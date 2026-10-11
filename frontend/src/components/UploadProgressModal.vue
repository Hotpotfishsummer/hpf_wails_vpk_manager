<template>
  <transition name="slide-card">
    <div
      v-if="store.uploadProgress.visible"
      :class="['upload-progress-widget', { 'is-collapsed': store.uploadProgress.collapsed, 'dock-shifted': hasBottomDock }]"
      role="region"
      aria-label="上传进度指示器"
    >
      <!-- 折叠态胶囊 (Pill Mode) -->
      <div v-if="store.uploadProgress.collapsed" class="pill-container" @click="store.toggleUploadProgressCollapsed">
        <div class="pill-status-icon">
          <SvgIcon v-if="isUploading" name="cloud-upload" :size="16" class="upload-icon-pulse" />
          <SvgIcon v-else-if="isCompleted" name="check-circle" :size="16" class="text-success" />
          <SvgIcon v-else-if="isCanceled" name="ban" :size="16" class="text-muted" />
          <SvgIcon v-else name="alert-triangle" :size="16" class="text-danger" />
        </div>

        <div class="pill-text">
          <span v-if="isUploading">
            {{ $t('uploading') }} ({{ store.uploadProgress.currentIndex }}/{{ store.uploadProgress.totalFiles }}) {{ overallPercent }}% • {{ formattedSpeed }}
          </span>
          <span v-else-if="isCompleted">
            {{ $t('uploadCompletedTitle') }} ({{ store.uploadProgress.successCount }}/{{ store.uploadProgress.totalFiles }})
          </span>
          <span v-else-if="isCanceled">
            {{ $t('uploadCanceled') }}
          </span>
          <span v-else>
            {{ $t('uploadFailed') }}
          </span>
        </div>

        <button
          class="pill-btn"
          @click.stop="store.toggleUploadProgressCollapsed"
          :title="$t('expand', '展开详情')"
        >
          <SvgIcon name="chevron-up" :size="13" />
        </button>

        <button
          v-if="!isUploading"
          class="pill-btn close"
          @click.stop="store.closeUploadProgress"
          :title="$t('close', '关闭')"
        >
          <SvgIcon name="x" :size="13" />
        </button>
      </div>

      <!-- 展开态卡片 (Detailed Card Mode) -->
      <div v-else class="card-container">
        <!-- 头部导航与操作 -->
        <div class="card-header">
          <div class="header-left">
            <div class="header-icon-box" :class="statusThemeClass">
              <SvgIcon v-if="isUploading" name="cloud-upload" :size="18" class="upload-icon-pulse" />
              <SvgIcon v-else-if="isCompleted" name="check-circle" :size="18" />
              <SvgIcon v-else-if="isCanceled" name="ban" :size="18" />
              <SvgIcon v-else name="alert-triangle" :size="18" />
            </div>

            <div class="header-title-box">
              <div class="title-row">
                <h4 class="card-title">{{ $t('uploadProgressTitle') }}</h4>
                <span class="status-badge" :class="statusThemeClass">
                  <span v-if="isUploading" class="pulse-dot"></span>
                  {{ statusBadgeText }}
                </span>
              </div>
              <p class="card-subtitle">
                {{ subtitleText }}
              </p>
            </div>
          </div>

          <div class="header-actions">
            <button
              class="header-icon-btn"
              @click="store.toggleUploadProgressCollapsed"
              :title="$t('collapse', '收起')"
            >
              <SvgIcon name="chevron-down" :size="15" />
            </button>
            <button
              class="header-icon-btn"
              @click="handleCloseOrCancel"
              :title="isUploading ? $t('cancelUpload') : $t('close', '关闭')"
            >
              <SvgIcon name="x" :size="15" />
            </button>
          </div>
        </div>

        <!-- 总体进度主卡片 -->
        <div class="overall-progress-box">
          <div class="progress-bar-wrap">
            <div class="progress-track">
              <div
                class="progress-fill"
                :class="statusThemeClass"
                :style="{ width: `${overallPercent}%` }"
              >
                <div class="progress-glow"></div>
              </div>
            </div>
            <span class="progress-percent-label">{{ overallPercent }}%</span>
          </div>

          <!-- 核心指标栅格: 已传大小 / 速率 / 剩余时间 -->
          <div class="metrics-grid">
            <div class="metric-item">
              <span class="metric-label">{{ $t('uploadTransferred') }}</span>
              <span class="metric-value">{{ formattedOverallTransferred }} / {{ formattedOverallTotal }}</span>
            </div>

            <div class="metric-item">
              <span class="metric-label">{{ $t('uploadSpeed') }}</span>
              <span class="metric-value highlight-cyan">{{ formattedSpeed }}</span>
            </div>

            <div class="metric-item">
              <span class="metric-label">{{ $t('uploadEta') }}</span>
              <span class="metric-value">{{ formattedEta }}</span>
            </div>
          </div>
        </div>

        <!-- 当前文件传输详情 (仅上传中或有正在传的文件时展示) -->
        <div v-if="isUploading && store.uploadProgress.currentFile" class="current-file-box">
          <div class="current-file-header">
            <div class="file-name-row">
              <SvgIcon name="box" :size="13" class="file-icon" />
              <span class="file-name" :title="store.uploadProgress.currentFile">
                {{ store.uploadProgress.currentFile }}
              </span>
            </div>
            <span class="file-percent">{{ currentFilePercent }}%</span>
          </div>

          <div class="mini-progress-track">
            <div
              class="mini-progress-fill"
              :style="{ width: `${currentFilePercent}%` }"
            ></div>
          </div>

          <div class="current-file-meta">
            <span>{{ formattedCurrentBytes }} / {{ formattedCurrentTotal }}</span>
            <span v-if="store.uploadProgress.totalFiles > 1">
              {{ $t('uploadedFileCount', { current: store.uploadProgress.currentIndex, total: store.uploadProgress.totalFiles }) }}
            </span>
          </div>
        </div>

        <!-- 失败文件列表展开区 (如有错误) -->
        <div v-if="store.uploadProgress.failedFiles.length > 0" class="failed-files-section">
          <div class="failed-header" @click="showFailedDetails = !showFailedDetails">
            <div class="failed-title">
              <SvgIcon name="alert-triangle" :size="13" class="text-danger" />
              <span>{{ $t('failedFilesList') }} ({{ store.uploadProgress.failedFiles.length }})</span>
            </div>
            <SvgIcon :name="showFailedDetails ? 'chevron-up' : 'chevron-down'" :size="12" />
          </div>

          <div v-if="showFailedDetails" class="failed-list">
            <div
              v-for="(item, idx) in store.uploadProgress.failedFiles"
              :key="idx"
              class="failed-item"
            >
              <span class="failed-filename" :title="item.file">{{ item.file }}</span>
              <span class="failed-reason" :title="item.error">{{ item.error }}</span>
            </div>
          </div>
        </div>

        <!-- 底部操作按钮 -->
        <div class="card-footer">
          <div class="footer-stats">
            <span class="stat-badge success">
              <SvgIcon name="check" :size="11" /> {{ store.uploadProgress.successCount }}
            </span>
            <span v-if="store.uploadProgress.failedCount > 0" class="stat-badge danger">
              <SvgIcon name="x" :size="11" /> {{ store.uploadProgress.failedCount }}
            </span>
          </div>

          <div class="footer-btns">
            <!-- 上传中: 取消按钮 -->
            <button
              v-if="isUploading"
              class="btn-cancel-upload"
              @click="handleCancelUpload"
              :disabled="store.uploadProgress.isCanceled"
            >
              <SvgIcon name="ban" :size="13" />
              <span>{{ store.uploadProgress.isCanceled ? $t('processing') : $t('cancelUpload') }}</span>
            </button>

            <!-- 结束状态: 关闭按钮 -->
            <button
              v-else
              class="btn-done"
              @click="store.closeUploadProgress"
            >
              <span>{{ $t('close', '关闭') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import SvgIcon from '@/components/common/SvgIcon.vue'

const store = useVpkManagerStore()
const { t: $t } = useI18n()

const showFailedDetails = ref(false)

// 是否有底部操作栏停靠（用于避免遮挡）
const hasBottomDock = computed(() => {
  const selectionCount = store.selectedLocalVpks.size + store.selectedWorkshopVpks.size
  return selectionCount > 0 && store.actionBarPosition === 'bottom'
})

const isUploading = computed(() => {
  return store.uploadProgress.stage === 'preparing' || store.uploadProgress.stage === 'uploading'
})

const isCompleted = computed(() => {
  return store.uploadProgress.stage === 'completed'
})

const isCanceled = computed(() => {
  return store.uploadProgress.stage === 'canceled' || store.uploadProgress.isCanceled
})

const isError = computed(() => {
  return store.uploadProgress.stage === 'error'
})

const statusThemeClass = computed(() => {
  if (isUploading.value) return 'theme-blue'
  if (isCompleted.value) {
    return store.uploadProgress.failedCount > 0 ? 'theme-warning' : 'theme-success'
  }
  if (isCanceled.value) return 'theme-muted'
  return 'theme-danger'
})

const statusBadgeText = computed(() => {
  if (store.uploadProgress.stage === 'preparing') return $t('uploadPreparing')
  if (isUploading.value) return $t('uploading')
  if (isCompleted.value) {
    return store.uploadProgress.failedCount > 0 ? '部分完成' : $t('uploadCompletedTitle')
  }
  if (isCanceled.value) return $t('uploadCanceled')
  return $t('uploadFailed')
})

const subtitleText = computed(() => {
  if (store.uploadProgress.stage === 'preparing') {
    return $t('uploadPreparing')
  }
  if (isUploading.value) {
    return $t('uploadedFileCount', {
      current: store.uploadProgress.currentIndex,
      total: store.uploadProgress.totalFiles,
    })
  }
  if (isCompleted.value) {
    return $t('uploadSummary', {
      success: store.uploadProgress.successCount,
      failed: store.uploadProgress.failedCount,
    })
  }
  if (isCanceled.value) {
    return $t('uploadCanceled')
  }
  return $t('uploadFailed')
})

const overallPercent = computed(() => {
  if (store.uploadProgress.overallTotal <= 0) {
    return isCompleted.value ? 100 : 0
  }
  const pct = Math.round((store.uploadProgress.overallBytes / store.uploadProgress.overallTotal) * 100)
  return Math.min(100, Math.max(0, pct))
})

const currentFilePercent = computed(() => {
  if (store.uploadProgress.totalBytes <= 0) return 0
  const pct = Math.round((store.uploadProgress.currentBytes / store.uploadProgress.totalBytes) * 100)
  return Math.min(100, Math.max(0, pct))
})

const formattedSpeed = computed(() => {
  if (!isUploading.value || store.uploadProgress.speed <= 0) {
    return isCompleted.value ? '--' : '0 KB/s'
  }
  return formatBytes(store.uploadProgress.speed) + '/s'
})

const formattedOverallTransferred = computed(() => {
  return formatBytes(store.uploadProgress.overallBytes)
})

const formattedOverallTotal = computed(() => {
  return formatBytes(store.uploadProgress.overallTotal)
})

const formattedCurrentBytes = computed(() => {
  return formatBytes(store.uploadProgress.currentBytes)
})

const formattedCurrentTotal = computed(() => {
  return formatBytes(store.uploadProgress.totalBytes)
})

const formattedEta = computed(() => {
  if (!isUploading.value) return '--'
  if (store.uploadProgress.speed <= 0) return $t('uploadEtaCalculating')

  const remainingBytes = store.uploadProgress.overallTotal - store.uploadProgress.overallBytes
  if (remainingBytes <= 0) return '0秒'

  const seconds = Math.ceil(remainingBytes / store.uploadProgress.speed)
  return formatDuration(seconds)
})

function formatBytes(bytes: number): string {
  if (!bytes || bytes <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB']
  let idx = 0
  let val = bytes
  while (val >= 1024 && idx < units.length - 1) {
    val /= 1024
    idx++
  }
  return `${val.toFixed(val >= 100 || idx === 0 ? 0 : 1)} ${units[idx]}`
}

function formatDuration(totalSeconds: number): string {
  if (totalSeconds < 60) {
    return `${totalSeconds}秒`
  }
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60
  if (minutes < 60) {
    return `${minutes}分${seconds > 0 ? `${seconds}秒` : ''}`
  }
  const hours = Math.floor(minutes / 60)
  const remMinutes = minutes % 60
  return `${hours}时${remMinutes}分`
}

async function handleCancelUpload() {
  if (await store.showConfirm($t('cancelUploadConfirm'))) {
    await store.cancelUpload()
  }
}

async function handleCloseOrCancel() {
  if (isUploading.value) {
    await handleCancelUpload()
  } else {
    store.closeUploadProgress()
  }
}
</script>

<style scoped>
.upload-progress-widget {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 120;
  max-width: 440px;
  width: calc(100vw - 48px);
  transition: bottom 0.25s ease, transform 0.25s ease, opacity 0.25s ease;
  user-select: none;
}

.upload-progress-widget.dock-shifted {
  bottom: 80px;
}

/* 折叠态胶囊 (Pill) */
.pill-container {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  padding: 8px 14px;
  background: rgba(15, 23, 42, 0.92);
  border: 1px solid rgba(56, 189, 248, 0.35);
  box-shadow: 0 10px 28px rgba(0, 0, 0, 0.45);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-radius: 9999px;
  cursor: pointer;
  transition: all 0.2s ease;
  float: right;
}

.pill-container:hover {
  background: rgba(30, 41, 59, 0.96);
  border-color: rgba(56, 189, 248, 0.6);
  transform: translateY(-2px);
}

.pill-status-icon {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.pill-text {
  font-size: 13px;
  font-weight: 500;
  color: #f1f5f9;
  white-space: nowrap;
}

.pill-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  color: #94a3b8;
  transition: all 0.15s ease;
}

.pill-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  color: #ffffff;
}

/* 展开态卡片 (Card) */
.card-container {
  background: rgba(15, 23, 42, 0.94);
  border: 1px solid rgba(148, 163, 184, 0.18);
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.55), 0 0 0 1px rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-radius: 14px;
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  gap: 14px;
  overflow: hidden;
}

/* 头部 */
.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.header-icon-box {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.header-icon-box.theme-blue {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
  border: 1px solid rgba(56, 189, 248, 0.3);
}

.header-icon-box.theme-success {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.header-icon-box.theme-warning {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.header-icon-box.theme-danger {
  background: rgba(244, 63, 94, 0.15);
  color: #f43f5e;
  border: 1px solid rgba(244, 63, 94, 0.3);
}

.header-icon-box.theme-muted {
  background: rgba(148, 163, 184, 0.15);
  color: #94a3b8;
  border: 1px solid rgba(148, 163, 184, 0.3);
}

.header-title-box {
  flex: 1;
  min-width: 0;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.card-title {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #f8fafc;
  line-height: 1.2;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 7px;
  border-radius: 9999px;
  font-size: 11px;
  font-weight: 500;
  line-height: 1;
}

.status-badge.theme-blue {
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
}

.status-badge.theme-success {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.status-badge.theme-warning {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}

.status-badge.theme-danger {
  background: rgba(244, 63, 94, 0.15);
  color: #f43f5e;
}

.status-badge.theme-muted {
  background: rgba(148, 163, 184, 0.15);
  color: #94a3b8;
}

.pulse-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: #38bdf8;
  animation: pulse-glow 1.5s infinite;
}

.card-subtitle {
  margin: 3px 0 0 0;
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 4px;
}

.header-icon-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  border-radius: 6px;
  color: #94a3b8;
  transition: all 0.15s ease;
}

.header-icon-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}

/* 总体进度条 */
.overall-progress-box {
  background: rgba(30, 41, 59, 0.45);
  border: 1px solid rgba(148, 163, 184, 0.1);
  border-radius: 10px;
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.progress-bar-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.progress-track {
  flex: 1;
  height: 8px;
  background: rgba(15, 23, 42, 0.7);
  border-radius: 9999px;
  overflow: hidden;
  position: relative;
}

.progress-fill {
  height: 100%;
  border-radius: 9999px;
  transition: width 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

.progress-fill.theme-blue {
  background: linear-gradient(90deg, #3b82f6, #06b6d4);
}

.progress-fill.theme-success {
  background: linear-gradient(90deg, #10b981, #34d399);
}

.progress-fill.theme-warning {
  background: linear-gradient(90deg, #f59e0b, #fbbf24);
}

.progress-fill.theme-danger {
  background: linear-gradient(90deg, #ef4444, #f87171);
}

.progress-fill.theme-muted {
  background: #64748b;
}

.progress-glow {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 14px;
  background: rgba(255, 255, 255, 0.45);
  filter: blur(2px);
}

.progress-percent-label {
  font-size: 13px;
  font-weight: 700;
  font-family: monospace;
  color: #f8fafc;
  min-width: 38px;
  text-align: right;
}

/* 关键指标栅格 */
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  border-top: 1px solid rgba(148, 163, 184, 0.08);
  padding-top: 8px;
}

.metric-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.metric-label {
  font-size: 11px;
  color: #64748b;
}

.metric-value {
  font-size: 12px;
  font-weight: 600;
  color: #e2e8f0;
  font-family: monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.highlight-cyan {
  color: #38bdf8;
}

/* 当前单文件详情卡片 */
.current-file-box {
  background: rgba(30, 41, 59, 0.35);
  border: 1px solid rgba(148, 163, 184, 0.08);
  border-radius: 8px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.current-file-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.file-name-row {
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.file-icon {
  color: #38bdf8;
  flex-shrink: 0;
}

.file-name {
  font-size: 12px;
  font-weight: 500;
  color: #f1f5f9;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.file-percent {
  font-size: 11px;
  font-weight: 600;
  color: #94a3b8;
  font-family: monospace;
}

.mini-progress-track {
  height: 4px;
  background: rgba(15, 23, 42, 0.6);
  border-radius: 9999px;
  overflow: hidden;
}

.mini-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #60a5fa, #38bdf8);
  border-radius: 9999px;
  transition: width 0.15s ease;
}

.current-file-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 11px;
  color: #64748b;
  font-family: monospace;
}

/* 失败文件展开区 */
.failed-files-section {
  background: rgba(244, 63, 94, 0.08);
  border: 1px solid rgba(244, 63, 94, 0.25);
  border-radius: 8px;
  overflow: hidden;
}

.failed-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  cursor: pointer;
  transition: background 0.15s ease;
}

.failed-header:hover {
  background: rgba(244, 63, 94, 0.12);
}

.failed-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 500;
  color: #fb7185;
}

.failed-list {
  max-height: 120px;
  overflow-y: auto;
  padding: 0 12px 8px 12px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.failed-item {
  display: flex;
  flex-direction: column;
  font-size: 11px;
  background: rgba(15, 23, 42, 0.5);
  padding: 6px 8px;
  border-radius: 4px;
}

.failed-filename {
  font-weight: 600;
  color: #fecdd3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.failed-reason {
  color: #94a3b8;
  font-size: 10px;
  margin-top: 2px;
}

/* 底部操作行 */
.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 4px;
}

.footer-stats {
  display: flex;
  align-items: center;
  gap: 6px;
}

.stat-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
}

.stat-badge.success {
  background: rgba(16, 185, 129, 0.12);
  color: #34d399;
}

.stat-badge.danger {
  background: rgba(244, 63, 94, 0.12);
  color: #f43f5e;
}

.footer-btns {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-cancel-upload {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  color: #f43f5e;
  background: rgba(244, 63, 94, 0.12);
  border: 1px solid rgba(244, 63, 94, 0.25);
  transition: all 0.15s ease;
  cursor: pointer;
}

.btn-cancel-upload:hover:not(:disabled) {
  background: rgba(244, 63, 94, 0.22);
  border-color: rgba(244, 63, 94, 0.4);
}

.btn-cancel-upload:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-done {
  padding: 6px 16px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  color: #ffffff;
  background: #3b82f6;
  transition: all 0.15s ease;
  cursor: pointer;
}

.btn-done:hover {
  background: #2563eb;
}

/* 动效 */
.upload-icon-pulse {
  animation: float-pulse 1.8s ease-in-out infinite;
}

@keyframes float-pulse {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
}

@keyframes pulse-glow {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.4; transform: scale(0.85); }
}

.text-success { color: #10b981; }
.text-danger { color: #f43f5e; }
.text-muted { color: #64748b; }

.slide-card-enter-active,
.slide-card-leave-active {
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.slide-card-enter-from,
.slide-card-leave-to {
  opacity: 0;
  transform: translateY(16px) scale(0.96);
}
</style>
