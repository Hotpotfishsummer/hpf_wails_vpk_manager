<template>
  <div :class="['action-bar-docked', `dock-${props.position}`]">
    <!-- 忙碌进度条 (若处于顶部停靠则置于下边缘，若处于底部停靠则置于上边缘) -->
    <div
      v-if="isAnyBusy && props.position === 'bottom'"
      class="dock-progress-bar"
      :class="{ 'is-clickable': store.isUploading }"
      @click="handleProgressBarClick"
      :title="store.isUploading ? $t('uploadProgressTitle') : ''"
    >
      <div class="dock-progress-fill" :style="{ width: progressWidth }"></div>
      <span class="dock-progress-text">{{ busyText }}</span>
    </div>

    <div class="dock-body">
      <!-- 左侧：选中统计与取消全选 -->
      <div class="dock-section dock-left">
        <div class="selection-badge">
          <span class="badge-dot"></span>
          <span class="badge-text">{{ selectionCountText }}</span>
        </div>

        <button
          class="dock-btn-clear"
          @click="handleClearSelection"
          :title="$t('deselectAll')"
        >
          <SvgIcon name="x" :size="13" />
          <span>{{ $t('deselectAll') }}</span>
        </button>
      </div>

      <div class="dock-divider"></div>

      <!-- 中间：核心操作按钮组 -->
      <div class="dock-section dock-actions">
        <!-- 导出文件 -->
        <button
          @click="handleExport"
          :disabled="store.isExporting"
          class="dock-btn btn-export"
          :class="{ 'btn-loading': store.isExporting }"
          :title="$t('exportFiles')"
        >
          <SvgIcon v-if="store.isExporting" name="loader" :size="14" :spin="true" />
          <SvgIcon v-else name="export" :size="14" />
          <span>{{ store.isExporting ? $t('exporting') : $t('exportFiles') }}</span>
        </button>

        <!-- 本地化工坊模组 -->
        <button
          v-if="store.selectedWorkshopVpks.size > 0 && store.selectedLocalVpks.size === 0"
          @click="handleLocalize"
          :disabled="store.isLocalizing"
          class="dock-btn btn-localize"
          :class="{ 'btn-loading': store.isLocalizing }"
          :title="$t('localizeTooltip')"
        >
          <SvgIcon v-if="store.isLocalizing" name="loader" :size="14" :spin="true" />
          <SvgIcon v-else name="box" :size="14" />
          <span>{{ store.isLocalizing ? $t('converting') : $t('localize') }}</span>
        </button>

        <!-- 上传到远程 SFTP -->
        <button
          v-if="store.isSftpActive && (store.selectedLocalVpks.size > 0 || store.selectedWorkshopVpks.size > 0)"
          @click="store.uploadLocalVpks"
          :disabled="store.isUploading"
          class="dock-btn btn-upload"
          :class="{ 'btn-loading': store.isUploading }"
          :title="$t('uploadToSftp')"
        >
          <SvgIcon v-if="store.isUploading" name="loader" :size="14" :spin="true" />
          <SvgIcon v-else name="cloud-upload" :size="14" />
          <span>{{ store.isUploading ? $t('uploading') : $t('uploadToSftp') }}</span>
        </button>

        <!-- 模组分组操作 -->
        <button
          @click="store.openGroupModal('assign')"
          class="dock-btn btn-group"
          :title="$t('groupManagement')"
        >
          <SvgIcon name="layers" :size="14" />
          <span>{{ $t('group') }}</span>
        </button>

        <!-- 快捷移出分组 (当选中的文件中包含已分组项目时显示) -->
        <button
          v-if="hasGroupedSelected"
          @click="handleRemoveFromGroup"
          class="dock-btn btn-ungroup"
          :title="$t('removeFromGroup')"
        >
          <SvgIcon name="folder-minus" :size="13" />
          <span>{{ $t('removeFromGroup') }}</span>
        </button>

        <!-- 禁用 -->
        <button
          @click="handleDisable"
          :disabled="store.isDeleting"
          class="dock-btn btn-disable"
          :title="$t('disable')"
        >
          <SvgIcon name="ban" :size="13" />
          <span>{{ $t('disable') }}</span>
        </button>

        <!-- 启用 -->
        <button
          @click="handleEnable"
          :disabled="store.isDeleting"
          class="dock-btn btn-enable"
          :title="$t('enable')"
        >
          <SvgIcon name="check" :size="13" />
          <span>{{ $t('enable') }}</span>
        </button>

        <!-- 彻底删除 -->
        <button
          @click="handleDelete"
          :disabled="store.isDeleting"
          class="dock-btn btn-delete"
          :title="$t('hardDelete')"
        >
          <SvgIcon name="trash" :size="13" />
          <span>{{ $t('hardDelete') }}</span>
        </button>
      </div>

      <div class="dock-divider"></div>

      <!-- 右侧：导出选项开关与停靠切换 -->
      <div class="dock-section dock-right">
        <label class="dock-checkbox" :title="$t('betterNamingHelp')">
          <input type="checkbox" v-model="store.useBetterNaming" />
          <span class="option-label">{{ $t('betterNaming') }}</span>
          <SvgIcon name="help-circle" :size="12" class="help-icon" />
        </label>

        <label class="dock-checkbox">
          <input type="checkbox" v-model="store.useAutoOpenExport" />
          <span class="option-label">{{ $t('autoOpenExport') }}</span>
        </label>

        <div class="dock-divider-mini"></div>

        <!-- 停靠位置切换按钮 -->
        <button
          class="dock-pos-toggle"
          @click="store.toggleActionBarPosition()"
          :title="positionTooltip"
        >
          <SvgIcon :name="props.position === 'bottom' ? 'chevron-up' : 'chevron-down'" :size="13" />
          <span class="toggle-text">{{ props.position === 'bottom' ? $t('dockToTop') : $t('dockToBottom') }}</span>
        </button>
      </div>
    </div>

    <!-- 忙碌进度条 (顶部停靠时置于下边缘) -->
    <div
      v-if="isAnyBusy && props.position === 'top'"
      class="dock-progress-bar"
      :class="{ 'is-clickable': store.isUploading }"
      @click="handleProgressBarClick"
      :title="store.isUploading ? $t('uploadProgressTitle') : ''"
    >
      <div class="dock-progress-fill" :style="{ width: progressWidth }"></div>
      <span class="dock-progress-text">{{ busyText }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import { getApi } from '@/composables/useApi'
import SvgIcon from '@/components/common/SvgIcon.vue'

const props = withDefaults(
  defineProps<{
    position?: 'bottom' | 'top'
  }>(),
  {
    position: 'bottom'
  }
)

const store = useVpkManagerStore()
const { t: $t } = useI18n()

const selectionCount = computed(() => {
  return store.selectedLocalVpks.size + store.selectedWorkshopVpks.size
})

const selectionCountText = computed(() => {
  return $t('selectionInfo', { count: selectionCount.value })
})

const positionTooltip = computed(() => {
  return props.position === 'bottom' ? $t('dockToTop') : $t('dockToBottom')
})

const isAnyBusy = computed(() => {
  return (
    store.isExporting ||
    store.isUploading ||
    store.isDownloading ||
    store.isImporting ||
    store.isLocalizing ||
    store.isDeleting
  )
})

const busyText = computed(() => {
  if (store.isExporting) {
    const { stage, done, total } = store.exportProgress
    if (stage === 'compressing') return $t('compressingArchive')
    if (stage === 'copying' && total > 0) return `${$t('packingFiles')} (${done}/${total})`
    return $t('packingFiles')
  }
  if (store.isUploading) {
    const p = store.uploadProgress
    if (p.totalFiles > 0) {
      const pct = p.overallTotal > 0 ? Math.round((p.overallBytes / p.overallTotal) * 100) : 0
      const speedStr = p.speed > 0 ? ` • ${formatSpeed(p.speed)}` : ''
      const filePart = p.currentFile ? `: ${p.currentFile}` : ''
      return `${$t('uploading')} (${p.currentIndex}/${p.totalFiles}) ${pct}%${filePart}${speedStr}`
    }
    return $t('uploadingLocal')
  }
  if (store.isDownloading) return $t('downloadingRemote')
  if (store.isImporting) return $t('importing')
  if (store.isLocalizing) return $t('converting')
  if (store.isDeleting) return $t('processing')
  return ''
})

const progressWidth = computed(() => {
  if (store.isExporting && store.exportProgress.total > 0) {
    const pct = Math.round((store.exportProgress.done / store.exportProgress.total) * 100)
    return `${pct}%`
  }
  if (store.isUploading && store.uploadProgress.overallTotal > 0) {
    const pct = Math.round((store.uploadProgress.overallBytes / store.uploadProgress.overallTotal) * 100)
    return `${Math.min(100, Math.max(0, pct))}%`
  }
  return '100%'
})

function formatSpeed(bytesPerSec: number): string {
  if (!bytesPerSec || bytesPerSec <= 0) return '0 KB/s'
  const units = ['B/s', 'KB/s', 'MB/s', 'GB/s']
  let idx = 0
  let val = bytesPerSec
  while (val >= 1024 && idx < units.length - 1) {
    val /= 1024
    idx++
  }
  return `${val.toFixed(val >= 100 || idx === 0 ? 0 : 1)} ${units[idx]}`
}

function handleProgressBarClick() {
  if (store.isUploading) {
    store.uploadProgress.visible = true
    store.uploadProgress.collapsed = false
  }
}

function handleClearSelection() {
  store.clearSelection()
}

async function handleExport() {
  try {
    if (store.selectedLocalVpks.size === 0 && store.selectedWorkshopVpks.size === 0) {
      store.showNotification($t('selectFiles'), 'error')
      return
    }

    let outputDir = ''
    try {
      const result = await getApi().get_downloads_directory()
      outputDir = result || ''
    } catch {
      outputDir = '.'
    }

    const result = await store.exportVpks(outputDir)
    if (result && result.exists) {
      if (await store.showConfirm($t('confirmOverwriteExport'))) {
        await store.exportVpks(outputDir, true)
      }
    }
  } catch (error) {
    store.showNotification(error instanceof Error ? error.message : $t('exportError'), 'error')
  }
}

async function handleLocalize() {
  try {
    await store.localizeVpks()
  } catch (error) {
    store.showNotification(error instanceof Error ? error.message : $t('localizeFailed'), 'error')
  }
}

async function handleDisable() {
  if (selectionCount.value === 0) return
  const pathsToDisable = [...store.selectedLocalVpks, ...store.selectedWorkshopVpks].filter(
    (path) => !path.endsWith('.disabled')
  )

  if (pathsToDisable.length > 0) {
    await store.deleteVpks(true, pathsToDisable)
  }
}

async function handleEnable() {
  if (selectionCount.value === 0) return
  const pathsToEnable = [...store.selectedLocalVpks, ...store.selectedWorkshopVpks].filter(
    (path) => path.endsWith('.disabled')
  )

  if (pathsToEnable.length > 0) {
    await store.deleteVpks(true, pathsToEnable)
  }
}

async function handleDelete() {
  const selectedCountVal = selectionCount.value
  if (selectedCountVal === 0) {
    store.showNotification($t('selectFiles'), 'error')
    return
  }

  const message = $t('confirmHardDelete') + ` (${selectedCountVal})`
  if (await store.showConfirm(message)) {
    await store.deleteVpks(false)
  }
}

const hasGroupedSelected = computed(() => {
  return store.getSelectedVpkPaths().some((p) => store.getFileGroup(p) !== undefined)
})

function handleRemoveFromGroup() {
  store.removeFilesFromGroups(store.getSelectedVpkPaths())
}
</script>

<style scoped>
/* 贴合固定工具条容器：纯色沉浸式深色底色，无毛玻璃马赛克模糊 */
.action-bar-docked {
  width: 100%;
  flex-shrink: 0;
  background: #0d1525;
  background: linear-gradient(180deg, #10192b 0%, #0c1220 100%);
  position: relative;
  z-index: 80;
  display: flex;
  flex-direction: column;
}

/* 底部贴合形态 */
.action-bar-docked.dock-bottom {
  border-top: 1px solid rgba(56, 189, 248, 0.22);
  box-shadow: 0 -6px 24px rgba(0, 0, 0, 0.5), inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

/* 顶部贴合形态 */
.action-bar-docked.dock-top {
  border-bottom: 1px solid rgba(56, 189, 248, 0.22);
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.5), inset 0 -1px 0 rgba(255, 255, 255, 0.05);
}

/* 进度条 */
.dock-progress-bar {
  width: 100%;
  height: 18px;
  background: rgba(0, 0, 0, 0.5);
  position: relative;
  display: flex;
  align-items: center;
  overflow: hidden;
}

.dock-progress-bar.is-clickable {
  cursor: pointer;
}

.dock-progress-bar.is-clickable:hover .dock-progress-text {
  text-decoration: underline;
  color: #38bdf8;
}

.dock-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #0284c7 0%, #38bdf8 100%);
  transition: width 0.3s ease;
}

.dock-progress-text {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  font-weight: 600;
  color: #fff;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.8);
}

/* 工具栏主体内容区 */
.dock-body {
  width: 100%;
  max-width: 1480px;
  margin: 0 auto;
  padding: 8px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 50px;
  overflow-x: auto;
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.15) transparent;
}

.dock-body::-webkit-scrollbar {
  height: 3px;
}

.dock-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 2px;
}

.dock-section {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

/* 左侧区域：已选统计与清空按钮 */
.dock-left {
  gap: 10px;
}

.selection-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 4px 12px;
  background: rgba(56, 189, 248, 0.12);
  border: 1px solid rgba(56, 189, 248, 0.3);
  border-radius: var(--radius-full);
}

.badge-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent-primary);
  box-shadow: 0 0 6px var(--accent-primary);
}

.badge-text {
  font-size: 12px;
  font-weight: 600;
  color: var(--accent-primary);
  white-space: nowrap;
}

.dock-btn-clear {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 500;
  color: var(--text-dim);
  background: transparent;
  border: 1px solid transparent;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.dock-btn-clear:hover {
  color: var(--text-main);
  background: rgba(255, 255, 255, 0.06);
  border-color: var(--border-subtle);
}

/* 分割线 */
.dock-divider {
  width: 1px;
  height: 22px;
  background: rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
}

.dock-divider-mini {
  width: 1px;
  height: 16px;
  background: rgba(255, 255, 255, 0.08);
  flex-shrink: 0;
}

/* 中间操作按钮组 */
.dock-actions {
  gap: 8px;
}

.dock-btn {
  height: 34px;
  padding: 0 13px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-weight: 600;
  transition: all var(--transition-fast);
  white-space: nowrap;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
  cursor: pointer;
}

.dock-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* 导出按钮 */
.btn-export {
  background: linear-gradient(135deg, #0284c7 0%, #0369a1 100%);
  color: #fff;
  box-shadow: 0 2px 6px rgba(2, 132, 199, 0.3);
  border: 1px solid rgba(56, 189, 248, 0.35);
}

.btn-export:hover:not(:disabled) {
  background: linear-gradient(135deg, #0369a1 0%, #075985 100%);
  box-shadow: 0 4px 12px rgba(2, 132, 199, 0.5);
  filter: brightness(1.08);
}

/* 本地化按钮 */
.btn-localize {
  background: rgba(139, 92, 246, 0.15);
  color: #c4b5fd;
  border: 1px solid rgba(139, 92, 246, 0.35);
}

.btn-localize:hover:not(:disabled) {
  background: rgba(139, 92, 246, 0.25);
  border-color: rgba(139, 92, 246, 0.55);
}

/* 上传按钮 */
.btn-upload {
  background: rgba(6, 182, 212, 0.15);
  color: #67e8f9;
  border: 1px solid rgba(6, 182, 212, 0.35);
}

.btn-upload:hover:not(:disabled) {
  background: rgba(6, 182, 212, 0.25);
  border-color: rgba(6, 182, 212, 0.55);
}

/* 分组按钮 */
.btn-group {
  background: rgba(14, 165, 233, 0.15);
  color: #38bdf8;
  border: 1px solid rgba(14, 165, 233, 0.35);
}

.btn-group:hover:not(:disabled) {
  background: rgba(14, 165, 233, 0.28);
  border-color: rgba(14, 165, 233, 0.6);
  color: #ffffff;
}

/* 移出分组按钮 */
.btn-ungroup {
  background: rgba(245, 158, 11, 0.12);
  color: #fcd34d;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.btn-ungroup:hover:not(:disabled) {
  background: rgba(245, 158, 11, 0.25);
  border-color: rgba(245, 158, 11, 0.55);
  color: #ffffff;
}

/* 禁用按钮 */
.btn-disable {
  background: rgba(245, 158, 11, 0.12);
  color: #fcd34d;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.btn-disable:hover:not(:disabled) {
  background: rgba(245, 158, 11, 0.22);
  border-color: rgba(245, 158, 11, 0.5);
}

/* 启用按钮 */
.btn-enable {
  background: rgba(16, 185, 129, 0.12);
  color: #6ee7b7;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.btn-enable:hover:not(:disabled) {
  background: rgba(16, 185, 129, 0.22);
  border-color: rgba(16, 185, 129, 0.5);
}

/* 删除按钮 */
.btn-delete {
  background: rgba(244, 63, 94, 0.12);
  color: #fda4af;
  border: 1px solid rgba(244, 63, 94, 0.3);
}

.btn-delete:hover:not(:disabled) {
  background: rgba(244, 63, 94, 0.22);
  border-color: rgba(244, 63, 94, 0.5);
}

/* 右侧选项区域：保证单行不换行 */
.dock-right {
  gap: 12px;
}

.dock-checkbox {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  user-select: none;
  white-space: nowrap;
  flex-shrink: 0;
  padding: 4px 6px;
  border-radius: var(--radius-sm);
  transition: background var(--transition-fast);
}

.dock-checkbox:hover {
  background: rgba(255, 255, 255, 0.04);
}

.dock-checkbox input {
  accent-color: var(--accent-primary);
  width: 14px;
  height: 14px;
  cursor: pointer;
  flex-shrink: 0;
}

.option-label {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-muted);
  white-space: nowrap;
  flex-shrink: 0;
}

.help-icon {
  color: var(--text-dim);
  flex-shrink: 0;
}

/* 停靠位置切换按钮 */
.dock-pos-toggle {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  height: 28px;
  padding: 0 10px;
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-subtle);
  color: var(--text-muted);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
  flex-shrink: 0;
}

.dock-pos-toggle:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-main);
  border-color: var(--border-medium);
}

.toggle-text {
  white-space: nowrap;
}
</style>
