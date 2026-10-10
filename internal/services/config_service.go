package services

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed autoexec.cfg.example
var autoexecTemplate string

type ConfigService struct{}

func NewConfigService() *ConfigService {
	return &ConfigService{}
}

func (c *ConfigService) ApplyAutoexecTemplate(l4d2Dir string) bool {
	if l4d2Dir == "" {
		return false
	}

	var cfgDir string
	if strings.ToLower(filepath.Base(l4d2Dir)) == "left4dead2" {
		cfgDir = filepath.Join(l4d2Dir, "cfg")
	} else {
		cfgDir = filepath.Join(l4d2Dir, "left4dead2", "cfg")
	}

	info, err := os.Stat(cfgDir)
	if err != nil || !info.IsDir() {
		log.Printf("cfg directory does not exist: %s", cfgDir)
		return false
	}

	autoexecPath := filepath.Join(cfgDir, "autoexec.cfg")
	bakPath := autoexecPath + ".bak"

	// If autoexec.cfg exists, check if update is needed
	if existingData, err := os.ReadFile(autoexecPath); err == nil {
		if string(existingData) == autoexecTemplate {
			log.Println("autoexec.cfg is already up-to-date")
			return true
		}
		// Backup existing file
		if err := os.WriteFile(bakPath, existingData, 0644); err != nil {
			log.Printf("Failed to backup autoexec.cfg: %v", err)
		} else {
			log.Printf("Backed up autoexec.cfg to %s", bakPath)
		}
	}

	// Write template to autoexec.cfg
	if err := os.WriteFile(autoexecPath, []byte(autoexecTemplate), 0644); err != nil {
		log.Printf("Failed to write autoexec.cfg: %v", err)
		return false
	}

	log.Printf("Applied autoexec template to %s", autoexecPath)
	return true
}
