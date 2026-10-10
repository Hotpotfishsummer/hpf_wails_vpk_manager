package vpk

import (
	"bytes"
	"encoding/binary"
	"os"
	"path/filepath"
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
