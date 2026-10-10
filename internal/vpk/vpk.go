package vpk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"

	"hpf_wails_vpk_manager/internal/models"
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

// KVNode represents a KeyValues tree node
type KVNode struct {
	Key      string
	Value    string
	Children []*KVNode
}

func tokenizeKV(text string) []string {
	var tokens []string
	n := len(text)
	i := 0
	for i < n {
		c := text[i]

		// Skip whitespace
		if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			i++
			continue
		}

		// Skip comments //
		if c == '/' && i+1 < n && text[i+1] == '/' {
			for i < n && text[i] != '\n' {
				i++
			}
			continue
		}

		// Braces
		if c == '{' || c == '}' {
			tokens = append(tokens, string(c))
			i++
			continue
		}

		// Quoted string
		if c == '"' {
			i++
			var buf []byte
			for i < n && text[i] != '"' {
				if text[i] == '\\' && i+1 < n {
					i++
					buf = append(buf, text[i])
				} else {
					buf = append(buf, text[i])
				}
				i++
			}
			if i < n && text[i] == '"' {
				i++
			}
			tokens = append(tokens, string(buf))
			continue
		}

		// Unquoted word
		start := i
		for i < n && text[i] != ' ' && text[i] != '\t' && text[i] != '\r' && text[i] != '\n' && text[i] != '{' && text[i] != '}' && text[i] != '"' {
			if text[i] == '/' && i+1 < n && text[i+1] == '/' {
				break
			}
			i++
		}
		if i > start {
			tokens = append(tokens, text[start:i])
		}
	}
	return tokens
}

func parseKVTokens(tokens []string) []*KVNode {
	idx := 0
	var parseBlock func() []*KVNode
	parseBlock = func() []*KVNode {
		var list []*KVNode
		for idx < len(tokens) {
			token := tokens[idx]
			if token == "}" {
				idx++
				break
			}
			if token == "{" {
				idx++
				continue
			}

			key := token
			idx++
			if idx >= len(tokens) {
				list = append(list, &KVNode{Key: key})
				break
			}

			nextToken := tokens[idx]
			if nextToken == "{" {
				idx++
				children := parseBlock()
				list = append(list, &KVNode{Key: key, Children: children})
			} else if nextToken == "}" {
				list = append(list, &KVNode{Key: key})
				idx++
				break
			} else {
				list = append(list, &KVNode{Key: key, Value: nextToken})
				idx++
			}
		}
		return list
	}
	return parseBlock()
}

func findKVNode(nodes []*KVNode, key string) *KVNode {
	for _, n := range nodes {
		if strings.EqualFold(n.Key, key) {
			return n
		}
	}
	return nil
}

func cleanChapterToken(token string, tokens map[string]string) string {
	if token == "" {
		return ""
	}
	// Try direct lookup
	if translated, ok := tokens[strings.ToLower(token)]; ok && translated != "" {
		return translated
	}
	trimmed := strings.TrimPrefix(token, "#")
	if translated, ok := tokens[strings.ToLower(trimmed)]; ok && translated != "" {
		return translated
	}
	// Fallback: cleanup valve prefixes
	if strings.HasPrefix(token, "#") {
		clean := strings.TrimPrefix(token, "#L4D360UI_Chapter_")
		clean = strings.TrimPrefix(clean, "#L4D360UI_")
		clean = strings.TrimPrefix(clean, "#")
		return clean
	}
	return token
}

func parseMissionMaps(content string, tokens map[string]string) []models.MapInfo {
	tokensList := tokenizeKV(content)
	if len(tokensList) == 0 {
		return nil
	}

	rootNodes := parseKVTokens(tokensList)
	// Find mission node (usually the root object or root list)
	var missionNode *KVNode
	for _, node := range rootNodes {
		if strings.EqualFold(node.Key, "mission") && len(node.Children) > 0 {
			missionNode = node
			break
		}
	}
	if missionNode == nil && len(rootNodes) > 0 {
		missionNode = rootNodes[0]
	}
	if missionNode == nil {
		return nil
	}

	modesNode := findKVNode(missionNode.Children, "modes")
	if modesNode == nil {
		return nil
	}

	// Priority: coop > campaign > versus > survival > scavenge > first available
	var targetMode *KVNode
	for _, modeName := range []string{"coop", "campaign", "versus", "survival", "scavenge"} {
		if m := findKVNode(modesNode.Children, modeName); m != nil && len(m.Children) > 0 {
			targetMode = m
			break
		}
	}
	if targetMode == nil && len(modesNode.Children) > 0 {
		targetMode = modesNode.Children[0]
	}
	if targetMode == nil {
		return nil
	}

	var maps []models.MapInfo
	for _, chapter := range targetMode.Children {
		var mapName, displayName string
		if len(chapter.Children) > 0 {
			for _, prop := range chapter.Children {
				if strings.EqualFold(prop.Key, "map") {
					mapName = prop.Value
				} else if strings.EqualFold(prop.Key, "displayname") {
					displayName = prop.Value
				}
			}
		} else if strings.EqualFold(chapter.Key, "map") {
			mapName = chapter.Value
		}

		mapName = strings.TrimSpace(mapName)
		mapName = strings.TrimSuffix(mapName, ".bsp")
		if mapName != "" {
			cleanedDisplay := cleanChapterToken(strings.TrimSpace(displayName), tokens)
			maps = append(maps, models.MapInfo{
				Name:        mapName,
				DisplayName: cleanedDisplay,
			})
		}
	}

	return maps
}

func naturalCompare(a, b string) bool {
	re := regexp.MustCompile(`(\d+|\D+)`)
	aParts := re.FindAllString(strings.ToLower(a), -1)
	bParts := re.FindAllString(strings.ToLower(b), -1)

	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aPart := aParts[i]
		bPart := bParts[i]

		aNum, aErr := strconv.Atoi(aPart)
		bNum, bErr := strconv.Atoi(bPart)

		if aErr == nil && bErr == nil {
			if aNum != bNum {
				return aNum < bNum
			}
		} else {
			if aPart != bPart {
				return aPart < bPart
			}
		}
	}
	return len(aParts) < len(bParts)
}

func isBackgroundMap(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasPrefix(lower, "background") ||
		strings.HasPrefix(lower, "l4d_background") ||
		strings.HasPrefix(lower, "l4d2_background") ||
		strings.Contains(lower, "_background")
}

// ExtractMaps parses all map names and chapter information from the VPK
func (v *File) ExtractMaps() []models.MapInfo {
	// 1. Collect all BSP files in VPK
	var bspMapNames []string
	seenBsp := make(map[string]bool)
	for _, entry := range v.entries {
		if strings.EqualFold(entry.Extension, "bsp") {
			name := strings.TrimSpace(entry.Filename)
			if name == "" {
				continue
			}
			lower := strings.ToLower(name)
			if !seenBsp[lower] {
				seenBsp[lower] = true
				bspMapNames = append(bspMapNames, name)
			}
		}
	}

	// 2. Load resource tokens if any (for translating chapter tokens like #L4D360UI_...)
	tokens := make(map[string]string)
	for key, entry := range v.entries {
		if strings.EqualFold(entry.Extension, "txt") && strings.HasPrefix(strings.ToLower(entry.Path), "resource") {
			if raw, err := v.ReadFile(key); err == nil {
				resKV := parseKVTokens(tokenizeKV(DecodeText(raw)))
				var walkTokens func(nodes []*KVNode)
				walkTokens = func(nodes []*KVNode) {
					for _, node := range nodes {
						if len(node.Children) > 0 {
							walkTokens(node.Children)
						} else if node.Value != "" {
							tokens[strings.ToLower(node.Key)] = node.Value
							tokens[strings.ToLower(strings.TrimPrefix(node.Key, "#"))] = node.Value
						}
					}
				}
				walkTokens(resKV)
			}
		}
	}

	// 3. Collect chapters from mission files
	var missionMaps []models.MapInfo
	for key, entry := range v.entries {
		if strings.EqualFold(entry.Extension, "txt") && (strings.HasPrefix(strings.ToLower(entry.Path), "missions") || strings.Contains(strings.ToLower(entry.Path), "missions")) {
			if raw, err := v.ReadFile(key); err == nil {
				parsed := parseMissionMaps(DecodeText(raw), tokens)
				if len(parsed) > 0 {
					missionMaps = append(missionMaps, parsed...)
				}
			}
		}
	}

	// 4. Merge mission maps and BSP maps
	var result []models.MapInfo
	used := make(map[string]bool)

	if len(missionMaps) > 0 {
		for _, m := range missionMaps {
			lower := strings.ToLower(m.Name)
			// If BSP exists or there are no BSPs in VPK (e.g. mission-only mod), include it
			if len(bspMapNames) == 0 || seenBsp[lower] {
				if !used[lower] {
					used[lower] = true
					result = append(result, m)
				}
			}
		}
	}

	// Append any remaining BSPs that were not mentioned in mission files
	sort.Slice(bspMapNames, func(i, j int) bool {
		return naturalCompare(bspMapNames[i], bspMapNames[j])
	})
	for _, bsp := range bspMapNames {
		lower := strings.ToLower(bsp)
		if !used[lower] {
			used[lower] = true
			result = append(result, models.MapInfo{
				Name: bsp,
			})
		}
	}

	// 5. Filter out background menu maps if there are other playable maps
	if len(result) > 1 {
		var filtered []models.MapInfo
		for _, m := range result {
			if !isBackgroundMap(m.Name) {
				filtered = append(filtered, m)
			}
		}
		if len(filtered) > 0 {
			result = filtered
		}
	}

	return result
}

// ExtractMetadata opens the VPK once and extracts both addon title and maps
func ExtractMetadata(vpkPath string) (*string, []models.MapInfo, error) {
	v, err := Open(vpkPath)
	if err != nil {
		return nil, nil, err
	}
	defer v.Close()

	var title *string
	if v.HasFile("addoninfo.txt") {
		if raw, err := v.ReadFile("addoninfo.txt"); err == nil {
			t := ParseAddonTitle(DecodeText(raw))
			if t != "" {
				title = &t
			}
		}
	}

	maps := v.ExtractMaps()
	return title, maps, nil
}

// ExtractMapNames opens the VPK and retrieves all map names
func ExtractMapNames(vpkPath string) ([]models.MapInfo, error) {
	v, err := Open(vpkPath)
	if err != nil {
		return nil, err
	}
	defer v.Close()

	return v.ExtractMaps(), nil
}
