package services

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"hpf_wails_vpk_manager/internal/models"
)

type StorageService struct {
	mu          sync.RWMutex
	storageDir  string
	storageFile string
	data        map[string]any
}

func NewStorageService() *StorageService {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	storageDir := filepath.Join(home, "Documents", "L4D2-VPK-Manager")
	_ = os.MkdirAll(storageDir, 0755)

	storageFile := filepath.Join(storageDir, "preferences.json")

	// Clean stale tmp file
	staleTmp := filepath.Join(storageDir, "preferences.json.tmp")
	_ = os.Remove(staleTmp)

	s := &StorageService{
		storageDir:  storageDir,
		storageFile: storageFile,
		data:        make(map[string]any),
	}
	s.load()
	return s
}

func (s *StorageService) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	bytes, err := os.ReadFile(s.storageFile)
	if err != nil {
		return
	}

	var parsed map[string]any
	if err := json.Unmarshal(bytes, &parsed); err != nil {
		log.Printf("Failed to decode preferences.json: %v", err)
		corruptFile := s.storageFile + ".corrupt"
		_ = os.Rename(s.storageFile, corruptFile)
		return
	}
	s.data = parsed
}

func (s *StorageService) save() {
	bytes, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal preferences: %v", err)
		return
	}

	tmpFile := s.storageFile + ".tmp"
	if err := os.WriteFile(tmpFile, bytes, 0644); err != nil {
		log.Printf("Failed to write preferences.tmp: %v", err)
		return
	}

	_ = os.Rename(tmpFile, s.storageFile)
}

func (s *StorageService) GetL4D2Directory() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if val, ok := s.data["l4d2_directory"]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (s *StorageService) SaveL4D2Directory(dir string) {
	s.mu.Lock()
	s.data["l4d2_directory"] = dir
	s.save()
	s.mu.Unlock()
}

func (s *StorageService) GetSftpSettings() models.SftpSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()

	def := models.SftpSettings{}
	val, ok := s.data["sftp_settings"]
	if !ok {
		return def
	}

	// Marshal/unmarshal to cleanly convert map[string]any to struct
	bytes, err := json.Marshal(val)
	if err != nil {
		return def
	}
	var res models.SftpSettings
	if err := json.Unmarshal(bytes, &res); err != nil {
		return def
	}
	return res
}

func (s *StorageService) SaveSftpSettings(settings models.SftpSettings) {
	s.mu.Lock()
	s.data["sftp_settings"] = settings
	s.save()
	s.mu.Unlock()
}

func (s *StorageService) GetStorageDir() string {
	return s.storageDir
}
