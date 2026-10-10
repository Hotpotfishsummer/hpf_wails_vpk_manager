import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
import i18n from '@/i18n'
import { getApi, callApi } from '@/composables/useApi'
import type { RemoteFile, DetectedGameClient, LaunchSettings, VpkGroup, AiSettings, GroupPredictResult } from '@/types/api'

export type { DetectedGameClient, LaunchSettings, VpkGroup, AiSettings, GroupPredictResult }

// 组件外使用 vue-i18n 的全局实例,保证通知文案跟随当前语言
const $t = i18n.global.t

export const GROUP_COLOR_PALETTE = [
  '#0ea5e9', // 天空蓝 Sky
  '#10b981', // 翡翠绿 Emerald
  '#8b5cf6', // 梦幻紫 Purple
  '#f59e0b', // 琥珀橙 Amber
  '#f43f5e', // 玫瑰红 Rose
  '#06b6d4', // 青碧色 Cyan
  '#ec4899', // 亮粉色 Pink
  '#6366f1', // 靛青色 Indigo
  '#14b8a6', // 青雀绿 Teal
  '#eab308', // 柠黄色 Yellow
  '#d946ef', // 魅惑洋红 Fuchsia
  '#fb923c', // 珊瑚橙 Coral
]

export function getVpkKey(vpkPath: string): string {
  if (!vpkPath) return ''
  const normalized = vpkPath.replace(/\\/g, '/')
  const filename = normalized.split('/').pop() || ''
  const clean = filename.toLowerCase().replace(/\.disabled$/i, '')
  if (normalized.toLowerCase().includes('/workshop/')) {
    return `workshop:${clean}`
  }
  return `local:${clean}`
}

export function hexToRgb(hex: string): string {
  const c = (hex || '').replace('#', '')
  if (c.length === 3) {
    const r = parseInt(c[0] + c[0], 16)
    const g = parseInt(c[1] + c[1], 16)
    const b = parseInt(c[2] + c[2], 16)
    return `${r}, ${g}, ${b}`
  }
  if (c.length === 6) {
    const r = parseInt(c.slice(0, 2), 16)
    const g = parseInt(c.slice(2, 4), 16)
    const b = parseInt(c.slice(4, 6), 16)
    return `${r}, ${g}, ${b}`
  }
  return '56, 189, 248' // fallback
}

export interface MapInfo {
  name: string
  display_name?: string
}

export interface VpkFile {
  name: string
  path: string
  size: number
  modified_time: number
  is_valid: boolean
  is_disabled: boolean
  thumbnail_path?: string
  addontitle?: string
  display_name: string
  note?: string
  maps?: MapInfo[]
}

export const useVpkManagerStore = defineStore('vpkManager', () => {
  const currentDirectory = ref<string>('')
  const localVpkFiles = ref<VpkFile[]>([])
  const workshopVpkFiles = ref<VpkFile[]>([])
  const searchQuery = ref<string>('')
  const filterStatus = ref<'all' | 'enabled' | 'disabled'>('all')
  const sortKey = ref<'name' | 'size' | 'date'>('name')
  const sortAsc = ref<boolean>(true)
  const viewMode = ref<'card' | 'compact'>('card')
  const actionBarPosition = ref<'bottom' | 'top'>('bottom')

  try {
    const saved = localStorage.getItem('vpk-view-mode')
    if (saved === 'card' || saved === 'compact') {
      viewMode.value = saved
    }
  } catch {}

  try {
    const savedPos = localStorage.getItem('vpk-actionbar-pos')
    if (savedPos === 'bottom' || savedPos === 'top') {
      actionBarPosition.value = savedPos
    }
  } catch {}

  function setViewMode(mode: 'card' | 'compact') {
    viewMode.value = mode
    try {
      localStorage.setItem('vpk-view-mode', mode)
    } catch {}
  }

  function setActionBarPosition(pos: 'bottom' | 'top') {
    actionBarPosition.value = pos
    try {
      localStorage.setItem('vpk-actionbar-pos', pos)
    } catch {}
  }

  function toggleActionBarPosition() {
    setActionBarPosition(actionBarPosition.value === 'bottom' ? 'top' : 'bottom')
  }

  // VPK 自定义分组状态
  const groups = ref<VpkGroup[]>([])
  const isGroupModalOpen = ref(false)
  const groupModalMode = ref<'assign' | 'edit'>('assign')
  const editingGroupId = ref<string>('')

  // 快速映射: stable vpkKey -> VpkGroup
  const fileGroupMap = computed(() => {
    const map = new Map<string, VpkGroup>()
    for (const grp of groups.value) {
      if (grp.file_keys) {
        for (const key of grp.file_keys) {
          map.set(key, grp)
        }
      }
    }
    return map
  })

  function getFileGroup(vpkPath: string): VpkGroup | undefined {
    return fileGroupMap.value.get(getVpkKey(vpkPath))
  }

  function matchesQuery(vpk: VpkFile, query: string): boolean {
    const q = query.trim().toLowerCase()
    if (!q) return true
    const grp = getFileGroup(vpk.path)
    return (
      vpk.display_name.toLowerCase().includes(q) ||
      vpk.name.toLowerCase().includes(q) ||
      (grp ? grp.name.toLowerCase().includes(q) : false) ||
      (vpk.note ? vpk.note.toLowerCase().includes(q) : false) ||
      (vpk.maps ? vpk.maps.some((m) => m.name.toLowerCase().includes(q) || (m.display_name ? m.display_name.toLowerCase().includes(q) : false)) : false)
    )
  }

  function processFiles(files: VpkFile[]): VpkFile[] {
    let list = files.filter((vpk) => matchesQuery(vpk, searchQuery.value))
    if (filterStatus.value === 'enabled') {
      list = list.filter((vpk) => !vpk.is_disabled)
    } else if (filterStatus.value === 'disabled') {
      list = list.filter((vpk) => vpk.is_disabled)
    }

    return [...list].sort((a, b) => {
      let cmp = 0
      if (sortKey.value === 'name') {
        cmp = a.display_name.localeCompare(b.display_name, 'zh-CN')
      } else if (sortKey.value === 'size') {
        cmp = a.size - b.size
      } else if (sortKey.value === 'date') {
        cmp = a.modified_time - b.modified_time
      }
      return sortAsc.value ? cmp : -cmp
    })
  }

  const filteredLocalVpkFiles = computed(() => processFiles(localVpkFiles.value))
  const filteredWorkshopVpkFiles = computed(() => processFiles(workshopVpkFiles.value))

  const totalModStats = computed(() => {
    const all = [...localVpkFiles.value, ...workshopVpkFiles.value]
    return {
      total: all.length,
      enabled: all.filter((v) => !v.is_disabled).length,
      disabled: all.filter((v) => v.is_disabled).length,
      local: localVpkFiles.value.length,
      workshop: workshopVpkFiles.value.length,
    }
  })
  const selectedLocalVpks = reactive(new Set<string>())
  const selectedWorkshopVpks = reactive(new Set<string>())
  const isScanning = ref(false)
  const isExporting = ref(false)
  const isDeleting = ref(false)
  const isImporting = ref(false)
  const isLocalizing = ref(false)
  const isSftpActive = ref(false)
  const isRemoteScanning = ref(false)
  const isDownloading = ref(false)
  const isUploading = ref(false)
  const isRestartingServer = ref(false)
  const useBetterNaming = ref(true) // 默认开启
  const useAutoOpenExport = ref(true) // 默认开启
  const exportProgress = reactive({ stage: '', done: 0, total: 0 })

  // SFTP 设置
  const sftpSettings = reactive({
    address: '',
    username: '',
    password: '',
    addons_path: ''
  })

  const remoteVpkFiles = ref<RemoteFile[]>([])
  const selectedRemoteVpks = reactive(new Set<string>())
  
  // 游戏目录检测状态
  const detectedClients = ref<DetectedGameClient[]>([])
  const isDetecting = ref(false)
  const isDeepDetecting = ref(false)
  const showDetectModal = ref(false)
  
  // 游戏启动状态与参数配置
  const isLaunchingGame = ref(false)
  const isCurrentSteamClient = ref(false)
  const launchSettings = reactive<LaunchSettings>({
    novid: true,
    chinese_lang: true,
    direct_launch: false,
    custom_args: '',
  })

  // OpenAI 兼容 API 设置
  const aiSettings = reactive<AiSettings>({
    api_key: '',
    base_url: 'https://api.openai.com/v1',
    model: 'gpt-4o-mini',
    temperature: 0.7,
    max_tokens: 500,
  })
  const isTestingAi = ref(false)
  const isPredictingGroupName = ref(false)
  
  // 消息提示状态
  const notification = reactive({
    message: '',
    type: 'info' as 'success' | 'error' | 'info'
  })

  const confirmState = reactive({ show: false, message: '', resolve: null as null | ((v: boolean) => void) })

  function showConfirm(message: string): Promise<boolean> {
    return new Promise((resolve) => {
      confirmState.message = message
      confirmState.show = true
      confirmState.resolve = resolve
    })
  }

  function resolveConfirm(value: boolean) {
    if (confirmState.resolve) {
      confirmState.resolve(value)
      confirmState.resolve = null
    }
    confirmState.show = false
  }

  function showNotification(message: string, type: 'success' | 'error' | 'info' = 'info') {
    notification.message = message
    notification.type = type
    
    // 如果是成功或提示消息，3秒后自动清除
    if (type !== 'error') {
      setTimeout(() => {
        if (notification.message === message) {
          notification.message = ''
        }
      }, 3000)
    }
  }

  function toggleLocalVpkSelection(path: string) {
    if (selectedLocalVpks.has(path)) {
      selectedLocalVpks.delete(path)
    } else {
      selectedLocalVpks.add(path)
    }
  }

  function toggleWorkshopVpkSelection(path: string) {
    if (selectedWorkshopVpks.has(path)) {
      selectedWorkshopVpks.delete(path)
    } else {
      selectedWorkshopVpks.add(path)
    }
  }

  function selectAllLocal(paths?: string[]) {
    const list = paths ?? localVpkFiles.value.map((v) => v.path)
    list.forEach((p) => selectedLocalVpks.add(p))
  }

  function deselectAllLocal(paths?: string[]) {
    if (paths) {
      paths.forEach((p) => selectedLocalVpks.delete(p))
    } else {
      selectedLocalVpks.clear()
    }
  }

  function selectAllWorkshop(paths?: string[]) {
    const list = paths ?? workshopVpkFiles.value.map((v) => v.path)
    list.forEach((p) => selectedWorkshopVpks.add(p))
  }

  function deselectAllWorkshop(paths?: string[]) {
    if (paths) {
      paths.forEach((p) => selectedWorkshopVpks.delete(p))
    } else {
      selectedWorkshopVpks.clear()
    }
  }

  function clearSelection() {
    selectedLocalVpks.clear()
    selectedWorkshopVpks.clear()
  }

  function getSelectedVpkPaths(): string[] {
    return [...Array.from(selectedLocalVpks), ...Array.from(selectedWorkshopVpks)]
  }

  /** 加载分组配置 */
  async function loadGroups() {
    try {
      const res = await callApi(api => api.get_vpk_groups(), [] as VpkGroup[])
      if (Array.isArray(res) && res.length > 0) {
        groups.value = res
        return
      }
    } catch {}
    // 回退尝试从 localStorage 恢复
    try {
      const cached = localStorage.getItem('vpk-groups')
      if (cached) {
        groups.value = JSON.parse(cached)
      }
    } catch {}
  }

  /** 持久化保存分组配置 */
  async function saveGroups() {
    try {
      localStorage.setItem('vpk-groups', JSON.stringify(groups.value))
    } catch {}
    try {
      await callApi(api => api.save_vpk_groups(groups.value), { success: true })
    } catch (e) {
      console.warn('Failed to persist groups to backend:', e)
    }
  }

  /** 新建分组并加入指定文件 */
  function createGroup(name: string, color?: string, filePaths: string[] = []): VpkGroup {
    const trimmed = name.trim() || `${$t('group', '分组')} ${groups.value.length + 1}`
    const chosenColor = color || GROUP_COLOR_PALETTE[groups.value.length % GROUP_COLOR_PALETTE.length]
    const id = `grp_${Date.now()}_${Math.random().toString(36).slice(2, 6)}`

    const newKeys = filePaths.map(p => getVpkKey(p)).filter(Boolean)

    // 从已有分组中移出这些 key (每个文件唯一属于一个组)
    if (newKeys.length > 0) {
      const keySet = new Set(newKeys)
      for (const grp of groups.value) {
        grp.file_keys = (grp.file_keys || []).filter(k => !keySet.has(k))
      }
    }

    const newGroup: VpkGroup = {
      id,
      name: trimmed,
      color: chosenColor,
      collapsed: false,
      file_keys: Array.from(new Set(newKeys)),
    }

    groups.value.push(newGroup)
    saveGroups()
    showNotification($t('groupCreatedSuccess', { name: trimmed, count: newKeys.length }), 'success')
    return newGroup
  }

  /** 将指定文件加入已有分组 */
  function addFilesToGroup(groupId: string, filePaths: string[]) {
    const targetGroup = groups.value.find(g => g.id === groupId)
    if (!targetGroup) return

    const newKeys = filePaths.map(p => getVpkKey(p)).filter(Boolean)
    if (newKeys.length === 0) return

    const keySet = new Set(newKeys)
    // 从其他分组中移除
    for (const grp of groups.value) {
      if (grp.id !== groupId) {
        grp.file_keys = (grp.file_keys || []).filter(k => !keySet.has(k))
      }
    }

    const currentKeys = new Set(targetGroup.file_keys || [])
    for (const k of newKeys) {
      currentKeys.add(k)
    }
    targetGroup.file_keys = Array.from(currentKeys)

    saveGroups()
    showNotification($t('groupAddedSuccess', { name: targetGroup.name, count: newKeys.length }), 'success')
  }

  /** 从分组中移出指定文件 (恢复为未分组) */
  function removeFilesFromGroups(filePaths: string[]) {
    const keysToRemove = new Set(filePaths.map(p => getVpkKey(p)).filter(Boolean))
    if (keysToRemove.size === 0) return

    let removedCount = 0
    for (const grp of groups.value) {
      const before = (grp.file_keys || []).length
      grp.file_keys = (grp.file_keys || []).filter(k => !keysToRemove.has(k))
      removedCount += before - grp.file_keys.length
    }

    saveGroups()
    showNotification($t('groupRemovedSuccess', { count: removedCount }), 'success')
  }

  /** 解散分组 */
  function dissolveGroup(groupId: string) {
    const idx = groups.value.findIndex(g => g.id === groupId)
    if (idx !== -1) {
      const name = groups.value[idx].name
      groups.value.splice(idx, 1)
      saveGroups()
      showNotification($t('groupDissolvedSuccess', { name }), 'success')
    }
  }

  /** 折叠/展开分组 */
  function toggleGroupCollapse(groupId: string) {
    const grp = groups.value.find(g => g.id === groupId)
    if (grp) {
      grp.collapsed = !grp.collapsed
      saveGroups()
    }
  }

  /** 更新分组名称与颜色 */
  function updateGroup(groupId: string, data: { name?: string; color?: string }) {
    const grp = groups.value.find(g => g.id === groupId)
    if (grp) {
      if (data.name !== undefined && data.name.trim()) {
        grp.name = data.name.trim()
      }
      if (data.color !== undefined) {
        grp.color = data.color
      }
      saveGroups()
      showNotification($t('groupUpdatedSuccess'), 'success')
    }
  }

  /** 组内全选或取消全选 */
  function selectGroupFiles(groupId: string, selectAll: boolean) {
    const grp = groups.value.find(g => g.id === groupId)
    if (!grp || !grp.file_keys) return

    const keySet = new Set(grp.file_keys)
    const updateSelection = (list: VpkFile[], set: Set<string>) => {
      for (const v of list) {
        if (keySet.has(getVpkKey(v.path))) {
          if (selectAll) {
            set.add(v.path)
          } else {
            set.delete(v.path)
          }
        }
      }
    }
    updateSelection(localVpkFiles.value, selectedLocalVpks)
    updateSelection(workshopVpkFiles.value, selectedWorkshopVpks)
  }

  /** 打开分组操作弹窗 */
  function openGroupModal(mode: 'assign' | 'edit' = 'assign', groupId?: string) {
    groupModalMode.value = mode
    editingGroupId.value = groupId || ''
    isGroupModalOpen.value = true
  }

  function closeGroupModal() {
    isGroupModalOpen.value = false
  }

  /** 扫描本地与创意工坊 VPK 列表(经后端 scan_vpk_files) */
  async function scanVpks(directory: string) {
    if (!directory) return

    isScanning.value = true
    showNotification('', 'info') // 清除之前的消息

    try {
      const result = await getApi().scan_vpk_files(directory)

      if (result.success) {
        currentDirectory.value = directory
        localVpkFiles.value = result.local_vpks || []
        workshopVpkFiles.value = result.workshop_vpks || []
        selectedLocalVpks.clear()
        selectedWorkshopVpks.clear()
        if (typeof window !== 'undefined' && new URLSearchParams(window.location.search).get('selected') === '1' && localVpkFiles.value.length > 0) {
          selectedLocalVpks.add(localVpkFiles.value[0].path)
        }
        console.log(`Successfully scanned ${localVpkFiles.value.length + workshopVpkFiles.value.length} files`)
        // 异步加载分组
        loadGroups().catch(() => {})
        // 异步检查是否为 Steam 客户端来源
        callApi(api => api.is_steam_client(directory), false).then(isSteam => {
          isCurrentSteamClient.value = isSteam
        })
      } else {
        showNotification(result.error || $t('scanFailed'), 'error')
      }
    } catch (error) {
      showNotification($t('scanError', { error }), 'error')
    } finally {
      isScanning.value = false
    }
  }

    /** 弹窗选择压缩包并导入到本地 addons(经后端 import_compressed_vpks) */
async function importVpks() {
    if (!currentDirectory.value) {
      showNotification($t('needDirectoryFirst'), 'error')
      return
    }

    isImporting.value = true
    showNotification('', 'info')

    try {
      const result = await getApi().import_compressed_vpks()

      if (result.success) {
        // 导入成功，刷新列表
        await scanVpks(currentDirectory.value)
        showNotification($t('importSuccessPacks', { count: result.imported_count }), 'success')
        return result
      } else {
        showNotification(result.error || $t('importFailed'), 'error')
        return result
      }
    } catch (error) {
      showNotification($t('importError', { error }), 'error')
      throw error
    } finally {
      isImporting.value = false
    }
  }

  /** 直接导入或拖入压缩包/VPK到本地 addons(经后端 import_files) */
  async function importFiles(filePaths: string[]) {
    if (!currentDirectory.value) {
      showNotification($t('needDirectoryFirst'), 'error')
      return null
    }
    if (!filePaths || filePaths.length === 0) {
      return null
    }

    isImporting.value = true
    showNotification($t('importing'), 'info')

    try {
      const result = await getApi().import_files(filePaths)
      if (result.success) {
        await scanVpks(currentDirectory.value)
        showNotification($t('importSuccessPacks', { count: result.imported_count }), 'success')
        return result
      } else if ((result.imported_count ?? 0) > 0) {
        await scanVpks(currentDirectory.value)
        showNotification($t('importDonePartial', {
          success: result.imported_count,
          failed: result.errors ? result.errors.length : 0,
        }), 'info')
        return result
      } else {
        showNotification(result.error || (result.errors && result.errors.join('; ')) || $t('importFailed'), 'error')
        return result
      }
    } catch (error) {
      showNotification($t('importError', { error }), 'error')
      throw error
    } finally {
      isImporting.value = false
    }
  }

    /** 将选中的创意工坊 VPK 本地化到 addons(可选 better_naming 重命名) */
async function localizeVpks() {
    const paths = Array.from(selectedWorkshopVpks)
    if (paths.length === 0) return

    isLocalizing.value = true
    showNotification('', 'info')

    try {
      const result = await getApi().localize_workshop_vpks(
        paths, 
        useBetterNaming.value
      )

      if ((result.success_count ?? 0) > 0) {
        // 本地化成功（至少部分），刷新列表
        selectedWorkshopVpks.clear()
        await scanVpks(currentDirectory.value)
        const type = result.success ? 'success' : 'info'
        let msg = result.success 
          ? $t('localizeSuccess', { count: result.success_count })
          : $t('localizePartial', { count: result.success_count })
        
        if (result.errors && result.errors.length > 0) {
          msg += $t('localizePartialFailed', { count: result.errors.length })
        }
        showNotification(msg, type)
        return result
      } else {
        const errorMsg = result.error || (result.errors && result.errors.length > 0 ? result.errors[0] : $t('localizeFailed'))
        showNotification(errorMsg, 'error')
        return result
      }
    } catch (error) {
      showNotification($t('localizeError', { error }), 'error')
      throw error
    } finally {
      isLocalizing.value = false
    }
  }

  /** 自动搜索 Left 4 Dead 2 客户端目录并弹窗展示 */
  async function detectGameDirectories(deepScan = false) {
    if (deepScan) {
      isDeepDetecting.value = true
    } else {
      isDetecting.value = true
    }
    showDetectModal.value = true

    try {
      const res = await getApi().detect_l4d2_directories(deepScan)
      if (res && res.clients) {
        detectedClients.value = res.clients
      } else {
        detectedClients.value = []
      }
    } catch (e) {
      console.error('Failed to detect L4D2 directories:', e)
      showNotification($t('detectNoClients'), 'error')
    } finally {
      isDetecting.value = false
      isDeepDetecting.value = false
    }
  }

  /** 选择并应用检测到的客户端目录 */
  async function selectDetectedClient(path: string) {
    if (!path) return
    showDetectModal.value = false
    await scanVpks(path)
  }

  function closeDetectModal() {
    showDetectModal.value = false
  }

  /** 启动游戏：Steam 客户端唤起 Steam，非 Steam 直接运行主程序 */
  async function launchGame() {
    if (!currentDirectory.value) {
      showNotification($t('needSelectL4d2', '请先选择 L4D2 安装目录'), 'error')
      return
    }

    isLaunchingGame.value = true
    try {
      const res = await getApi().launch_game()
      if (res.success) {
        showNotification(res.message || $t('gameLaunchedSuccess'), 'success')
      } else {
        showNotification(res.error || $t('gameLaunchFailed'), 'error')
      }
    } catch (e) {
      showNotification($t('gameLaunchError', { error: String(e) }), 'error')
    } finally {
      setTimeout(() => {
        isLaunchingGame.value = false
      }, 1000)
    }
  }

  /** 应用启动初始化:恢复上次有效目录或自动探测 */
  async function initialize() {
    console.log('Initializing VpkManager store...')
    try {
      // 加载上次目录
      const lastDir = await getApi().get_last_directory()
      let hasValidLastDir = false
      if (lastDir) {
        const isValid = await callApi(api => api.validate_l4d2_directory(lastDir), false)
        if (isValid) {
          hasValidLastDir = true
          console.log('Last directory found and valid:', lastDir)
          await scanVpks(lastDir)
        }
      }

      if (!hasValidLastDir) {
        console.log('No valid last directory stored, auto detecting...')
        await detectGameDirectories(false)
      }

      // 加载 SFTP 设置
      const settings = await getApi().get_sftp_settings()
      if (settings) {
        Object.assign(sftpSettings, settings)
      }

      // 加载启动参数设置
      const launch = await getApi().get_launch_settings()
      if (launch) {
        Object.assign(launchSettings, launch)
      }

      // 加载自定义分组配置
      await loadGroups()

      // 加载 AI 大模型配置
      await loadAiSettings()
    } catch (error) {
      console.warn('Failed to initialize settings:', error)
    }
  }

  /** 保存启动参数设置到偏好文件 */
  async function saveLaunchSettings() {
    try {
      const result = await getApi().save_launch_settings({ ...launchSettings })
      if (!result.success) {
        showNotification(result.error || '保存启动设置失败', 'error')
      }
      return result
    } catch (error) {
      showNotification(`保存启动设置错误: ${error}`, 'error')
      throw error
    }
  }

  /** 保存 SFTP 设置到偏好文件 */
  async function saveSftpSettings() {
    try {
      const result = await getApi().save_sftp_settings({ ...sftpSettings })
      if (!result.success) {
        showNotification(result.error || '保存 SFTP 设置失败', 'error')
      }
      return result
    } catch (error) {
      showNotification(`保存 SFTP 设置错误: ${error}`, 'error')
      throw error
    }
  }

  /** 加载 OpenAI 兼容大模型 API 设置 */
  async function loadAiSettings() {
    try {
      const settings = await callApi(api => api.get_ai_settings(), null)
      if (settings) {
        Object.assign(aiSettings, settings)
      }
    } catch (e) {
      console.warn('Failed to load AI settings:', e)
    }
  }

  /** 保存 OpenAI 兼容大模型 API 设置 */
  async function saveAiSettings() {
    try {
      const result = await getApi().save_ai_settings({ ...aiSettings })
      if (!result.success) {
        showNotification(result.error || $t('saveAiSettingsFailed', '保存 AI 设置失败'), 'error')
      } else {
        showNotification(result.message || $t('aiSettingsSaved', 'AI 设置已保存'), 'success')
      }
      return result
    } catch (error) {
      showNotification(`保存 AI 设置错误: ${error}`, 'error')
      throw error
    }
  }

  /** 测试 OpenAI 兼容 API 连接 */
  async function testAiConnection() {
    isTestingAi.value = true
    try {
      const result = await getApi().test_ai_connection({ ...aiSettings })
      if (result.success) {
        showNotification(result.message || $t('aiConnectSuccess', 'AI 服务连接成功！'), 'success')
      } else {
        showNotification(result.error || $t('aiConnectFailed', 'AI 服务连接失败'), 'error')
      }
      return result
    } catch (error) {
      showNotification(`测试连接错误: ${error}`, 'error')
      return { success: false, error: String(error) }
    } finally {
      isTestingAi.value = false
    }
  }

  /** 智能预测/推断选中文件的分组名称 */
  async function predictGroupName(files?: VpkFile[]): Promise<{ name: string; source: 'ai' | 'heuristic'; reason?: string }> {
    isPredictingGroupName.value = true
    try {
      let targetFiles = files
      if (!targetFiles || targetFiles.length === 0) {
        const selectedPaths = new Set(getSelectedVpkPaths())
        targetFiles = [...localVpkFiles.value, ...workshopVpkFiles.value].filter(f => selectedPaths.has(f.path))
      }
      const result = await callApi(api => api.predict_group_name(targetFiles as any), {
        success: true,
        predicted_name: '新模组分组',
        source: 'heuristic' as const,
      })
      return {
        name: result.predicted_name || '新模组分组',
        source: result.source || 'heuristic',
        reason: result.reason,
      }
    } catch (error) {
      console.warn('Failed to predict group name:', error)
      return { name: '新模组分组', source: 'heuristic' }
    } finally {
      isPredictingGroupName.value = false
    }
  }

  function toggleRemoteVpkSelection(path: string) {
    if (selectedRemoteVpks.has(path)) {
      selectedRemoteVpks.delete(path)
    } else {
      selectedRemoteVpks.add(path)
    }
  }

  function selectAllRemote(paths?: string[]) {
    const list = paths ?? remoteVpkFiles.value.map(v => v.path)
    list.forEach(p => {
      selectedRemoteVpks.add(p)
    })
  }

  function deselectAllRemote(paths?: string[]) {
    if (paths) {
      paths.forEach(p => selectedRemoteVpks.delete(p))
    } else {
      selectedRemoteVpks.clear()
    }
  }

    /** 列出远端 addons 目录的 VPK 文件 */
async function scanRemoteVpks() {
    isRemoteScanning.value = true
    try {
      const result = await getApi().list_remote_vpks()
      if (result.success) {
        remoteVpkFiles.value = result.files || []
        selectedRemoteVpks.clear() // 扫描后清除之前的选中
      } else {
        showNotification(result.error || $t('remoteScanFailed'), 'error')
      }
    } catch (error) {
      showNotification($t('remoteScanError', { error }), 'error')
    } finally {
      isRemoteScanning.value = false
    }
  }

    /** 删除单个远端 VPK(含关联图片) */
async function deleteRemoteVpk(path: string) {
    try {
      const result = await getApi().delete_remote_vpk(path)
      if (result.success) {
        showNotification($t('remoteDeleted'), 'success')
        await scanRemoteVpks()
      } else {
        showNotification(result.error || $t('remoteDeleteFailed'), 'error')
      }
    } catch (error) {
      showNotification($t('remoteDeleteError', { error }), 'error')
    }
  }

    /** 批量删除选中的远端 VPK(含关联图片) */
async function deleteMultipleRemoteVpks() {
    const paths = Array.from(selectedRemoteVpks)
    if (paths.length === 0) return
    
    isDeleting.value = true
    try {
      const result = await getApi().delete_multiple_remote_vpks(paths)
      if (result.success) {
        showNotification($t('remoteDeletedSelected'), 'success')
      } else if ((result.success_count ?? 0) > 0) {
        showNotification($t('remoteDeletePartial', { success: result.success_count ?? 0, failed: result.error_count ?? 0 }), 'info')
      } else {
        showNotification($t('deleteActionFailed'), 'error')
      }
      await scanRemoteVpks()
    } catch (error) {
      showNotification($t('deleteActionError', { error }), 'error')
    } finally {
      isDeleting.value = false
    }
  }

    /** 远端禁用:重命名为 .disabled(含关联图片) */
async function disableRemoteVpks() {
    const paths = Array.from(selectedRemoteVpks).filter(p => !p.endsWith('.disabled'))
    if (paths.length === 0) return
    
    isDeleting.value = true
    try {
      let successCount = 0
      for (const path of paths) {
        const result = await getApi().rename_remote_vpk(path, `${path}.disabled`)
        if (result.success) successCount++
      }
      showNotification($t('remoteDisabled', { count: successCount }), 'success')
      await scanRemoteVpks()
    } catch (error) {
      showNotification($t('remoteDisableError', { error }), 'error')
    } finally {
      isDeleting.value = false
    }
  }

    /** 远端启用:移除 .disabled 后缀(含关联图片) */
async function enableRemoteVpks() {
    const paths = Array.from(selectedRemoteVpks).filter(p => p.endsWith('.disabled'))
    if (paths.length === 0) return
    
    isDeleting.value = true
    try {
      let successCount = 0
      for (const path of paths) {
        const result = await getApi().rename_remote_vpk(path, path.substring(0, path.length - 9))
        if (result.success) successCount++
      }
      showNotification($t('remoteEnabled', { count: successCount }), 'success')
      await scanRemoteVpks()
    } catch (error) {
      showNotification($t('remoteEnableError', { error }), 'error')
    } finally {
      isDeleting.value = false
    }
  }

    /** 下载选中的远端 VPK 及关联图片到本地 addons */
async function downloadRemoteVpks() {
    const paths = Array.from(selectedRemoteVpks)
    if (paths.length === 0) return
    
    isDownloading.value = true
    showNotification($t('downloadingRemote'), 'info')
    try {
      const result = await getApi().download_remote_vpks(paths)
      if (result.success) {
        showNotification($t('downloadSuccess', { count: result.downloaded_count }), 'success')
        // 刷新本地列表
        if (currentDirectory.value) {
          await scanVpks(currentDirectory.value)
        }
      } else {
        showNotification(result.error || $t('downloadFailed'), 'error')
      }
    } catch (error) {
      showNotification($t('downloadError', { error }), 'error')
    } finally {
      isDownloading.value = false
    }
  }

    /** 上传选中的本地 VPK 及关联图片到远端 addons */
async function uploadLocalVpks() {
    const paths = getSelectedVpkPaths()
    if (paths.length === 0) return
    
    isUploading.value = true
    showNotification($t('uploadingLocal'), 'info')
    try {
      const result = await getApi().upload_local_vpks(
        paths, 
        useBetterNaming.value
      )
      if (result.success) {
        showNotification($t('uploadSuccess', { count: result.uploaded_count }), 'success')
        // 刷新远程列表
        if (isSftpActive.value) {
          await scanRemoteVpks()
        }
      } else {
        showNotification(result.error || $t('uploadFailed'), 'error')
      }
    } catch (error) {
      showNotification($t('uploadError', { error }), 'error')
    } finally {
      isUploading.value = false
    }
  }

  /** 重启远端 L4D2 服务器容器 */
  async function restartL4d2Server() {
    isRestartingServer.value = true
    showNotification($t('restartingServer'), 'info')
    try {
      const result = await getApi().restart_l4d2_server()
      if (result.success) {
        const duration = result.duration_s ? ` (${result.duration_s}s)` : ''
        showNotification($t('restartServerSuccess') + duration, 'success')
      } else {
        showNotification(result.error || $t('restartServerFailed'), 'error')
      }
    } catch (error) {
      showNotification($t('restartServerError', { error }), 'error')
    } finally {
      isRestartingServer.value = false
    }
  }

  async function openAddonsFolder() {
    if (!currentDirectory.value) return
    const path = `${currentDirectory.value}/left4dead2/addons`
    await callApi(api => api.open_directory(path), false)
  }

  async function openWorkshopFolder() {
    if (!currentDirectory.value) return
    const path = `${currentDirectory.value}/left4dead2/addons/workshop`
    await callApi(api => api.open_directory(path), false)
  }

  async function exportVpks(outputDir: string, overwrite = false) {
    const paths = getSelectedVpkPaths()
    if (paths.length === 0) {
      showNotification($t('selectFiles'), 'error')
      return
    }

    isExporting.value = true
    showNotification('', 'info')

    try {
      const result = await getApi().export_vpks(
        paths,
        outputDir,
        useBetterNaming.value,
        overwrite
      )

      if (result.success) {
        // Clear selection after export
        selectedLocalVpks.clear()
        selectedWorkshopVpks.clear()
        showNotification($t('exportDone'), 'success')
        if (useAutoOpenExport.value && outputDir) {
          await callApi(api => api.open_directory(outputDir), false)
        }
      } else if (result.exists) {
        // 文件已存在，由调用方决定如何处理（前端不自动弹出通知）
      } else {
        showNotification(result.error || $t('exportError'), 'error')
      }

      return result
    } catch (error) {
      showNotification($t('exportErrorMsg', { error }), 'error')
    } finally {
      isExporting.value = false
    }
  }

  async function deleteVpks(disable: boolean = true, targetPaths?: string[]) {
    const paths = targetPaths || getSelectedVpkPaths()
    if (paths.length === 0) {
      showNotification($t('selectFiles'), 'error')
      return
    }

    isDeleting.value = true
    showNotification('', 'info')

    try {
      // Use batch delete method for better performance
      const result = await getApi().delete_multiple_vpk_files(paths, disable)
      
      if (!result.success && (result.error_count ?? 0) > 0) {
        showNotification($t('deletePartialFailed', { count: result.error_count ?? 0 }), 'error')
      } else if (result.success) {
        // Deleted successfully, clear selections that were processed
        if (targetPaths) {
          targetPaths.forEach(p => {
            selectedLocalVpks.delete(p)
            selectedWorkshopVpks.delete(p)
          })
        } else {
          selectedLocalVpks.clear()
          selectedWorkshopVpks.clear()
        }
        showNotification(disable ? $t('statusUpdated') : $t('deleteSuccess'), 'success')
      }

      // Always rescan after deletion attempt to update UI
      await scanVpks(currentDirectory.value)
    } catch (error) {
      showNotification($t('deleteActionError', { error }), 'error')
      // Still try to rescan even if there was an error
      try {
        await scanVpks(currentDirectory.value)
      } catch (rescanError) {
        console.error('Error rescanning after deletion:', rescanError)
      }
    } finally {
      isDeleting.value = false
    }
  }

  async function getLastDirectory() {
    try {
      const lastDir = await getApi().get_last_directory()
      if (lastDir) {
        return lastDir
      }
    } catch (error) {
      console.error('Error getting last directory:', error)
    }
    return ''
  }

  async function updateVpkNote(vpkPath: string, note: string): Promise<boolean> {
    const res = await callApi(api => api.set_vpk_note(vpkPath, note), { success: false, error: '调用失败' })
    if (res.success) {
      const cleanNote = note.trim()
      const updateList = (list: VpkFile[]) => {
        const item = list.find(v => v.path === vpkPath)
        if (item) {
          item.note = cleanNote
          const rawTitle = item.addontitle && item.addontitle.trim() ? item.addontitle.trim() : item.name
          item.display_name = cleanNote ? `${rawTitle} [${cleanNote}]` : rawTitle
        }
      }
      updateList(localVpkFiles.value)
      updateList(workshopVpkFiles.value)
      showNotification($t('noteSaved', '备注已保存'), 'success')
      return true
    } else {
      showNotification(res.error || $t('saveNoteFailed', '保存备注失败'), 'error')
      return false
    }
  }

  return {
    currentDirectory,
    localVpkFiles,
    workshopVpkFiles,
    searchQuery,
    filterStatus,
    sortKey,
    sortAsc,
    viewMode,
    setViewMode,
    actionBarPosition,
    setActionBarPosition,
    toggleActionBarPosition,
    totalModStats,
    filteredLocalVpkFiles,
    filteredWorkshopVpkFiles,
    selectedLocalVpks,
    selectedWorkshopVpks,
    selectedRemoteVpks,
    isScanning,
    isExporting,
    isDeleting,
    isImporting,
    isLocalizing,
    isSftpActive,
    isRemoteScanning,
    isDownloading,
    isUploading,
    isRestartingServer,
    useBetterNaming,
    useAutoOpenExport,
    exportProgress,
    sftpSettings,
    remoteVpkFiles,
    notification,
    showNotification,
    showConfirm,
    resolveConfirm,
    confirmState,
    toggleLocalVpkSelection,
    toggleWorkshopVpkSelection,
    toggleRemoteVpkSelection,
    selectAllLocal,
    deselectAllLocal,
    selectAllWorkshop,
    deselectAllWorkshop,
    selectAllRemote,
    deselectAllRemote,
    clearSelection,
    getSelectedVpkPaths,
    scanVpks,
    importVpks,
    importFiles,
    localizeVpks,
    saveSftpSettings,
    scanRemoteVpks,
    deleteRemoteVpk,
    deleteMultipleRemoteVpks,
    disableRemoteVpks,
    enableRemoteVpks,
    downloadRemoteVpks,
    uploadLocalVpks,
    restartL4d2Server,
    initialize,
    openAddonsFolder,
    openWorkshopFolder,
    exportVpks,
    deleteVpks,
    getLastDirectory,
    updateVpkNote,
    detectedClients,
    isDetecting,
    isDeepDetecting,
    showDetectModal,
    detectGameDirectories,
    selectDetectedClient,
    closeDetectModal,
    isLaunchingGame,
    isCurrentSteamClient,
    launchGame,
    launchSettings,
    saveLaunchSettings,
    groups,
    isGroupModalOpen,
    groupModalMode,
    editingGroupId,
    fileGroupMap,
    getFileGroup,
    loadGroups,
    saveGroups,
    createGroup,
    addFilesToGroup,
    removeFilesFromGroups,
    dissolveGroup,
    toggleGroupCollapse,
    updateGroup,
    selectGroupFiles,
    openGroupModal,
    closeGroupModal,
    aiSettings,
    isTestingAi,
    isPredictingGroupName,
    loadAiSettings,
    saveAiSettings,
    testAiConnection,
    predictGroupName,
  }
})
