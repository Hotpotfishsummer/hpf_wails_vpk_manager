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
		stat, err := os.Stat(p)
		if err != nil {
			errors = append(errors, fmt.Sprintf("文件不存在: %s", p))
			continue
		}

		if stat.IsDir() {
			count, dirErrs := s.importDirectory(p, targetAddonsDir)
			importedCount += count
			errors = append(errors, dirErrs...)
			continue
		}

		lower := strings.ToLower(p)
		var importErr error
		if strings.HasSuffix(lower, ".tar.zst") || strings.HasSuffix(lower, ".zst") {
			importErr = s.extractTarZst(p, targetAddonsDir)
		} else if strings.HasSuffix(lower, ".zip") {
			importErr = s.extractZip(p, targetAddonsDir)
		} else if strings.HasSuffix(lower, ".7z") {
			importErr = s.extract7z(p, targetAddonsDir)
		} else if s.isVpkOrImage(p) {
			importErr = s.copyFileToDir(p, targetAddonsDir)
		} else {
			errors = append(errors, fmt.Sprintf("不支持的文件格式: %s", filepath.Base(p)))
			continue
		}

		if importErr != nil {
			errors = append(errors, fmt.Sprintf("导入 %s 失败: %v", filepath.Base(p), importErr))
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

func (s *ImportService) importDirectory(dirPath, targetDir string) (int, []string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, []string{fmt.Sprintf("读取目录 %s 失败: %v", filepath.Base(dirPath), err)}
	}

	count := 0
	var errors []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		fullPath := filepath.Join(dirPath, name)
		lower := strings.ToLower(name)

		var itemErr error
		if strings.HasSuffix(lower, ".tar.zst") || strings.HasSuffix(lower, ".zst") {
			itemErr = s.extractTarZst(fullPath, targetDir)
		} else if strings.HasSuffix(lower, ".zip") {
			itemErr = s.extractZip(fullPath, targetDir)
		} else if strings.HasSuffix(lower, ".7z") {
			itemErr = s.extract7z(fullPath, targetDir)
		} else if s.isVpkOrImage(name) {
			itemErr = s.copyFileToDir(fullPath, targetDir)
		} else {
			continue
		}

		if itemErr != nil {
			errors = append(errors, fmt.Sprintf("导入 %s 失败: %v", name, itemErr))
		} else {
			count++
		}
	}
	return count, errors
}

func (s *ImportService) copyFileToDir(srcPath, targetDir string) error {
	baseName := filepath.Base(srcPath)
	destPath := filepath.Join(targetDir, baseName)
	return copyFile(srcPath, destPath)
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
