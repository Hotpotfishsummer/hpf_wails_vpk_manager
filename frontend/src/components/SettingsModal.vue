<template>
  <transition name="fade">
    <div v-if="show" class="modal-overlay" @click.self="emit('update:show', false)">
      <div class="modal-content settings-modal">
        <div class="modal-header">
          <div class="title-with-icon">
            <SvgIcon name="settings" :size="18" class="header-icon" />
            <h2 class="modal-title">{{ $t('settingsTitle') }}</h2>
          </div>
          <button class="modal-close-btn" @click="emit('update:show', false)">
            <SvgIcon name="x" :size="16" />
          </button>
        </div>

        <div class="modal-body">
          <!-- 常规设置 -->
          <div class="settings-section">
            <h3 class="section-subtitle">{{ $t('generalSettings') }}</h3>
            <!-- 操作栏停靠位置设置 -->
            <div class="settings-item action-item">
              <div class="setting-text">
                <span class="setting-title">{{ $t('actionBarPosition') }}</span>
                <p class="setting-item-desc">{{ $t('actionBarPositionDesc') }}</p>
              </div>
              <div class="pos-choice-group">
                <button
                  type="button"
                  class="btn-pos-choice"
                  :class="{ active: store.actionBarPosition === 'bottom' }"
                  @click="store.setActionBarPosition('bottom')"
                >
                  <SvgIcon name="chevron-down" :size="12" />
                  <span>{{ $t('dockPositionBottom') }}</span>
                </button>
                <button
                  type="button"
                  class="btn-pos-choice"
                  :class="{ active: store.actionBarPosition === 'top' }"
                  @click="store.setActionBarPosition('top')"
                >
                  <SvgIcon name="chevron-up" :size="12" />
                  <span>{{ $t('dockPositionTop') }}</span>
                </button>
              </div>
            </div>

            <div class="settings-item action-item">
              <div class="setting-text">
                <span class="setting-title">{{ $t('resetAppData') }}</span>
                <p class="setting-item-desc">{{ $t('resetAppDataDesc') }}</p>
              </div>
              <button
                class="btn-reset-data"
                @click="handleResetData"
                :disabled="isResetting"
              >
                {{ isResetting ? $t('resetting') : $t('executeNow') }}
              </button>
            </div>
          </div>

          <!-- SFTP 设置 -->
          <div class="settings-section">
            <div class="section-header-row">
              <h3 class="section-subtitle">{{ $t('sftpSettings') }}</h3>
              <p class="section-desc">{{ $t('sftpHelp') }}</p>
            </div>

            <div class="form-grid">
              <div class="settings-field">
                <label class="setting-label">{{ $t('sftpAddress') }}</label>
                <input
                  type="text"
                  v-model="store.sftpSettings.address"
                  @change="store.saveSftpSettings()"
                  class="setting-input"
                  placeholder="example.com:22"
                />
              </div>

              <div class="settings-field">
                <label class="setting-label">{{ $t('sftpUsername') }}</label>
                <input
                  type="text"
                  v-model="store.sftpSettings.username"
                  @change="store.saveSftpSettings()"
                  class="setting-input"
                  placeholder="root"
                />
              </div>

              <div class="settings-field">
                <label class="setting-label">{{ $t('sftpPassword') }}</label>
                <input
                  type="password"
                  v-model="store.sftpSettings.password"
                  @change="store.saveSftpSettings()"
                  class="setting-input"
                  placeholder="••••••••"
                />
              </div>

              <div class="settings-field full-width">
                <label class="setting-label">{{ $t('sftpAddonsPath') }}</label>
                <input
                  type="text"
                  v-model="store.sftpSettings.addons_path"
                  @change="store.saveSftpSettings()"
                  class="setting-input"
                  placeholder="/home/user/l4d2/left4dead2/addons"
                />
              </div>
            </div>

            <!-- SFTP 连接测试 -->
            <div class="connection-test-row">
              <button
                class="btn-test-conn"
                :disabled="isTestingConnection"
                @click="handleTestConnection"
              >
                <SvgIcon v-if="isTestingConnection" name="loader" :size="14" :spin="true" />
                <SvgIcon v-else name="zap" :size="14" />
                <span>{{ isTestingConnection ? $t('sftpTesting') : $t('sftpTestConnection') }}</span>
              </button>

              <span
                v-if="connectionTestResult"
                class="test-result-badge"
                :class="connectionTestResult.ok ? 'test-ok' : 'test-fail'"
              >
                {{ connectionTestResult.message }}
              </span>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn-close-modal" @click="emit('update:show', false)">
            {{ $t('close') }}
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import { getApi } from '@/composables/useApi'
import SvgIcon from '@/components/common/SvgIcon.vue'

defineProps<{ show: boolean }>()
const emit = defineEmits<{ (e: 'update:show', value: boolean): void }>()

const store = useVpkManagerStore()
const { t: $t } = useI18n()

const isResetting = ref(false)
const isTestingConnection = ref(false)
const connectionTestResult = ref<{ ok: boolean; message: string } | null>(null)

async function handleResetData() {
  const confirmed = await store.showConfirm($t('confirmResetAppData'))
  if (!confirmed) return

  isResetting.value = true
  try {
    await getApi().reset_app_data()
    store.showNotification($t('resetAppDataSuccess'), 'success')
  } catch (error) {
    store.showNotification(
      error instanceof Error ? error.message : $t('resetAppDataFailed'),
      'error'
    )
  } finally {
    isResetting.value = false
  }
}

async function handleTestConnection() {
  isTestingConnection.value = true
  connectionTestResult.value = null
  try {
    const res = await getApi().list_remote_vpks()
    if (res.success) {
      connectionTestResult.value = {
        ok: true,
        message: $t('sftpConnectSuccess', { count: res.files ? res.files.length : 0 }),
      }
    } else {
      connectionTestResult.value = {
        ok: false,
        message: res.error || $t('sftpConnectFailed'),
      }
    }
  } catch (e: any) {
    connectionTestResult.value = {
      ok: false,
      message: e?.message || $t('sftpConnectFailed'),
    }
  } finally {
    isTestingConnection.value = false
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(3, 7, 18, 0.78);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.modal-content {
  background: #151e30 !important;
  border-radius: var(--radius-lg);
  border: 1.5px solid #384c6c;
  width: 100%;
  max-width: 560px;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 70px rgba(0, 0, 0, 0.9), 0 0 0 1px rgba(56, 189, 248, 0.25), 0 0 40px rgba(56, 189, 248, 0.12);
  animation: modal-scale 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  overflow: hidden;
}

@keyframes modal-scale {
  from {
    transform: scale(0.94);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}

.modal-header {
  padding: 16px 22px;
  background: #101726;
  border-bottom: 1px solid #28374d;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.title-with-icon {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon {
  color: var(--accent-primary);
}

.modal-title {
  margin: 0;
  font-size: 16px;
  font-weight: 700;
  color: #fff;
}

.modal-close-btn {
  appearance: none;
  -webkit-appearance: none;
  width: 30px;
  height: 30px;
  border-radius: var(--radius-sm);
  background: #1e293b;
  border: 1px solid #334460;
  color: #cbd5e1;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.modal-close-btn:hover {
  background: #ef4444;
  border-color: #ef4444;
  color: #ffffff;
}

.modal-body {
  padding: 20px 22px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-section {
  background: #101826;
  border: 1px solid #253348;
  border-radius: 10px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.section-header-row {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.section-subtitle {
  font-size: 13px;
  font-weight: 700;
  color: #38bdf8;
  letter-spacing: 0.02em;
  text-transform: uppercase;
}

.section-desc {
  font-size: 12px;
  color: var(--text-muted);
}

.action-item {
  background: #162032;
  border: 1px solid #2b3b52;
  border-radius: var(--radius-md);
  padding: 14px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.setting-text {
  flex: 1;
}

.setting-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-main);
  display: block;
}

.setting-item-desc {
  font-size: 12px;
  color: var(--text-dim);
  margin-top: 3px;
}

.btn-reset-data {
  appearance: none;
  -webkit-appearance: none;
  height: 32px;
  padding: 0 16px;
  border-radius: var(--radius-sm);
  background: #991b1b;
  border: 1px solid #dc2626;
  color: #ffffff !important;
  font-weight: 700;
  font-size: 12px;
  flex-shrink: 0;
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(220, 38, 38, 0.3);
  transition: all var(--transition-fast);
}

.btn-reset-data:hover:not(:disabled) {
  background: #b91c1c;
  border-color: #ef4444;
}

.pos-choice-group {
  display: flex;
  background: #0f1624;
  border: 1px solid #28374d;
  border-radius: var(--radius-sm);
  padding: 3px;
  gap: 3px;
  flex-shrink: 0;
}

.btn-pos-choice {
  appearance: none;
  -webkit-appearance: none;
  height: 28px;
  padding: 0 12px;
  border-radius: 4px;
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all var(--transition-fast);
}

.btn-pos-choice.active {
  background: linear-gradient(135deg, #0284c7 0%, #0369a1 100%);
  color: #fff;
  box-shadow: 0 2px 6px rgba(2, 132, 199, 0.4);
}

.btn-pos-choice:hover:not(.active) {
  color: var(--text-main);
  background: rgba(255, 255, 255, 0.06);
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.settings-field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.settings-field.full-width {
  grid-column: span 2;
}

.setting-label {
  font-size: 12px;
  font-weight: 600;
  color: #cbd5e1;
}

.setting-input {
  background: #0b111e !important;
  border: 1px solid #2d3e58 !important;
  border-radius: var(--radius-sm);
  padding: 8px 12px;
  font-size: 13px;
  color: #ffffff !important;
  width: 100%;
  outline: none;
  transition: all var(--transition-fast);
}

.setting-input:focus {
  border-color: #38bdf8 !important;
  box-shadow: 0 0 0 2px rgba(56, 189, 248, 0.35) !important;
}

.connection-test-row {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 2px;
}

.btn-test-conn {
  appearance: none;
  -webkit-appearance: none;
  height: 32px;
  padding: 0 16px;
  border-radius: var(--radius-sm);
  background: #0369a1;
  border: 1px solid #0284c7;
  color: #ffffff !important;
  font-size: 12px;
  font-weight: 700;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  box-shadow: 0 2px 6px rgba(2, 132, 199, 0.3);
  transition: all var(--transition-fast);
}

.btn-test-conn:hover:not(:disabled) {
  background: #0284c7;
  border-color: #38bdf8;
}

.test-result-badge {
  font-size: 12px;
  padding: 4px 10px;
  border-radius: var(--radius-sm);
}

.test-result-badge.test-ok {
  background: rgba(16, 185, 129, 0.2);
  color: #34d399;
  border: 1px solid rgba(16, 185, 129, 0.4);
}

.test-result-badge.test-fail {
  background: rgba(244, 63, 94, 0.2);
  color: #fb7185;
  border: 1px solid rgba(244, 63, 94, 0.4);
}

.modal-footer {
  padding: 14px 22px;
  background: #101726;
  border-top: 1px solid #28374d;
  display: flex;
  justify-content: flex-end;
}

.btn-close-modal {
  appearance: none;
  -webkit-appearance: none;
  height: 34px;
  padding: 0 24px;
  border-radius: var(--radius-sm);
  background: #1e293b;
  border: 1px solid #3b4d6b;
  color: #ffffff !important;
  font-weight: 600;
  font-size: 13px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-close-modal:hover {
  background: #2b3952;
  border-color: #4f6488;
}
</style>
