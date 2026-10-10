package services

import (
	"path"
	"path/filepath"
	"strings"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
)

type UploadService struct {
	metaService *VpkMetadataService
	sftpService *SftpService
}

func NewUploadService(meta *VpkMetadataService, sftp *SftpService) *UploadService {
	return &UploadService{
		metaService: meta,
		sftpService: sftp,
	}
}

func (s *UploadService) Upload(
	localFilePaths []string,
	remoteAddonsPath string,
	config models.SftpSettings,
	betterNaming bool,
) models.UploadResult {
	uploadedCount := 0

	for _, lPath := range localFilePaths {
		isDisabled := strings.HasSuffix(strings.ToLower(lPath), ".vpk.disabled")
		baseName := s.resolveBaseName(lPath, betterNaming)

		suffix := ".vpk"
		if isDisabled {
			suffix = ".vpk.disabled"
		}

		rPath := path.Join(remoteAddonsPath, baseName+suffix)
		if err := s.sftpService.UploadFile(config, lPath, rPath); err == nil {
			uploadedCount++
			s.uploadAssociatedImages(lPath, baseName, remoteAddonsPath, config)
		}
	}

	return models.UploadResult{
		ApiResult:     models.ApiResult{Success: true},
		UploadedCount: uploadedCount,
	}
}

func (s *UploadService) resolveBaseName(localPath string, betterNaming bool) string {
	name := filepath.Base(localPath)
	lower := strings.ToLower(name)

	var fallback string
	if strings.HasSuffix(lower, ".vpk.disabled") {
		fallback = name[:len(name)-len(".vpk.disabled")]
	} else {
		fallback = strings.TrimSuffix(name, filepath.Ext(name))
	}

	if !betterNaming {
		return fallback
	}

	title := s.metaService.GetAddonTitle(localPath)
	var titleStr string
	if title != nil {
		titleStr = *title
	}
	return naming.BuildExportBaseName(titleStr, name, fallback)
}

func (s *UploadService) uploadAssociatedImages(localPath, baseName, remoteAddonsPath string, config models.SftpSettings) {
	uploadedExts := make(map[string]bool)
	matches := naming.FindAssociatedImages(filepath.Dir(localPath), filepath.Base(localPath), true)

	for _, m := range matches {
		extLower := strings.ToLower(m.Extension)
		if !uploadedExts[extLower] {
			remoteImgPath := path.Join(remoteAddonsPath, baseName+m.Extension)
			if err := s.sftpService.UploadFile(config, m.FullPath, remoteImgPath); err == nil {
				uploadedExts[extLower] = true
			}
		}
	}
}
