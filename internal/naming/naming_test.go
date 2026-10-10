package naming

import (
	"testing"
)

func TestExtractWorkshopID(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"12345678.vpk", "12345678"},
		{"workshop_123456789_mod.vpk", "123456789"},
		{"short123.vpk", ""},
		{"no_digits.vpk", ""},
	}
	for _, tc := range tests {
		actual := ExtractWorkshopID(tc.filename)
		if actual != tc.expected {
			t.Errorf("ExtractWorkshopID(%q) = %q; want %q", tc.filename, actual, tc.expected)
		}
	}
}

func TestSanitizeAddonTitle(t *testing.T) {
	tests := []struct {
		title    string
		expected string
	}{
		{"My Cool Mod", "My_Cool_Mod"},
		{"a  b??c", "a_b_c"},
		{"  x  ", "_x_"},
		{"x\ty", "x_y"},
		{"Trailing.Dot.", "Trailing_Dot"},
	}
	for _, tc := range tests {
		actual := SanitizeAddonTitle(tc.title)
		if actual != tc.expected {
			t.Errorf("SanitizeAddonTitle(%q) = %q; want %q", tc.title, actual, tc.expected)
		}
	}
}

func TestBuildExportBaseName(t *testing.T) {
	tests := []struct {
		title        string
		filename     string
		fallbackStem string
		expected     string
	}{
		{"Survivor Skin", "12345678.vpk", "12345678", "Survivor_Skin-12345678"},
		{"", "12345678.vpk", "12345678", "12345678"},
		{"Weapon Mod", "ak47.vpk", "ak47", "Weapon_Mod"},
		{"???*", "backup.vpk", "backup", "backup"},
	}
	for _, tc := range tests {
		actual := BuildExportBaseName(tc.title, tc.filename, tc.fallbackStem)
		if actual != tc.expected {
			t.Errorf("BuildExportBaseName(%q, %q, %q) = %q; want %q", tc.title, tc.filename, tc.fallbackStem, actual, tc.expected)
		}
	}
}
