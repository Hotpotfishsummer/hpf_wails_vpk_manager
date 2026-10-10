package services

import (
	"strings"
	"testing"

	"hpf_wails_vpk_manager/internal/models"
)

func TestLauncherService(t *testing.T) {
	detector := NewDetectorService()
	launcher := NewLauncherService(detector, nil)

	// 测试 Steam 路径判断
	steamPath := `C:\SteamLibrary\steamapps\common\Left 4 Dead 2`
	isSteam := launcher.IsSteamClient(steamPath)
	t.Logf("IsSteamClient(%s) = %v", steamPath, isSteam)
	if !isSteam {
		t.Errorf("Expected true for steamPath, got %v", isSteam)
	}

	nonSteamPath := `D:\Games\Left4Dead2`
	isNonSteam := launcher.IsSteamClient(nonSteamPath)
	t.Logf("IsSteamClient(%s) = %v", nonSteamPath, isNonSteam)
	if isNonSteam {
		t.Errorf("Expected false for nonSteamPath, got %v", isNonSteam)
	}

	// 测试启动参数构造
	settings := models.LaunchSettings{
		NoVid:        true,
		ChineseLang:  true,
		DirectLaunch: false,
		CustomArgs:   "-console -windowed -noborder",
	}
	args := launcher.BuildLaunchArgs(settings)
	argsStr := strings.Join(args, " ")
	expectedPrefix := "-novid -language schinese -console -windowed -noborder"
	if argsStr != expectedPrefix {
		t.Errorf("Expected args '%s', got '%s'", expectedPrefix, argsStr)
	}
}
