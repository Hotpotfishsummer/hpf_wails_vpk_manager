package vpk

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func buildTestVpk(t *testing.T, dir, filename string, addoninfoBytes []byte) string {
	vpkPath := filepath.Join(dir, filename)

	var treeBuf bytes.Buffer
	if addoninfoBytes != nil {
		treeBuf.WriteString("txt\x00")
		treeBuf.WriteString(" \x00")
		treeBuf.WriteString("addoninfo\x00")

		_ = binary.Write(&treeBuf, binary.LittleEndian, uint32(0))                   // CRC
		_ = binary.Write(&treeBuf, binary.LittleEndian, uint16(len(addoninfoBytes))) // PreloadBytes
		_ = binary.Write(&treeBuf, binary.LittleEndian, uint16(0x7fff))              // ArchiveIndex
		_ = binary.Write(&treeBuf, binary.LittleEndian, uint32(0))                   // EntryOffset
		_ = binary.Write(&treeBuf, binary.LittleEndian, uint32(len(addoninfoBytes))) // EntryLength
		_ = binary.Write(&treeBuf, binary.LittleEndian, uint16(0xffff))              // Terminator
		treeBuf.Write(addoninfoBytes)

		treeBuf.WriteByte(0) // End of filenames
		treeBuf.WriteByte(0) // End of paths
	}
	treeBuf.WriteByte(0) // End of extensions

	var headerBuf bytes.Buffer
	_ = binary.Write(&headerBuf, binary.LittleEndian, uint32(Signature))
	_ = binary.Write(&headerBuf, binary.LittleEndian, uint32(1)) // Version 1
	_ = binary.Write(&headerBuf, binary.LittleEndian, uint32(treeBuf.Len()))

	f, err := os.Create(vpkPath)
	if err != nil {
		t.Fatalf("failed to create test VPK: %v", err)
	}
	defer f.Close()

	_, _ = f.Write(headerBuf.Bytes())
	_, _ = f.Write(treeBuf.Bytes())

	return vpkPath
}

func TestExtractAddonTitle(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("utf8 title", func(t *testing.T) {
		content := []byte("\"addonversion\"\t\"1.0\"\n\"addontitle\"\t\"Real Test Mod\"\n")
		p := buildTestVpk(t, tempDir, "mod1.vpk", content)
		title, err := ExtractAddonTitle(p)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if title != "Real Test Mod" {
			t.Errorf("got %q, want 'Real Test Mod'", title)
		}
	})

	t.Run("chinese utf8 title", func(t *testing.T) {
		content := []byte("addontitle \"中文模组标题\"\n")
		p := buildTestVpk(t, tempDir, "mod2.vpk", content)
		title, err := ExtractAddonTitle(p)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if title != "中文模组标题" {
			t.Errorf("got %q, want '中文模组标题'", title)
		}
	})

	t.Run("gbk encoded title", func(t *testing.T) {
		gbkContent, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte("addontitle \"GBK模组\"\n"))
		p := buildTestVpk(t, tempDir, "mod3.vpk", gbkContent)
		title, err := ExtractAddonTitle(p)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if title != "GBK模组" {
			t.Errorf("got %q, want 'GBK模组'", title)
		}
	})

	t.Run("no addoninfo returns error", func(t *testing.T) {
		p := buildTestVpk(t, tempDir, "empty.vpk", nil)
		_, err := ExtractAddonTitle(p)
		if err == nil {
			t.Errorf("expected error for empty vpk, got nil")
		}
	})
}

func buildMultiFileVpk(t *testing.T, dir, filename string, files map[string][]byte) string {
	vpkPath := filepath.Join(dir, filename)

	// Group files by ext -> path -> filename
	type fileEntry struct {
		name string
		data []byte
	}
	grouped := make(map[string]map[string][]fileEntry)

	for fullPath, data := range files {
		ext := filepath.Ext(fullPath)
		ext = strings.TrimPrefix(ext, ".")
		dirPath := filepath.Dir(fullPath)
		dirPath = strings.ReplaceAll(dirPath, "\\", "/")
		if dirPath == "." {
			dirPath = " "
		}
		base := filepath.Base(fullPath)
		name := strings.TrimSuffix(base, "."+ext)

		if grouped[ext] == nil {
			grouped[ext] = make(map[string][]fileEntry)
		}
		grouped[ext][dirPath] = append(grouped[ext][dirPath], fileEntry{name: name, data: data})
	}

	var treeBuf bytes.Buffer
	for ext, paths := range grouped {
		treeBuf.WriteString(ext + "\x00")
		for p, entries := range paths {
			treeBuf.WriteString(p + "\x00")
			for _, e := range entries {
				treeBuf.WriteString(e.name + "\x00")
				_ = binary.Write(&treeBuf, binary.LittleEndian, uint32(0))             // CRC
				_ = binary.Write(&treeBuf, binary.LittleEndian, uint16(len(e.data)))   // PreloadBytes
				_ = binary.Write(&treeBuf, binary.LittleEndian, uint16(0x7fff))        // ArchiveIndex
				_ = binary.Write(&treeBuf, binary.LittleEndian, uint32(0))             // EntryOffset
				_ = binary.Write(&treeBuf, binary.LittleEndian, uint32(len(e.data)))   // EntryLength
				_ = binary.Write(&treeBuf, binary.LittleEndian, uint16(0xffff))        // Terminator
				treeBuf.Write(e.data)
			}
			treeBuf.WriteByte(0) // End of filenames
		}
		treeBuf.WriteByte(0) // End of paths
	}
	treeBuf.WriteByte(0) // End of extensions

	var headerBuf bytes.Buffer
	_ = binary.Write(&headerBuf, binary.LittleEndian, uint32(Signature))
	_ = binary.Write(&headerBuf, binary.LittleEndian, uint32(1))
	_ = binary.Write(&headerBuf, binary.LittleEndian, uint32(treeBuf.Len()))

	f, err := os.Create(vpkPath)
	if err != nil {
		t.Fatalf("failed to create multi-file test VPK: %v", err)
	}
	defer f.Close()

	_, _ = f.Write(headerBuf.Bytes())
	_, _ = f.Write(treeBuf.Bytes())

	return vpkPath
}

func TestExtractMapNames(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("campaign with mission file and bsp maps", func(t *testing.T) {
		missionContent := `"mission" {
			"DisplayTitle" "Custom Campaign"
			"modes" {
				"coop" {
					"1" {
						"Map" "c1m1_hotel"
						"DisplayName" "旅馆"
					}
					"2" {
						"Map" "c1m2_streets"
						"DisplayName" "街道"
					}
				}
			}
		}`
		files := map[string][]byte{
			"maps/c1m1_hotel.bsp":       []byte("fake-bsp-1"),
			"maps/c1m2_streets.bsp":     []byte("fake-bsp-2"),
			"maps/background01.bsp":     []byte("fake-bg"),
			"missions/campaign.txt":     []byte(missionContent),
			"addoninfo.txt":             []byte("addontitle \"测试战役模组\"\n"),
		}

		p := buildMultiFileVpk(t, tempDir, "campaign.vpk", files)
		title, maps, err := ExtractMetadata(p)
		if err != nil {
			t.Fatalf("ExtractMetadata failed: %v", err)
		}
		if title == nil || *title != "测试战役模组" {
			t.Errorf("title = %v, want '测试战役模组'", title)
		}
		if len(maps) != 2 {
			t.Fatalf("expected 2 maps (background filtered), got %d: %+v", len(maps), maps)
		}
		if maps[0].Name != "c1m1_hotel" || maps[0].DisplayName != "旅馆" {
			t.Errorf("maps[0] = %+v, want c1m1_hotel / 旅馆", maps[0])
		}
		if maps[1].Name != "c1m2_streets" || maps[1].DisplayName != "街道" {
			t.Errorf("maps[1] = %+v, want c1m2_streets / 街道", maps[1])
		}
	})

	t.Run("survival standalone map without mission file", func(t *testing.T) {
		files := map[string][]byte{
			"maps/l4d_rooftop_survival.bsp": []byte("fake-bsp"),
		}
		p := buildMultiFileVpk(t, tempDir, "survival.vpk", files)
		maps, err := ExtractMapNames(p)
		if err != nil {
			t.Fatalf("ExtractMapNames failed: %v", err)
		}
		if len(maps) != 1 {
			t.Fatalf("expected 1 map, got %d", len(maps))
		}
		if maps[0].Name != "l4d_rooftop_survival" {
			t.Errorf("maps[0].Name = %q, want 'l4d_rooftop_survival'", maps[0].Name)
		}
	})

	t.Run("non-map addon has no maps", func(t *testing.T) {
		files := map[string][]byte{
			"models/v_models/v_rifle.mdl": []byte("fake-mdl"),
			"addoninfo.txt":               []byte("addontitle \"AK47 Skin\"\n"),
		}
		p := buildMultiFileVpk(t, tempDir, "weapon.vpk", files)
		maps, err := ExtractMapNames(p)
		if err != nil {
			t.Fatalf("ExtractMapNames failed: %v", err)
		}
		if len(maps) != 0 {
			t.Errorf("expected 0 maps for weapon mod, got %d", len(maps))
		}
	})
}

