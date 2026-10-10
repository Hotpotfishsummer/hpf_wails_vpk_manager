package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
	"hpf_wails_vpk_manager/internal/services"
)

// App struct
type App struct {
	ctx              context.Context
	storageService   *services.StorageService
	configService    *services.ConfigService
	metaService      *services.VpkMetadataService
	scanService      *services.ScanService
	exportService    *services.ExportService
	importService    *services.ImportService
	localizeService  *services.LocalizeService
	sftpService      *services.SftpService
	uploadService    *services.UploadService
	downloadService  *services.DownloadService
	restarterService *services.RestarterService
}

// NewApp creates a new App application struct
func NewApp() *App {
	storage := services.NewStorageService()
	config := services.NewConfigService()
	meta := services.NewVpkMetadataService()
	scan := services.NewScanService(meta)
	export := services.NewExportService(meta)
	importSvc := services.NewImportService()
	localize := services.NewLocalizeService(meta)
	sftp := services.NewSftpService()
	upload := services.NewUploadService(meta, sftp)
	download := services.NewDownloadService(sftp)
	restarter := services.NewRestarterService()

	return &App{
		storageService:   storage,
		configService:    config,
		metaService:      meta,
		scanService:      scan,
		exportService:    export,
		importService:    importSvc,
		localizeService:  localize,
		sftpService:      sftp,
		uploadService:    upload,
		downloadService:  download,
		restarterService: restarter,
	}
}

// startup is called when the app starts. The context is saved so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	a.sftpService.Close()
}

// GetLastDirectory returns the last chosen L4D2 directory
func (a *App) GetLastDirectory() string {
	return a.storageService.GetL4D2Directory()
}

// BrowseDirectory opens a native directory dialog
func (a *App) BrowseDirectory() string {
	initial := a.storageService.GetL4D2Directory()
	if initial == "" {
		if home, err := os.UserHomeDir(); err == nil {
			initial = home
		}
	}

	dir, err := wailsRuntime.OpenDirectoryDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title:            "选择 L4D2 安装目录",
		DefaultDirectory: initial,
	})
	if err != nil {
		log.Printf("OpenDirectoryDialog failed: %v", err)
		return ""
	}
	return dir
}

// GetDownloadsDirectory returns user's downloads directory
func (a *App) GetDownloadsDirectory() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	downloads := filepath.Join(home, "Downloads")
	if info, err := os.Stat(downloads); err == nil && info.IsDir() {
		return downloads
	}
	return home
}

// ScanVpkFiles scans directory for local and workshop VPKs
func (a *App) ScanVpkFiles(directory string) models.ScanResult {
	if directory == "" {
		return models.ScanResult{
			ApiResult: models.ApiResult{Success: false, Error: "未指定目录"},
		}
	}

	if _, err := os.Stat(directory); err != nil {
		return models.ScanResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("目录不存在: %s", directory)},
		}
	}

	a.storageService.SaveL4D2Directory(directory)
	a.configService.ApplyAutoexecTemplate(directory)

	return a.scanService.Scan(directory)
}

// ExportVpks exports chosen VPKs to a tar.zst archive
func (a *App) ExportVpks(filePaths []string, outputDir string, betterNaming bool, overwrite bool) models.ExportResult {
	progressCb := func(stage string, done, total int) {
		wailsRuntime.EventsEmit(a.ctx, "export-progress", map[string]any{
			"stage": stage,
			"done":  done,
			"total": total,
		})
	}

	return a.exportService.ExportVpks(filePaths, outputDir, betterNaming, overwrite, progressCb)
}

// ImportCompressedVpks opens a native file dialog and imports archives
func (a *App) ImportCompressedVpks() models.ImportResult {
	l4d2Dir := a.storageService.GetL4D2Directory()
	if l4d2Dir == "" {
		return models.ImportResult{
			ApiResult: models.ApiResult{Success: false, Error: "请先选择 L4D2 安装目录"},
		}
	}

	addonsDir := filepath.Join(l4d2Dir, "left4dead2", "addons")
	if _, err := os.Stat(addonsDir); err != nil {
		return models.ImportResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("addons 目录不存在: %s", addonsDir)},
		}
	}

	files, err := wailsRuntime.OpenMultipleFilesDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: "选择 VPK 压缩包",
		Filters: []wailsRuntime.FileFilter{
			{
				DisplayName: "Addon Archives (*.zip;*.7z;*.zst)",
				Pattern:     "*.zip;*.7z;*.zst;*.tar.zst",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})
	if err != nil || len(files) == 0 {
		return models.ImportResult{
			ApiResult: models.ApiResult{Success: true, Message: "用户取消了选择"},
		}
	}

	return a.importService.ImportCompressedFiles(files, addonsDir)
}

// LocalizeWorkshopVpks copies workshop VPKs to local addons and moves originals to .trash
func (a *App) LocalizeWorkshopVpks(filePaths []string, betterNaming bool) models.LocalizeResult {
	l4d2Dir := a.storageService.GetL4D2Directory()
	if l4d2Dir == "" {
		return models.LocalizeResult{
			ApiResult: models.ApiResult{Success: false, Error: "未找到 L4D2 目录"},
		}
	}
	return a.localizeService.Localize(filePaths, l4d2Dir, betterNaming)
}

// OpenDirectory opens a directory in Windows Explorer
func (a *App) OpenDirectory(dirPath string) bool {
	if dirPath == "" {
		return false
	}
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	target := dirPath
	if !info.IsDir() {
		target = filepath.Dir(dirPath)
	}

	cmd := exec.Command("explorer", target)
	return cmd.Start() == nil
}

func disableOrDeleteSingle(filePath string, disable bool) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	base := strings.TrimSuffix(filePath, filepath.Ext(filePath))

	if disable {
		var newPath string
		if strings.HasSuffix(strings.ToLower(filePath), ".disabled") {
			newPath = filePath[:len(filePath)-len(".disabled")]
		} else {
			newPath = filePath + ".disabled"
		}

		newBase := strings.TrimSuffix(newPath, filepath.Ext(newPath))
		for _, ext := range naming.AssociatedSyncExtensions {
			oldImg := base + ext
			newImg := newBase + ext
			if _, err := os.Stat(oldImg); err == nil {
				_ = os.Remove(newImg)
				_ = os.Rename(oldImg, newImg)
			}
		}

		_ = os.Remove(newPath)
		return os.Rename(filePath, newPath)
	}

	// Delete
	for _, ext := range naming.AssociatedSyncExtensions {
		img := base + ext
		if _, err := os.Stat(img); err == nil {
			_ = os.Remove(img)
		}
	}
	return os.Remove(filePath)
}

// DeleteVpkFile deletes or disables a single VPK file
func (a *App) DeleteVpkFile(filePath string, disable bool) models.ApiResult {
	if err := disableOrDeleteSingle(filePath, disable); err != nil {
		return models.ApiResult{Success: false, Error: err.Error()}
	}
	action := "删除"
	if disable {
		action = "禁用"
	}
	return models.ApiResult{Success: true, Message: fmt.Sprintf("文件%s成功", action)}
}

// DeleteMultipleVpkFiles deletes or disables multiple VPK files
func (a *App) DeleteMultipleVpkFiles(filePaths []string, disable bool) models.BatchDeleteResult {
	successCount := 0
	var errors []string

	for _, p := range filePaths {
		if err := disableOrDeleteSingle(p, disable); err != nil {
			errors = append(errors, fmt.Sprintf("处理 %s 失败: %v", filepath.Base(p), err))
		} else {
			successCount++
		}
	}

	return models.BatchDeleteResult{
		ApiResult:    models.ApiResult{Success: len(errors) == 0},
		SuccessCount: successCount,
		ErrorCount:   len(errors),
		Errors:       errors,
	}
}

// EnableVpkFile enables a disabled VPK file
func (a *App) EnableVpkFile(filePath string) models.ApiResult {
	if _, err := os.Stat(filePath); err != nil {
		return models.ApiResult{Success: false, Error: fmt.Sprintf("文件不存在: %s", filePath)}
	}

	if !strings.HasSuffix(strings.ToLower(filePath), ".disabled") {
		return models.ApiResult{Success: false, Error: "文件未被禁用"}
	}

	newPath := filePath[:len(filePath)-len(".disabled")]
	if _, err := os.Stat(newPath); err == nil {
		return models.ApiResult{Success: false, Error: fmt.Sprintf("目标文件已存在: %s", newPath)}
	}

	if err := os.Rename(filePath, newPath); err != nil {
		return models.ApiResult{Success: false, Error: err.Error()}
	}

	return models.ApiResult{Success: true, Message: "文件启用成功"}
}

// ResetAppData resets app data keeping preferences.json and exits
func (a *App) ResetAppData() bool {
	dir := a.storageService.GetStorageDir()
	entries, err := os.ReadDir(dir)
	if err == nil {
		for _, entry := range entries {
			if entry.Name() != "preferences.json" {
				_ = os.RemoveAll(filepath.Join(dir, entry.Name()))
			}
		}
	}
	wailsRuntime.Quit(a.ctx)
	return true
}

// LogFromJs logs a message from the JavaScript frontend
func (a *App) LogFromJs(level, message string) {
	log.Printf("[Frontend/%s] %s", level, message)
}

// GetSftpSettings returns saved SFTP configuration
func (a *App) GetSftpSettings() models.SftpSettings {
	return a.storageService.GetSftpSettings()
}

// SaveSftpSettings saves SFTP configuration
func (a *App) SaveSftpSettings(settings models.SftpSettings) models.ApiResult {
	a.storageService.SaveSftpSettings(settings)
	return models.ApiResult{Success: true}
}

// ListRemoteVpks lists VPKs on remote SFTP server
func (a *App) ListRemoteVpks() models.RemoteListResult {
	files, err := a.sftpService.ListVpkFiles(a.storageService.GetSftpSettings())
	if err != nil {
		return models.RemoteListResult{
			ApiResult: models.ApiResult{Success: false, Error: err.Error()},
		}
	}
	return models.RemoteListResult{
		ApiResult: models.ApiResult{Success: true},
		Files:     files,
	}
}

// DeleteRemoteVpk deletes a VPK on remote SFTP server
func (a *App) DeleteRemoteVpk(filePath string) models.ApiResult {
	err := a.sftpService.DeleteFileWithImages(a.storageService.GetSftpSettings(), filePath)
	if err != nil {
		return models.ApiResult{Success: false, Error: err.Error()}
	}
	return models.ApiResult{Success: true}
}

// DeleteMultipleRemoteVpks deletes multiple VPKs on remote SFTP server
func (a *App) DeleteMultipleRemoteVpks(filePaths []string) models.BatchDeleteResult {
	cfg := a.storageService.GetSftpSettings()
	successCount := 0
	var errors []string

	for _, p := range filePaths {
		if err := a.sftpService.DeleteFileWithImages(cfg, p); err != nil {
			errors = append(errors, p)
		} else {
			successCount++
		}
	}

	return models.BatchDeleteResult{
		ApiResult:    models.ApiResult{Success: len(errors) == 0},
		SuccessCount: successCount,
		ErrorCount:   len(errors),
		Errors:       errors,
	}
}

// RenameRemoteVpk renames a VPK on remote SFTP server
func (a *App) RenameRemoteVpk(oldPath, newPath string) models.ApiResult {
	err := a.sftpService.RenameFileWithImages(a.storageService.GetSftpSettings(), oldPath, newPath)
	if err != nil {
		return models.ApiResult{Success: false, Error: err.Error()}
	}
	return models.ApiResult{Success: true}
}

// DownloadRemoteVpks downloads VPKs from remote SFTP server
func (a *App) DownloadRemoteVpks(remoteFilePaths []string) models.DownloadResult {
	l4d2Dir := a.storageService.GetL4D2Directory()
	if l4d2Dir == "" {
		return models.DownloadResult{
			ApiResult: models.ApiResult{Success: false, Error: "未配置本地 L4D2 目录"},
		}
	}
	localAddonsDir := filepath.Join(l4d2Dir, "left4dead2", "addons")
	return a.downloadService.Download(remoteFilePaths, localAddonsDir, a.storageService.GetSftpSettings())
}

// UploadLocalVpks uploads local VPKs to remote SFTP server
func (a *App) UploadLocalVpks(localFilePaths []string, betterNaming bool) models.UploadResult {
	cfg := a.storageService.GetSftpSettings()
	if cfg.AddonsPath == "" {
		return models.UploadResult{
			ApiResult: models.ApiResult{Success: false, Error: "未配置远程 addons 路径"},
		}
	}
	return a.uploadService.Upload(localFilePaths, cfg.AddonsPath, cfg, betterNaming)
}

// OpenExternalUrl opens a URL in system default browser
func (a *App) OpenExternalUrl(url string) {
	wailsRuntime.BrowserOpenURL(a.ctx, url)
}

// RestartL4d2Server requests restart of remote L4D2 container service
func (a *App) RestartL4d2Server() models.RestartResult {
	return a.restarterService.RestartServer()
}
