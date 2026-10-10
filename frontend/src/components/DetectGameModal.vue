<template>
  <transition name="fade">
    <div v-if="store.showDetectModal" class="modal-overlay" @click.self="handleClose">
      <div class="modal-content detect-modal">
        <div class="modal-header">
          <div class="title-with-icon">
            <SvgIcon name="compass" :size="20" class="header-icon" />
            <div>
              <h2 class="modal-title">{{ $t('detectModalTitle') }}</h2>
              <p class="modal-subtitle">{{ $t('detectModalSubtitle') }}</p>
            </div>
          </div>
          <button class="modal-close-btn" @click="handleClose" :title="$t('close')">
            <SvgIcon name="x" :size="16" />
          </button>
        </div>

        <div class="modal-body">
          <!-- 正在扫描动画状态 -->
          <div v-if="store.isDetecting || store.isDeepDetecting" class="scanning-state-box">
            <div class="radar-scan-wrapper">
              <div class="radar-pulse"></div>
              <div class="radar-pulse-2"></div>
              <div class="radar-core">
                <SvgIcon name="compass" :size="32" :spin="true" class="radar-icon" />
              </div>
            </div>
            <h3 class="scanning-title">
              {{ store.isDeepDetecting ? $t('deepScanning') : $t('detecting') }}
            </h3>
            <p class="scanning-desc">{{ $t('detectingDesc') }}</p>
          </div>

          <!-- 扫描结果列表 -->
          <div v-else-if="store.detectedClients.length > 0" class="clients-list-container">
            <div class="list-stats-bar">
              <span class="stats-count">
                {{ $t('clientFoundCount', { count: store.detectedClients.length }) }}
              </span>
              <span class="stats-hint">点击卡片选中要管理的客户端</span>
            </div>

            <div class="clients-grid">
              <div
                v-for="client in store.detectedClients"
                :key="client.path"
                class="client-card"
                :class="{ active: selectedClientPath === client.path }"
                @click="selectedClientPath = client.path"
              >
                <div class="client-card-header">
                  <div class="radio-indicator">
                    <span class="radio-inner" v-if="selectedClientPath === client.path"></span>
                  </div>

                  <span class="source-badge" :class="getSourceBadgeClass(client.source_name)">
                    <SvgIcon :name="getSourceIcon(client.source_name)" :size="12" />
                    <span>{{ getSourceLabel(client) }}</span>
                  </span>

                  <span class="folder-name">{{ getFolderName(client.path) }}</span>
                </div>

                <div class="client-path-row">
                  <span class="client-path" :title="client.path">{{ client.path }}</span>
                  <button
                    class="mini-copy-btn"
                    @click.stop="handleCopy(client.path)"
                    :title="$t('copyPath')"
                  >
                    <SvgIcon :name="copiedPath === client.path ? 'check' : 'copy'" :size="13" />
                  </button>
                </div>

                <div class="client-card-meta">
                  <span class="meta-tag mod-tag">
                    <SvgIcon name="box" :size="12" />
                    <span>{{ $t('modCountText', { count: client.mod_count }) }}</span>
                  </span>

                  <span class="meta-tag workshop-tag" :class="{ disabled: !client.has_workshop }">
                    <SvgIcon name="globe" :size="12" />
                    <span>{{ client.has_workshop ? $t('hasWorkshop') : $t('noWorkshop') }}</span>
                  </span>

                  <span v-if="client.last_modified" class="meta-tag date-tag">
                    <span>{{ formatDate(client.last_modified) }}</span>
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 未扫描到任何客户端 -->
          <div v-else class="empty-state-box">
            <div class="empty-icon-wrap">
              <SvgIcon name="folder" :size="40" class="empty-icon" />
            </div>
            <h3 class="empty-title">{{ $t('detectNoClients') }}</h3>
            <p class="empty-desc">{{ $t('detectNoClientsTip') }}</p>
          </div>
        </div>

        <div class="modal-footer">
          <div class="footer-left-actions">
            <!-- 全盘深度扫描按钮 -->
            <button
              type="button"
              class="btn-deep-scan"
              :disabled="store.isDetecting || store.isDeepDetecting"
              @click="handleDeepScan"
              :title="$t('deepScanTip')"
            >
              <SvgIcon name="hard-drive" :size="14" />
              <span>{{ store.isDeepDetecting ? $t('deepScanning') : $t('deepScan') }}</span>
            </button>

            <!-- 手动浏览选择按钮 -->
            <button
              type="button"
              class="btn-browse-fallback"
              @click="handleManualBrowse"
            >
              <SvgIcon name="folder-open" :size="14" />
              <span>{{ $t('manualBrowse') }}</span>
            </button>
          </div>

          <div class="footer-right-actions">
            <button type="button" class="btn-cancel" @click="handleClose">
              {{ $t('cancel') }}
            </button>

            <button
              type="button"
              class="btn-confirm-select"
              :disabled="!selectedClientPath || store.isDetecting || store.isDeepDetecting"
              @click="handleConfirmSelect"
            >
              <SvgIcon name="check" :size="15" />
              <span>{{ $t('confirmSelectClient') }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import { getApi } from '@/composables/useApi'
import SvgIcon from '@/components/common/SvgIcon.vue'
import type { DetectedGameClient } from '@/types/api'

const store = useVpkManagerStore()
const { t: $t } = useI18n()

const selectedClientPath = ref<string>('')
const copiedPath = ref<string>('')

// 默认选中第一个候选客户端（通常为模组数最多或官方版本）
watch(
  () => store.detectedClients,
  (clients) => {
    if (clients && clients.length > 0) {
      if (!selectedClientPath.value || !clients.some((c) => c.path === selectedClientPath.value)) {
        selectedClientPath.value = clients[0].path
      }
    } else {
      selectedClientPath.value = ''
    }
  },
  { immediate: true }
)

function handleClose() {
  store.closeDetectModal()
}

async function handleConfirmSelect() {
  if (!selectedClientPath.value) return
  await store.selectDetectedClient(selectedClientPath.value)
}

async function handleDeepScan() {
  await store.detectGameDirectories(true)
}

async function handleManualBrowse() {
  try {
    const selected = await getApi().browse_directory()
    if (selected) {
      store.closeDetectModal()
      await store.scanVpks(selected)
    }
  } catch (e) {
    console.error('Manual browse failed:', e)
  }
}

async function handleCopy(path: string) {
  try {
    await navigator.clipboard.writeText(path)
    copiedPath.value = path
    setTimeout(() => {
      copiedPath.value = ''
    }, 2000)
  } catch (e) {
    console.error('Copy failed:', e)
  }
}

function getFolderName(fullPath: string): string {
  if (!fullPath) return ''
  const parts = fullPath.replace(/\\/g, '/').split('/')
  return parts.filter(Boolean).pop() || fullPath
}

function getSourceLabel(client: DetectedGameClient): string {
  if (client.source_name === 'running') return $t('clientSourceRunning')
  if (client.source_name === 'steam') return $t('clientSourceSteam')
  if (client.source_name === 'registry') return $t('clientSourceRegistry')
  if (client.source_name === 'shortcut') return $t('clientSourceShortcut')
  if (client.source_name === 'common') return $t('clientSourceCommon')
  if (client.source_name === 'disk') return $t('clientSourceDisk')
  return client.source || $t('clientSourceCommon')
}

function getSourceIcon(sourceName: string): string {
  switch (sourceName) {
    case 'running':
      return 'zap'
    case 'steam':
      return 'globe'
    case 'registry':
      return 'settings'
    case 'shortcut':
      return 'external-link'
    case 'disk':
      return 'hard-drive'
    default:
      return 'folder'
  }
}

function getSourceBadgeClass(sourceName: string): string {
  switch (sourceName) {
    case 'running':
      return 'badge-running'
    case 'steam':
      return 'badge-steam'
    case 'registry':
      return 'badge-registry'
    case 'shortcut':
      return 'badge-shortcut'
    case 'disk':
      return 'badge-disk'
    default:
      return 'badge-common'
  }
}

function formatDate(timestamp: number): string {
  if (!timestamp) return ''
  const d = new Date(timestamp * 1000)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.75);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.detect-modal {
  width: 100%;
  max-width: 680px;
  max-height: 85vh;
  background: var(--bg-card, #131d2e);
  border: 1px solid var(--border-medium, #2a3a52);
  border-radius: var(--radius-lg, 12px);
  display: flex;
  flex-direction: column;
  box-shadow: 0 16px 40px rgba(0, 0, 0, 0.6);
  animation: modal-pop 0.25s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

@keyframes modal-pop {
  from {
    opacity: 0;
    transform: scale(0.96) translateY(8px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-subtle, #1e2c40);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(0, 0, 0, 0.2);
}

.title-with-icon {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon {
  color: var(--accent-primary, #38bdf8);
}

.modal-title {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.3px;
}

.modal-subtitle {
  margin: 3px 0 0;
  font-size: 12px;
  color: var(--text-dim, #94a3b8);
}

.modal-close-btn {
  background: transparent;
  border: none;
  color: var(--text-dim, #94a3b8);
  width: 32px;
  height: 32px;
  border-radius: var(--radius-sm, 6px);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast, 0.15s ease);
}

.modal-close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.modal-body {
  padding: 18px 20px;
  overflow-y: auto;
  flex: 1;
  min-height: 240px;
  display: flex;
  flex-direction: column;
}

/* 扫描中雷达动画 */
.scanning-state-box {
  margin: auto;
  padding: 30px 20px;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.radar-scan-wrapper {
  position: relative;
  width: 90px;
  height: 90px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}

.radar-pulse,
.radar-pulse-2 {
  position: absolute;
  inset: 0;
  border-radius: 50%;
  border: 2px solid rgba(56, 189, 248, 0.5);
  animation: pulse-ring 2s cubic-bezier(0.215, 0.61, 0.355, 1) infinite;
}

.radar-pulse-2 {
  animation-delay: 1s;
}

@keyframes pulse-ring {
  0% {
    transform: scale(0.6);
    opacity: 0.8;
  }
  100% {
    transform: scale(1.6);
    opacity: 0;
  }
}

.radar-core {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, #0284c7 0%, #0369a1 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 20px rgba(56, 189, 248, 0.4);
  z-index: 2;
}

.radar-icon {
  color: #fff;
}

.scanning-title {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
}

.scanning-desc {
  margin: 8px 0 0;
  font-size: 13px;
  color: var(--text-dim, #94a3b8);
  max-width: 380px;
}

/* 列表展示 */
.clients-list-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.list-stats-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 12px;
}

.stats-count {
  font-weight: 700;
  color: var(--accent-primary, #38bdf8);
}

.stats-hint {
  color: var(--text-dim, #94a3b8);
}

.clients-grid {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.client-card {
  background: rgba(15, 23, 42, 0.6);
  border: 1px solid var(--border-subtle, #1e2c40);
  border-radius: var(--radius-md, 8px);
  padding: 12px 14px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  gap: 8px;
  transition: all var(--transition-fast, 0.15s ease);
}

.client-card:hover {
  background: rgba(30, 41, 59, 0.7);
  border-color: rgba(56, 189, 248, 0.4);
  transform: translateY(-1px);
}

.client-card.active {
  background: rgba(14, 165, 233, 0.1);
  border-color: var(--accent-primary, #38bdf8);
  box-shadow: 0 0 14px rgba(56, 189, 248, 0.2);
}

.client-card-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.radio-indicator {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  border: 2px solid var(--border-medium, #334155);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all 0.15s ease;
}

.client-card.active .radio-indicator {
  border-color: var(--accent-primary, #38bdf8);
}

.radio-inner {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--accent-primary, #38bdf8);
  box-shadow: 0 0 6px var(--accent-primary, #38bdf8);
}

.folder-name {
  font-weight: 700;
  font-size: 14px;
  color: #fff;
}

.source-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.2px;
}

.badge-running {
  background: rgba(34, 197, 94, 0.16);
  color: #4ade80;
  border: 1px solid rgba(34, 197, 94, 0.4);
}

.badge-steam {
  background: rgba(14, 165, 233, 0.16);
  color: #38bdf8;
  border: 1px solid rgba(14, 165, 233, 0.4);
}

.badge-registry {
  background: rgba(168, 85, 247, 0.16);
  color: #c084fc;
  border: 1px solid rgba(168, 85, 247, 0.4);
}

.badge-shortcut {
  background: rgba(249, 115, 22, 0.16);
  color: #fb923c;
  border: 1px solid rgba(249, 115, 22, 0.4);
}

.badge-disk {
  background: rgba(234, 179, 8, 0.16);
  color: #facc15;
  border: 1px solid rgba(234, 179, 8, 0.4);
}

.badge-common {
  background: rgba(148, 163, 184, 0.16);
  color: #cbd5e1;
  border: 1px solid rgba(148, 163, 184, 0.3);
}

.client-path-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(0, 0, 0, 0.3);
  padding: 5px 10px;
  border-radius: 4px;
  gap: 8px;
}

.client-path {
  font-size: 12px;
  color: var(--text-dim, #94a3b8);
  font-family: Consolas, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.mini-copy-btn {
  background: transparent;
  border: none;
  color: var(--text-dim, #94a3b8);
  cursor: pointer;
  padding: 3px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: all var(--transition-fast, 0.15s ease);
}

.mini-copy-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.client-card-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.meta-tag {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: var(--text-dim, #94a3b8);
}

.mod-tag {
  color: #38bdf8;
  font-weight: 600;
}

.workshop-tag.disabled {
  color: #64748b;
}

.date-tag {
  color: #64748b;
  margin-left: auto;
}

/* 空状态 */
.empty-state-box {
  margin: auto;
  text-align: center;
  padding: 30px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.empty-icon-wrap {
  width: 70px;
  height: 70px;
  border-radius: 50%;
  background: rgba(239, 68, 68, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
}

.empty-icon {
  color: #f87171;
}

.empty-title {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
}

.empty-desc {
  margin: 8px 0 0;
  font-size: 13px;
  color: var(--text-dim, #94a3b8);
  max-width: 440px;
  line-height: 1.5;
}

/* 弹窗底部操作区 */
.modal-footer {
  padding: 14px 20px;
  border-top: 1px solid var(--border-subtle, #1e2c40);
  background: rgba(0, 0, 0, 0.25);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.footer-left-actions,
.footer-right-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-deep-scan {
  height: 36px;
  padding: 0 14px;
  border-radius: var(--radius-sm, 6px);
  background: rgba(234, 179, 8, 0.1);
  color: #facc15;
  border: 1px solid rgba(234, 179, 8, 0.3);
  font-size: 13px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all var(--transition-fast, 0.15s ease);
}

.btn-deep-scan:hover:not(:disabled) {
  background: rgba(234, 179, 8, 0.2);
  border-color: rgba(234, 179, 8, 0.5);
}

.btn-browse-fallback {
  height: 36px;
  padding: 0 14px;
  border-radius: var(--radius-sm, 6px);
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-main, #f1f5f9);
  border: 1px solid var(--border-medium, #334155);
  font-size: 13px;
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all var(--transition-fast, 0.15s ease);
}

.btn-browse-fallback:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: var(--text-dim, #94a3b8);
}

.btn-cancel {
  height: 36px;
  padding: 0 14px;
  border-radius: var(--radius-sm, 6px);
  background: transparent;
  color: var(--text-dim, #94a3b8);
  border: 1px solid transparent;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast, 0.15s ease);
}

.btn-cancel:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #fff;
}

.btn-confirm-select {
  height: 36px;
  padding: 0 18px;
  border-radius: var(--radius-sm, 6px);
  background: linear-gradient(135deg, #0284c7 0%, #0369a1 100%);
  color: #fff;
  border: 1px solid rgba(56, 189, 248, 0.4);
  font-size: 13px;
  font-weight: 700;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  box-shadow: 0 2px 10px rgba(2, 132, 199, 0.4);
  transition: all var(--transition-fast, 0.15s ease);
}

.btn-confirm-select:hover:not(:disabled) {
  filter: brightness(1.15);
  box-shadow: 0 0 14px rgba(56, 189, 248, 0.6);
}

.btn-confirm-select:disabled,
.btn-deep-scan:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  filter: grayscale(0.5);
}

/* 淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
