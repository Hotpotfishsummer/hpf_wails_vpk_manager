<template>
  <header class="app-header" :class="{ 'collapsed': collapsed }">
    <div class="header-content">
      <div class="header-left">
        <div class="logo-box">
          <SvgIcon name="gamepad" :size="20" stroke-width="2.2" />
        </div>
        <div class="title-group">
          <div class="title-row">
            <h1 class="app-title">{{ $t('appTitle') }}</h1>
            <span v-if="totalModsCount > 0" class="mods-count-badge">
              {{ totalModsCount }} {{ $t('modsLoaded', 'MODS') }}
            </span>
          </div>
          <p class="app-subtitle">Left 4 Dead 2 VPK Manager</p>
        </div>
      </div>

      <div class="header-right">
        <!-- 语言切换 -->
        <div class="language-pill">
          <button
            @click="switchLanguage('en')"
            :class="{ active: currentLocale === 'en' }"
            class="lang-pill-btn"
          >
            EN
          </button>
          <button
            @click="switchLanguage('zh')"
            :class="{ active: currentLocale === 'zh' }"
            class="lang-pill-btn"
          >
            中文
          </button>
        </div>

        <!-- SFTP 远程分屏切换 -->
        <button
          @click="store.isSftpActive = !store.isSftpActive"
          class="header-icon-btn sftp-toggle-btn"
          :class="{ active: store.isSftpActive }"
          :title="store.isSftpActive ? $t('sftpDisable') : $t('sftpEnable')"
        >
          <SvgIcon name="globe" :size="16" />
          <span class="btn-label">{{ $t('remoteSync', '远程同步') }}</span>
          <span v-if="store.isSftpActive" class="status-dot online"></span>
        </button>

        <!-- 设置按钮 -->
        <button
          @click="$emit('open-settings')"
          class="header-icon-btn settings-btn"
          :title="$t('settings')"
        >
          <SvgIcon name="settings" :size="16" />
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import SvgIcon from '@/components/common/SvgIcon.vue'

defineProps<{ collapsed: boolean }>()
defineEmits<{ (e: 'open-settings'): void }>()

const store = useVpkManagerStore()
const { locale, t: $t } = useI18n()

const currentLocale = computed(() => locale.value)
const totalModsCount = computed(() => store.localVpkFiles.length + store.workshopVpkFiles.length)

const STORAGE_KEY = 'locale'

function saveLocale(lang: string) {
  try {
    localStorage.setItem(STORAGE_KEY, lang)
  } catch {
    // 忽略不可写 localStorage
  }
}

function switchLanguage(lang: string) {
  locale.value = lang
  saveLocale(lang)
}
</script>

<style scoped>
.app-header {
  background: rgba(11, 15, 25, 0.85);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border-subtle);
  padding: 14px 28px;
  z-index: 100;
  position: relative;
  transition: all var(--transition-normal);
}

.app-header.collapsed {
  padding: 8px 28px;
}

.header-content {
  max-width: 1360px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.logo-box {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md);
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2) 0%, rgba(139, 92, 246, 0.2) 100%);
  border: 1px solid rgba(56, 189, 248, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-primary);
  box-shadow: 0 0 16px var(--accent-primary-glow);
}

.title-group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.app-title {
  font-size: 17px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.02em;
  line-height: 1.2;
}

.mods-count-badge {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  background: rgba(56, 189, 248, 0.15);
  color: var(--accent-primary);
  border: 1px solid rgba(56, 189, 248, 0.3);
}

.app-subtitle {
  font-size: 11px;
  color: var(--text-dim);
  font-weight: 500;
  line-height: 1;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

/* 语言切换胶囊 */
.language-pill {
  display: flex;
  background: #141c2c;
  padding: 3px;
  border-radius: var(--radius-full);
  border: 1px solid var(--border-medium);
  gap: 2px;
}

.lang-pill-btn {
  appearance: none;
  -webkit-appearance: none;
  background: transparent;
  border: none;
  font-size: 11px;
  font-weight: 600;
  padding: 4px 10px;
  border-radius: var(--radius-full);
  color: var(--text-muted);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.lang-pill-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.08);
}

.lang-pill-btn.active {
  background: #0284c7 !important;
  color: #ffffff !important;
  font-weight: 700;
  box-shadow: 0 1px 6px rgba(2, 132, 199, 0.4);
}

/* 图标按钮 */
.header-icon-btn {
  appearance: none;
  -webkit-appearance: none;
  height: 32px;
  padding: 0 12px;
  border-radius: var(--radius-md);
  background: #182235;
  border: 1px solid var(--border-medium);
  color: #cbd5e1;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.header-icon-btn:hover {
  background: #223048;
  color: #fff;
  border-color: var(--border-strong);
}

.header-icon-btn.active {
  background: rgba(56, 189, 248, 0.18);
  color: var(--accent-primary);
  border-color: var(--accent-primary);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--accent-success);
  box-shadow: 0 0 8px var(--accent-success);
}

.settings-btn {
  width: 32px;
  padding: 0;
  justify-content: center;
}
</style>
