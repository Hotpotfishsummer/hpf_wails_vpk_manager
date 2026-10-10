package services

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"hpf_wails_vpk_manager/internal/naming"
)

type NotesService struct {
	mu       sync.RWMutex
	filePath string
	notes    map[string]string
}

func NewNotesService(storageDir ...string) *NotesService {
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

	filePath := filepath.Join(dir, "vpk_notes.json")
	_ = os.Remove(filePath + ".tmp")

	s := &NotesService{
		filePath: filePath,
		notes:    make(map[string]string),
	}
	s.load()
	return s
}

func (s *NotesService) load() {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return
	}

	var parsed map[string]string
	if err := json.Unmarshal(data, &parsed); err != nil {
		log.Printf("Failed to decode vpk_notes.json: %v", err)
		return
	}
	s.notes = parsed
}

func (s *NotesService) save() error {
	data, err := json.MarshalIndent(s.notes, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal notes: %v", err)
		return err
	}

	tmpFile := s.filePath + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		log.Printf("Failed to write notes.tmp: %v", err)
		return err
	}

	return os.Rename(tmpFile, s.filePath)
}

func (s *NotesService) GetKey(vpkPath string) string {
	filename := filepath.Base(vpkPath)
	lower := strings.ToLower(filename)
	if strings.HasSuffix(lower, ".disabled") {
		filename = filename[:len(filename)-len(".disabled")]
	}

	workshopID := naming.ExtractWorkshopID(filename)
	if workshopID != "" {
		return "workshop:" + workshopID
	}
	return "local:" + strings.ToLower(filename)
}

func (s *NotesService) GetNote(vpkPath string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.GetKey(vpkPath)
	if note, ok := s.notes[key]; ok && strings.TrimSpace(note) != "" {
		return strings.TrimSpace(note)
	}

	// Fallback: match by base filename directly
	base := strings.ToLower(filepath.Base(vpkPath))
	if strings.HasSuffix(base, ".disabled") {
		base = base[:len(base)-len(".disabled")]
	}
	if note, ok := s.notes[base]; ok && strings.TrimSpace(note) != "" {
		return strings.TrimSpace(note)
	}

	return ""
}

func (s *NotesService) SetNote(vpkPath string, note string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := s.GetKey(vpkPath)
	trimmed := strings.TrimSpace(note)
	if trimmed == "" {
		delete(s.notes, key)
	} else {
		s.notes[key] = trimmed
	}
	return s.save()
}

func (s *NotesService) GetAllNotes() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make(map[string]string, len(s.notes))
	for k, v := range s.notes {
		res[k] = v
	}
	return res
}
