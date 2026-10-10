package services

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/klauspost/compress/zstd"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
)

var illegalArchiveChars = regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)

type ProgressCallback func(stage string, done int, total int)

type ExportService struct {
	metaService  *VpkMetadataService
	notesService *NotesService
}

func NewExportService(meta *VpkMetadataService, notes ...*NotesService) *ExportService {
	var n *NotesService
	if len(notes) > 0 {
		n = notes[0]
	}
	return &ExportService{metaService: meta, notesService: n}
}

func (s *ExportService) GenerateArchiveName(filePaths []string, titles []string) string {
	var names []string
	if len(titles) > 0 {
		for _, t := range titles {
			if strings.TrimSpace(t) != "" {
				names = append(names, strings.TrimSpace(t))
			}
		}
	}

	for i := len(names); i < len(filePaths); i++ {
		stem := strings.TrimSuffix(filepath.Base(filePaths[i]), filepath.Ext(filePaths[i]))
		names = append(names, stem)
	}

	count := len(names)
	limit := 3
	if limit > count {
		limit = count
	}

	archiveName := strings.Join(names[:limit], "-")
	if count > 3 {
		archiveName = fmt.Sprintf("%s-and-%d-more", archiveName, count-3)
	}

	// Clean illegal chars
	archiveName = illegalArchiveChars.ReplaceAllString(archiveName, "_")
	archiveName = strings.ReplaceAll(archiveName, " ", "_")

	if len(archiveName) > 100 {
		archiveName = archiveName[:97] + "..."
	}

	if strings.Trim(archiveName, "._-") == "" || archiveName == "..." {
		archiveName = fmt.Sprintf("vpk_export_%d", time.Now().Unix())
	}

	return archiveName
}

func (s *ExportService) ExportVpks(
	filePaths []string,
	outputDir string,
	betterNaming bool,
	overwrite bool,
	progress ProgressCallback,
) models.ExportResult {
	if len(filePaths) == 0 {
		return models.ExportResult{
			ApiResult: models.ApiResult{Success: false, Error: "未选择要导出的文件"},
		}
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return models.ExportResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("创建输出目录失败: %v", err)},
		}
	}

	// Gather titles
	titles := make([]string, len(filePaths))
	for i, p := range filePaths {
		if t := s.metaService.GetAddonTitle(p); t != nil {
			titles[i] = *t
		}
	}

	archiveName := s.GenerateArchiveName(filePaths, titles)
	archivePath := filepath.Join(outputDir, archiveName+".tar.zst")

	if _, err := os.Stat(archivePath); err == nil && !overwrite {
		return models.ExportResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("目标文件已存在: %s", archivePath)},
			Exists:    true,
		}
	}

	tempDir, err := os.MkdirTemp(outputDir, ".vpk_temp_*")
	if err != nil {
		return models.ExportResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("创建临时目录失败: %v", err)},
		}
	}
	defer os.RemoveAll(tempDir)

	// Copy files to temp directory
	for i, srcPath := range filePaths {
		if progress != nil {
			progress("copying", i+1, len(filePaths))
		}

		info, err := os.Stat(srcPath)
		if err != nil || info.IsDir() {
			continue
		}

		filename := filepath.Base(srcPath)
		stem := strings.TrimSuffix(filename, filepath.Ext(filename))
		baseName := stem

		if betterNaming {
			var title string
			if i < len(titles) {
				title = titles[i]
			}
			var note string
			if s.notesService != nil {
				note = s.notesService.GetNote(srcPath)
			}
			baseName = naming.BuildExportBaseNameWithNote(title, filename, stem, note)
		}

		// Copy VPK file
		destVpk := filepath.Join(tempDir, baseName+".vpk")
		if err := copyFile(srcPath, destVpk); err != nil {
			return models.ExportResult{
				ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("复制文件失败: %v", err)},
			}
		}

		// Copy associated images
		for _, img := range naming.FindAssociatedImages(filepath.Dir(srcPath), filename, false) {
			destImg := filepath.Join(tempDir, baseName+img.Extension)
			_ = copyFile(img.FullPath, destImg)
		}
	}

	if progress != nil {
		progress("compressing", 0, 0)
	}

	// Create zstd tar archive
	outFile, err := os.Create(archivePath)
	if err != nil {
		return models.ExportResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("创建输出文件失败: %v", err)},
		}
	}
	defer outFile.Close()

	concurrency := runtime.NumCPU() - 1
	if concurrency < 1 {
		concurrency = 1
	}

	zstdWriter, err := zstd.NewWriter(outFile,
		zstd.WithEncoderConcurrency(concurrency),
		zstd.WithEncoderLevel(zstd.SpeedDefault),
	)
	if err != nil {
		return models.ExportResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("初始化 zstd 编码器失败: %v", err)},
		}
	}
	defer zstdWriter.Close()

	tarWriter := tar.NewWriter(zstdWriter)
	defer tarWriter.Close()

	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return models.ExportResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("读取临时目录失败: %v", err)},
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		itemPath := filepath.Join(tempDir, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			continue
		}
		header.Name = entry.Name()

		if err := tarWriter.WriteHeader(header); err != nil {
			return models.ExportResult{
				ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("写入 tar 头失败: %v", err)},
			}
		}

		f, err := os.Open(itemPath)
		if err != nil {
			continue
		}
		_, err = io.Copy(tarWriter, f)
		f.Close()
		if err != nil {
			return models.ExportResult{
				ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("写入归档数据失败: %v", err)},
			}
		}
	}

	_ = tarWriter.Close()
	_ = zstdWriter.Close()
	_ = outFile.Close()

	return models.ExportResult{
		ApiResult: models.ApiResult{
			Success: true,
			Message: fmt.Sprintf("导出完成: %s", archivePath),
		},
		FilePath: archivePath,
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
