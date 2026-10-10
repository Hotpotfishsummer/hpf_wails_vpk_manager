<template>
  <div id="app" class="app-container">
    <!-- 顶部应用栏 -->
    <AppHeader :collapsed="isScrolled" @open-settings="showSettings = true" />

    <!-- 全局设置弹窗 -->
    <SettingsModal v-model:show="showSettings" />

    <!-- 确认操作对话框 -->
    <ConfirmDialog />

    <!-- 全局状态通知 -->
    <transition name="slide-fade">
      <StatusBanner v-if="store.notification.message" />
    </transition>

    <!-- 顶部贴合固定操作栏 (当有选定文件且设置为顶部停靠时显示) -->
    <transition name="dock-top-slide">
      <ActionBar
        v-if="hasSelection && store.actionBarPosition === 'top'"
        position="top"
      />
    </transition>

    <!-- 主工作区 (支持 SFTP 左右拖拽分屏) -->
    <main
      ref="mainContainerRef"
      class="app-main"
      :class="{ 'split-view': store.isSftpActive }"
      @scroll="handleScroll"
    >
      <!-- 左侧：本地与工坊模组管理面板 -->
      <div
        class="left-panel-wrapper"
        :style="store.isSftpActive ? { width: `${leftPanelPercent}%` } : {}"
      >
        <div class="panel-inner">
          <!-- 目录选择与状态卡片 -->
          <section class="directory-section">
            <DirectorySelector />
          </section>

          <!-- 扫描加载中状态 -->
          <div v-if="store.isScanning" class="loading-state">
            <SvgIcon name="loader" :size="36" :spin="true" class="main-spinner" />
            <p class="loading-text">{{ $t('loading') }}...</p>
          </div>

          <!-- 模组列表 -->
          <template v-else-if="hasAnyFiles">
            <section class="files-section">
              <VpkFilesList />
            </section>
          </template>

          <!-- 空状态 (未找到或未选择目录) -->
          <section v-else class="empty-state">
            <div class="empty-card">
              <div class="empty-icon-wrap">
                <SvgIcon name="folder" :size="48" />
              </div>
              <h2 class="empty-heading">{{ $t("noVpkFilesHeading") }}</h2>
              <p class="empty-desc">{{ $t("selectDirWithVpks") }}</p>
              <button class="empty-browse-btn" @click="handleEmptyBrowse">
                <SvgIcon name="folder-open" :size="16" />
                <span>{{ $t('browseFolder') }}</span>
              </button>
            </div>
          </section>

          <!-- 页脚 -->
          <footer class="app-footer">
            <p class="footer-text">© 2026 L4D2 VPK Manager • Powered by Go + Wails v2 + Vue 3</p>
          </footer>
        </div>
      </div>

      <!-- 可拖拽分屏分割线 -->
      <div
        v-if="store.isSftpActive"
        class="split-resizer"
        @mousedown="startResizing"
        title="拖动调整左右面板比例"
      >
        <div class="resizer-handle"></div>
      </div>

      <!-- 右侧：SFTP 远程模组管理面板 -->
      <div
        v-if="store.isSftpActive"
        class="right-panel-wrapper"
        :style="{ width: `${100 - leftPanelPercent}%` }"
      >
        <SftpBrowser />
      </div>
    </main>

    <!-- 底部贴合固定操作栏 (默认底部停靠) -->
    <transition name="dock-bottom-slide">
      <ActionBar
        v-if="hasSelection && store.actionBarPosition !== 'top'"
        position="bottom"
      />
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import { usePyWebViewApi, getApi } from '@/composables/useApi'
import AppHeader from '@/components/AppHeader.vue'
import StatusBanner from '@/components/StatusBanner.vue'
import ActionBar from '@/components/ActionBar.vue'
import DirectorySelector from '@/components/DirectorySelector.vue'
import VpkFilesList from '@/components/VpkFilesList.vue'
import SettingsModal from '@/components/SettingsModal.vue'
import ConfirmDialog from '@/components/ConfirmDialog.vue'
import SftpBrowser from '@/components/SftpBrowser.vue'
import SvgIcon from '@/components/common/SvgIcon.vue'

const store = useVpkManagerStore()
const { t: $t } = useI18n()
const { waitForApiReady } = usePyWebViewApi()

const isScrolled = ref(false)
const showSettings = ref(false)
const mainContainerRef = ref<HTMLElement | null>(null)

// 分屏比例 (左侧占比百分比)
const leftPanelPercent = ref(58)
const isResizing = ref(false)

try {
  const savedRatio = localStorage.getItem('vpk-split-ratio')
  if (savedRatio) {
    const parsed = parseFloat(savedRatio)
    if (parsed >= 30 && parsed <= 75) {
      leftPanelPercent.value = parsed
    }
  }
} catch {}

function startResizing(e: MouseEvent) {
  e.preventDefault()
  isResizing.value = true
  document.body.style.cursor = 'col-resize'
  document.body.style.userSelect = 'none'

  const onMouseMove = (moveEvent: MouseEvent) => {
    if (!isResizing.value || !mainContainerRef.value) return
    const containerRect = mainContainerRef.value.getBoundingClientRect()
    const relativeX = moveEvent.clientX - containerRect.left
    const newPercent = (relativeX / containerRect.width) * 100
    if (newPercent >= 30 && newPercent <= 75) {
      leftPanelPercent.value = Math.round(newPercent)
    }
  }

  const onMouseUp = () => {
    isResizing.value = false
    document.body.style.cursor = ''
    document.body.style.userSelect = ''
    window.removeEventListener('mousemove', onMouseMove)
    window.removeEventListener('mouseup', onMouseUp)
    try {
      localStorage.setItem('vpk-split-ratio', leftPanelPercent.value.toString())
    } catch {}
  }

  window.addEventListener('mousemove', onMouseMove)
  window.addEventListener('mouseup', onMouseUp)
}

function handleScroll(event: Event) {
  const target = event.target as HTMLElement
  isScrolled.value = target.scrollTop > 20
}

async function handleEmptyBrowse() {
  try {
    const selected = await getApi().browse_directory()
    if (selected) {
      await store.scanVpks(selected)
    }
  } catch (e) {
    console.error('Failed to browse directory', e)
  }
}

onMounted(async () => {
  const ready = await waitForApiReady(5000)
  if (ready) {
    await store.initialize()
  }
})

const hasSelection = computed(() => {
  return store.selectedLocalVpks.size > 0 || store.selectedWorkshopVpks.size > 0
})

const hasAnyFiles = computed(() => {
  return store.localVpkFiles.length > 0 || store.workshopVpkFiles.length > 0
})
</script>

<style scoped>
.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100vw;
  background: var(--bg-canvas);
  color: var(--text-main);
  overflow: hidden;
  position: relative;
}

/* 主区域 */
.app-main {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  position: relative;
}

.app-main.split-view {
  overflow: hidden;
}

.left-panel-wrapper {
  flex: 1;
  overflow-y: auto;
  min-width: 320px;
  display: flex;
  flex-direction: column;
}

.panel-inner {
  max-width: 1360px;
  margin: 0 auto;
  padding: 20px 28px;
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 18px;
  min-height: 100%;
}

/* 分屏分割线 */
.split-resizer {
  width: 8px;
  background: rgba(255, 255, 255, 0.02);
  border-left: 1px solid var(--border-subtle);
  border-right: 1px solid var(--border-subtle);
  cursor: col-resize;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background var(--transition-fast);
  z-index: 10;
  flex-shrink: 0;
}

.split-resizer:hover {
  background: rgba(56, 189, 248, 0.2);
  border-color: rgba(56, 189, 248, 0.4);
}

.resizer-handle {
  width: 2px;
  height: 24px;
  background: var(--text-dim);
  border-radius: 2px;
  transition: background var(--transition-fast);
}

.split-resizer:hover .resizer-handle {
  background: var(--accent-primary);
}

.right-panel-wrapper {
  min-width: 280px;
  height: 100%;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 0;
  gap: 16px;
}

.main-spinner {
  color: var(--accent-primary);
}

.loading-text {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
}

/* 空状态卡片 */
.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
}

.empty-card {
  max-width: 440px;
  padding: 36px 28px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: var(--shadow-md);
}

.empty-icon-wrap {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  background: rgba(56, 189, 248, 0.08);
  border: 1px solid rgba(56, 189, 248, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-primary);
  margin-bottom: 20px;
}

.empty-heading {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 8px;
}

.empty-desc {
  font-size: 13px;
  color: var(--text-dim);
  margin-bottom: 22px;
  line-height: 1.5;
}

.empty-browse-btn {
  height: 38px;
  padding: 0 20px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, #0284c7 0%, #0369a1 100%);
  color: white;
  font-size: 13px;
  font-weight: 600;
  box-shadow: 0 4px 12px rgba(2, 132, 199, 0.4);
  border: 1px solid rgba(56, 189, 248, 0.4);
}

.empty-browse-btn:hover {
  background: linear-gradient(135deg, #0369a1 0%, #075985 100%);
  box-shadow: 0 6px 18px rgba(2, 132, 199, 0.6);
}

/* 页脚 */
.app-footer {
  margin-top: auto;
  padding: 20px 0 10px;
  text-align: center;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.footer-text {
  font-size: 12px;
  color: var(--text-dim);
}

/* 动效 */
.dock-top-slide-enter-active,
.dock-top-slide-leave-active {
  transition: transform 0.22s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.2s ease, max-height 0.22s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

.dock-top-slide-enter-from,
.dock-top-slide-leave-to {
  opacity: 0;
  transform: translateY(-100%);
  max-height: 0;
}

.dock-top-slide-enter-to,
.dock-top-slide-leave-from {
  opacity: 1;
  transform: translateY(0);
  max-height: 80px;
}

.dock-bottom-slide-enter-active,
.dock-bottom-slide-leave-active {
  transition: transform 0.22s cubic-bezier(0.16, 1, 0.3, 1), opacity 0.2s ease, max-height 0.22s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

.dock-bottom-slide-enter-from,
.dock-bottom-slide-leave-to {
  opacity: 0;
  transform: translateY(100%);
  max-height: 0;
}

.dock-bottom-slide-enter-to,
.dock-bottom-slide-leave-from {
  opacity: 1;
  transform: translateY(0);
  max-height: 80px;
}

.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.2s ease;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
