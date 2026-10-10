package models

// MapInfo represents a Left 4 Dead 2 map / chapter inside a VPK
type MapInfo struct {
	Name        string `json:"name"`                  // e.g. "c1m1_hotel"
	DisplayName string `json:"display_name,omitempty"` // e.g. "Hotel" / "旅馆"
}

// VpkFile represents metadata for a local or workshop VPK file
type VpkFile struct {
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	Size          int64     `json:"size"`
	ModifiedTime  float64   `json:"modified_time"`
	IsValid       bool      `json:"is_valid"`
	IsDisabled    bool      `json:"is_disabled"`
	ThumbnailPath *string   `json:"thumbnail_path,omitempty"`
	AddonTitle    *string   `json:"addontitle,omitempty"`
	DisplayName   string    `json:"display_name"`
	Note          string    `json:"note,omitempty"`
	Maps          []MapInfo `json:"maps,omitempty"`
}

// ApiResult is the standard response structure
type ApiResult struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

// ScanResult contains local and workshop VPK files
type ScanResult struct {
	ApiResult
	LocalVpks    []VpkFile `json:"local_vpks"`
	WorkshopVpks []VpkFile `json:"workshop_vpks"`
}

// ImportResult contains statistics of an archive import operation
type ImportResult struct {
	ApiResult
	ImportedCount int      `json:"imported_count"`
	Errors        []string `json:"errors,omitempty"`
}

// LocalizeResult contains statistics of a workshop localization operation
type LocalizeResult struct {
	ApiResult
	SuccessCount int      `json:"success_count"`
	Errors       []string `json:"errors,omitempty"`
}

// ExportResult contains the exported file path and optional existence check flag
type ExportResult struct {
	ApiResult
	FilePath string `json:"file_path,omitempty"`
	Exists   bool   `json:"exists,omitempty"`
}

// RemoteFile represents a file on the remote SFTP server
type RemoteFile struct {
	Name         string  `json:"name"`
	Path         string  `json:"path"`
	Size         float64 `json:"size"` // MB
	ModifiedTime float64 `json:"modified_time"`
	IsDisabled   bool    `json:"is_disabled"`
}

// RemoteListResult contains remote SFTP files list
type RemoteListResult struct {
	ApiResult
	Files []RemoteFile `json:"files,omitempty"`
}

// BatchDeleteResult contains statistics for batch delete operations
type BatchDeleteResult struct {
	ApiResult
	SuccessCount int      `json:"success_count"`
	ErrorCount   int      `json:"error_count"`
	Errors       []string `json:"errors,omitempty"`
}

// DownloadResult contains download count
type DownloadResult struct {
	ApiResult
	DownloadedCount int `json:"downloaded_count"`
}

// UploadResult contains upload count
type UploadResult struct {
	ApiResult
	UploadedCount int `json:"uploaded_count"`
}

// RestartResult contains L4D2 server restart status
type RestartResult struct {
	ApiResult
	DurationS float64 `json:"duration_s,omitempty"`
}

// SftpSettings represents SFTP connection config
type SftpSettings struct {
	Address    string `json:"address"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AddonsPath string `json:"addons_path"`
}
