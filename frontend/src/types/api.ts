import type { VpkFile, MapInfo } from '@/stores/vpk-manager'
export type { MapInfo, VpkFile }

/**
 * 后端 API 返回结果的类型定义。
 * 字段与 backend 各 service/api 的返回 dict 一一对应;
 * 标注为可选是因为失败分支通常只带 success/error。
 */

export interface ApiResult {
  success: boolean
  error?: string
  message?: string
}

/** scan_vpk_files */
export interface ScanResult extends ApiResult {
  local_vpks: VpkFile[]
  workshop_vpks: VpkFile[]
}

/** import_compressed_vpks */
export interface ImportResult extends ApiResult {
  imported_count?: number
  errors?: string[] | null
}

/** localize_workshop_vpks */
export interface LocalizeResult extends ApiResult {
  success_count?: number
  errors?: string[] | null
}

/** export_vpks */
export interface ExportResult extends ApiResult {
  file_path?: string
  /** 目标归档已存在且未指定 overwrite 时为 true,由前端决定是否覆盖重试 */
  exists?: boolean
}

/** 远程 SFTP 文件条目(sftp_service.list_vpk_files) */
export interface RemoteFile {
  name: string
  path: string
  size: number
  modified_time: number
  is_disabled: boolean
}

/** list_remote_vpks */
export interface RemoteListResult extends ApiResult {
  files?: RemoteFile[]
}

/** delete_multiple_vpk_files(本地批量)与 delete_multiple_remote_vpks(远程批量) */
export interface BatchDeleteResult extends ApiResult {
  success_count?: number
  error_count?: number
  errors?: string[]
}

/** download_remote_vpks */
export interface DownloadResult extends ApiResult {
  downloaded_count?: number
}

/** upload_local_vpks */
export interface UploadFailedItem {
  file: string
  error: string
}

export interface UploadProgressEvent {
  current_index: number
  total_files: number
  current_file: string
  current_bytes: number
  total_bytes: number
  overall_bytes: number
  overall_total: number
  speed: number
  success_count: number
  failed_count: number
  stage: 'preparing' | 'uploading' | 'completed' | 'error' | 'canceled'
}

export interface UploadResult extends ApiResult {
  uploaded_count?: number
  total_count?: number
  failed_files?: UploadFailedItem[]
  canceled?: boolean
}

/** 自动探测到的 L4D2 客户端 */
export interface DetectedGameClient {
  path: string
  exe_path: string
  source: string
  source_name: string
  mod_count: number
  has_workshop: boolean
  last_modified: number
}

/** detect_l4d2_directories */
export interface DetectGameResult extends ApiResult {
  clients: DetectedGameClient[]
}

/** 游戏启动设置选项 */
export interface LaunchSettings {
  novid: boolean
  chinese_lang: boolean
  direct_launch: boolean
  custom_args: string
  server_connect_cmd?: string
  network_rates_cmd?: string
}

/** VPK 自定义分组 */
export interface VpkGroup {
  id: string
  name: string
  color: string
  collapsed?: boolean
  file_keys: string[]
}

/** OpenAI 兼容大模型 API 设置 */
export interface AiSettings {
  api_key: string
  base_url: string
  model: string
  temperature: number
  max_tokens: number
}

export interface AiChatMessage {
  role: string
  content: string
}

export interface AiChatRequest {
  model: string
  messages: AiChatMessage[]
  temperature?: number
  max_tokens?: number
}

export interface AiResult extends ApiResult {
  content?: string
}

export interface GroupPredictResult extends ApiResult {
  predicted_name: string
  reason?: string
  source: 'ai' | 'heuristic'
}

