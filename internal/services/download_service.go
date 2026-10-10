package services

import (
	"path"
	"path/filepath"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
)

type DownloadService struct {
	sftpService *SftpService
}

func NewDownloadService(sftp *SftpService) *DownloadService {
	return &DownloadService{sftpService: sftp}
}

func (s *DownloadService) Download(
	remoteFilePaths []string,
	localAddonsDir string,
	config models.SftpSettings,
) models.DownloadResult {
	downloadedCount := 0

	for _, rPath := range remoteFilePaths {
		fileName := path.Base(rPath)
		lPath := filepath.Join(localAddonsDir, fileName)

		if err := s.sftpService.DownloadFile(config, rPath, lPath); err != nil {
			continue
		}
		downloadedCount++

		baseName := naming.VpkStemForAssociation(fileName)
		remoteDir := path.Dir(rPath)

		for _, ext := range naming.ImageExtensions {
			rImgPath := path.Join(remoteDir, baseName+ext)
			if s.sftpService.FileExists(config, rImgPath) {
				lImgPath := filepath.Join(localAddonsDir, baseName+ext)
				_ = s.sftpService.DownloadFile(config, rImgPath, lImgPath)
			}
		}
	}

	return models.DownloadResult{
		ApiResult:       models.ApiResult{Success: true},
		DownloadedCount: downloadedCount,
	}
}
