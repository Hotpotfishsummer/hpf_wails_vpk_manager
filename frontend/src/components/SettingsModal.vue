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

          <!-- 游戏启动设置 -->
          <div class="settings-section">
            <div class="section-header-row">
              <h3 class="section-subtitle">{{ $t('launchSettings') }}</h3>
              <p class="section-desc">{{ $t('launchSettingsDesc') }}</p>
            </div>

            <!-- 跳过片头动画 (-novid) -->
            <div class="settings-item action-item">
              <div class="setting-text">
                <span class="setting-title">{{ $t('launchNoVid') }}</span>
                <p class="setting-item-desc">{{ $t('launchNoVidDesc') }}</p>
              </div>
              <label class="switch-control">
                <input
                  type="checkbox"
                  v-model="store.launchSettings.novid"
                  @change="store.saveLaunchSettings()"
                />
                <span class="switch-slider"></span>
              </label>
            </div>

            <!-- 强制中文启动 (-language schinese) -->
            <div class="settings-item action-item">
              <div class="setting-text">
                <span class="setting-title">{{ $t('launchChineseLang') }}</span>
                <p class="setting-item-desc">{{ $t('launchChineseLangDesc') }}</p>
              </div>
              <label class="switch-control">
                <input
                  type="checkbox"
                  v-model="store.launchSettings.chinese_lang"
                  @change="store.saveLaunchSettings()"
                />
                <span class="switch-slider"></span>
              </label>
            </div>

            <!-- 优先本地直接运行 (不通过 Steam 协议) -->
            <div class="settings-item action-item">
              <div class="setting-text">
                <span class="setting-title">{{ $t('launchDirectMode') }}</span>
                <p class="setting-item-desc">{{ $t('launchDirectModeDesc') }}</p>
              </div>
              <label class="switch-control">
                <input
                  type="checkbox"
                  v-model="store.launchSettings.direct_launch"
                  @change="store.saveLaunchSettings()"
                />
                <span class="switch-slider"></span>
              </label>
            </div>

            <!-- 自定义额外启动参数 -->
            <div class="settings-field full-width">
              <label class="setting-label">{{ $t('launchCustomArgs') }}</label>
              <input
                type="text"
                v-model="store.launchSettings.custom_args"
                @change="store.saveLaunchSettings()"
                class="setting-input"
                :placeholder="$t('launchCustomArgsPlaceholder')"
              />
              <span class="field-hint">{{ $t('launchCustomArgsDesc') }}</span>
            </div>

            <!-- 快捷复制服务器指令设置 -->
            <div class="settings-field full-width">
              <div class="field-label-row">
                <label class="setting-label">{{ $t('serverConnectCmd') }}</label>
                <button
                  type="button"
                  class="btn-inline-copy"
                  @click="handleTestCopyCmd"
                  :title="$t('testCopyServerTooltip')"
                >
                  <SvgIcon :name="testCopied ? 'check' : 'copy'" :size="12" />
                  <span>{{ testCopied ? $t('copied') : $t('copyServer') }}</span>
                </button>
              </div>
              <input
                type="text"
                v-model="store.launchSettings.server_connect_cmd"
                @change="store.saveLaunchSettings()"
                class="setting-input"
                placeholder="connect 192.168.1.200"
              />
              <span class="field-hint">{{ $t('serverConnectCmdDesc') }}</span>
            </div>

            <!-- 修改网络参数设置 (Tickrate 100 匹配) -->
            <div class="settings-field full-width">
              <div class="field-label-row">
                <div class="label-with-presets">
                  <label class="setting-label">{{ $t('networkRatesCmd') }}</label>
                  <div class="preset-tag-group">
                    <button
                      type="button"
                      class="btn-preset-tag"
                      @click="applyNetworkPreset('100')"
                      :title="$t('preset100TickTooltip')"
                    >
                      {{ $t('preset100Tick') }}
                    </button>
                    <button
                      type="button"
                      class="btn-preset-tag"
                      @click="applyNetworkPreset('60')"
                      :title="$t('preset60TickTooltip')"
                    >
                      {{ $t('preset60Tick') }}
                    </button>
                    <button
                      type="button"
                      class="btn-preset-tag"
                      @click="applyNetworkPreset('30')"
                      :title="$t('preset30TickTooltip')"
                    >
                      {{ $t('preset30Tick') }}
                    </button>
                  </div>
                </div>
                <button
                  type="button"
                  class="btn-inline-copy btn-inline-network"
                  @click="handleTestCopyNetwork"
                  :title="$t('testCopyNetworkTooltip')"
                >
                  <SvgIcon :name="testNetworkCopied ? 'check' : 'copy'" :size="12" />
                  <span>{{ testNetworkCopied ? $t('copied') : $t('networkParams') }}</span>
                </button>
              </div>
              <textarea
                v-model="store.launchSettings.network_rates_cmd"
                @change="store.saveLaunchSettings()"
                class="setting-textarea font-mono"
                rows="4"
                :placeholder="DEFAULT_NETWORK_RATES"
              ></textarea>
              <span class="field-hint">{{ $t('networkRatesCmdDesc') }}</span>
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

          <!-- AI 大模型接入 (OpenAI 兼容) -->
          <div class="settings-section">
            <div class="section-header-row">
              <h3 class="section-subtitle">{{ $t('aiSettingsTitle') }}</h3>
              <p class="section-desc">{{ $t('aiSettingsDesc') }}</p>
            </div>

            <div class="form-grid">
              <div class="settings-field full-width">
                <label class="setting-label">{{ $t('aiBaseUrl') }}</label>
                <input
                  type="text"
                  v-model="store.aiSettings.base_url"
                  @change="store.saveAiSettings()"
                  class="setting-input"
                  placeholder="https://api.openai.com/v1"
                />
                <span class="field-hint">{{ $t('aiBaseUrlHint') }}</span>
              </div>

              <div class="settings-field full-width">
                <label class="setting-label">{{ $t('aiApiKey') }}</label>
                <input
                  type="password"
                  v-model="store.aiSettings.api_key"
                  @change="store.saveAiSettings()"
                  class="setting-input"
                  placeholder="sk-••••••••"
                />
              </div>

              <div class="settings-field">
                <label class="setting-label">{{ $t('aiModel') }}</label>
                <input
                  type="text"
                  v-model="store.aiSettings.model"
                  @change="store.saveAiSettings()"
                  class="setting-input"
                  placeholder="gpt-4o-mini"
                />
              </div>

              <div class="settings-field">
                <label class="setting-label">{{ $t('aiTemperature') }}</label>
                <input
                  type="number"
                  step="0.1"
                  min="0"
                  max="2"
                  v-model.number="store.aiSettings.temperature"
                  @change="store.saveAiSettings()"
                  class="setting-input"
                  placeholder="0.7"
                />
              </div>

              <div class="settings-field">
                <label class="setting-label">{{ $t('aiMaxTokens') }}</label>
                <input
                  type="number"
                  step="50"
                  min="50"
                  max="4096"
                  v-model.number="store.aiSettings.max_tokens"
                  @change="store.saveAiSettings()"
                  class="setting-input"
                  placeholder="500"
                />
              </div>
            </div>

            <!-- AI 连接测试 -->
            <div class="connection-test-row">
              <button
                class="btn-test-conn"
                :disabled="isTestingAi"
                @click="handleTestAiConnection"
              >
                <SvgIcon v-if="isTestingAi" name="loader" :size="14" :spin="true" />
                <SvgIcon v-else name="sparkles" :size="14" />
                <span>{{ isTestingAi ? $t('aiTesting') : $t('aiTestConnection') }}</span>
              </button>

              <span
                v-if="aiTestResult"
                class="test-result-badge"
                :class="aiTestResult.ok ? 'test-ok' : 'test-fail'"
              >
                {{ aiTestResult.message }}
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
import { useVpkManagerStore, DEFAULT_NETWORK_RATES } from '@/stores/vpk-manager'
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

const isTestingAi = ref(false)
const aiTestResult = ref<{ ok: boolean; message: string } | null>(null)

const testCopied = ref(false)
const testNetworkCopied = ref(false)

async function handleTestCopyCmd() {
  const ok = await store.copyServerConnectCmd()
  if (ok) {
    testCopied.value = true
    setTimeout(() => {
      testCopied.value = false
    }, 2000)
  }
}

async function handleTestCopyNetwork() {
  const ok = await store.copyNetworkRatesCmd()
  if (ok) {
    testNetworkCopied.value = true
    setTimeout(() => {
      testNetworkCopied.value = false
    }, 2000)
  }
}

function applyNetworkPreset(preset: '100' | '60' | '30') {
  if (preset === '100') {
    store.launchSettings.network_rates_cmd = `cl_updaterate 100; cl_cmdrate 100; cl_interp 0.01; cl_interp_ratio 1`
  } else if (preset === '60') {
    store.launchSettings.network_rates_cmd = `cl_updaterate 60; cl_cmdrate 60; cl_interp 0.016667; cl_interp_ratio 1`
  } else if (preset === '30') {
    store.launchSettings.network_rates_cmd = `cl_updaterate 30; cl_cmdrate 30; cl_interp 0.033333; cl_interp_ratio 1`
  }
  store.saveLaunchSettings()
}

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

async function handleTestAiConnection() {
  if (!store.aiSettings.api_key?.trim()) {
    aiTestResult.value = {
      ok: false,
      message: $t('aiKeyRequired'),
    }
    return
  }
  isTestingAi.value = true
  aiTestResult.value = null
  try {
    const res = await store.testAiConnection()
    if (res.success) {
      aiTestResult.value = {
        ok: true,
        message: res.message || $t('aiConnectSuccess'),
      }
    } else {
      aiTestResult.value = {
        ok: false,
        message: res.error || $t('aiConnectFailed'),
      }
    }
  } catch (e: any) {
    aiTestResult.value = {
      ok: false,
      message: e?.message || $t('aiConnectFailed'),
    }
  } finally {
    isTestingAi.value = false
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

/* 优雅的拨动开关组件 */
.switch-control {
  position: relative;
  display: inline-block;
  width: 44px;
  height: 24px;
  flex-shrink: 0;
  cursor: pointer;
}

.switch-control input {
  opacity: 0;
  width: 0;
  height: 0;
}

.switch-slider {
  position: absolute;
  inset: 0;
  background-color: #1e293b;
  border: 1.5px solid #334155;
  border-radius: 24px;
  transition: all 0.22s cubic-bezier(0.4, 0, 0.2, 1);
}

.switch-slider::before {
  position: absolute;
  content: "";
  height: 16px;
  width: 16px;
  left: 2px;
  bottom: 2px;
  background-color: #94a3b8;
  border-radius: 50%;
  transition: all 0.22s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.4);
}

.switch-control input:checked + .switch-slider {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  border-color: #34d399;
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.35);
}

.switch-control input:checked + .switch-slider::before {
  transform: translateX(20px);
  background-color: #ffffff;
}

.switch-control:hover .switch-slider {
  border-color: #475569;
}

.switch-control input:checked:hover + .switch-slider {
  border-color: #6ee7b7;
}

.field-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.field-label-row .setting-label {
  margin-bottom: 0;
}

.btn-inline-copy {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: rgba(56, 189, 248, 0.12);
  color: var(--accent-primary);
  border: 1px solid rgba(56, 189, 248, 0.35);
  padding: 3px 8px;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-inline-copy:hover {
  background: rgba(56, 189, 248, 0.22);
  border-color: rgba(56, 189, 248, 0.6);
  box-shadow: 0 0 10px var(--accent-primary-glow);
}

.btn-inline-network {
  background: rgba(245, 158, 11, 0.12);
  color: #fbbf24;
  border-color: rgba(245, 158, 11, 0.35);
}

.btn-inline-network:hover {
  background: rgba(245, 158, 11, 0.22);
  border-color: rgba(245, 158, 11, 0.6);
  box-shadow: 0 0 10px rgba(245, 158, 11, 0.35);
}

.label-with-presets {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.preset-tag-group {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.btn-preset-tag {
  appearance: none;
  background: rgba(245, 158, 11, 0.1);
  color: #fbbf24;
  border: 1px solid rgba(245, 158, 11, 0.25);
  padding: 2px 7px;
  border-radius: var(--radius-sm);
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-preset-tag:hover {
  background: rgba(245, 158, 11, 0.22);
  border-color: rgba(245, 158, 11, 0.5);
  box-shadow: 0 0 8px rgba(245, 158, 11, 0.3);
}

.setting-textarea {
  width: 100%;
  background: #0b111e !important;
  border: 1px solid #2d3e58 !important;
  border-radius: var(--radius-sm);
  color: #e2e8f0;
  padding: 10px 12px;
  font-size: 12px;
  resize: vertical;
  line-height: 1.6;
  min-height: 92px;
  transition: all var(--transition-fast);
  box-sizing: border-box;
}

.setting-textarea:focus {
  outline: none;
  border-color: #fbbf24 !important;
  box-shadow: 0 0 0 2px rgba(245, 158, 11, 0.25) !important;
}

.font-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
}

.field-hint {
  font-size: 11px;
  color: var(--text-dim);
  margin-top: 2px;
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
