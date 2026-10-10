<template>
  <div class="sftp-browser">
    <!-- 头部工具栏 -->
    <div class="browser-header">
      <!-- 顶行：标题、路径与核心操作 -->
      <div class="header-top-row">
        <div class="header-title-box">
          <SvgIcon name="globe" :size="16" class="sftp-icon" />
          <h3 class="header-title">{{ $t('remoteAddons') }}</h3>
          <span class="path-badge" :title="store.sftpSettings.addons_path">
            {{ store.sftpSettings.addons_path || $t('pathNotConfigured') }}
          </span>
        </div>

        <div class="header-actions">
          <!-- 刷新远程文件列表 -->
          <button
            @click="store.scanRemoteVpks()"
            class="header-btn btn-refresh"
            :disabled="store.isRemoteScanning"
            :title="$t('refresh')"
          >
            <SvgIcon name="refresh" :size="14" :spin="store.isRemoteScanning" />
            <span class="btn-text">{{ $t('refresh') }}</span>
          </button>

          <!-- 重启远程服务器 -->
          <button
            @click="handleRestartServer"
            class="header-btn btn-restart"
            :disabled="store.isRestartingServer"
            :title="$t('restartServerTooltip')"
          >
            <SvgIcon v-if="store.isRestartingServer" name="loader" :size="14" :spin="true" />
            <SvgIcon v-else name="zap" :size="14" />
            <span class="btn-text">{{ store.isRestartingServer ? $t('restartingServer') : $t('restartServer') }}</span>
          </button>

          <!-- 关闭分屏 -->
          <button
            @click="store.isSftpActive = false"
            class="header-btn btn-close"
            :title="$t('close', '关闭分屏')"
          >
            <SvgIcon name="x" :size="14" />
          </button>
        </div>
      </div>

      <!-- 次行：搜索过滤与排序 -->
      <div class="header-sub-row">
        <div class="search-wrap">
          <SvgIcon name="search" :size="13" class="search-icon" />
          <input
            v-model="remoteSearchQuery"
            type="text"
            class="remote-search-input"
            :placeholder="$t('searchFilesPlaceholder')"
          />
          <button
            v-if="remoteSearchQuery"
            class="clear-search-btn"
            @click="remoteSearchQuery = ''"
          >
            <SvgIcon name="x" :size="12" />
          </button>
        </div>

        <div class="sort-controls">
          <button
            v-for="key in sortKeys"
            :key="key"
            class="sort-btn"
            :class="{ active: sortField === key }"
            @click="toggleSort(key)"
          >
            {{ $t(sortLabels[key]) }}
            <span v-if="sortField === key" class="sort-indicator">{{ sortAsc ? '▲' : '▼' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 内容区域 -->
    <div class="browser-content">
      <!-- 多选操作栏 -->
      <div v-if="store.remoteVpkFiles.length > 0" class="selection-actions">
        <div class="selection-info">
          <label class="custom-checkbox" @click.stop>
            <input
              type="checkbox"
              :checked="isAllSelected"
              @change="toggleSelectAll"
            />
            <span class="checkbox-mark" :class="{ 'is-checked': isAllSelected }">
              <SvgIcon v-if="isAllSelected" name="check" :size="12" stroke-width="3.2" class="check-icon" />
            </span>
            <span class="count-text">
              {{ store.selectedRemoteVpks.size > 0 ? $t('selectionInfo', { count: store.selectedRemoteVpks.size }) : $t('selectAll') }}
            </span>
          </label>
        </div>

        <div class="batch-action-btns" v-if="store.selectedRemoteVpks.size > 0">
          <button
            @click="store.downloadRemoteVpks"
            class="mini-action-btn download"
            :title="$t('syncToLocal')"
            :disabled="store.isDownloading"
          >
            <SvgIcon v-if="store.isDownloading" name="loader" :size="13" :spin="true" />
            <SvgIcon v-else name="import" :size="13" />
          </button>
          <button @click="handleBatchEnable" class="mini-action-btn enable" :title="$t('enable')">
            <SvgIcon name="check" :size="13" />
          </button>
          <button @click="handleBatchDisable" class="mini-action-btn disable" :title="$t('disable')">
            <SvgIcon name="ban" :size="13" />
          </button>
          <button @click="handleBatchDelete" class="mini-action-btn delete" :title="$t('delete')">
            <SvgIcon name="trash" :size="13" />
          </button>
        </div>
      </div>

      <!-- 加载中 -->
      <div v-if="store.isRemoteScanning" class="remote-loading">
        <SvgIcon name="loader" :size="24" :spin="true" class="spinner-icon" />
        <p class="loading-label">{{ $t('loading') }}...</p>
      </div>

      <!-- 空状态 -->
      <div v-else-if="store.remoteVpkFiles.length === 0" class="remote-empty">
        <SvgIcon name="folder" :size="32" class="empty-folder-icon" />
        <p>{{ $t('noRemoteFiles') }}</p>
      </div>

      <!-- 远程列表 -->
      <div v-else class="remote-list">
        <div
          v-for="file in sortedRemoteFiles"
          :key="file.path"
          class="remote-item"
          :class="{
            'is-disabled': file.is_disabled,
            'is-selected': store.selectedRemoteVpks.has(file.path),
          }"
          role="button"
          tabindex="0"
          @click="store.toggleRemoteVpkSelection(file.path)"
        >
          <label class="custom-checkbox" @click.stop>
            <input
              type="checkbox"
              :checked="store.selectedRemoteVpks.has(file.path)"
              @change="store.toggleRemoteVpkSelection(file.path)"
            />
            <span class="checkbox-mark" :class="{ 'is-checked': store.selectedRemoteVpks.has(file.path) }">
              <SvgIcon v-if="store.selectedRemoteVpks.has(file.path)" name="check" :size="12" stroke-width="3.2" class="check-icon" />
            </span>
          </label>

          <div class="item-icon-box">
            <SvgIcon :name="file.is_disabled ? 'ban' : 'box'" :size="14" />
          </div>

          <div class="item-details">
            <span class="item-name" :title="file.name">{{ file.name }}</span>
            <span class="item-meta">{{ file.size }} MB • {{ formatDate(file.modified_time) }}</span>
          </div>

          <div class="item-status">
            <span v-if="file.is_disabled" class="status-pill disabled">已禁用</span>
            <span v-else class="status-pill active">在线</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import type { RemoteFile } from '@/types/api'
import SvgIcon from '@/components/common/SvgIcon.vue'

type SortField = 'name' | 'size' | 'modified_time'
const sortKeys: SortField[] = ['name', 'size', 'modified_time']
const sortField = ref<SortField>('name')
const sortAsc = ref(true)
const remoteSearchQuery = ref('')
const sortLabels: Record<SortField, string> = {
  name: 'sortByName',
  size: 'sortBySize',
  modified_time: 'sortByDate',
}

function toggleSort(key: SortField) {
  if (sortField.value === key) {
    sortAsc.value = !sortAsc.value
  } else {
    sortField.value = key
    sortAsc.value = true
  }
}

const store = useVpkManagerStore()
const { t: $t } = useI18n()

onMounted(async () => {
  if (store.sftpSettings.address && store.remoteVpkFiles.length === 0) {
    await store.scanRemoteVpks()
  }
})

const filteredRemoteFiles = computed(() => {
  if (!remoteSearchQuery.value) return store.remoteVpkFiles
  const q = remoteSearchQuery.value.toLowerCase()
  return store.remoteVpkFiles.filter((f) => f.name.toLowerCase().includes(q))
})

const sortedRemoteFiles = computed(() => {
  return [...filteredRemoteFiles.value].sort((a, b) => {
    let cmp = 0
    if (sortField.value === 'name') cmp = a.name.localeCompare(b.name)
    else if (sortField.value === 'size') cmp = a.size - b.size
    else if (sortField.value === 'modified_time') cmp = a.modified_time - b.modified_time
    return sortAsc.value ? cmp : -cmp
  })
})

const isAllSelected = computed(() => {
  return (
    sortedRemoteFiles.value.length > 0 &&
    sortedRemoteFiles.value.every((f) => store.selectedRemoteVpks.has(f.path))
  )
})

function toggleSelectAll() {
  if (isAllSelected.value) {
    store.deselectAllRemote()
  } else {
    store.selectAllRemote(sortedRemoteFiles.value.map((f) => f.path))
  }
}

async function handleRestartServer() {
  const confirmed = await store.showConfirm($t('confirmRestartServer'))
  if (confirmed) {
    await store.restartL4d2Server()
  }
}

async function handleBatchEnable() {
  await store.enableRemoteVpks()
}

async function handleBatchDisable() {
  await store.disableRemoteVpks()
}

async function handleBatchDelete() {
  const count = store.selectedRemoteVpks.size
  if (count === 0) return
  const confirmed = await store.showConfirm($t('confirmBatchDeleteRemote', { count }))
  if (confirmed) {
    await store.deleteMultipleRemoteVpks()
  }
}

function formatDate(timestamp: number): string {
  if (!timestamp) return ''
  const d = new Date(timestamp * 1000)
  return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours().toString().padStart(2, '0')}:${d
    .getMinutes()
    .toString()
    .padStart(2, '0')}`
}
</script>

<style scoped>
.sftp-browser {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--bg-surface);
  border-left: 1px solid var(--border-subtle);
  overflow: hidden;
}

.browser-header {
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.25);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.header-top-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.header-title-box {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  flex: 1;
}

.sftp-icon {
  color: var(--accent-primary);
  flex-shrink: 0;
}

.header-title {
  font-size: 14px;
  font-weight: 700;
  color: #fff;
  white-space: nowrap;
}

.path-badge {
  font-size: 11px;
  font-family: monospace;
  color: var(--text-dim);
  background: rgba(255, 255, 255, 0.04);
  padding: 2px 6px;
  border-radius: 4px;
  border: 1px solid var(--border-subtle);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 160px;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.header-btn {
  height: 28px;
  padding: 0 10px;
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border-subtle);
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
}

.header-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-restart {
  background: rgba(245, 158, 11, 0.12);
  color: var(--accent-warning);
  border-color: rgba(245, 158, 11, 0.3);
}

.btn-restart:hover:not(:disabled) {
  background: rgba(245, 158, 11, 0.22);
}

.btn-close {
  width: 28px;
  padding: 0;
  justify-content: center;
}

/* 次行 */
.header-sub-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  padding: 0 8px;
  height: 28px;
  gap: 6px;
}

.search-icon {
  color: var(--text-dim);
}

.remote-search-input {
  flex: 1;
  background: transparent !important;
  border: none !important;
  padding: 0 !important;
  font-size: 12px;
  color: var(--text-main);
}

.remote-search-input:focus {
  box-shadow: none !important;
}

.clear-search-btn {
  width: 16px;
  height: 16px;
  border-radius: 50%;
  color: var(--text-dim);
  display: flex;
  align-items: center;
  justify-content: center;
}

.sort-controls {
  display: flex;
  background: rgba(0, 0, 0, 0.25);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  padding: 2px;
}

.sort-btn {
  height: 24px;
  padding: 0 6px;
  font-size: 11px;
  color: var(--text-dim);
  border-radius: 3px;
}

.sort-btn.active {
  background: rgba(255, 255, 255, 0.1);
  color: var(--accent-primary);
  font-weight: 600;
}

.sort-indicator {
  font-size: 8px;
  margin-left: 2px;
}

/* 内容区 */
.browser-content {
  flex: 1;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.selection-actions {
  padding: 8px 14px;
  background: rgba(56, 189, 248, 0.06);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.count-text {
  font-size: 12px;
  font-weight: 600;
  color: var(--accent-primary);
  margin-left: 6px;
}

.batch-action-btns {
  display: flex;
  align-items: center;
  gap: 4px;
}

.mini-action-btn {
  appearance: none;
  -webkit-appearance: none;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-sm);
  background: #182337;
  border: 1px solid var(--border-medium);
  color: #cbd5e1;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.mini-action-btn:hover {
  background: #22314a;
  color: #fff;
  border-color: var(--border-strong);
}

.mini-action-btn.enable:hover {
  background: rgba(16, 185, 129, 0.25);
  color: var(--accent-success);
  border-color: rgba(16, 185, 129, 0.5);
}

.mini-action-btn.disable:hover {
  background: rgba(245, 158, 11, 0.25);
  color: var(--accent-warning);
  border-color: rgba(245, 158, 11, 0.5);
}

.mini-action-btn.delete:hover {
  background: rgba(244, 63, 94, 0.25);
  color: var(--accent-danger);
  border-color: rgba(244, 63, 94, 0.5);
}

.remote-loading,
.remote-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-dim);
  padding: 40px;
  font-size: 13px;
}

.spinner-icon {
  color: var(--accent-primary);
}

.empty-folder-icon {
  opacity: 0.5;
}

.remote-list {
  padding: 8px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.remote-item {
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.remote-item:hover {
  background: var(--bg-card-hover);
  border-color: var(--border-medium);
}

.remote-item.is-selected {
  background: rgba(56, 189, 248, 0.08);
  border-color: var(--border-active);
}

.remote-item.is-disabled {
  opacity: 0.6;
}

.item-icon-box {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.04);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-dim);
  flex-shrink: 0;
}

.item-details {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.item-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-main);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-meta {
  font-size: 11px;
  color: var(--text-dim);
}

.status-pill {
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: var(--radius-full);
}

.status-pill.active {
  background: rgba(16, 185, 129, 0.15);
  color: var(--accent-success);
}

.status-pill.disabled {
  background: rgba(245, 158, 11, 0.15);
  color: var(--accent-warning);
}

/* 复选框样式 */
.custom-checkbox {
  position: relative;
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
}

.custom-checkbox input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.checkbox-mark {
  width: 17px;
  height: 17px;
  border-radius: 4px;
  border: 1.5px solid #475569;
  background: #141c2c;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
  flex-shrink: 0;
}

.custom-checkbox:hover .checkbox-mark {
  border-color: #38bdf8;
  background: #192336;
}

.custom-checkbox input:checked + .checkbox-mark,
.checkbox-mark.is-checked {
  background: #0284c7 !important;
  border-color: #38bdf8 !important;
  color: #ffffff !important;
  box-shadow: 0 0 8px rgba(56, 189, 248, 0.45);
}

.check-icon {
  color: #ffffff !important;
}
</style>
