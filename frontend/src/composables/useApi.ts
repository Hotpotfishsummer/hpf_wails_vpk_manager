import { ref, onMounted } from 'vue'
import * as AppApi from '../../wailsjs/go/main/App'
import { EventsOn } from '../../wailsjs/runtime/runtime'
import type {
  ApiResult,
  BatchDeleteResult,
  DownloadResult,
  ExportResult,
  ImportResult,
  LocalizeResult,
  RemoteListResult,
  ScanResult,
  UploadResult,
} from '@/types/api'

// 监听 Wails 后端推送的导出进度事件
EventsOn('export-progress', (data: any) => {
  if (typeof (window as any).__onExportProgress === 'function') {
    ;(window as any).__onExportProgress(data)
  }
})

// Wails API 适配器（保持与原 PyWebViewApi 完全一致的签名和返回值）
const wailsAdapter = {
  get_last_directory: () => AppApi.GetLastDirectory(),
  browse_directory: async () => {
    const res = await AppApi.BrowseDirectory()
    return res || null
  },
  get_downloads_directory: () => AppApi.GetDownloadsDirectory(),
  scan_vpk_files: (directory: string) => AppApi.ScanVpkFiles(directory) as Promise<ScanResult>,
  export_vpks: (file_paths: string[], output_dir: string, better_naming = false, overwrite = false) =>
    AppApi.ExportVpks(file_paths, output_dir, better_naming, overwrite) as Promise<ExportResult>,
  import_compressed_vpks: () => AppApi.ImportCompressedVpks() as Promise<ImportResult>,
  localize_workshop_vpks: (file_paths: string[], better_naming = false) =>
    AppApi.LocalizeWorkshopVpks(file_paths, better_naming) as Promise<LocalizeResult>,
  open_directory: (path: string) => AppApi.OpenDirectory(path),
  delete_vpk_file: (file_path: string, disable = true) =>
    AppApi.DeleteVpkFile(file_path, disable) as Promise<ApiResult>,
  delete_multiple_vpk_files: (file_paths: string[], disable = true) =>
    AppApi.DeleteMultipleVpkFiles(file_paths, disable) as Promise<BatchDeleteResult>,
  enable_vpk_file: (file_path: string) => AppApi.EnableVpkFile(file_path) as Promise<ApiResult>,
  reset_app_data: () => AppApi.ResetAppData(),
  log_from_js: (level: string, message: string) => AppApi.LogFromJs(level, message),
  get_sftp_settings: () => AppApi.GetSftpSettings() as Promise<Record<string, string>>,
  save_sftp_settings: (settings: Record<string, string>) =>
    AppApi.SaveSftpSettings(settings as any) as Promise<ApiResult>,
  list_remote_vpks: () => AppApi.ListRemoteVpks() as Promise<RemoteListResult>,
  delete_remote_vpk: (file_path: string) => AppApi.DeleteRemoteVpk(file_path) as Promise<ApiResult>,
  delete_multiple_remote_vpks: (file_paths: string[]) =>
    AppApi.DeleteMultipleRemoteVpks(file_paths) as Promise<BatchDeleteResult>,
  rename_remote_vpk: (old_path: string, new_path: string) =>
    AppApi.RenameRemoteVpk(old_path, new_path) as Promise<ApiResult>,
  download_remote_vpks: (remote_file_paths: string[]) =>
    AppApi.DownloadRemoteVpks(remote_file_paths) as Promise<DownloadResult>,
  upload_local_vpks: (local_file_paths: string[], better_naming = false) =>
    AppApi.UploadLocalVpks(local_file_paths, better_naming) as Promise<UploadResult>,
  open_external_url: (url: string) => AppApi.OpenExternalUrl(url),
  restart_l4d2_server: () => AppApi.RestartL4d2Server() as Promise<ApiResult & { duration_s?: number }>,
}

export type PyWebViewApi = typeof wailsAdapter

declare global {
  interface Window {
    pywebview?: {
      api: PyWebViewApi
    }
    __onExportProgress?: (payload: { stage: string; done: number; total: number }) => void
  }
}

// 注入到 window.pywebview.api 保证兼容性
if (typeof window !== 'undefined') {
  ;(window as any).pywebview = { api: wailsAdapter }
}

export function getApi(): PyWebViewApi {
  return wailsAdapter
}

export async function callApi<T>(
  operation: (api: PyWebViewApi) => Promise<T>,
  fallback: T
): Promise<T> {
  try {
    return await operation(getApi())
  } catch (error) {
    console.error('API 调用失败:', error)
    return fallback
  }
}

export function usePyWebViewApi() {
  const isApiReady = ref(true)

  const checkApiReady = () => {
    isApiReady.value = true
    return true
  }

  const waitForApiReady = async (): Promise<boolean> => {
    isApiReady.value = true
    return true
  }

  onMounted(() => {
    checkApiReady()
  })

  return { isApiReady, checkApiReady, waitForApiReady }
}
