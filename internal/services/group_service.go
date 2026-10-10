package services

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"

	"hpf_wails_vpk_manager/internal/models"
)

type GroupService struct {
	mu       sync.RWMutex
	filePath string
	groups   []models.VpkGroup
}

func NewGroupService(storageDir ...string) *GroupService {
	var dir string
	if len(storageDir) > 0 && storageDir[0] != "" {
		dir = storageDir[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		dir = filepath.Join(home, "Documents", "L4D2-VPK-Manager")
	}
	_ = os.MkdirAll(dir, 0755)

	filePath := filepath.Join(dir, "vpk_groups.json")
	_ = os.Remove(filePath + ".tmp")

	s := &GroupService{
		filePath: filePath,
		groups:   make([]models.VpkGroup, 0),
	}
	s.load()
	return s
}

func (s *GroupService) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}

	var parsed []models.VpkGroup
	if err := json.Unmarshal(data, &parsed); err != nil {
		log.Printf("Failed to decode vpk_groups.json: %v", err)
		return
	}
	if parsed == nil {
		parsed = make([]models.VpkGroup, 0)
	}
	s.groups = parsed
}

func (s *GroupService) save() error {
	data, err := json.MarshalIndent(s.groups, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal groups: %v", err)
		return err
	}

	tmpFile := s.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		log.Printf("Failed to write groups.tmp: %v", err)
		return err
	}

	return os.Rename(tmpFile, s.filePath)
}

func (s *GroupService) GetGroups() []models.VpkGroup {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]models.VpkGroup, len(s.groups))
	copy(res, s.groups)
	return res
}

func (s *GroupService) SaveGroups(groups []models.VpkGroup) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if groups == nil {
		groups = make([]models.VpkGroup, 0)
	}
	s.groups = groups
	return s.save()
}
