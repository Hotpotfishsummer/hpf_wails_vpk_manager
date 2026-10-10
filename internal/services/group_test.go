package services

import (
	"testing"

	"hpf_wails_vpk_manager/internal/models"
)

func TestGroupService(t *testing.T) {
	tempDir := t.TempDir()
	svc := NewGroupService(tempDir)

	// Initially empty
	groups := svc.GetGroups()
	if len(groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(groups))
	}

	// Save groups
	newGroups := []models.VpkGroup{
		{
			ID:        "grp_1",
			Name:      "初音未来套装",
			Color:     "#0ea5e9",
			Collapsed: false,
			FileKeys:  []string{"local:miku_zoey.vpk", "local:miku_voice.vpk"},
		},
		{
			ID:        "grp_2",
			Name:      "高清武器包",
			Color:     "#10b981",
			Collapsed: true,
			FileKeys:  []string{"local:hd_weapons.vpk"},
		},
	}

	if err := svc.SaveGroups(newGroups); err != nil {
		t.Fatalf("SaveGroups failed: %v", err)
	}

	loaded := svc.GetGroups()
	if len(loaded) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(loaded))
	}
	if loaded[0].Name != "初音未来套装" || loaded[0].Color != "#0ea5e9" || len(loaded[0].FileKeys) != 2 {
		t.Errorf("unexpected group 0: %+v", loaded[0])
	}

	// Reload from disk with fresh service instance
	svc2 := NewGroupService(tempDir)
	reloaded := svc2.GetGroups()
	if len(reloaded) != 2 {
		t.Fatalf("reloaded service expected 2 groups, got %d", len(reloaded))
	}
	if reloaded[1].Name != "高清武器包" || !reloaded[1].Collapsed {
		t.Errorf("unexpected reloaded group 1: %+v", reloaded[1])
	}
}
