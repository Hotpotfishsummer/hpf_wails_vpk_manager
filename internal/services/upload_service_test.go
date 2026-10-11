package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"hpf_wails_vpk_manager/internal/models"
)

func TestUploadService_EmptyFiles(t *testing.T) {
	svc := NewUploadService(nil, nil)
	res := svc.Upload(context.Background(), nil, "/remote/addons", models.SftpSettings{}, false, nil)
	if !res.Success {
		t.Errorf("expected success for empty file list, got error: %s", res.Error)
	}
	if res.UploadedCount != 0 {
		t.Errorf("expected uploaded count 0, got %d", res.UploadedCount)
	}
}

func TestUploadService_ResolveBaseName(t *testing.T) {
	tempDir := t.TempDir()
	notes := NewNotesService(tempDir)
	meta := NewVpkMetadataService(filepath.Join(tempDir, "cache"))
	meta.SetNotesService(notes)
	svc := NewUploadService(meta, nil, notes)

	localPath := filepath.Join(tempDir, "123456789.vpk")
	_ = os.WriteFile(localPath, []byte("fake"), 0644)
	_ = notes.SetNote(localPath, "好玩地图")

	// 不开启 betterNaming
	base := svc.resolveBaseName(localPath, false)
	if base != "123456789" {
		t.Errorf("resolveBaseName(betterNaming=false) = %q, want '123456789'", base)
	}

	// 开启 betterNaming
	baseBetter := svc.resolveBaseName(localPath, true)
	if baseBetter != "好玩地图-123456789" {
		t.Errorf("resolveBaseName(betterNaming=true) = %q, want '好玩地图-123456789'", baseBetter)
	}

	// disabled 文件
	disabledPath := filepath.Join(tempDir, "mymod.vpk.disabled")
	baseDisabled := svc.resolveBaseName(disabledPath, false)
	if baseDisabled != "mymod" {
		t.Errorf("resolveBaseName(disabled) = %q, want 'mymod'", baseDisabled)
	}
}

func TestUploadService_ContextCanceled(t *testing.T) {
	tempDir := t.TempDir()
	svc := NewUploadService(nil, NewSftpService())

	testFile := filepath.Join(tempDir, "test.vpk")
	_ = os.WriteFile(testFile, []byte("test content"), 0644)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 立刻取消

	var reportedStages []string
	res := svc.Upload(ctx, []string{testFile}, "/addons", models.SftpSettings{}, false, func(evt models.UploadProgressEvent) {
		reportedStages = append(reportedStages, evt.Stage)
	})

	if !res.Canceled {
		t.Errorf("expected res.Canceled to be true")
	}
	if res.Success {
		t.Errorf("expected res.Success to be false when canceled")
	}

	// 验证最终上报的状态为 canceled
	if len(reportedStages) == 0 {
		t.Errorf("expected at least one progress event")
	} else if reportedStages[len(reportedStages)-1] != "canceled" {
		t.Errorf("expected last stage to be 'canceled', got %q", reportedStages[len(reportedStages)-1])
	}
}
