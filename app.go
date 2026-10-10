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
	notesService     *services.NotesService
	metaService      *services.VpkMetadataService
	scanService      *services.ScanService
	exportService    *services.ExportService
	importService    *services.ImportService
	localizeService  *services.LocalizeService
	sftpService      *services.SftpService
	uploadService    *services.UploadService
	downloadService  *services.DownloadService
	restarterService *services.RestarterService
	detectorService  *services.DetectorService
	launcherService  *services.LauncherService
	groupService     *services.GroupService
	aiService        *services.AiService
}

// NewApp creates a new App application struct
func NewApp() *App {
	storage := services.NewStorageService()
	config := services.NewConfigService()
	notes := services.NewNotesService()
	group := services.NewGroupService()
	ai := services.NewAiService()
	meta := services.NewVpkMetadataService()
	meta.SetNotesService(notes)
	scan := services.NewScanService(meta)
	export := services.NewExportService(meta, notes)
	importSvc := services.NewImportService()
	localize := services.NewLocalizeService(meta, notes)
	sftp := services.NewSftpService()
	upload := services.NewUploadService(meta, sftp, notes)
	download := services.NewDownloadService(sftp)
	restarter := services.NewRestarterService()
	detector := services.NewDetectorService()
	launcher := services.NewLauncherService(detector, storage)

	return &App{
		storageService:   storage,
		configService:    config,
		notesService:     notes,
		groupService:     group,
		aiService:        ai,
		metaService:      meta,
		scanService:      scan,
		exportService:    export,
		importService:    importSvc,
		localizeService:  localize,
		sftpService:      sftp,
		uploadService:    upload,
		downloadService:  download,
		restarterService: restarter,
		detectorService:  detector,
		launcherService:  launcher,
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

// DetectL4D2Directories automatically searches for Left 4 Dead 2 directories across processes, registry, shortcuts, and disks
func (a *App) DetectL4D2Directories(deepScan bool) models.DetectGameResult {
	clients := a.detectorService.DetectClients(deepScan)
	return models.DetectGameResult{
		ApiResult: models.ApiResult{
			Success: true,
		},
		Clients: clients,
	}
}

// ValidateL4D2Directory checks if a path is a valid L4D2 directory
func (a *App) ValidateL4D2Directory(dir string) bool {
	valid, _, _ := a.detectorService.ValidateL4D2Directory(dir)
	return valid
}

// LaunchGame launches the game via Steam (if Steam client) or directly via executable
func (a *App) LaunchGame() models.ApiResult {
	l4d2Dir := a.storageService.GetL4D2Directory()
	return a.launcherService.LaunchGame(l4d2Dir)
}

// IsSteamClient checks if a directory or current directory is from Steam
func (a *App) IsSteamClient(directory string) bool {
	if directory == "" {
		directory = a.storageService.GetL4D2Directory()
	}
	return a.launcherService.IsSteamClient(directory)
}

// GetLaunchSettings returns game startup options
func (a *App) GetLaunchSettings() models.LaunchSettings {
	return a.storageService.GetLaunchSettings()
}

// SaveLaunchSettings saves game startup options
func (a *App) SaveLaunchSettings(settings models.LaunchSettings) models.ApiResult {
	a.storageService.SaveLaunchSettings(settings)
	return models.ApiResult{Success: true}
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

// ImportFiles imports chosen archive, VPK, or folder files directly into the addons directory
func (a *App) ImportFiles(filePaths []string) models.ImportResult {
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

	if len(filePaths) == 0 {
		return models.ImportResult{
			ApiResult: models.ApiResult{Success: false, Error: "未检测到拖入的文件"},
		}
	}

	return a.importService.ImportCompressedFiles(filePaths, addonsDir)
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

// SetVpkNote sets a custom note for a local or workshop VPK file
func (a *App) SetVpkNote(filePath, note string) models.ApiResult {
	if filePath == "" {
		return models.ApiResult{Success: false, Error: "文件路径不能为空"}
	}
	if err := a.notesService.SetNote(filePath, note); err != nil {
		return models.ApiResult{Success: false, Error: err.Error()}
	}
	return models.ApiResult{Success: true, Message: "备注保存成功"}
}

// GetVpkNote retrieves the custom note for a VPK file
func (a *App) GetVpkNote(filePath string) string {
	return a.notesService.GetNote(filePath)
}

// GetAllVpkNotes retrieves all notes mapped by key
func (a *App) GetAllVpkNotes() map[string]string {
	return a.notesService.GetAllNotes()
}

// GetVpkGroups retrieves all saved VPK groups
func (a *App) GetVpkGroups() []models.VpkGroup {
	return a.groupService.GetGroups()
}

// SaveVpkGroups persists all VPK groups
func (a *App) SaveVpkGroups(groups []models.VpkGroup) models.ApiResult {
	if err := a.groupService.SaveGroups(groups); err != nil {
		return models.ApiResult{Success: false, Error: err.Error()}
	}
	return models.ApiResult{Success: true, Message: "分组保存成功"}
}

// GetAiSettings retrieves the saved OpenAI API settings
func (a *App) GetAiSettings() models.AiSettings {
	return a.storageService.GetAiSettings()
}

// SaveAiSettings persists the OpenAI API settings
func (a *App) SaveAiSettings(settings models.AiSettings) models.ApiResult {
	a.storageService.SaveAiSettings(settings)
	return models.ApiResult{Success: true, Message: "AI 配置保存成功"}
}

// TestAiConnection tests the OpenAI API configuration
func (a *App) TestAiConnection(settings models.AiSettings) models.ApiResult {
	return a.aiService.TestConnection(settings)
}

// PredictGroupName predicts a concise group name for selected VPK files
func (a *App) PredictGroupName(files []models.VpkFile) models.GroupPredictResult {
	settings := a.storageService.GetAiSettings()
	return a.aiService.PredictGroupName(files, settings)
}

// AiChat performs a generic chat completion with user-configured OpenAI settings
func (a *App) AiChat(req models.AiChatRequest) models.AiResult {
	settings := a.storageService.GetAiSettings()
	reply, err := a.aiService.ChatCompletion(a.ctx, req, settings)
	if err != nil {
		return models.AiResult{
			ApiResult: models.ApiResult{Success: false, Error: err.Error()},
		}
	}
	return models.AiResult{
		ApiResult: models.ApiResult{Success: true},
		Content:   reply,
	}
}

