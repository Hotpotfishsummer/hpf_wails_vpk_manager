package vpk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
)

const (
	Signature = 0x55aa1234
)

type Entry struct {
	Extension    string
	Path         string
	Filename     string
	CRC          uint32
	PreloadBytes uint16
	ArchiveIndex uint16
	EntryOffset  uint32
	EntryLength  uint32
	Terminator   uint16
	PreloadData  []byte
}

type File struct {
	reader     io.ReadSeeker
	closer     io.Closer
	version    uint32
	headerSize int64
	treeSize   uint32
	entries    map[string]Entry // lowercase "path/filename.ext" -> Entry
}

func readNullString(r io.Reader) (string, error) {
	var buf []byte
	b := make([]byte, 1)
	for {
		_, err := r.Read(b)
		if err != nil {
			return "", err
		}
		if b[0] == 0 {
			break
		}
		buf = append(buf, b[0])
	}
	return string(buf), nil
}

// Open opens a VPK file and reads its header and directory tree
func Open(vpkPath string) (*File, error) {
	f, err := os.Open(vpkPath)
	if err != nil {
		return nil, err
	}

	var sig, ver, treeSize uint32
	if err := binary.Read(f, binary.LittleEndian, &sig); err != nil {
		f.Close()
		return nil, err
	}
	if sig != Signature {
		f.Close()
		return nil, fmt.Errorf("invalid VPK signature: 0x%x", sig)
	}

	if err := binary.Read(f, binary.LittleEndian, &ver); err != nil {
		f.Close()
		return nil, err
	}
	if err := binary.Read(f, binary.LittleEndian, &treeSize); err != nil {
		f.Close()
		return nil, err
	}

	var headerSize int64
	if ver == 1 {
		headerSize = 12
	} else if ver == 2 {
		headerSize = 28
		// skip the remaining 16 bytes of v2 header
		if _, err := f.Seek(headerSize, io.SeekStart); err != nil {
			f.Close()
			return nil, err
		}
	} else {
		f.Close()
		return nil, fmt.Errorf("unsupported VPK version: %d", ver)
	}

	treeBytes := make([]byte, treeSize)
	if _, err := io.ReadFull(f, treeBytes); err != nil {
		f.Close()
		return nil, fmt.Errorf("failed to read VPK directory tree: %w", err)
	}

	tr := bytes.NewReader(treeBytes)
	entries := make(map[string]Entry)

	for {
		ext, err := readNullString(tr)
		if err != nil || ext == "" {
			break
		}

		for {
			path, err := readNullString(tr)
			if err != nil || path == "" {
				break
			}

			for {
				name, err := readNullString(tr)
				if err != nil || name == "" {
					break
				}

				var e Entry
				e.Extension = ext
				e.Path = path
				e.Filename = name

				if err := binary.Read(tr, binary.LittleEndian, &e.CRC); err != nil {
					break
				}
				if err := binary.Read(tr, binary.LittleEndian, &e.PreloadBytes); err != nil {
					break
				}
				if err := binary.Read(tr, binary.LittleEndian, &e.ArchiveIndex); err != nil {
					break
				}
				if err := binary.Read(tr, binary.LittleEndian, &e.EntryOffset); err != nil {
					break
				}
				if err := binary.Read(tr, binary.LittleEndian, &e.EntryLength); err != nil {
					break
				}
				if err := binary.Read(tr, binary.LittleEndian, &e.Terminator); err != nil {
					break
				}

				if e.PreloadBytes > 0 {
					e.PreloadData = make([]byte, e.PreloadBytes)
					if _, err := io.ReadFull(tr, e.PreloadData); err != nil {
						break
					}
				}

				// Normalize full path key
				var fullPath string
				trimmedPath := strings.Trim(path, " /\\")
				if trimmedPath == "" || trimmedPath == " " {
					fullPath = fmt.Sprintf("%s.%s", name, ext)
				} else {
					fullPath = fmt.Sprintf("%s/%s.%s", trimmedPath, name, ext)
				}
				entries[strings.ToLower(fullPath)] = e
			}
		}
	}

	return &File{
		reader:     f,
		closer:     f,
		version:    ver,
		headerSize: headerSize,
		treeSize:   treeSize,
		entries:    entries,
	}, nil
}

// Close closes the underlying VPK file
func (v *File) Close() error {
	if v.closer != nil {
		return v.closer.Close()
	}
	return nil
}

// ReadFile reads the content of a file from the VPK
func (v *File) ReadFile(filename string) ([]byte, error) {
	normKey := strings.ToLower(strings.ReplaceAll(filename, "\\", "/"))
	entry, ok := v.entries[normKey]
	if !ok {
		return nil, os.ErrNotExist
	}

	var result []byte
	if entry.PreloadBytes > 0 {
		result = append(result, entry.PreloadData...)
	}

	remaining := int(entry.EntryLength) - int(entry.PreloadBytes)
	if remaining > 0 && entry.ArchiveIndex == 0x7fff {
		offset := v.headerSize + int64(v.treeSize) + int64(entry.EntryOffset)
		if _, err := v.reader.Seek(offset, io.SeekStart); err != nil {
			return nil, err
		}
		data := make([]byte, remaining)
		if _, err := io.ReadFull(v.reader, data); err != nil {
			return nil, err
		}
		result = append(result, data...)
	}

	return result, nil
}

// HasFile checks if a file exists in the VPK
func (v *File) HasFile(filename string) bool {
	normKey := strings.ToLower(strings.ReplaceAll(filename, "\\", "/"))
	_, ok := v.entries[normKey]
	return ok
}

// DecodeText tries multiple encodings to decode raw text bytes
func DecodeText(raw []byte) string {
	if utf8.Valid(raw) {
		return string(raw)
	}

	// Try GBK
	if decoded, err := simplifiedchinese.GBK.NewDecoder().Bytes(raw); err == nil {
		return string(decoded)
	}

	// Try UTF-16LE
	if len(raw) >= 2 {
		u16le := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
		if decoded, err := u16le.NewDecoder().Bytes(raw); err == nil {
			return string(decoded)
		}
	}

	// Try Big5
	if decoded, err := traditionalchinese.Big5.NewDecoder().Bytes(raw); err == nil {
		return string(decoded)
	}

	// Try ISO-8859-1 (Latin1)
	if decoded, err := charmap.ISO8859_1.NewDecoder().Bytes(raw); err == nil {
		return string(decoded)
	}

	// Fallback to UTF-8 with replacement
	return strings.ToValidUTF8(string(raw), "\uFFFD")
}

// ParseAddonTitle parses addontitle from addoninfo.txt text content
func ParseAddonTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "\"addontitle\"") || strings.HasPrefix(lower, "addontitle") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				return parts[3]
			} else if len(parts) >= 2 {
				return parts[1]
			}
		}
	}
	return ""
}

// ExtractAddonTitle opens the VPK and retrieves the addontitle from addoninfo.txt
func ExtractAddonTitle(vpkPath string) (string, error) {
	v, err := Open(vpkPath)
	if err != nil {
		return "", err
	}
	defer v.Close()

	if !v.HasFile("addoninfo.txt") {
		return "", fmt.Errorf("addoninfo.txt not found")
	}

	raw, err := v.ReadFile("addoninfo.txt")
	if err != nil {
		return "", err
	}

	content := DecodeText(raw)
	title := ParseAddonTitle(content)
	if title == "" {
		return "", fmt.Errorf("addontitle key not found in addoninfo.txt")
	}
	return title, nil
}
