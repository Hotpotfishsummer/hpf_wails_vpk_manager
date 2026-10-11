package services

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
)

type UploadProgressCallback func(evt models.UploadProgressEvent)

type UploadService struct {
	metaService  *VpkMetadataService
	sftpService  *SftpService
	notesService *NotesService
}

func NewUploadService(meta *VpkMetadataService, sftp *SftpService, notes ...*NotesService) *UploadService {
	var n *NotesService
	if len(notes) > 0 {
		n = notes[0]
	}
	return &UploadService{
		metaService:  meta,
		sftpService:  sftp,
		notesService: n,
	}
}

type fileUploadTask struct {
	localPath  string
	remotePath string
	size       int64
}

type vpkUploadItem struct {
	mainTask    fileUploadTask
	imageTasks  []fileUploadTask
	totalSize   int64
	displayName string
}

func (s *UploadService) Upload(
	ctx context.Context,
	localFilePaths []string,
	remoteAddonsPath string,
	config models.SftpSettings,
	betterNaming bool,
	progressCb UploadProgressCallback,
) models.UploadResult {
	if ctx == nil {
		ctx = context.Background()
	}

	totalFiles := len(localFilePaths)
	if totalFiles == 0 {
		return models.UploadResult{
			ApiResult: models.ApiResult{Success: true},
		}
	}

	// 1. 预先解析所有任务并统计总字节大小
	var items []vpkUploadItem
	var overallTotal int64

	for _, lPath := range localFilePaths {
		isDisabled := strings.HasSuffix(strings.ToLower(lPath), ".vpk.disabled")
		baseName := s.resolveBaseName(lPath, betterNaming)

		suffix := ".vpk"
		if isDisabled {
			suffix = ".vpk.disabled"
		}

		rPath := path.Join(remoteAddonsPath, baseName+suffix)
		var mainSize int64
		if fi, err := os.Stat(lPath); err == nil {
			mainSize = fi.Size()
		}

		item := vpkUploadItem{
			mainTask: fileUploadTask{
				localPath:  lPath,
				remotePath: rPath,
				size:       mainSize,
			},
			totalSize:   mainSize,
			displayName: baseName + suffix,
		}

		// 检查关联图片
		uploadedExts := make(map[string]bool)
		matches := naming.FindAssociatedImages(filepath.Dir(lPath), filepath.Base(lPath), true)
		for _, m := range matches {
			extLower := strings.ToLower(m.Extension)
			if !uploadedExts[extLower] {
				uploadedExts[extLower] = true
				var imgSize int64
				if fi, err := os.Stat(m.FullPath); err == nil {
					imgSize = fi.Size()
				}
				imgTask := fileUploadTask{
					localPath:  m.FullPath,
					remotePath: path.Join(remoteAddonsPath, baseName+m.Extension),
					size:       imgSize,
				}
				item.imageTasks = append(item.imageTasks, imgTask)
				item.totalSize += imgSize
			}
		}

		overallTotal += item.totalSize
		items = append(items, item)
	}

	// 2. 状态变量
	var (
		overallBytes  int64
		uploadedCount int
		failedCount   int
		failedFiles   []models.UploadFailedItem
		canceled      bool
		lastTime      = time.Now()
		lastBytes     int64
		currentSpeed  int64
		mu            sync.Mutex
	)

	report := func(currIdx int, currFile string, currBytes, currTotal int64, stage string) {
		if progressCb == nil {
			return
		}
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		elapsed := now.Sub(lastTime).Seconds()
		if elapsed >= 0.5 {
			bytesDiff := overallBytes - lastBytes
			if bytesDiff < 0 {
				bytesDiff = 0
			}
			instSpeed := int64(float64(bytesDiff) / elapsed)
			if currentSpeed == 0 {
				currentSpeed = instSpeed
			} else {
				// EMA 平滑处理网速波动
				currentSpeed = int64(float64(currentSpeed)*0.7 + float64(instSpeed)*0.3)
			}
			lastBytes = overallBytes
			lastTime = now
		}

		evt := models.UploadProgressEvent{
			CurrentIndex: currIdx,
			TotalFiles:   totalFiles,
			CurrentFile:  currFile,
			CurrentBytes: currBytes,
			TotalBytes:   currTotal,
			OverallBytes: overallBytes,
			OverallTotal: overallTotal,
			Speed:        currentSpeed,
			SuccessCount: uploadedCount,
			FailedCount:  failedCount,
			Stage:        stage,
		}
		progressCb(evt)
	}

	// 发送初始准备中事件
	report(0, "", 0, 0, "preparing")

	// 3. 逐个文件上传
	for i, item := range items {
		select {
		case <-ctx.Done():
			canceled = true
			break
		default:
		}
		if canceled {
			break
		}

		currIdx := i + 1
		var itemBytesTransferred int64
		itemTotal := item.totalSize
		vpkFailed := false

		// 上传主 VPK 文件
		report(currIdx, item.displayName, 0, itemTotal, "uploading")

		var lastMainTransferred int64
		err := s.sftpService.UploadFileWithProgress(ctx, config, item.mainTask.localPath, item.mainTask.remotePath, func(transferred, total int64) {
			delta := transferred - lastMainTransferred
			lastMainTransferred = transferred

			mu.Lock()
			overallBytes += delta
			itemBytesTransferred += delta
			mu.Unlock()

			report(currIdx, item.displayName, itemBytesTransferred, itemTotal, "uploading")
		})

		if err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				canceled = true
				break
			}
			vpkFailed = true
			failedCount++
			failedFiles = append(failedFiles, models.UploadFailedItem{
				File:  filepath.Base(item.mainTask.localPath),
				Error: err.Error(),
			})
			// 补偿未传完的字节计算，保持总体进度条一致
			remainingItemBytes := itemTotal - itemBytesTransferred
			if remainingItemBytes > 0 {
				mu.Lock()
				overallBytes += remainingItemBytes
				mu.Unlock()
			}
			continue
		}

		// 上传关联图片（如果有）
		for _, imgTask := range item.imageTasks {
			select {
			case <-ctx.Done():
				canceled = true
				break
			default:
			}
			if canceled {
				break
			}

			var lastImgTransferred int64
			_ = s.sftpService.UploadFileWithProgress(ctx, config, imgTask.localPath, imgTask.remotePath, func(transferred, total int64) {
				delta := transferred - lastImgTransferred
				lastImgTransferred = transferred

				mu.Lock()
				overallBytes += delta
				itemBytesTransferred += delta
				mu.Unlock()

				report(currIdx, item.displayName, itemBytesTransferred, itemTotal, "uploading")
			})
		}

		if canceled {
			break
		}

		if !vpkFailed {
			uploadedCount++
		}
	}

	// 4. 发送最终状态事件
	finalStage := "completed"
	if canceled {
		finalStage = "canceled"
	} else if uploadedCount == 0 && failedCount > 0 {
		finalStage = "error"
	}
	report(totalFiles, "", overallTotal, overallTotal, finalStage)

	return models.UploadResult{
		ApiResult: models.ApiResult{
			Success: !canceled && (failedCount == 0 || uploadedCount > 0),
			Error: func() string {
				if canceled {
					return "上传已被用户取消"
				}
				if failedCount > 0 && uploadedCount == 0 {
					return "所有文件上传均失败"
				}
				return ""
			}(),
		},
		UploadedCount: uploadedCount,
		TotalCount:    totalFiles,
		FailedFiles:   failedFiles,
		Canceled:      canceled,
	}
}

func (s *UploadService) resolveBaseName(localPath string, betterNaming bool) string {
	name := filepath.Base(localPath)
	lower := strings.ToLower(name)

	var fallback string
	if strings.HasSuffix(lower, ".vpk.disabled") {
		fallback = name[:len(name)-len(".vpk.disabled")]
	} else {
		fallback = strings.TrimSuffix(name, filepath.Ext(name))
	}

	if !betterNaming {
		return fallback
	}

	title := s.metaService.GetAddonTitle(localPath)
	var titleStr string
	if title != nil {
		titleStr = *title
	}
	var note string
	if s.notesService != nil {
		note = s.notesService.GetNote(localPath)
	}
	return naming.BuildExportBaseNameWithNote(titleStr, name, fallback, note)
}
