package naming

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	illegalCharsRegex = regexp.MustCompile(`[^a-zA-Z0-9_-]`)
	multiUnderscores  = regexp.MustCompile(`_{2,}`)
	workshopIDRegex   = regexp.MustCompile(`(\d{8,})`)
)

const maxBaseNameLen = 120

// ImageExtensions supported for VPK thumbnails and associated assets
var ImageExtensions = []string{".jpg", ".jpeg", ".png", ".tga", ".bmp", ".gif"}

// AssociatedSyncExtensions used when disabling, deleting, or renaming VPK files
var AssociatedSyncExtensions = []string{".jpg", ".png", ".jpeg", ".txt"}

type ImageMatch struct {
	FullPath  string
	Extension string
}

// ExtractWorkshopID extracts an 8+ digit ID from the filename if present
func ExtractWorkshopID(filename string) string {
	matches := workshopIDRegex.FindStringSubmatch(filename)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// VpkStemForAssociation calculates the stem used for thumbnail/association matching
func VpkStemForAssociation(filename string) string {
	name := filename
	lower := strings.ToLower(name)
	if strings.HasSuffix(lower, ".disabled") {
		name = name[:len(name)-len(".disabled")]
	}
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext)
}

// SanitizeAddonTitle cleans a title into L4D2-safe characters (ASCII alphanumeric, underscore, hyphen)
func SanitizeAddonTitle(title string) string {
	if title == "" {
		return ""
	}
	s := strings.TrimSpace(strings.ReplaceAll(title, " ", "_"))
	s = strings.TrimRight(s, ".")
	s = illegalCharsRegex.ReplaceAllString(s, "_")
	s = multiUnderscores.ReplaceAllString(s, "_")
	return strings.TrimRight(s, ".")
}

// BuildExportBaseName generates a sanitized base filename based on addontitle and workshop ID
func BuildExportBaseName(addontitle, filename, fallbackStem string) string {
	workshopID := ExtractWorkshopID(filename)

	safeFallback := SanitizeAddonTitle(fallbackStem)
	if len(safeFallback) > maxBaseNameLen {
		safeFallback = safeFallback[:maxBaseNameLen]
	}

	nonIDPart := safeFallback
	if workshopID != "" {
		nonIDPart = strings.ReplaceAll(nonIDPart, workshopID, "")
	}
	if strings.Trim(nonIDPart, "_-") == "" {
		if workshopID != "" {
			safeFallback = workshopID
		} else if strings.Trim(safeFallback, "_-") == "" {
			safeFallback = "addon"
		}
	}

	if addontitle == "" {
		return safeFallback
	}

	baseName := SanitizeAddonTitle(addontitle)
	if len(baseName) > maxBaseNameLen {
		baseName = baseName[:maxBaseNameLen]
	}
	if strings.Trim(baseName, "_-") == "" {
		baseName = safeFallback
	}

	if workshopID != "" && !strings.Contains(baseName, workshopID) {
		baseName = baseName + "-" + workshopID
	}
	return baseName
}

// FindAssociatedImages scans the directory for images associated with the VPK file (case-insensitive)
func FindAssociatedImages(dir, vpkFilename string, stripVpkSuffix bool) []ImageMatch {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	stem := strings.ToLower(filepath.Base(vpkFilename))
	if strings.HasSuffix(stem, ".disabled") {
		stem = stem[:len(stem)-len(".disabled")]
	}
	ext := strings.ToLower(filepath.Ext(stem))
	stem = strings.TrimSuffix(stem, ext)

	if stripVpkSuffix && strings.HasSuffix(stem, ".vpk") {
		stem = strings.TrimSuffix(stem, ".vpk")
	}

	allowedExts := make(map[string]bool)
	for _, e := range ImageExtensions {
		allowedExts[strings.ToLower(e)] = true
	}

	seenExts := make(map[string]bool)
	var matches []ImageMatch

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		entryName := entry.Name()
		entryLower := strings.ToLower(entryName)
		entryExt := strings.ToLower(filepath.Ext(entryLower))

		if !allowedExts[entryExt] {
			continue
		}

		entryStem := strings.TrimSuffix(entryLower, entryExt)
		if entryStem == stem {
			if !seenExts[entryExt] {
				seenExts[entryExt] = true
				matches = append(matches, ImageMatch{
					FullPath:  filepath.Join(dir, entryName),
					Extension: filepath.Ext(entryName),
				})
			}
		}
	}

	return matches
}
