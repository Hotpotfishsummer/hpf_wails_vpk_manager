package services

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/klauspost/compress/zstd"

	"hpf_wails_vpk_manager/internal/models"
	"hpf_wails_vpk_manager/internal/naming"
)

type ImportService struct{}

func NewImportService() *ImportService {
	return &ImportService{}
}

func (s *ImportService) isVpkOrImage(filename string) bool {
	lower := strings.ToLower(filename)
	if strings.HasSuffix(lower, ".vpk") || strings.HasSuffix(lower, ".vpk.disabled") {
		return true
	}
	for _, ext := range naming.ImageExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func (s *ImportService) ImportCompressedFiles(filePaths []string, targetAddonsDir string) models.ImportResult {
	targetInfo, err := os.Stat(targetAddonsDir)
	if err != nil || !targetInfo.IsDir() {
		return models.ImportResult{
			ApiResult: models.ApiResult{Success: false, Error: fmt.Sprintf("目标目录不存在: %s", targetAddonsDir)},
		}
	}

	importedCount := 0
	var errors []string

	for _, p := range filePaths {
		if _, err := os.Stat(p); err != nil {
			errors = append(errors, fmt.Sprintf("文件不存在: %s", p))
			continue
		}

		lower := strings.ToLower(p)
		var err error
		if strings.HasSuffix(lower, ".tar.zst") {
			err = s.extractTarZst(p, targetAddonsDir)
		} else if strings.HasSuffix(lower, ".zip") {
			err = s.extractZip(p, targetAddonsDir)
		} else if strings.HasSuffix(lower, ".7z") {
			err = s.extract7z(p, targetAddonsDir)
		} else {
			errors = append(errors, fmt.Sprintf("不支持的文件格式: %s", filepath.Base(p)))
			continue
		}

		if err != nil {
			errors = append(errors, fmt.Sprintf("解压 %s 失败: %v", filepath.Base(p), err))
		} else {
			importedCount++
		}
	}

	return models.ImportResult{
		ApiResult: models.ApiResult{
			Success: len(errors) == 0,
		},
		ImportedCount: importedCount,
		Errors:        errors,
	}
}

func (s *ImportService) extractZip(archivePath, targetDir string) error {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if !s.isVpkOrImage(f.Name) {
			continue
		}

		baseName := filepath.Base(f.Name)
		if baseName == "" || baseName == "." {
			continue
		}

		destPath := filepath.Join(targetDir, baseName)
		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ImportService) extractTarZst(archivePath, targetDir string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	zstdReader, err := zstd.NewReader(f)
	if err != nil {
		return err
	}
	defer zstdReader.Close()

	tarReader := tar.NewReader(zstdReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeDir {
			continue
		}
		if !s.isVpkOrImage(header.Name) {
			continue
		}

		baseName := filepath.Base(header.Name)
		if baseName == "" || baseName == "." {
			continue
		}

		destPath := filepath.Join(targetDir, baseName)
		outFile, err := os.Create(destPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, tarReader)
		outFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ImportService) extract7z(archivePath, targetDir string) error {
	r, err := sevenzip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		if !s.isVpkOrImage(f.Name) {
			continue
		}

		baseName := filepath.Base(f.Name)
		if baseName == "" || baseName == "." {
			continue
		}

		destPath := filepath.Join(targetDir, baseName)
		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.Create(destPath)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
