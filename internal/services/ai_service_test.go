package services

import (
	"testing"

	"hpf_wails_vpk_manager/internal/models"
)

func TestFormatEndpoint(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "https://api.openai.com/v1/chat/completions"},
		{"https://api.openai.com/v1", "https://api.openai.com/v1/chat/completions"},
		{"https://api.openai.com/v1/", "https://api.openai.com/v1/chat/completions"},
		{"https://api.deepseek.com", "https://api.deepseek.com/v1/chat/completions"},
		{"https://api.deepseek.com/", "https://api.deepseek.com/v1/chat/completions"},
		{"https://my-proxy.com/v1/chat/completions", "https://my-proxy.com/v1/chat/completions"},
	}

	for _, tc := range tests {
		actual := formatEndpoint(tc.input)
		if actual != tc.expected {
			t.Errorf("formatEndpoint(%q) = %q; expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestCleanPredictedName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"生还者角色合集"`, "生还者角色合集"},
		{"【特感皮肤包】\n额外解释", "特感皮肤包"},
		{"`武器音效合集`", "武器音效合集"},
		{"  “战役地图包”  ", "战役地图包"},
	}

	for _, tc := range tests {
		actual := cleanPredictedName(tc.input)
		if actual != tc.expected {
			t.Errorf("cleanPredictedName(%q) = %q; expected %q", tc.input, actual, tc.expected)
		}
	}
}

func TestPredictGroupNameHeuristic(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		res := PredictGroupNameHeuristic([]models.VpkFile{})
		if res != "新模组分组" {
			t.Errorf("expected '新模组分组', got %q", res)
		}
	})

	t.Run("character zoey skin", func(t *testing.T) {
		files := []models.VpkFile{
			{Name: "zoey_casual_outfit.vpk", DisplayName: "Zoey Casual"},
			{Name: "zoey_hair_mod.vpk", DisplayName: "Zoey Hair"},
		}
		res := PredictGroupNameHeuristic(files)
		if res != "佐伊皮肤套装" {
			t.Errorf("expected '佐伊皮肤套装', got %q", res)
		}
	})

	t.Run("weapons ak47", func(t *testing.T) {
		files := []models.VpkFile{
			{Name: "cod_ak47_sound.vpk", DisplayName: "AK47 Sound Replacer"},
			{Name: "m16_reanimation.vpk", DisplayName: "M16 Animation"},
		}
		res := PredictGroupNameHeuristic(files)
		if res != "武器枪械模组包" {
			t.Errorf("expected '武器枪械模组包', got %q", res)
		}
	})

	t.Run("infected witch tank", func(t *testing.T) {
		files := []models.VpkFile{
			{Name: "cute_witch.vpk", DisplayName: "Witch Model"},
			{Name: "hulk_tank.vpk", DisplayName: "Tank Sound"},
		}
		res := PredictGroupNameHeuristic(files)
		if res != "Tank坦克特感模组" {
			t.Errorf("expected 'Tank坦克特感模组', got %q", res)
		}
	})

	t.Run("campaign maps", func(t *testing.T) {
		files := []models.VpkFile{
			{Name: "suicide_blitz_part1.vpk", DisplayName: "Suicide Blitz", Maps: []models.MapInfo{{Name: "c1m1"}, {Name: "c1m2"}}},
			{Name: "suicide_blitz_part2.vpk", DisplayName: "Suicide Blitz Part 2"},
		}
		res := PredictGroupNameHeuristic(files)
		if res != "战役地图模组包" {
			t.Errorf("expected '战役地图模组包', got %q", res)
		}
	})

	t.Run("common stem fallback", func(t *testing.T) {
		files := []models.VpkFile{
			{Name: "genshin_pack_01.vpk", DisplayName: "genshin_pack_01"},
			{Name: "genshin_pack_02.vpk", DisplayName: "genshin_pack_02"},
		}
		res := PredictGroupNameHeuristic(files)
		if res != "genshin模组包" {
			t.Errorf("expected 'genshin模组包', got %q", res)
		}
	})
}
