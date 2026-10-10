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
		{"  x  ", "x"},
		{"x\ty", "x_y"},
		{"Trailing.Dot.", "Trailing_Dot"},
		{"全自动连喷音效（超大声！特别好用）", "全自动连喷音效_超大声_特别好用"},
		{"【生化危机】吉尔-人物模型", "生化危机_吉尔-人物模型"},
		{"AK47: 涂装/二次元", "AK47_涂装_二次元"},
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
		{"二次元皮肤", "99887766.vpk", "99887766", "二次元皮肤-99887766"},
	}
	for _, tc := range tests {
		actual := BuildExportBaseName(tc.title, tc.filename, tc.fallbackStem)
		if actual != tc.expected {
			t.Errorf("BuildExportBaseName(%q, %q, %q) = %q; want %q", tc.title, tc.filename, tc.fallbackStem, actual, tc.expected)
		}
	}
}

func TestBuildExportBaseNameWithNote(t *testing.T) {
	tests := []struct {
		title        string
		filename     string
		fallbackStem string
		note         string
		expected     string
	}{
		// 无标题，纯数字文件，有中文备注
		{"", "123456789.vpk", "123456789", "全自动连喷", "全自动连喷-123456789"},
		// 只有英文标题，有中文备注
		{"Auto Shotgun", "123456789.vpk", "123456789", "连喷音效", "Auto_Shotgun_连喷音效-123456789"},
		// 备注中带有各种中文标点与特殊符号
		{"M16", "234567890.vpk", "234567890", "【好用！】消音+换弹（自制）", "M16_好用_消音_换弹_自制-234567890"},
		// 本地文件无 workshop ID
		{"Weapon", "custom.vpk", "custom", "雷神涂装", "Weapon_雷神涂装"},
		// 备注和标题重复时避免冗余
		{"雷神涂装", "custom.vpk", "custom", "雷神涂装", "雷神涂装"},
	}
	for _, tc := range tests {
		actual := BuildExportBaseNameWithNote(tc.title, tc.filename, tc.fallbackStem, tc.note)
		if actual != tc.expected {
			t.Errorf("BuildExportBaseNameWithNote(%q, %q, %q, %q) = %q; want %q", tc.title, tc.filename, tc.fallbackStem, tc.note, actual, tc.expected)
		}
	}
}
