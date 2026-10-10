package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows"

	"hpf_wails_vpk_manager/internal/models"
)

type LauncherService struct {
	detectorService *DetectorService
	storageService  *StorageService
}

func NewLauncherService(detector *DetectorService, storage *StorageService) *LauncherService {
	return &LauncherService{
		detectorService: detector,
		storageService:  storage,
	}
}

// IsSteamClient 判断指定的 L4D2 目录是否来源于 Steam
func (s *LauncherService) IsSteamClient(directory string) bool {
	if directory == "" {
		return false
	}
	clean := filepath.Clean(directory)
	cleanLower := strings.ToLower(clean)

	// 1. 路径中包含 steamapps/common 或 steamapps\common
	if strings.Contains(cleanLower, "steamapps\\common") || strings.Contains(cleanLower, "steamapps/common") {
		return true
	}

	// 2. 检查父级的 steamapps 目录下是否存在 appmanifest_550.acf
	parent := filepath.Dir(clean)
	grandParent := filepath.Dir(parent)
	manifest := filepath.Join(grandParent, "appmanifest_550.acf")
	if fi, err := os.Stat(manifest); err == nil && !fi.IsDir() {
		return true
	}

	return false
}

// BuildLaunchArgs 构建本地启动游戏的命令行参数
func (s *LauncherService) BuildLaunchArgs(settings models.LaunchSettings) []string {
	var args []string
	if settings.NoVid {
		args = append(args, "-novid")
	}
	if settings.ChineseLang {
		args = append(args, "-language", "schinese")
	}
	if trimmed := strings.TrimSpace(settings.CustomArgs); trimmed != "" {
		args = append(args, strings.Fields(trimmed)...)
	}
	return args
}

// LaunchGame 启动游戏（若为 Steam 来源且未开启直接启动模式则唤起 Steam，否则直接运行主程序）
func (s *LauncherService) LaunchGame(directory string) models.ApiResult {
	if directory == "" {
		return models.ApiResult{
			Success: false,
			Error:   "请先选择 L4D2 游戏安装目录",
		}
	}

	valid, finalDir, exePath := s.detectorService.ValidateL4D2Directory(directory)
	if !valid {
		return models.ApiResult{
			Success: false,
			Error:   fmt.Sprintf("目录不是有效的 Left 4 Dead 2 客户端: %s", directory),
		}
	}

	settings := models.DefaultLaunchSettings()
	if s.storageService != nil {
		settings = s.storageService.GetLaunchSettings()
	}

	// 1. 若为 Steam 客户端来源且未指定优先本地直接运行，通过 Steam 协议启动 (AppID: 550)
	if s.IsSteamClient(finalDir) && !settings.DirectLaunch {
		verbPtr, err := windows.UTF16PtrFromString("open")
		if err != nil {
			return models.ApiResult{Success: false, Error: err.Error()}
		}
		urlPtr, err := windows.UTF16PtrFromString("steam://rungameid/550")
		if err != nil {
			return models.ApiResult{Success: false, Error: err.Error()}
		}

		err = windows.ShellExecute(0, verbPtr, urlPtr, nil, nil, windows.SW_SHOWNORMAL)
		if err != nil {
			// 若 ShellExecute 失败，回退为直接启动主程序
			return s.launchDirect(finalDir, exePath, settings)
		}

		return models.ApiResult{
			Success: true,
			Message: "已通过 Steam 唤起启动 Left 4 Dead 2",
		}
	}

	// 2. 若为非 Steam 客户端或开启了优先本地直接运行，直接运行游戏主程序
	return s.launchDirect(finalDir, exePath, settings)
}

func (s *LauncherService) launchDirect(dir, exePath string, settings models.LaunchSettings) models.ApiResult {
	if fi, err := os.Stat(exePath); err != nil || fi.IsDir() {
		return models.ApiResult{
			Success: false,
			Error:   fmt.Sprintf("未找到游戏主程序: %s", exePath),
		}
	}

	args := s.BuildLaunchArgs(settings)

	cmd := exec.Command(exePath, args...)
	cmd.Dir = dir
	cmd.SysProcAttr = &windows.SysProcAttr{
		CreationFlags: windows.DETACHED_PROCESS,
	}

	if err := cmd.Start(); err != nil {
		return models.ApiResult{
			Success: false,
			Error:   fmt.Sprintf("启动游戏失败: %v", err),
		}
	}

	argMsg := ""
	if len(args) > 0 {
		argMsg = fmt.Sprintf(" (启动参数: %s)", strings.Join(args, " "))
	}

	return models.ApiResult{
		Success: true,
		Message: fmt.Sprintf("已成功启动游戏: %s%s", filepath.Base(exePath), argMsg),
	}
}
