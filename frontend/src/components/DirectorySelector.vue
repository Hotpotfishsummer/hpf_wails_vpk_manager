<template>
  <div class="directory-selector">
    <div class="selector-card">
      <div class="path-display-box">
        <div class="folder-badge">
          <SvgIcon name="folder" :size="16" />
        </div>
        <input
          id="directory-path"
          v-model="store.currentDirectory"
          type="text"
          :placeholder="$t('enterFolderPath')"
          readonly
          class="path-input"
        />
        <div class="path-actions" v-if="store.currentDirectory">
          <button
            class="path-mini-btn"
            @click="store.refreshVpks()"
            :disabled="store.isScanning"
            :title="$t('refreshTooltip')"
          >
            <SvgIcon name="refresh" :size="14" :spin="store.isScanning" />
          </button>
          <button
            class="path-mini-btn"
            @click="handleCopyPath"
            :title="$t('copyPath', '复制路径')"
          >
            <SvgIcon :name="copied ? 'check' : 'copy'" :size="14" />
          </button>
          <button
            class="path-mini-btn"
            @click="store.openAddonsFolder"
            :title="$t('openInFileManager')"
          >
            <SvgIcon name="external-link" :size="14" />
          </button>
        </div>
      </div>

      <div class="action-buttons-group">
        <!-- 导入模组按钮 -->
        <button
          @click="handleImport"
          :disabled="!store.currentDirectory || store.isImporting"
          class="btn-import-box"
          :class="{ 'btn-loading': store.isImporting }"
          :title="!store.currentDirectory ? $t('needSelectL4d2') : $t('importFilesTooltip')"
        >
          <SvgIcon v-if="store.isImporting" name="loader" :size="15" :spin="true" />
          <SvgIcon v-else name="import" :size="15" />
          <span>{{ store.isImporting ? $t('importing') : $t('importFiles') }}</span>
        </button>

        <!-- 刷新模组列表按钮 -->
        <button
          @click="store.refreshVpks()"
          :disabled="!store.currentDirectory || store.isScanning"
          class="btn-refresh-box"
          :class="{ 'btn-loading': store.isScanning }"
          :title="!store.currentDirectory ? $t('needSelectL4d2') : $t('refreshTooltip')"
        >
          <SvgIcon name="refresh" :size="14" :spin="store.isScanning" />
          <span>{{ store.isScanning ? $t('refreshing') : $t('refresh') }}</span>
        </button>

        <!-- 自动搜索客户端按钮 -->
        <button
          @click="handleAutoDetect"
          :disabled="store.isDetecting || store.isDeepDetecting"
          class="btn-detect-box"
          :class="{ 'btn-loading': store.isDetecting || store.isDeepDetecting }"
          :title="$t('autoDetectTooltip')"
        >
          <SvgIcon v-if="store.isDetecting || store.isDeepDetecting" name="loader" :size="15" :spin="true" />
          <SvgIcon v-else name="compass" :size="15" />
          <span>{{ store.isDetecting || store.isDeepDetecting ? $t('detecting') : $t('autoDetect') }}</span>
        </button>

        <!-- 浏览并选择安装目录按钮 -->
        <button
          @click="handleBrowse"
          :disabled="isLoading"
          class="btn-browse-box"
          :class="{ 'btn-loading': isLoading }"
          :title="isLoading ? $t('browsing') : $t('browseFolder')"
        >
          <SvgIcon v-if="isLoading" name="loader" :size="15" :spin="true" />
          <SvgIcon v-else name="folder-open" :size="15" />
          <span>{{ isLoading ? $t('browsing') : $t('browseFolder') }}</span>
        </button>

        <!-- 复制服务器连接指令按钮 -->
        <button
          @click="handleCopyServer"
          class="btn-copy-server-box"
          :title="$t('copyServerTooltip', { cmd: store.launchSettings.server_connect_cmd || 'connect 192.168.1.200' })"
        >
          <SvgIcon :name="serverCopied ? 'check' : 'server'" :size="14" />
          <span>{{ serverCopied ? $t('copied') : $t('copyServer') }}</span>
        </button>

        <!-- 修改网络参数按钮 (复制 100 Tick 网络参数) -->
        <button
          @click="handleCopyNetworkRates"
          class="btn-copy-network-box"
          :title="$t('copyNetworkTooltip', { cmd: formatConsoleCommands(store.launchSettings.network_rates_cmd || DEFAULT_NETWORK_RATES) })"
        >
          <SvgIcon :name="networkCopied ? 'check' : 'sliders'" :size="14" />
          <span>{{ networkCopied ? $t('copied') : $t('networkParams') }}</span>
        </button>

        <!-- 启动游戏按钮 -->
        <button
          @click="store.launchGame"
          :disabled="!store.currentDirectory || store.isLaunchingGame"
          class="btn-launch-box"
          :class="{ 'btn-loading': store.isLaunchingGame }"
          :title="launchTooltip"
        >
          <SvgIcon v-if="store.isLaunchingGame" name="loader" :size="14" :spin="true" />
          <SvgIcon v-else name="play" :size="13" />
          <span>{{ store.isLaunchingGame ? $t('launchingGame') : $t('launchGame') }}</span>
          <span v-if="store.currentDirectory && isSteamLaunchMode" class="steam-tag">{{ $t('launchTagSteam') }}</span>
          <span v-else-if="store.currentDirectory" class="direct-tag">{{ $t('launchTagDirect') }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useVpkManagerStore, DEFAULT_NETWORK_RATES, formatConsoleCommands } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import { ref, computed } from 'vue'
import { getApi } from '@/composables/useApi'
import SvgIcon from '@/components/common/SvgIcon.vue'

const { t: $t } = useI18n()
const store = useVpkManagerStore()
const isLoading = ref(false)
const copied = ref(false)
const serverCopied = ref(false)
const networkCopied = ref(false)
let networkCopyTimer: number | null = null

const isSteamLaunchMode = computed(() => {
  return store.isCurrentSteamClient && !store.launchSettings.direct_launch
})

const launchTooltip = computed(() => {
  if (!store.currentDirectory) return $t('needSelectL4d2')
  return isSteamLaunchMode.value ? $t('launchGameTooltipSteam') : $t('launchGameTooltipDirect')
})

async function handleCopyPath() {
  if (!store.currentDirectory) return
  try {
    await navigator.clipboard.writeText(store.currentDirectory)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (e) {
    console.error('Failed to copy path', e)
  }
}

async function handleCopyServer() {
  const success = await store.copyServerConnectCmd()
  if (success) {
    serverCopied.value = true
    setTimeout(() => {
      serverCopied.value = false
    }, 2000)
  }
}

async function handleCopyNetworkRates() {
  const success = await store.copyNetworkRatesCmd()
  if (success) {
    networkCopied.value = true
    if (networkCopyTimer) clearTimeout(networkCopyTimer)
    networkCopyTimer = window.setTimeout(() => {
      networkCopied.value = false
    }, 2000)
  }
}

async function handleImport() {
  if (store.isImporting || !store.currentDirectory) return

  try {
    const result = await store.importVpks()
    if (result) {
      if ((result.imported_count ?? 0) > 0) {
        const type = result.success ? 'success' : 'info'
        const msg = result.success
          ? $t('importSuccess', { count: result.imported_count })
          : $t('importDonePartial', {
              success: result.imported_count,
              failed: result.errors ? result.errors.length : 0,
            })
        store.showNotification(msg, type)
      } else if (result.error) {
        store.showNotification($t('importFailedWithError', { error: result.error }), 'error')
      }
    }
  } catch (error) {
    console.error('导入失败:', error)
  }
}

async function handleBrowse() {
  if (isLoading.value) return

  try {
    isLoading.value = true
    const selectedPath = await getApi().browse_directory()
    if (selectedPath) {
      await store.scanVpks(selectedPath)
    }
  } catch (error) {
    store.showNotification(
      $t('browseError', { error: error instanceof Error ? error.message : String(error) }),
      'error'
    )
    console.error('选择文件夹失败:', error)
  } finally {
    isLoading.value = false
  }
}

function handleAutoDetect() {
  store.detectGameDirectories(false)
}
</script>

<style scoped>
.directory-selector {
  width: 100%;
}

.selector-card {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  padding: 10px 14px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: var(--shadow-sm);
  transition: border-color var(--transition-fast);
}

.selector-card:hover {
  border-color: var(--border-medium);
}

.path-display-box {
  flex: 1;
  display: flex;
  align-items: center;
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 4px 10px;
  gap: 8px;
  min-width: 0;
}

.folder-badge {
  color: var(--accent-primary);
  display: flex;
  align-items: center;
}

.path-input {
  flex: 1;
  background: transparent !important;
  border: none !important;
  padding: 4px 0 !important;
  font-size: 13px;
  color: var(--text-main);
  box-shadow: none !important;
  min-width: 0;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
  cursor: default;
}

.path-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.path-mini-btn {
  width: 26px;
  height: 26px;
  border-radius: var(--radius-sm);
  background: transparent !important;
  border: none !important;
  color: var(--text-dim);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.path-mini-btn:hover {
  background: rgba(255, 255, 255, 0.1) !important;
  color: #fff;
}

.action-buttons-group {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
  flex-wrap: wrap;
}

/* 导入按钮 */
.btn-import-box {
  height: 36px;
  padding: 0 14px;
  border-radius: var(--radius-md);
  background: rgba(16, 185, 129, 0.12);
  color: var(--accent-success);
  border: 1px solid rgba(16, 185, 129, 0.3);
  font-weight: 600;
  font-size: 13px;
}

.btn-import-box:hover:not(:disabled) {
  background: rgba(16, 185, 129, 0.22);
  border-color: rgba(16, 185, 129, 0.5);
  box-shadow: 0 0 12px var(--accent-success-glow);
}

/* 刷新列表按钮 */
.btn-refresh-box {
  height: 36px;
  padding: 0 14px;
  border-radius: var(--radius-md);
  background: rgba(6, 182, 212, 0.12);
  color: #22d3ee;
  border: 1px solid rgba(6, 182, 212, 0.35);
  font-weight: 600;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.btn-refresh-box:hover:not(:disabled) {
  background: rgba(6, 182, 212, 0.22);
  border-color: rgba(6, 182, 212, 0.6);
  box-shadow: 0 0 12px rgba(6, 182, 212, 0.3);
  transform: translateY(-1px);
}

.btn-refresh-box:active:not(:disabled) {
  transform: translateY(0);
}

.btn-refresh-box:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 自动探测按钮 */
.btn-detect-box {
  height: 36px;
  padding: 0 14px;
  border-radius: var(--radius-md);
  background: rgba(168, 85, 247, 0.12);
  color: #c084fc;
  border: 1px solid rgba(168, 85, 247, 0.35);
  font-weight: 600;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-detect-box:hover:not(:disabled) {
  background: rgba(168, 85, 247, 0.22);
  border-color: rgba(168, 85, 247, 0.6);
  box-shadow: 0 0 12px rgba(168, 85, 247, 0.3);
}

/* 浏览文件夹主按钮 */
.btn-browse-box {
  height: 36px;
  padding: 0 16px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, #0284c7 0%, #0369a1 100%);
  color: #fff;
  border: 1px solid rgba(56, 189, 248, 0.4);
  font-weight: 600;
  font-size: 13px;
  box-shadow: 0 2px 8px rgba(2, 132, 199, 0.3);
}

.btn-browse-box:hover:not(:disabled) {
  background: linear-gradient(135deg, #0369a1 0%, #075985 100%);
  box-shadow: 0 4px 14px rgba(2, 132, 199, 0.5);
}

/* 复制服务器地址按钮 */
.btn-copy-server-box {
  height: 36px;
  padding: 0 14px;
  border-radius: var(--radius-md);
  background: rgba(56, 189, 248, 0.12);
  color: var(--accent-primary);
  border: 1px solid rgba(56, 189, 248, 0.35);
  font-weight: 600;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.btn-copy-server-box:hover {
  background: rgba(56, 189, 248, 0.22);
  border-color: rgba(56, 189, 248, 0.6);
  box-shadow: 0 0 12px var(--accent-primary-glow);
  transform: translateY(-1px);
}

.btn-copy-server-box:active {
  transform: translateY(0);
}

/* 修改网络参数按钮 */
.btn-copy-network-box {
  height: 36px;
  padding: 0 14px;
  border-radius: var(--radius-md);
  background: rgba(245, 158, 11, 0.12);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.35);
  font-weight: 600;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.btn-copy-network-box:hover {
  background: rgba(245, 158, 11, 0.22);
  border-color: rgba(245, 158, 11, 0.6);
  box-shadow: 0 0 12px rgba(245, 158, 11, 0.35);
  transform: translateY(-1px);
}

.btn-copy-network-box:active {
  transform: translateY(0);
}

/* 启动游戏主按钮 */
.btn-launch-box {
  height: 36px;
  padding: 0 16px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: #fff;
  border: 1px solid rgba(52, 211, 153, 0.4);
  font-weight: 700;
  font-size: 13px;
  box-shadow: 0 2px 10px rgba(16, 185, 129, 0.35);
  display: flex;
  align-items: center;
  gap: 7px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-launch-box:hover:not(:disabled) {
  filter: brightness(1.15);
  box-shadow: 0 0 16px rgba(16, 185, 129, 0.6);
  transform: translateY(-1px);
}

.btn-launch-box:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  filter: grayscale(0.5);
}

.steam-tag {
  font-size: 10px;
  padding: 1px 5px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 3px;
  color: #a7f3d0;
  font-weight: 800;
  letter-spacing: 0.5px;
}

.direct-tag {
  font-size: 10px;
  padding: 1px 5px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 3px;
  color: #fed7aa;
  font-weight: 800;
  letter-spacing: 0.5px;
}
</style>
