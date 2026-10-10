package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
	"hpf_wails_vpk_manager/internal/vpk"
)

type VpkCachedMetadata struct {
	AddonTitle  *string `json:"addontitle,omitempty"`
	SourceMtime float64 `json:"source_mtime,omitempty"`
	SourceSize  int64   `json:"source_size,omitempty"`
}

type VpkMetadataService struct {
	mu        sync.RWMutex
	configDir string
}

func NewVpkMetadataService(configDir ...string) *VpkMetadataService {
	var dir string
	if len(configDir) > 0 && configDir[0] != "" {
		dir = configDir[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		dir = filepath.Join(home, "Documents", "L4D2-VPK-Manager", ".vpk_config")
	}

	_ = os.MkdirAll(dir, 0755)

	// Clean stale .json.tmp files
	if entries, err := os.ReadDir(dir); err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".json.tmp") {
				_ = os.Remove(filepath.Join(dir, entry.Name()))
			}
		}
	}

	return &VpkMetadataService{configDir: dir}
}

func (m *VpkMetadataService) getJsonPath(vpkPath string) string {
	stem := strings.TrimSuffix(filepath.Base(vpkPath), filepath.Ext(vpkPath))
	return filepath.Join(m.configDir, stem+".json")
}

func (m *VpkMetadataService) LoadMetadata(vpkPath string) *VpkCachedMetadata {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jsonPath := m.getJsonPath(vpkPath)
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil
	}

	var meta VpkCachedMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil
	}
	return &meta
}

func (m *VpkMetadataService) SaveMetadata(vpkPath string, meta *VpkCachedMetadata) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	jsonPath := m.getJsonPath(vpkPath)
	tmpPath := jsonPath + ".tmp"

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return false
	}

	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return false
	}

	if err := os.Rename(tmpPath, jsonPath); err != nil {
		_ = os.Remove(tmpPath)
		return false
	}
	return true
}

func (m *VpkMetadataService) isCacheFresh(cached *VpkCachedMetadata, info os.FileInfo) bool {
	if cached == nil {
		return false
	}
	if cached.SourceSize == 0 || cached.SourceMtime == 0 {
		return false
	}
	mtimeFloat := float64(info.ModTime().UnixNano()) / 1e9
	// Check size and mtime (allowing small float diff)
	return cached.SourceSize == info.Size() && (cached.SourceMtime == mtimeFloat || int64(cached.SourceMtime) == info.ModTime().Unix())
}

func (m *VpkMetadataService) GetAddonTitle(vpkPath string) *string {
	info, err := os.Stat(vpkPath)
	if err != nil {
		// File might not exist, check cache
		if cached := m.LoadMetadata(vpkPath); cached != nil {
			return cached.AddonTitle
		}
		return nil
	}

	cached := m.LoadMetadata(vpkPath)
	if cached != nil && m.isCacheFresh(cached, info) {
		return cached.AddonTitle
	}

	// Extract title from VPK
	titleStr, err := vpk.ExtractAddonTitle(vpkPath)
	var title *string
	if err == nil && titleStr != "" {
		title = &titleStr
	}

	newMeta := &VpkCachedMetadata{
		AddonTitle:  title,
		SourceMtime: float64(info.ModTime().UnixNano()) / 1e9,
		SourceSize:  info.Size(),
	}
	m.SaveMetadata(vpkPath, newMeta)
	return title
}

func (m *VpkMetadataService) GetVpkMetadata(vpkPath string) (*models.VpkFile, error) {
	info, err := os.Stat(vpkPath)
	if err != nil {
		return nil, err
	}

	filename := filepath.Base(vpkPath)
	lowerName := strings.ToLower(filename)
	isDisabled := strings.HasSuffix(lowerName, ".disabled")

	addonTitle := m.GetAddonTitle(vpkPath)

	displayName := filename
	if addonTitle != nil && *addonTitle != "" {
		displayName = *addonTitle
	}

	// Look for thumbnail image
	var thumbnailURL *string
	parentDir := filepath.Dir(vpkPath)
	stem := naming.VpkStemForAssociation(filename)

	for _, ext := range []string{".jpg", ".png", ".jpeg"} {
		thumbPath := filepath.Join(parentDir, stem+ext)
		if thumbBytes, err := os.ReadFile(thumbPath); err == nil {
			mime := "image/jpeg"
			if ext == ".png" {
				mime = "image/png"
			}
			encoded := base64.StdEncoding.EncodeToString(thumbBytes)
			dataURL := fmt.Sprintf("data:%s;base64,%s", mime, encoded)
			thumbnailURL = &dataURL
			break
		}
	}

	return &models.VpkFile{
		Name:          filename,
		Path:          vpkPath,
		Size:          info.Size(),
		ModifiedTime:  float64(info.ModTime().UnixNano()) / 1e9,
		IsValid:       true,
		IsDisabled:    isDisabled,
		ThumbnailPath: thumbnailURL,
		AddonTitle:    addonTitle,
		DisplayName:   displayName,
	}, nil
}
