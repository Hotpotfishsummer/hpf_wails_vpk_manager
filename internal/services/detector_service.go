package services

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"

	"hpf_wails_vpk_manager/internal/models"
)

var (
	// Windows 路径正则，用于从快捷方式或文本中匹配可能的 exe 路径
	exePathRegex = regexp.MustCompile(`(?i)[a-zA-Z]:\\[^<>"|?*\r\n\x00-\x1f]+\.exe`)

	// 磁盘扫描跳过的目录（小写）
	skipDirNames = map[string]bool{
		"$recycle.bin":              true,
		"system volume information": true,
		"windows":                   true,
		"appdata":                   true,
		"programdata":               true,
		"node_modules":              true,
		".git":                      true,
		".svn":                      true,
		"recovery":                  true,
		"windowsapps":               true,
		"$winreagent":               true,
	}
)

type DetectorService struct{}

func NewDetectorService() *DetectorService {
	return &DetectorService{}
}

// DetectClients 自动探测所有可用的 Left 4 Dead 2 客户端目录
func (s *DetectorService) DetectClients(deepScan bool) []models.DetectedGameClient {
	var mu sync.Mutex
	detectedMap := make(map[string]models.DetectedGameClient) // key 为小写规范化路径

	sourcePriority := map[string]int{
		"running":  100,
		"steam":    80,
		"registry": 60,
		"shortcut": 40,
		"common":   20,
		"disk":     10,
	}

	addCandidate := func(candidatePath string, source string, sourceName string) {
		if candidatePath == "" {
			return
		}
		cleanPath := filepath.Clean(candidatePath)
		valid, finalDir, exePath := s.ValidateL4D2Directory(cleanPath)
		if !valid {
			return
		}

		key := strings.ToLower(finalDir)
		mu.Lock()
		defer mu.Unlock()
		if existing, exists := detectedMap[key]; exists {
			if sourcePriority[sourceName] > sourcePriority[existing.SourceName] {
				existing.Source = source
				existing.SourceName = sourceName
				detectedMap[key] = existing
			}
			return
		}

		client := s.buildClientInfo(finalDir, exePath, source, sourceName)
		detectedMap[key] = client
	}

	// 阶段 1：快速多源探测（并行执行）
	var wg sync.WaitGroup

	// 1.1 正在运行的进程
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, path := range s.detectRunningProcesses() {
			addCandidate(path, "正在运行的进程", "running")
		}
	}()

	// 1.2 Windows 注册表（卸载项、运行缓存 MuiCache、防火墙规则）
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, path := range s.detectRegistry() {
			addCandidate(path, "系统注册表/记录", "registry")
		}
	}()

	// 1.3 Steam 客户端与库目录（Steam 专属探测，作为多源之一）
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, path := range s.detectSteamLibraries() {
			addCandidate(path, "Steam 游戏库", "steam")
		}
	}()

	// 1.4 桌面与开始菜单快捷方式
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, path := range s.detectShortcuts() {
			addCandidate(path, "桌面/开始菜单快捷方式", "shortcut")
		}
	}()

	// 1.5 所有固定驱动器常见路径启发式探测
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, path := range s.detectCommonPaths() {
			addCandidate(path, "常见安装目录", "common")
		}
	}()

	wg.Wait()

	// 阶段 2：如果显式要求深度全盘扫描，或者快速扫描完全没有找到任何客户端
	if deepScan || len(detectedMap) == 0 {
		for _, path := range s.detectDeepDrives(8 * time.Second) {
			addCandidate(path, "磁盘智能扫描", "disk")
		}
	}

	var results []models.DetectedGameClient
	for _, client := range detectedMap {
		results = append(results, client)
	}

	// 智能排序：优先显示正在运行或模组数更多的主力版本
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			pI := sourcePriority[results[i].SourceName]
			pJ := sourcePriority[results[j].SourceName]
			// 如果其中一个是正在运行，置顶
			if results[i].SourceName == "running" && results[j].SourceName != "running" {
				continue
			}
			if results[j].SourceName == "running" && results[i].SourceName != "running" {
				results[i], results[j] = results[j], results[i]
				continue
			}
			// 否则按模组数量倒序，若模组数相同则按来源优先级倒序
			if results[j].ModCount > results[i].ModCount ||
				(results[j].ModCount == results[i].ModCount && pJ > pI) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results
}

// ValidateL4D2Directory 校验目录是否是合法的 L4D2 客户端根目录
// 同时兼容用户可能不小心传入了 "left4dead2" 子目录的情况（会自动上移一层）
func (s *DetectorService) ValidateL4D2Directory(dir string) (bool, string, string) {
	fi, err := os.Stat(dir)
	if err != nil || !fi.IsDir() {
		return false, "", ""
	}

	// 如果传入的目录名本身就叫 left4dead2，检查其父目录是否为根目录
	base := strings.ToLower(filepath.Base(dir))
	if base == "left4dead2" {
		parent := filepath.Dir(dir)
		if valid, finalDir, exe := s.checkDirValid(parent); valid {
			return true, finalDir, exe
		}
	}

	return s.checkDirValid(dir)
}

func (s *DetectorService) checkDirValid(dir string) (bool, string, string) {
	// 1. 检查是否存在主程序 left4dead2.exe (或 l4d2.exe)
	exeName := "left4dead2.exe"
	exePath := filepath.Join(dir, exeName)
	if fi, err := os.Stat(exePath); err != nil || fi.IsDir() {
		// 备用：检查 l4d2.exe
		altPath := filepath.Join(dir, "l4d2.exe")
		if fi2, err2 := os.Stat(altPath); err2 == nil && !fi2.IsDir() {
			exePath = altPath
		} else {
			return false, "", ""
		}
	}

	// 2. 检查是否存在子目录 left4dead2
	gameSubDir := filepath.Join(dir, "left4dead2")
	subFi, err := os.Stat(gameSubDir)
	if err != nil || !subFi.IsDir() {
		return false, "", ""
	}

	// 3. 检查特征文件：gameinfo.txt 或 pak01_dir.vpk 或 addons 目录
	gameinfoPath := filepath.Join(gameSubDir, "gameinfo.txt")
	pakPath := filepath.Join(gameSubDir, "pak01_dir.vpk")
	addonsPath := filepath.Join(gameSubDir, "addons")

	hasGameInfo := fileExists(gameinfoPath)
	hasPak := fileExists(pakPath)
	hasAddons := dirExists(addonsPath)

	if !hasGameInfo && !hasPak && !hasAddons {
		return false, "", ""
	}

	return true, dir, exePath
}

func (s *DetectorService) buildClientInfo(dir, exePath, source, sourceName string) models.DetectedGameClient {
	client := models.DetectedGameClient{
		Path:       dir,
		ExePath:    exePath,
		Source:     source,
		SourceName: sourceName,
	}

	if fi, err := os.Stat(exePath); err == nil {
		client.LastModified = fi.ModTime().Unix()
	}

	addonsDir := filepath.Join(dir, "left4dead2", "addons")
	if dirExists(addonsDir) {
		workshopDir := filepath.Join(addonsDir, "workshop")
		client.HasWorkshop = dirExists(workshopDir)

		// 快速统计 addons 下模组数量
		modCount := 0
		if entries, err := os.ReadDir(addonsDir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					name := strings.ToLower(entry.Name())
					if strings.HasSuffix(name, ".vpk") || strings.HasSuffix(name, ".vpk.disabled") {
						modCount++
					}
				}
			}
		}
		if client.HasWorkshop {
			if wsEntries, err := os.ReadDir(workshopDir); err == nil {
				for _, entry := range wsEntries {
					if !entry.IsDir() {
						name := strings.ToLower(entry.Name())
						if strings.HasSuffix(name, ".vpk") || strings.HasSuffix(name, ".vpk.disabled") {
							modCount++
						}
					}
				}
			}
		}
		client.ModCount = modCount
	}

	return client
}

// 1.1 正在运行的进程探测
func (s *DetectorService) detectRunningProcesses() []string {
	var paths []string
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))
	if err := windows.Process32First(snapshot, &entry); err != nil {
		return nil
	}

	for {
		exeName := strings.ToLower(windows.UTF16ToString(entry.ExeFile[:]))
		if exeName == "left4dead2.exe" || exeName == "l4d2.exe" {
			hProc, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, entry.ProcessID)
			if err == nil {
				var buf [windows.MAX_PATH]uint16
				size := uint32(len(buf))
				if err := windows.QueryFullProcessImageName(hProc, 0, &buf[0], &size); err == nil {
					fullPath := windows.UTF16ToString(buf[:size])
					paths = append(paths, filepath.Dir(fullPath))
				}
				windows.CloseHandle(hProc)
			}
		}

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			break
		}
	}
	return paths
}

// 1.2 系统注册表探测
func (s *DetectorService) detectRegistry() []string {
	var paths []string

	// A. 卸载项 (涵盖各类非 Steam 安装包、绿色安装包)
	uninstallKeys := []struct {
		root registry.Key
		path string
	}{
		{registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall`},
		{registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Uninstall`},
	}

	for _, uk := range uninstallKeys {
		k, err := registry.OpenKey(uk.root, uk.path, registry.READ)
		if err != nil {
			continue
		}
		subKeys, err := k.ReadSubKeyNames(-1)
		k.Close()
		if err != nil {
			continue
		}

		for _, sub := range subKeys {
			sk, err := registry.OpenKey(uk.root, uk.path+`\`+sub, registry.READ)
			if err != nil {
				continue
			}

			displayName, _, _ := sk.GetStringValue("DisplayName")
			installLocation, _, _ := sk.GetStringValue("InstallLocation")
			displayIcon, _, _ := sk.GetStringValue("DisplayIcon")
			sk.Close()

			lowerName := strings.ToLower(displayName)
			if strings.Contains(lowerName, "left 4 dead") || strings.Contains(lowerName, "求生之路") || strings.Contains(lowerName, "l4d2") {
				if installLocation != "" {
					paths = append(paths, installLocation)
				}
			}
			if displayIcon != "" && (strings.Contains(strings.ToLower(displayIcon), "left4dead2.exe") || strings.Contains(strings.ToLower(displayIcon), "l4d2.exe")) {
				iconPath := strings.Trim(displayIcon, "\"")
				if idx := strings.Index(strings.ToLower(iconPath), ".exe"); idx != -1 {
					iconPath = iconPath[:idx+4]
					paths = append(paths, filepath.Dir(iconPath))
				}
			}
		}
	}

	// B. MuiCache 记录（凡是在 Windows 上运行过的客户端必有记录，免安装绿色版神器）
	muiKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Classes\Local Settings\Software\Microsoft\Windows\Shell\MuiCache`, registry.READ)
	if err == nil {
		valNames, _ := muiKey.ReadValueNames(-1)
		muiKey.Close()
		for _, val := range valNames {
			lower := strings.ToLower(val)
			if strings.Contains(lower, "left4dead2.exe") || strings.Contains(lower, "l4d2.exe") {
				if idx := strings.Index(lower, ".exe"); idx != -1 {
					exePath := val[:idx+4]
					paths = append(paths, filepath.Dir(exePath))
				}
			}
		}
	}

	// C. Compatibility Assistant 运行记录
	compatKey, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows NT\CurrentVersion\AppCompatFlags\Compatibility Assistant\Store`, registry.READ)
	if err == nil {
		valNames, _ := compatKey.ReadValueNames(-1)
		compatKey.Close()
		for _, val := range valNames {
			lower := strings.ToLower(val)
			if strings.Contains(lower, "left4dead2.exe") || strings.Contains(lower, "l4d2.exe") {
				paths = append(paths, filepath.Dir(val))
			}
		}
	}

	// D. 防火墙联网规则（曾联网或对战过的程序）
	fwKey, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Services\SharedAccess\Parameters\FirewallPolicy\FirewallRules`, registry.READ)
	if err == nil {
		valNames, _ := fwKey.ReadValueNames(-1)
		for _, valName := range valNames {
			ruleData, _, err := fwKey.GetStringValue(valName)
			if err == nil && (strings.Contains(strings.ToLower(ruleData), "left4dead2.exe") || strings.Contains(strings.ToLower(ruleData), "l4d2.exe")) {
				// 格式如: ...|App=C:\path\to\left4dead2.exe|...
				if match := exePathRegex.FindString(ruleData); match != "" {
					paths = append(paths, filepath.Dir(match))
				}
			}
		}
		fwKey.Close()
	}

	return paths
}

// 1.3 Steam 库探测
func (s *DetectorService) detectSteamLibraries() []string {
	var paths []string

	var steamPath string
	if k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.READ); err == nil {
		steamPath, _, _ = k.GetStringValue("SteamPath")
		k.Close()
	}
	if steamPath == "" {
		if k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, registry.READ); err == nil {
			steamPath, _, _ = k.GetStringValue("InstallPath")
			k.Close()
		}
	}

	if steamPath != "" {
		steamPath = filepath.Clean(steamPath)
		// 默认库
		paths = append(paths, filepath.Join(steamPath, "steamapps", "common", "Left 4 Dead 2"))

		// 解析 libraryfolders.vdf
		vdfPath := filepath.Join(steamPath, "steamapps", "libraryfolders.vdf")
		if data, err := os.ReadFile(vdfPath); err == nil {
			vdfContent := string(data)
			re := regexp.MustCompile(`"path"\s+"([^"]+)"`)
			matches := re.FindAllStringSubmatch(vdfContent, -1)
			for _, match := range matches {
				if len(match) > 1 {
					libPath := strings.ReplaceAll(match[1], `\\`, `\`)
					paths = append(paths, filepath.Join(libPath, "steamapps", "common", "Left 4 Dead 2"))
				}
			}
		}
	}

	return paths
}

// 1.4 桌面与开始菜单快捷方式（快速有界搜索）
func (s *DetectorService) detectShortcuts() []string {
	var paths []string
	home, _ := os.UserHomeDir()

	// 桌面快捷方式（仅扫描桌面本身）
	desktopDirs := []string{
		filepath.Join(home, "Desktop"),
		`C:\Users\Public\Desktop`,
	}
	for _, dDir := range desktopDirs {
		entries, err := os.ReadDir(dDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() && strings.ToLower(filepath.Ext(e.Name())) == ".lnk" {
				fullPath := filepath.Join(dDir, e.Name())
				if data, err := os.ReadFile(fullPath); err == nil {
					matches := exePathRegex.FindAllString(string(data), -1)
					for _, m := range matches {
						lower := strings.ToLower(m)
						if strings.HasSuffix(lower, "left4dead2.exe") || strings.HasSuffix(lower, "l4d2.exe") {
							paths = append(paths, filepath.Dir(m))
						}
					}
				}
			}
		}
	}

	// 开始菜单快捷方式（最多深入 2 层）
	menuDirs := []string{
		filepath.Join(home, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs"),
		`C:\ProgramData\Microsoft\Windows\Start Menu\Programs`,
	}
	for _, mDir := range menuDirs {
		s.scanLnkShallow(mDir, 0, 2, &paths)
	}

	return paths
}

func (s *DetectorService) scanLnkShallow(dir string, depth, maxDepth int, out *[]string) {
	if depth > maxDepth {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			lower := strings.ToLower(name)
			if skipDirNames[lower] || strings.HasPrefix(lower, ".") {
				continue
			}
			s.scanLnkShallow(filepath.Join(dir, name), depth+1, maxDepth, out)
		} else if strings.ToLower(filepath.Ext(name)) == ".lnk" {
			// 仅当快捷方式名称包含特征词或位于游戏目录时深入读取，避免读取无关办公软件lnk
			lower := strings.ToLower(name)
			if strings.Contains(lower, "left") || strings.Contains(lower, "l4d") || strings.Contains(lower, "求生") || depth > 0 {
				fullPath := filepath.Join(dir, name)
				if data, err := os.ReadFile(fullPath); err == nil {
					matches := exePathRegex.FindAllString(string(data), -1)
					for _, m := range matches {
						mLower := strings.ToLower(m)
						if strings.HasSuffix(mLower, "left4dead2.exe") || strings.HasSuffix(mLower, "l4d2.exe") {
							*out = append(*out, filepath.Dir(m))
						}
					}
				}
			}
		}
	}
}

// 1.5 常见安装目录与所有固定磁盘探测
func (s *DetectorService) detectCommonPaths() []string {
	var paths []string
	fixedDrives := s.getFixedDrives()

	commonRelativeDirs := []string{
		`Left 4 Dead 2`,
		`Left4Dead2`,
		`L4D2`,
		`求生之路2`,
		`求生之路`,
		`Games\Left 4 Dead 2`,
		`Games\Left4Dead2`,
		`Games\L4D2`,
		`Games\求生之路2`,
		`Game\Left 4 Dead 2`,
		`Game\Left4Dead2`,
		`Game\L4D2`,
		`SteamLibrary\steamapps\common\Left 4 Dead 2`,
		`Program Files (x86)\Steam\steamapps\common\Left 4 Dead 2`,
		`Program Files\Steam\steamapps\common\Left 4 Dead 2`,
		`Program Files (x86)\Left 4 Dead 2`,
		`Program Files\Left 4 Dead 2`,
		`WeGameApps\Left 4 Dead 2`,
		`Dbg\Left 4 Dead 2`,
	}

	for _, drive := range fixedDrives {
		for _, rel := range commonRelativeDirs {
			paths = append(paths, filepath.Join(drive, rel))
		}
	}

	// 检查用户的 Downloads 和 Desktop 的一级子目录
	home, _ := os.UserHomeDir()
	if home != "" {
		for _, userDir := range []string{filepath.Join(home, "Downloads"), filepath.Join(home, "Desktop")} {
			entries, err := os.ReadDir(userDir)
			if err == nil {
				for _, entry := range entries {
					if entry.IsDir() {
						lower := strings.ToLower(entry.Name())
						if strings.Contains(lower, "left") || strings.Contains(lower, "l4d") || strings.Contains(lower, "求生") {
							paths = append(paths, filepath.Join(userDir, entry.Name()))
						}
					}
				}
			}
		}
	}

	return paths
}

// 阶段 2：全盘智能剪枝深度扫描（多协程并发，限制深度至 4 层，带超时保护）
func (s *DetectorService) detectDeepDrives(timeout time.Duration) []string {
	fixedDrives := s.getFixedDrives()
	if len(fixedDrives) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var mu sync.Mutex
	var foundPaths []string
	var wg sync.WaitGroup

	for _, drive := range fixedDrives {
		wg.Add(1)
		go func(drv string) {
			defer wg.Done()
			var scanDir func(dir string, depth int)
			scanDir = func(dir string, depth int) {
				select {
				case <-ctx.Done():
					return
				default:
				}

				if depth > 4 {
					return
				}

				entries, err := os.ReadDir(dir)
				if err != nil {
					return
				}

				hasExe := false
				var subDirs []string

				for _, e := range entries {
					name := e.Name()
					lowerName := strings.ToLower(name)
					if e.IsDir() {
						if skipDirNames[lowerName] || strings.HasPrefix(lowerName, "$") || strings.HasPrefix(lowerName, ".") {
							continue
						}
						subDirs = append(subDirs, filepath.Join(dir, name))
					} else {
						if lowerName == "left4dead2.exe" || lowerName == "l4d2.exe" {
							hasExe = true
						}
					}
				}

				if hasExe {
					if valid, finalDir, _ := s.ValidateL4D2Directory(dir); valid {
						mu.Lock()
						foundPaths = append(foundPaths, finalDir)
						mu.Unlock()
						return // 找到了该游戏目录，不再扫描内部子目录
					}
				}

				for _, sub := range subDirs {
					scanDir(sub, depth+1)
				}
			}

			scanDir(drv, 0)
		}(drive)
	}

	wg.Wait()
	return foundPaths
}

func (s *DetectorService) getFixedDrives() []string {
	var drives []string
	mask, err := windows.GetLogicalDrives()
	if err != nil {
		return []string{`C:\`}
	}

	for i := 0; i < 26; i++ {
		if mask&(1<<i) != 0 {
			root := string(rune('A'+i)) + `:\`
			driveType := windows.GetDriveType(windows.StringToUTF16Ptr(root))
			if driveType == windows.DRIVE_FIXED {
				drives = append(drives, root)
			}
		}
	}
	if len(drives) == 0 {
		drives = append(drives, `C:\`)
	}
	return drives
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}
