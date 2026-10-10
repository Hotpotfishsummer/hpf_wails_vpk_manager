package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
)

type LocalizeService struct {
	metaService  *VpkMetadataService
	notesService *NotesService
}

func NewLocalizeService(meta *VpkMetadataService, notes ...*NotesService) *LocalizeService {
	var n *NotesService
	if len(notes) > 0 {
		n = notes[0]
	}
	return &LocalizeService{metaService: meta, notesService: n}
}

func (s *LocalizeService) Localize(filePaths []string, l4d2Dir string, betterNaming bool) models.LocalizeResult {
	addonsDir := filepath.Join(l4d2Dir, "left4dead2", "addons")
	_ = os.MkdirAll(addonsDir, 0755)

	trashDir := filepath.Join(addonsDir, "workshop", ".trash")
	_ = os.MkdirAll(trashDir, 0755)

	successCount := 0
	var errors []string

	for _, srcPathStr := range filePaths {
		srcPath, err := filepath.Abs(srcPathStr)
		if err != nil {
			errors = append(errors, fmt.Sprintf("路径无效: %s", srcPathStr))
			continue
		}

		if _, err := os.Stat(srcPath); err != nil {
			errors = append(errors, fmt.Sprintf("文件不存在: %s", filepath.Base(srcPathStr)))
			continue
		}

		isInWorkshop := strings.Contains(strings.ToLower(srcPath), "workshop")
		baseName := s.resolveBaseName(srcPath, betterNaming)

		destVpk := filepath.Join(addonsDir, baseName+".vpk")
		if err := copyFile(srcPath, destVpk); err != nil {
			errors = append(errors, fmt.Sprintf("复制 %s 失败: %v", filepath.Base(srcPath), err))
			continue
		}

		// Copy associated images
		var associatedImages []string
		for _, img := range naming.FindAssociatedImages(filepath.Dir(srcPath), filepath.Base(srcPath), false) {
			destImg := filepath.Join(addonsDir, baseName+img.Extension)
			_ = copyFile(img.FullPath, destImg)
			associatedImages = append(associatedImages, img.FullPath)
		}

		// If in workshop, move original files to .trash
		if isInWorkshop {
			s.moveToTrash(srcPath, trashDir)
			for _, imgPath := range associatedImages {
				if _, err := os.Stat(imgPath); err == nil {
					s.moveToTrash(imgPath, trashDir)
				}
			}
		}

		successCount++
	}

	return models.LocalizeResult{
		ApiResult: models.ApiResult{
			Success: len(errors) == 0,
		},
		SuccessCount: successCount,
		Errors:       errors,
	}
}

func (s *LocalizeService) resolveBaseName(srcPath string, betterNaming bool) string {
	name := filepath.Base(srcPath)
	lower := strings.ToLower(name)

	var fallback string
	if strings.HasSuffix(lower, ".disabled") {
		withoutDisabled := name[:len(name)-len(".disabled")]
		fallback = strings.TrimSuffix(withoutDisabled, filepath.Ext(withoutDisabled))
	} else {
		fallback = strings.TrimSuffix(name, filepath.Ext(name))
	}

	if !betterNaming {
		return fallback
	}

	title := s.metaService.GetAddonTitle(srcPath)
	var titleStr string
	if title != nil {
		titleStr = *title
	}
	var note string
	if s.notesService != nil {
		note = s.notesService.GetNote(srcPath)
	}
	return naming.BuildExportBaseNameWithNote(titleStr, name, fallback, note)
}

func (s *LocalizeService) moveToTrash(filePath, trashDir string) {
	name := filepath.Base(filePath)
	trashPath := filepath.Join(trashDir, name)
	_ = os.Remove(trashPath)
	_ = os.Rename(filePath, trashPath)
}
