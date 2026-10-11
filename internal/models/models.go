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

// UploadFailedItem represents a file that failed to upload and its error reason
type UploadFailedItem struct {
	File  string `json:"file"`
	Error string `json:"error"`
}

// UploadResult contains upload count, total count, failed files, and cancel status
type UploadResult struct {
	ApiResult
	UploadedCount int                `json:"uploaded_count"`
	TotalCount    int                `json:"total_count"`
	FailedFiles   []UploadFailedItem `json:"failed_files,omitempty"`
	Canceled      bool               `json:"canceled,omitempty"`
}

// UploadProgressEvent represents real-time progress of SFTP upload
type UploadProgressEvent struct {
	CurrentIndex int    `json:"current_index"` // 当前文件序号 (从 1 开始)
	TotalFiles   int    `json:"total_files"`   // 总选中的 VPK 文件数
	CurrentFile  string `json:"current_file"`  // 当前正在上传的文件名
	CurrentBytes int64  `json:"current_bytes"` // 当前文件已上传字节
	TotalBytes   int64  `json:"total_bytes"`   // 当前文件总字节
	OverallBytes int64  `json:"overall_bytes"` // 所有文件累计已上传字节
	OverallTotal int64  `json:"overall_total"` // 所有文件累计总字节
	Speed        int64  `json:"speed"`         // 当前上传速度 (字节/秒)
	SuccessCount int    `json:"success_count"` // 已成功上传的文件数
	FailedCount  int    `json:"failed_count"`  // 已失败的文件数
	Stage        string `json:"stage"`         // "preparing" | "uploading" | "completed" | "error" | "canceled"
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

// DetectedGameClient represents a detected Left 4 Dead 2 installation
type DetectedGameClient struct {
	Path         string `json:"path"`
	ExePath      string `json:"exe_path"`
	Source       string `json:"source"`
	SourceName   string `json:"source_name"`
	ModCount     int    `json:"mod_count"`
	HasWorkshop  bool   `json:"has_workshop"`
	LastModified int64  `json:"last_modified"`
}

// DetectGameResult contains detected clients list
type DetectGameResult struct {
	ApiResult
	Clients []DetectedGameClient `json:"clients"`
}

// LaunchSettings represents game startup parameters configuration
type LaunchSettings struct {
	NoVid        bool   `json:"novid"`          // 跳过片头动画 (-novid)
	ChineseLang  bool   `json:"chinese_lang"`   // 强制中文启动 (-language schinese)
	DirectLaunch bool   `json:"direct_launch"`  // 强制直接启动主程序 (跳过 Steam 唤起)
	CustomArgs   string `json:"custom_args"`    // 自定义额外参数
}

func DefaultLaunchSettings() LaunchSettings {
	return LaunchSettings{
		NoVid:        true,
		ChineseLang:  true,
		DirectLaunch: false,
		CustomArgs:   "",
	}
}

// VpkGroup represents a user-defined group of VPK files
type VpkGroup struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Color     string   `json:"color"`
	Collapsed bool     `json:"collapsed"`
	FileKeys  []string `json:"file_keys"`
}

// AiSettings represents configuration for OpenAI-compatible API
type AiSettings struct {
	ApiKey      string  `json:"api_key"`
	BaseURL     string  `json:"base_url"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

func DefaultAiSettings() AiSettings {
	return AiSettings{
		ApiKey:      "",
		BaseURL:     "https://api.openai.com/v1",
		Model:       "gpt-4o-mini",
		Temperature: 0.7,
		MaxTokens:   500,
	}
}

// AiChatMessage represents chat message for completions
type AiChatMessage struct {
	Role    string `json:"role"` // "system", "user", "assistant"
	Content string `json:"content"`
}

// AiChatRequest represents OpenAI chat completions request payload
type AiChatRequest struct {
	Model       string          `json:"model"`
	Messages    []AiChatMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

// AiResult represents generic result returned to frontend
type AiResult struct {
	ApiResult
	Content string `json:"content,omitempty"`
}

// GroupPredictResult represents result of group prediction
type GroupPredictResult struct {
	ApiResult
	PredictedName string `json:"predicted_name"`
	Reason        string `json:"reason,omitempty"`
	Source        string `json:"source"` // "ai" or "heuristic"
}

