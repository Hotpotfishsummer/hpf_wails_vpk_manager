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

	// 测试 LaunchSettings 默认值与 ServerConnectCmd 及 NetworkRatesCmd
	launchDef := storage.GetLaunchSettings()
	if launchDef.ServerConnectCmd != "connect 192.168.1.200" {
		t.Errorf("Expected default ServerConnectCmd to be 'connect 192.168.1.200', got %q", launchDef.ServerConnectCmd)
	}
	if launchDef.NetworkRatesCmd != models.DefaultNetworkRatesCmd {
		t.Errorf("Expected default NetworkRatesCmd to be %q, got %q", models.DefaultNetworkRatesCmd, launchDef.NetworkRatesCmd)
	}

	// 保存并读取自定义服务器连接指令与网络参数
	launchDef.ServerConnectCmd = "connect 10.0.0.1:27016"
	launchDef.NetworkRatesCmd = "cl_updaterate 60; cl_cmdrate 60; cl_interp 0.016667; cl_interp_ratio 1"
	storage.SaveLaunchSettings(launchDef)
	launchSaved := storage.GetLaunchSettings()
	if launchSaved.ServerConnectCmd != "connect 10.0.0.1:27016" {
		t.Errorf("Expected ServerConnectCmd to be 'connect 10.0.0.1:27016', got %q", launchSaved.ServerConnectCmd)
	}
	if launchSaved.NetworkRatesCmd != "cl_updaterate 60; cl_cmdrate 60; cl_interp 0.016667; cl_interp_ratio 1" {
		t.Errorf("Expected NetworkRatesCmd to be custom, got %q", launchSaved.NetworkRatesCmd)
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

	// Test Import direct VPK and folder
	directVpk := filepath.Join(tempDir, "direct.vpk")
	_ = os.WriteFile(directVpk, []byte("direct vpk content"), 0644)
	directRes := importSvc.ImportCompressedFiles([]string{directVpk}, targetAddons)
	if !directRes.Success || directRes.ImportedCount != 1 {
		t.Fatalf("Import direct VPK failed: %+v", directRes)
	}
	if _, err := os.Stat(filepath.Join(targetAddons, "direct.vpk")); err != nil {
		t.Errorf("direct vpk missing in target")
	}

	// Test Import Directory
	subFolder := filepath.Join(tempDir, "mod_pack_folder")
	_ = os.MkdirAll(subFolder, 0755)
	_ = os.WriteFile(filepath.Join(subFolder, "folder_mod.vpk"), []byte("folder vpk"), 0644)
	_ = os.WriteFile(filepath.Join(subFolder, "folder_mod.png"), []byte("folder png"), 0644)
	dirRes := importSvc.ImportCompressedFiles([]string{subFolder}, targetAddons)
	if !dirRes.Success || dirRes.ImportedCount != 2 {
		t.Fatalf("Import directory failed: %+v", dirRes)
	}
	if _, err := os.Stat(filepath.Join(targetAddons, "folder_mod.vpk")); err != nil {
		t.Errorf("folder_mod.vpk missing in target")
	}
}

func TestNotesService(t *testing.T) {
	tempDir := t.TempDir()
	notesSvc := NewNotesService(tempDir)

	workshopVpk := "C:\\steam\\addons\\workshop\\123456789.vpk"
	localVpk := "C:\\steam\\addons\\my_mod.vpk"

	if err := notesSvc.SetNote(workshopVpk, "连喷（超大声！）"); err != nil {
		t.Fatalf("SetNote failed: %v", err)
	}
	if err := notesSvc.SetNote(localVpk, "AK47-雷神涂装"); err != nil {
		t.Fatalf("SetNote failed: %v", err)
	}

	// 验证按 workshopID 获取
	note1 := notesSvc.GetNote(workshopVpk)
	if note1 != "连喷（超大声！）" {
		t.Errorf("GetNote workshop = %q, want '连喷（超大声！）'", note1)
	}

	// 即使路径变到 addons，根据 workshopID 也能认出
	relocatedVpk := "C:\\steam\\addons\\123456789.vpk"
	noteRelocated := notesSvc.GetNote(relocatedVpk)
	if noteRelocated != "连喷（超大声！）" {
		t.Errorf("GetNote relocated = %q, want '连喷（超大声！）'", noteRelocated)
	}

	// 即使被重命名为 .disabled，依然能读出备注
	disabledVpk := "C:\\steam\\addons\\my_mod.vpk.disabled"
	noteDisabled := notesSvc.GetNote(disabledVpk)
	if noteDisabled != "AK47-雷神涂装" {
		t.Errorf("GetNote disabled = %q, want 'AK47-雷神涂装'", noteDisabled)
	}

	// 重新加载测试持久化
	notesSvc2 := NewNotesService(tempDir)
	if n := notesSvc2.GetNote(workshopVpk); n != "连喷（超大声！）" {
		t.Errorf("Persisted note = %q", n)
	}
}

func TestExportWithBetterNamingAndNotes(t *testing.T) {
	tempDir := t.TempDir()
	notesSvc := NewNotesService(tempDir)
	metaSvc := NewVpkMetadataService(filepath.Join(tempDir, "cache"))
	metaSvc.SetNotesService(notesSvc)
	exportSvc := NewExportService(metaSvc, notesSvc)

	srcDir := filepath.Join(tempDir, "src")
	_ = os.MkdirAll(srcDir, 0755)
	vpkPath := filepath.Join(srcDir, "123456789.vpk")
	_ = os.WriteFile(vpkPath, []byte("fake vpk"), 0644)

	// 设置带特殊中文符号的备注
	_ = notesSvc.SetNote(vpkPath, "【全自动连喷】（自制音效！）")

	outDir := filepath.Join(tempDir, "out")
	res := exportSvc.ExportVpks([]string{vpkPath}, outDir, true, true, nil)
	if !res.Success {
		t.Fatalf("ExportVpks failed: %s", res.Error)
	}

	// 解压并验证归档内的文件名保留汉字且去掉了中文特殊符号
	targetAddons := filepath.Join(tempDir, "extracted")
	_ = os.MkdirAll(targetAddons, 0755)
	importSvc := NewImportService()
	impRes := importSvc.ImportCompressedFiles([]string{res.FilePath}, targetAddons)
	if !impRes.Success {
		t.Fatalf("ImportCompressedFiles failed: %+v", impRes.Errors)
	}

	expectedName := "全自动连喷_自制音效-123456789.vpk"
	extractedFile := filepath.Join(targetAddons, expectedName)
	if _, err := os.Stat(extractedFile); err != nil {
		t.Errorf("Expected exported file %s not found, checking directory contents", expectedName)
		entries, _ := os.ReadDir(targetAddons)
		for _, e := range entries {
			t.Logf("Found file: %s", e.Name())
		}
	}
}
