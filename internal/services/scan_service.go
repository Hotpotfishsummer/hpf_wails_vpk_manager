package services

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"hpf_wails_vpk_manager/internal/models"
)

type ScanService struct {
	metaService *VpkMetadataService
}

func NewScanService(meta *VpkMetadataService) *ScanService {
	return &ScanService{metaService: meta}
}

func (s *ScanService) Scan(directory string) models.ScanResult {
	addonsDir := filepath.Join(directory, "left4dead2", "addons")
	workshopDir := filepath.Join(addonsDir, "workshop")

	var result models.ScanResult
	result.Success = true

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		result.LocalVpks = s.scanDir(addonsDir)
	}()

	go func() {
		defer wg.Done()
		result.WorkshopVpks = s.scanDir(workshopDir)
	}()

	wg.Wait()
	return result
}

func (s *ScanService) scanDir(dir string) []models.VpkFile {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return []models.VpkFile{}
	}

	type scanJob struct {
		path string
	}

	var jobs []scanJob
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		lower := strings.ToLower(name)
		if strings.HasSuffix(lower, ".vpk") || strings.HasSuffix(lower, ".vpk.disabled") {
			jobs = append(jobs, scanJob{path: filepath.Join(dir, name)})
		}
	}

	if len(jobs) == 0 {
		return []models.VpkFile{}
	}

	results := make([]models.VpkFile, len(jobs))
	valid := make([]bool, len(jobs))

	// Run concurrent metadata worker pool
	workers := 8
	if len(jobs) < workers {
		workers = len(jobs)
	}

	var wg sync.WaitGroup
	jobChan := make(chan int, len(jobs))

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobChan {
				job := jobs[idx]
				meta, err := s.metaService.GetVpkMetadata(job.path)
				if err == nil && meta != nil {
					results[idx] = *meta
					valid[idx] = true
				}
			}
		}()
	}

	for i := range jobs {
		jobChan <- i
	}
	close(jobChan)
	wg.Wait()

	var output []models.VpkFile
	for i, ok := range valid {
		if ok {
			output = append(output, results[i])
		}
	}
	return output
}
