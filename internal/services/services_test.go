package services

import (
	"os"
	"path/filepath"
	"testing"

	"hpf_wails_vpk_manager/internal/models"
)

func TestStorageService(t *testing.T) {
	tempDir := t.TempDir()
	storage := &StorageService{
		storageDir:  tempDir,
		storageFile: filepath.Join(tempDir, "preferences.json"),
		data:        make(map[string]any),
	}

	storage.SaveL4D2Directory("D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2")
	if dir := storage.GetL4D2Directory(); dir != "D:\\SteamLibrary\\steamapps\\common\\Left 4 Dead 2" {
		t.Errorf("GetL4D2Directory = %q", dir)
	}

	storage.SaveSftpSettings(models.SftpSettings{
		Address:  "127.0.0.1:22",
		Username: "testuser",
	})
	cfg := storage.GetSftpSettings()
	if cfg.Address != "127.0.0.1:22" || cfg.Username != "testuser" {
		t.Errorf("GetSftpSettings = %+v", cfg)
	}
}

func TestConfigService(t *testing.T) {
	tempDir := t.TempDir()
	l4d2Dir := filepath.Join(tempDir, "L4D2")
	cfgDir := filepath.Join(l4d2Dir, "left4dead2", "cfg")
	_ = os.MkdirAll(cfgDir, 0755)

	autoexecPath := filepath.Join(cfgDir, "autoexec.cfg")
	_ = os.WriteFile(autoexecPath, []byte("// existing old config"), 0644)

	configSvc := NewConfigService()
	ok := configSvc.ApplyAutoexecTemplate(l4d2Dir)
	if !ok {
		t.Fatalf("ApplyAutoexecTemplate failed")
	}

	// Verify backup was created
	bakPath := autoexecPath + ".bak"
	if _, err := os.Stat(bakPath); err != nil {
		t.Errorf("backup file %s not found", bakPath)
	}

	// Verify autoexec was updated
	data, _ := os.ReadFile(autoexecPath)
	if string(data) != autoexecTemplate {
		t.Errorf("autoexec content does not match template")
	}
}

func TestExportAndImport(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "src")
	_ = os.MkdirAll(srcDir, 0755)

	vpk1 := filepath.Join(srcDir, "mod1.vpk")
	img1 := filepath.Join(srcDir, "mod1.jpg")
	_ = os.WriteFile(vpk1, []byte("fake vpk content"), 0644)
	_ = os.WriteFile(img1, []byte("fake image content"), 0644)

	metaSvc := NewVpkMetadataService(filepath.Join(tempDir, "cache"))
	exportSvc := NewExportService(metaSvc)

	outDir := filepath.Join(tempDir, "out")
	res := exportSvc.ExportVpks([]string{vpk1}, outDir, false, true, nil)
	if !res.Success {
		t.Fatalf("ExportVpks failed: %s", res.Error)
	}

	if _, err := os.Stat(res.FilePath); err != nil {
		t.Fatalf("exported archive %s does not exist", res.FilePath)
	}

	// Test Import
	targetAddons := filepath.Join(tempDir, "target_addons")
	_ = os.MkdirAll(targetAddons, 0755)

	importSvc := NewImportService()
	impRes := importSvc.ImportCompressedFiles([]string{res.FilePath}, targetAddons)
	if !impRes.Success {
		t.Fatalf("ImportCompressedFiles failed: %+v", impRes.Errors)
	}
	if impRes.ImportedCount != 1 {
		t.Errorf("ImportedCount = %d, want 1", impRes.ImportedCount)
	}

	// Check extracted files
	extractedVpk := filepath.Join(targetAddons, "mod1.vpk")
	extractedImg := filepath.Join(targetAddons, "mod1.jpg")

	if _, err := os.Stat(extractedVpk); err != nil {
		t.Errorf("extracted vpk missing: %s", extractedVpk)
	}
	if _, err := os.Stat(extractedImg); err != nil {
		t.Errorf("extracted image missing: %s", extractedImg)
	}
}
