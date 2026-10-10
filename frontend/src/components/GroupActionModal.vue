<template>
  <Teleport to="body">
    <transition name="modal-fade">
      <div
        v-if="store.isGroupModalOpen"
        class="modal-backdrop"
        @click.self="store.closeGroupModal()"
        @keydown.esc="store.closeGroupModal()"
      >
        <div class="modal-dialog group-modal" role="dialog" aria-modal="true">
          <!-- 弹窗标题栏 -->
          <div class="modal-header">
            <div class="header-title-wrap">
              <div class="header-icon-box">
                <SvgIcon :name="isEditMode ? 'edit' : 'layers'" :size="18" />
              </div>
              <div>
                <h3 class="modal-title">
                  {{ isEditMode ? $t('editGroup') : $t('groupManagement') }}
                </h3>
                <p class="modal-subtitle">
                  {{ isEditMode ? targetGroupName : targetFileCountText }}
                </p>
              </div>
            </div>
            <button
              class="modal-close-btn"
              :title="$t('cancel')"
              @click="store.closeGroupModal()"
            >
              <SvgIcon name="x" :size="15" />
            </button>
          </div>

          <!-- 模式一：为选中的文件设置/加入分组 (assign) -->
          <div v-if="!isEditMode" class="modal-body">
            <!-- 标签页切换：新建分组 vs 加入已有分组 -->
            <div class="tab-switcher">
              <button
                class="tab-btn"
                :class="{ active: activeTab === 'create' }"
                @click="activeTab = 'create'"
              >
                <SvgIcon name="folder-plus" :size="14" />
                <span>{{ $t('createGroup') }}</span>
              </button>
              <button
                class="tab-btn"
                :class="{ active: activeTab === 'existing' }"
                :disabled="store.groups.length === 0"
                @click="activeTab = 'existing'"
              >
                <SvgIcon name="layers" :size="14" />
                <span>{{ $t('addToGroup') }}</span>
                <span v-if="store.groups.length > 0" class="tab-count-badge">{{ store.groups.length }}</span>
              </button>
            </div>

            <!-- Tab 1: 新建分组 -->
            <div v-if="activeTab === 'create'" class="tab-content">
              <div class="form-group">
                <div class="form-label-row">
                  <label class="form-label">{{ $t('groupName') }}</label>
                  <div
                    v-if="predictionSource"
                    class="prediction-badge"
                    :class="predictionSource === 'ai' ? 'badge-ai' : 'badge-heuristic'"
                    :title="predictionReason || (predictionSource === 'ai' ? $t('aiPredicted') : $t('rulePredicted'))"
                  >
                    <SvgIcon :name="predictionSource === 'ai' ? 'sparkles' : 'cpu'" :size="12" />
                    <span>{{ predictionSource === 'ai' ? $t('aiPredicted') : $t('rulePredicted') }}</span>
                  </div>
                </div>
                <div class="input-with-action">
                  <input
                    ref="nameInputRef"
                    v-model="newGroupName"
                    type="text"
                    class="form-input"
                    :placeholder="$t('groupNamePlaceholder')"
                    @keydown.enter="handleCreateGroup"
                  />
                  <button
                    type="button"
                    class="btn-predict-action"
                    :disabled="isPredicting"
                    :title="$t('rePredictTooltip')"
                    @click="runPrediction"
                  >
                    <SvgIcon v-if="isPredicting" name="loader" :size="13" :spin="true" />
                    <SvgIcon v-else name="sparkles" :size="13" />
                    <span>{{ isPredicting ? $t('predicting') : $t('autoPredict') }}</span>
                  </button>
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">{{ $t('groupColor') }}</label>
                <div class="color-palette-grid">
                  <button
                    v-for="color in GROUP_COLOR_PALETTE"
                    :key="color"
                    type="button"
                    class="color-swatch-btn"
                    :class="{ 'is-selected': selectedColor === color }"
                    :style="{ '--swatch-color': color }"
                    @click="selectedColor = color"
                  >
                    <span class="color-swatch-circle" :style="{ backgroundColor: color }"></span>
                    <SvgIcon
                      v-if="selectedColor === color"
                      name="check"
                      :size="12"
                      stroke-width="3"
                      class="swatch-check"
                    />
                  </button>
                </div>
              </div>

              <!-- 预览小卡片 -->
              <div class="group-preview-card" :style="{ borderColor: selectedColor, background: previewBg }">
                <span class="preview-dot" :style="{ backgroundColor: selectedColor }"></span>
                <span class="preview-name">{{ newGroupName || $t('groupNamePlaceholder') }}</span>
                <span class="preview-count">({{ selectedFiles.length }} {{ $t('modsLoaded') }})</span>
              </div>
            </div>

            <!-- Tab 2: 加入已有分组 -->
            <div v-else class="tab-content">
              <div v-if="store.groups.length === 0" class="empty-groups-hint">
                <SvgIcon name="info" :size="18" />
                <span>{{ $t('noGroups') }}</span>
              </div>
              <div v-else class="existing-groups-list custom-scrollbar">
                <div
                  v-for="grp in store.groups"
                  :key="grp.id"
                  class="existing-group-item"
                  :class="{ 'is-selected': chosenExistingGroupId === grp.id }"
                  :style="{ '--grp-color': grp.color }"
                  @click="chosenExistingGroupId = grp.id"
                >
                  <div class="item-left">
                    <span class="group-color-dot" :style="{ backgroundColor: grp.color }"></span>
                    <span class="group-item-name">{{ grp.name }}</span>
                  </div>
                  <div class="item-right">
                    <span class="group-item-count">
                      {{ (grp.file_keys || []).length }} {{ $t('modsLoaded') }}
                    </span>
                    <span v-if="chosenExistingGroupId === grp.id" class="check-indicator">
                      <SvgIcon name="check" :size="13" />
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 快捷提示：移出当前分组选项 -->
            <div v-if="hasGroupedSelectedFiles" class="ungroup-quick-bar">
              <span class="ungroup-hint-text">
                <SvgIcon name="info" :size="13" />
                已勾选的文件中包含已分组项
              </span>
              <button class="ungroup-quick-btn" @click="handleRemoveSelectedFromGroup">
                <SvgIcon name="folder-minus" :size="13" />
                <span>{{ $t('removeFromGroup') }}</span>
              </button>
            </div>
          </div>

          <!-- 模式二：编辑已有分组 (edit) -->
          <div v-else class="modal-body">
            <div class="form-group">
              <label class="form-label">{{ $t('groupName') }}</label>
              <input
                ref="editNameInputRef"
                v-model="editGroupName"
                type="text"
                class="form-input"
                :placeholder="$t('groupNamePlaceholder')"
                @keydown.enter="handleSaveEdit"
              />
            </div>

            <div class="form-group">
              <label class="form-label">{{ $t('groupColor') }}</label>
              <div class="color-palette-grid">
                <button
                  v-for="color in GROUP_COLOR_PALETTE"
                  :key="color"
                  type="button"
                  class="color-swatch-btn"
                  :class="{ 'is-selected': editColor === color }"
                  :style="{ '--swatch-color': color }"
                  @click="editColor = color"
                >
                  <span class="color-swatch-circle" :style="{ backgroundColor: color }"></span>
                  <SvgIcon
                    v-if="editColor === color"
                    name="check"
                    :size="12"
                    stroke-width="3"
                    class="swatch-check"
                  />
                </button>
              </div>
            </div>

            <div class="edit-preview-card" :style="{ borderColor: editColor, background: editPreviewBg }">
              <span class="preview-dot" :style="{ backgroundColor: editColor }"></span>
              <span class="preview-name">{{ editGroupName || targetGroupName }}</span>
              <span class="preview-count">({{ editingGroupFilesCount }} {{ $t('modsLoaded') }})</span>
            </div>
          </div>

          <!-- 底部操作按钮 -->
          <div class="modal-footer">
            <div v-if="isEditMode" class="footer-left">
              <button class="btn-dissolve" @click="handleDissolveGroup">
                <SvgIcon name="trash" :size="13" />
                <span>{{ $t('dissolveGroup') }}</span>
              </button>
            </div>
            <div v-else class="footer-left"></div>

            <div class="footer-right">
              <button class="btn-cancel" @click="store.closeGroupModal()">
                {{ $t('cancel') }}
              </button>

              <button
                v-if="!isEditMode && activeTab === 'create'"
                class="btn-primary"
                @click="handleCreateGroup"
              >
                <SvgIcon name="folder-plus" :size="14" />
                <span>{{ $t('createGroupAndAdd') }}</span>
              </button>

              <button
                v-else-if="!isEditMode && activeTab === 'existing'"
                class="btn-primary"
                :disabled="!chosenExistingGroupId"
                @click="handleAddToExistingGroup"
              >
                <SvgIcon name="check" :size="14" />
                <span>{{ $t('addToGroup') }}</span>
              </button>

              <button
                v-else-if="isEditMode"
                class="btn-primary"
                @click="handleSaveEdit"
              >
                <SvgIcon name="check" :size="14" />
                <span>{{ $t('saveGroup') }}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { useVpkManagerStore, GROUP_COLOR_PALETTE, hexToRgb, getVpkKey } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import SvgIcon from '@/components/common/SvgIcon.vue'

const store = useVpkManagerStore()
const { t: $t } = useI18n()

const nameInputRef = ref<HTMLInputElement | null>(null)
const editNameInputRef = ref<HTMLInputElement | null>(null)

const activeTab = ref<'create' | 'existing'>('create')
const newGroupName = ref('')
const selectedColor = ref(GROUP_COLOR_PALETTE[0])

const isPredicting = ref(false)
const predictionSource = ref<'ai' | 'heuristic' | ''>('')
const predictionReason = ref('')

async function runPrediction() {
  if (isPredicting.value) return
  isPredicting.value = true
  try {
    const res = await store.predictGroupName()
    if (res && res.name) {
      newGroupName.value = res.name
      predictionSource.value = res.source
      predictionReason.value = res.reason || ''
    }
  } catch (e) {
    console.warn('Prediction failed:', e)
  } finally {
    isPredicting.value = false
  }
}

watch(activeTab, (tab) => {
  if (tab === 'create' && !newGroupName.value) {
    runPrediction()
  }
})

const editGroupName = ref('')
const editColor = ref(GROUP_COLOR_PALETTE[0])
const chosenExistingGroupId = ref<string>('')

const isEditMode = computed(() => store.groupModalMode === 'edit')

const selectedFiles = computed(() => store.getSelectedVpkPaths())

const targetFileCountText = computed(() => {
  return `${selectedFiles.value.length} 个选中的模组文件`
})

const editingGroup = computed(() => {
  if (!store.editingGroupId) return null
  return store.groups.find(g => g.id === store.editingGroupId) || null
})

const targetGroupName = computed(() => {
  return editingGroup.value ? editingGroup.value.name : ''
})

const editingGroupFilesCount = computed(() => {
  return editingGroup.value?.file_keys?.length || 0
})

const previewBg = computed(() => {
  const rgb = hexToRgb(selectedColor.value)
  return `rgba(${rgb}, 0.12)`
})

const editPreviewBg = computed(() => {
  const rgb = hexToRgb(editColor.value)
  return `rgba(${rgb}, 0.12)`
})

const hasGroupedSelectedFiles = computed(() => {
  return selectedFiles.value.some(p => store.getFileGroup(p) !== undefined)
})

// 当弹窗打开时初始化表单
watch(
  () => store.isGroupModalOpen,
  (open) => {
    if (open) {
      if (store.groupModalMode === 'edit') {
        if (editingGroup.value) {
          editGroupName.value = editingGroup.value.name
          editColor.value = editingGroup.value.color
        }
        nextTick(() => {
          editNameInputRef.value?.focus()
        })
      } else {
        // assign mode
        activeTab.value = store.groups.length > 0 ? 'existing' : 'create'
        chosenExistingGroupId.value = store.groups.length > 0 ? store.groups[0].id : ''
        // Pick an unused or next color
        const nextColorIdx = store.groups.length % GROUP_COLOR_PALETTE.length
        selectedColor.value = GROUP_COLOR_PALETTE[nextColorIdx]
        newGroupName.value = ''
        predictionSource.value = ''
        predictionReason.value = ''
        if (activeTab.value === 'create') {
          runPrediction()
        }
        nextTick(() => {
          nameInputRef.value?.focus()
        })
      }
    }
  }
)

function handleCreateGroup() {
  const files = selectedFiles.value
  const name = newGroupName.value.trim()
  store.createGroup(name, selectedColor.value, files)
  store.closeGroupModal()
}

function handleAddToExistingGroup() {
  if (!chosenExistingGroupId.value) return
  const files = selectedFiles.value
  store.addFilesToGroup(chosenExistingGroupId.value, files)
  store.closeGroupModal()
}

function handleRemoveSelectedFromGroup() {
  const files = selectedFiles.value
  store.removeFilesFromGroups(files)
  store.closeGroupModal()
}

function handleSaveEdit() {
  if (!store.editingGroupId) return
  store.updateGroup(store.editingGroupId, {
    name: editGroupName.value.trim(),
    color: editColor.value,
  })
  store.closeGroupModal()
}

async function handleDissolveGroup() {
  if (!editingGroup.value) return
  const confirmed = await store.showConfirm(
    $t('confirmDissolveGroup', {
      name: editingGroup.value.name,
      count: editingGroupFilesCount.value,
    })
  )
  if (confirmed) {
    store.dissolveGroup(editingGroup.value.id)
    store.closeGroupModal()
  }
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.72);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
}

.modal-dialog {
  width: 100%;
  max-width: 480px;
  background: var(--bg-surface);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-modal);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  animation: modal-pop 0.2s cubic-bezier(0.16, 1, 0.3, 1);
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
  padding: 18px 20px;
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.02);
}

.header-title-wrap {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-icon-box {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  background: rgba(56, 189, 248, 0.12);
  border: 1px solid rgba(56, 189, 248, 0.3);
  color: var(--accent-primary);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-main);
  line-height: 1.2;
}

.modal-subtitle {
  font-size: 12px;
  color: var(--text-dim);
  margin-top: 2px;
}

.modal-close-btn {
  width: 30px;
  height: 30px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-dim);
  background: transparent;
  transition: all var(--transition-fast);
}

.modal-close-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-main);
}

.modal-body {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-height: 70vh;
  overflow-y: auto;
}

/* 标签页切换 */
.tab-switcher {
  display: flex;
  background: var(--bg-input);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 3px;
  gap: 3px;
}

.tab-btn {
  flex: 1;
  padding: 8px 12px;
  font-size: 13px;
  font-weight: 500;
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: all var(--transition-fast);
}

.tab-btn:hover:not(:disabled) {
  color: var(--text-main);
}

.tab-btn.active {
  background: var(--bg-card);
  color: var(--accent-primary);
  font-weight: 600;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.tab-btn:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.tab-count-badge {
  font-size: 11px;
  padding: 1px 6px;
  border-radius: var(--radius-full);
  background: rgba(56, 189, 248, 0.15);
  color: var(--accent-primary);
}

.tab-content {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* 表单输入控件 */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.form-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.prediction-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  font-weight: 500;
  line-height: 1.2;
}

.prediction-badge.badge-ai {
  background: rgba(168, 85, 247, 0.16);
  color: #c084fc;
  border: 1px solid rgba(168, 85, 247, 0.35);
}

.prediction-badge.badge-heuristic {
  background: rgba(56, 189, 248, 0.14);
  color: #38bdf8;
  border: 1px solid rgba(56, 189, 248, 0.3);
}

.input-with-action {
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn-predict-action {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 0 13px;
  height: 38px;
  background: rgba(168, 85, 247, 0.12);
  border: 1px solid rgba(168, 85, 247, 0.35);
  border-radius: var(--radius-md);
  color: #c084fc;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  flex-shrink: 0;
  transition: all var(--transition-fast);
}

.btn-predict-action:hover:not(:disabled) {
  background: rgba(168, 85, 247, 0.22);
  border-color: #c084fc;
  color: #ffffff;
  box-shadow: 0 0 10px rgba(168, 85, 247, 0.25);
}

.btn-predict-action:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.form-input {
  width: 100%;
  padding: 10px 14px;
  background: var(--bg-input);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  color: var(--text-main);
  font-size: 13px;
  outline: none;
  transition: all var(--transition-fast);
}

.form-input:focus {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 1px var(--border-focus), 0 0 8px rgba(56, 189, 248, 0.2);
}

/* 调色盘 Swatches */
.color-palette-grid {
  display: grid;
  grid-template-columns: repeat(6, 1fr);
  gap: 8px;
}

.color-swatch-btn {
  width: 100%;
  aspect-ratio: 1;
  border-radius: var(--radius-sm);
  background: var(--bg-input);
  border: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  transition: all var(--transition-fast);
}

.color-swatch-btn:hover {
  transform: scale(1.08);
  border-color: var(--swatch-color);
}

.color-swatch-circle {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.4);
}

.color-swatch-btn.is-selected {
  border-color: #ffffff;
  box-shadow: 0 0 0 1px #ffffff, 0 0 10px var(--swatch-color);
  background: rgba(255, 255, 255, 0.08);
}

.swatch-check {
  position: absolute;
  color: #ffffff;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.8));
}

/* 预览卡片 */
.group-preview-card,
.edit-preview-card {
  padding: 10px 14px;
  border-radius: var(--radius-md);
  border: 1px solid;
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  transition: all var(--transition-fast);
}

.preview-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.preview-name {
  font-weight: 600;
  color: var(--text-main);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-count {
  font-size: 11px;
  color: var(--text-dim);
  margin-left: auto;
}

/* 已有分组列表 */
.existing-groups-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  max-height: 220px;
  overflow-y: auto;
  padding-right: 2px;
}

.existing-group-item {
  padding: 10px 12px;
  background: var(--bg-input);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.existing-group-item:hover {
  background: var(--bg-card-hover);
  border-color: var(--grp-color);
}

.existing-group-item.is-selected {
  border-color: var(--grp-color);
  background: rgba(255, 255, 255, 0.05);
  box-shadow: 0 0 0 1px var(--grp-color);
}

.item-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.group-color-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  flex-shrink: 0;
}

.group-item-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-main);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.group-item-count {
  font-size: 11px;
  color: var(--text-dim);
}

.check-indicator {
  color: var(--accent-primary);
}

.empty-groups-hint {
  padding: 30px 16px;
  text-align: center;
  color: var(--text-dim);
  font-size: 13px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

/* 移出分组快速操作条 */
.ungroup-quick-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: rgba(245, 158, 11, 0.08);
  border: 1px solid rgba(245, 158, 11, 0.25);
  border-radius: var(--radius-md);
  font-size: 12px;
}

.ungroup-hint-text {
  color: var(--accent-warning);
  display: flex;
  align-items: center;
  gap: 6px;
}

.ungroup-quick-btn {
  padding: 4px 10px;
  font-size: 12px;
  background: rgba(245, 158, 11, 0.18);
  border: 1px solid rgba(245, 158, 11, 0.4);
  color: var(--accent-warning);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.ungroup-quick-btn:hover {
  background: rgba(245, 158, 11, 0.3);
  color: #ffffff;
}

/* 底部操作 */
.modal-footer {
  padding: 14px 20px;
  border-top: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.01);
}

.footer-left {
  display: flex;
  align-items: center;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.btn-cancel {
  padding: 8px 14px;
  font-size: 13px;
  color: var(--text-muted);
  border-radius: var(--radius-sm);
}

.btn-cancel:hover {
  background: rgba(255, 255, 255, 0.06);
  color: var(--text-main);
}

.btn-primary {
  padding: 8px 16px;
  background: var(--accent-primary);
  color: #0b111e;
  font-size: 13px;
  font-weight: 600;
  border-radius: var(--radius-sm);
  display: inline-flex;
  align-items: center;
  gap: 6px;
  transition: all var(--transition-fast);
}

.btn-primary:hover:not(:disabled) {
  background: var(--accent-primary-hover);
  box-shadow: 0 0 12px var(--accent-primary-glow);
}

.btn-primary:disabled {
  opacity: 0.45;
  cursor: not-allowed;
}

.btn-dissolve {
  padding: 8px 12px;
  font-size: 12px;
  color: var(--accent-danger);
  background: rgba(244, 63, 94, 0.1);
  border: 1px solid rgba(244, 63, 94, 0.25);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  gap: 6px;
  transition: all var(--transition-fast);
}

.btn-dissolve:hover {
  background: rgba(244, 63, 94, 0.22);
  color: #ffffff;
  border-color: rgba(244, 63, 94, 0.5);
}

.modal-fade-enter-active,
.modal-fade-leave-active {
  transition: opacity 0.2s ease;
}

.modal-fade-enter-from,
.modal-fade-leave-to {
  opacity: 0;
}
</style>
