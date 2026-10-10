<template>
  <div
    class="vpk-file-item"
    :class="{
      'is-disabled': vpk.is_disabled,
      'is-selected': isSelected,
      'is-compact': store.viewMode === 'compact',
      'is-grouped': !!group,
    }"
    :style="itemStyle"
    :title="cardHoverTitle"
    @click="handleCardClick"
  >
    <!-- 紧凑列表模式 -->
    <div v-if="store.viewMode === 'compact'" class="compact-row">
      <label class="custom-checkbox" @click.stop>
        <input
          type="checkbox"
          :checked="isSelected"
          @change="$emit('toggle')"
        />
        <span class="checkbox-mark" :class="{ 'is-checked': isSelected }">
          <SvgIcon v-if="isSelected" name="check" :size="12" stroke-width="3.2" class="check-icon" />
        </span>
      </label>

      <!-- 迷你缩略图 / 状态图标 -->
      <div class="mini-thumb" @click.stop="openLightbox">
        <img
          v-if="vpk.thumbnail_path"
          :src="vpk.thumbnail_path"
          class="mini-thumb-img"
          alt=""
          loading="lazy"
        />
        <span v-else class="mini-thumb-fallback">VPK</span>
      </div>

      <!-- 标题与文件名 -->
      <div class="compact-info">
        <div class="compact-title-wrap">
          <span class="compact-title" :title="vpk.display_name">{{ vpk.display_name }}</span>

          <!-- 地图标签与悬浮气泡 -->
          <div
            v-if="vpk.maps && vpk.maps.length > 0"
            class="map-badge-wrap"
            @click.stop
            @mouseenter="openMapPopover"
            @mouseleave="closeMapPopover"
          >
            <span class="map-pill" :class="{ 'is-active': isMapPopoverOpen }">
              <SvgIcon name="map" :size="10" />
              <span class="map-text">{{ formatMapBadgeLabel(vpk.maps) }}</span>
            </span>

            <transition name="popover-fade">
              <div v-if="isMapPopoverOpen" class="map-popover" @click.stop>
                <div class="map-popover-header">
                  <div class="map-popover-title">
                    <SvgIcon name="map" :size="13" class="map-header-icon" />
                    <span>{{ $t('mapInternalNames', '地图内部名称') }}</span>
                    <span class="map-count-tag">({{ vpk.maps.length }})</span>
                  </div>
                  <button
                    v-if="vpk.maps.length > 1"
                    class="copy-all-btn"
                    @click.stop="copyAllMapCommands"
                    :title="$t('copyAllMapCommands', '复制全部建图指令')"
                  >
                    <SvgIcon name="copy" :size="11" />
                    <span>{{ $t('copyAllMapCommands', '全部复制') }}</span>
                  </button>
                </div>
                <div class="map-popover-hint">
                  {{ $t('mapConsoleHint', '可在游戏控制台输入 map <代码> 启动') }}
                </div>
                <div class="map-popover-list custom-scrollbar">
                  <div
                    v-for="(m, idx) in vpk.maps"
                    :key="m.name"
                    class="map-popover-item"
                    :class="{ 'is-copied': copiedMap === m.name }"
                    @click.stop="copyMapCommand(m.name)"
                    :title="$t('copyMapCommand', '点击复制: map ') + m.name"
                  >
                    <div class="map-item-main">
                      <span class="map-item-idx">{{ idx + 1 }}.</span>
                      <code class="map-item-code">{{ m.name }}</code>
                      <span v-if="m.display_name" class="map-item-name">({{ m.display_name }})</span>
                    </div>
                    <div class="map-item-action">
                      <span v-if="copiedMap === m.name" class="copied-badge">
                        <SvgIcon name="check" :size="11" />
                        {{ $t('mapCommandCopied', '已复制') }}
                      </span>
                      <button v-else class="copy-icon-btn" :title="$t('copyMapCommand', '复制指令')">
                        <SvgIcon name="copy" :size="12" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </transition>
          </div>

          <span
            v-if="vpk.note"
            class="note-pill"
            :title="$t('customNote') + ': ' + vpk.note"
            @click.stop="openNoteModal"
          >
            <SvgIcon name="edit" :size="10" />
            <span class="note-text">{{ vpk.note }}</span>
          </span>

          <!-- 分组标签 -->
          <span
            v-if="group"
            class="group-pill"
            :style="{ '--badge-color': group.color }"
            :title="$t('group') + ': ' + group.name"
            @click.stop="openGroupAction"
          >
            <span class="group-pill-dot" :style="{ backgroundColor: group.color }"></span>
            <span class="group-pill-text">{{ group.name }}</span>
            <button
              class="group-pill-remove"
              :title="$t('removeFromGroup')"
              @click.stop="removeThisFromGroup"
            >
              <SvgIcon name="x" :size="9" />
            </button>
          </span>
        </div>
        <span class="compact-filename" :title="vpk.name">{{ vpk.name }}</span>
      </div>

      <!-- 状态标签 -->
      <div class="compact-status">
        <span v-if="vpk.is_disabled" class="status-pill disabled">
          {{ $t('disabledStatus', '已禁用') }}
        </span>
        <span v-else class="status-pill active">
          {{ $t('enabledStatus', '已启用') }}
        </span>
      </div>

      <!-- 文件大小 -->
      <span class="compact-size">{{ formatSize(vpk.size) }}</span>

      <!-- 快捷操作按钮 -->
      <div class="item-actions compact-actions" @click.stop>
        <button
          v-if="group"
          class="action-btn-circle ungroup"
          :title="$t('removeFromGroup')"
          @click="removeThisFromGroup"
        >
          <SvgIcon name="folder-minus" :size="13" />
        </button>
        <button
          v-else
          class="action-btn-circle group-add"
          :title="$t('addToGroup')"
          @click="openGroupAction"
        >
          <SvgIcon name="layers" :size="13" />
        </button>
        <button
          class="action-btn-circle note"
          :title="$t('editNote')"
          @click="openNoteModal"
        >
          <SvgIcon name="edit" :size="13" />
        </button>
        <button
          v-if="workshopId"
          @click="openWorkshop"
          class="action-btn-circle workshop"
          :title="$t('openWorkshopPage')"
        >
          <SvgIcon name="external-link" :size="13" />
        </button>
        <button
          v-if="vpk.is_disabled"
          class="action-btn-circle enable"
          :title="$t('enable')"
          @click="$emit('enable', vpk.path)"
        >
          <SvgIcon name="check" :size="13" />
        </button>
        <button
          v-else
          class="action-btn-circle disable"
          :title="$t('disable')"
          @click="$emit('disable', vpk.path)"
        >
          <SvgIcon name="ban" :size="13" />
        </button>
        <button
          class="action-btn-circle delete"
          :title="$t('hardDelete')"
          @click="handleDelete"
        >
          <SvgIcon name="trash" :size="13" />
        </button>
      </div>
    </div>

    <!-- 标准精致卡片模式 -->
    <div v-else class="card-content">
      <label class="custom-checkbox" @click.stop>
        <input
          type="checkbox"
          :checked="isSelected"
          @change="$emit('toggle')"
        />
        <span class="checkbox-mark" :class="{ 'is-checked': isSelected }">
          <SvgIcon v-if="isSelected" name="check" :size="12" stroke-width="3.2" class="check-icon" />
        </span>
      </label>

      <!-- 缩略图容器 (支持点击大图预览) -->
      <div class="thumbnail-wrapper" @click.stop="openLightbox" :title="$t('clickToPreview', '点击查看大图')">
        <img
          v-if="vpk.thumbnail_path"
          :src="vpk.thumbnail_path"
          class="thumbnail-image"
          :alt="vpk.display_name"
          loading="lazy"
          decoding="async"
        />
        <div v-else class="thumbnail-placeholder">
          <span class="vpk-tag">VPK</span>
        </div>
        <div v-if="vpk.thumbnail_path" class="thumbnail-zoom-hint">
          <SvgIcon name="eye" :size="14" />
        </div>
      </div>

      <!-- 模组元信息 -->
      <div class="file-info">
        <div class="title-row">
          <h3 class="file-name" :title="vpk.display_name">{{ vpk.display_name }}</h3>
          <span v-if="vpk.is_disabled" class="status-pill disabled">
            {{ $t('disabledStatus', '已禁用') }}
          </span>
          <span v-else class="status-pill active">
            {{ $t('enabledStatus', '已启用') }}
          </span>
        </div>

        <div class="meta-row">
          <span class="file-filename" :title="vpk.path">{{ vpk.name }}</span>
          <span class="dot-divider">•</span>
          <span class="file-size-badge">{{ formatSize(vpk.size) }}</span>

          <!-- 地图标签与悬浮气泡 -->
          <div
            v-if="vpk.maps && vpk.maps.length > 0"
            class="map-badge-wrap"
            @click.stop
            @mouseenter="openMapPopover"
            @mouseleave="closeMapPopover"
          >
            <span class="map-badge" :class="{ 'is-active': isMapPopoverOpen }">
              <SvgIcon name="map" :size="11" />
              <span class="map-text">{{ formatMapBadgeLabel(vpk.maps) }}</span>
            </span>

            <transition name="popover-fade">
              <div v-if="isMapPopoverOpen" class="map-popover" @click.stop>
                <div class="map-popover-header">
                  <div class="map-popover-title">
                    <SvgIcon name="map" :size="13" class="map-header-icon" />
                    <span>{{ $t('mapInternalNames', '地图内部名称') }}</span>
                    <span class="map-count-tag">({{ vpk.maps.length }})</span>
                  </div>
                  <button
                    v-if="vpk.maps.length > 1"
                    class="copy-all-btn"
                    @click.stop="copyAllMapCommands"
                    :title="$t('copyAllMapCommands', '复制全部建图指令')"
                  >
                    <SvgIcon name="copy" :size="11" />
                    <span>{{ $t('copyAllMapCommands', '全部复制') }}</span>
                  </button>
                </div>
                <div class="map-popover-hint">
                  {{ $t('mapConsoleHint', '可在游戏控制台输入 map <代码> 启动') }}
                </div>
                <div class="map-popover-list custom-scrollbar">
                  <div
                    v-for="(m, idx) in vpk.maps"
                    :key="m.name"
                    class="map-popover-item"
                    :class="{ 'is-copied': copiedMap === m.name }"
                    @click.stop="copyMapCommand(m.name)"
                    :title="$t('copyMapCommand', '点击复制: map ') + m.name"
                  >
                    <div class="map-item-main">
                      <span class="map-item-idx">{{ idx + 1 }}.</span>
                      <code class="map-item-code">{{ m.name }}</code>
                      <span v-if="m.display_name" class="map-item-name">({{ m.display_name }})</span>
                    </div>
                    <div class="map-item-action">
                      <span v-if="copiedMap === m.name" class="copied-badge">
                        <SvgIcon name="check" :size="11" />
                        {{ $t('mapCommandCopied', '已复制') }}
                      </span>
                      <button v-else class="copy-icon-btn" :title="$t('copyMapCommand', '复制指令')">
                        <SvgIcon name="copy" :size="12" />
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </transition>
          </div>

          <span
            v-if="vpk.note"
            class="note-pill"
            :title="$t('customNote') + ': ' + vpk.note"
            @click.stop="openNoteModal"
          >
            <SvgIcon name="edit" :size="11" />
            <span class="note-text">{{ vpk.note }}</span>
          </span>

          <!-- 分组标签 -->
          <span
            v-if="group"
            class="group-badge"
            :style="{ '--badge-color': group.color }"
            :title="$t('group') + ': ' + group.name"
            @click.stop="openGroupAction"
          >
            <span class="group-badge-dot" :style="{ backgroundColor: group.color }"></span>
            <span class="group-badge-text">{{ group.name }}</span>
            <button
              class="group-badge-remove"
              :title="$t('removeFromGroup')"
              @click.stop="removeThisFromGroup"
            >
              <SvgIcon name="x" :size="9" />
            </button>
          </span>
        </div>

        <div class="tag-row" v-if="workshopId">
          <button
            @click.stop="openWorkshop"
            class="workshop-link-btn"
            :title="$t('viewInWorkshop')"
          >
            <SvgIcon name="external-link" :size="12" />
            <span>{{ $t("openWorkshopPage") }}</span>
          </button>
        </div>
      </div>

      <!-- 操作按钮组 -->
      <div class="item-actions" @click.stop>
        <button
          v-if="group"
          class="action-btn-circle ungroup"
          :title="$t('removeFromGroup')"
          @click="removeThisFromGroup"
        >
          <SvgIcon name="folder-minus" :size="15" />
        </button>
        <button
          v-else
          class="action-btn-circle group-add"
          :title="$t('addToGroup')"
          @click="openGroupAction"
        >
          <SvgIcon name="layers" :size="15" />
        </button>
        <button
          class="action-btn-circle note"
          :title="$t('editNote')"
          @click="openNoteModal"
        >
          <SvgIcon name="edit" :size="15" />
        </button>
        <button
          v-if="vpk.is_disabled"
          class="action-btn-circle enable"
          :title="$t('enable')"
          @click="$emit('enable', vpk.path)"
        >
          <SvgIcon name="check" :size="15" />
        </button>
        <button
          v-else
          class="action-btn-circle disable"
          :title="$t('disable')"
          @click="$emit('disable', vpk.path)"
        >
          <SvgIcon name="ban" :size="15" />
        </button>
        <button
          class="action-btn-circle delete"
          :title="$t('hardDelete')"
          @click="handleDelete"
        >
          <SvgIcon name="trash" :size="15" />
        </button>
      </div>
    </div>

    <!-- 灯箱大图预览 -->
    <ImageLightbox
      :is-open="isLightboxOpen"
      :src="vpk.thumbnail_path || ''"
      :title="vpk.display_name"
      @close="isLightboxOpen = false"
    />

    <!-- 编辑附加备注模态框 -->
    <Teleport to="body">
      <div v-if="isNoteModalOpen" class="note-modal-overlay" @click.self="closeNoteModal">
        <div class="note-modal" role="dialog" aria-modal="true" @click.stop>
          <div class="note-modal-header">
            <div class="note-modal-title">
              <SvgIcon name="edit" :size="16" />
              <span>{{ $t('editNote', '附加备注') }}</span>
            </div>
            <button class="note-modal-close" @click="closeNoteModal">
              <SvgIcon name="x" :size="14" />
            </button>
          </div>

          <div class="note-modal-body">
            <div class="note-modal-target">
              <span class="target-name">{{ vpk.display_name }}</span>
              <span class="target-file">{{ vpk.name }}</span>
            </div>
            <input
              ref="noteInputRef"
              v-model="tempNote"
              type="text"
              class="note-modal-input"
              :placeholder="$t('notePlaceholder')"
              @keydown.enter="saveNote"
              @keydown.esc="closeNoteModal"
            />
            <p class="note-modal-hint">
              备注将独立保存在配置文件中，并在列表显示和开启“更好的命名”导出时自动追加。
            </p>
          </div>

          <div class="note-modal-footer">
            <button
              v-if="vpk.note"
              class="note-btn clear"
              @click="clearNote"
            >
              {{ $t('clearNote', '清空备注') }}
            </button>
            <div class="footer-right">
              <button class="note-btn cancel" @click="closeNoteModal">
                {{ $t('cancel', '取消') }}
              </button>
              <button class="note-btn confirm" @click="saveNote">
                {{ $t('saveNote', '保存备注') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { VpkFile, MapInfo, useVpkManagerStore, hexToRgb } from '@/stores/vpk-manager'
import { callApi } from '@/composables/useApi'
import { useI18n } from 'vue-i18n'
import SvgIcon from '@/components/common/SvgIcon.vue'
import ImageLightbox from '@/components/common/ImageLightbox.vue'

const props = defineProps<{
  vpk: VpkFile
  isSelected: boolean
}>()

const emit = defineEmits<{
  toggle: []
  enable: [path: string]
  disable: [path: string]
  delete: [path: string]
}>()

const store = useVpkManagerStore()
const { t: $t } = useI18n()
const isLightboxOpen = ref(false)

const group = computed(() => store.getFileGroup(props.vpk.path))

const itemStyle = computed(() => {
  if (!group.value) return {}
  const rgb = hexToRgb(group.value.color)
  return {
    '--group-color': group.value.color,
    '--group-rgb': rgb,
  }
})

function removeThisFromGroup() {
  store.removeFilesFromGroups([props.vpk.path])
}

function openGroupAction() {
  if (!group.value) {
    if (!props.isSelected) {
      emit('toggle')
    }
    store.openGroupModal('assign')
  } else {
    store.openGroupModal('edit', group.value.id)
  }
}

const isNoteModalOpen = ref(false)
const tempNote = ref('')
const noteInputRef = ref<HTMLInputElement | null>(null)

function openNoteModal() {
  tempNote.value = props.vpk.note || ''
  isNoteModalOpen.value = true
  setTimeout(() => {
    noteInputRef.value?.focus()
    noteInputRef.value?.select()
  }, 60)
}

function closeNoteModal() {
  isNoteModalOpen.value = false
}

async function saveNote() {
  await store.updateVpkNote(props.vpk.path, tempNote.value)
  closeNoteModal()
}

async function clearNote() {
  tempNote.value = ''
  await store.updateVpkNote(props.vpk.path, '')
  closeNoteModal()
}

function handleCardClick() {
  emit('toggle')
}

function openLightbox() {
  if (props.vpk.thumbnail_path) {
    isLightboxOpen.value = true
  }
}

async function handleDelete() {
  const confirmed = await store.showConfirm($t('confirmHardDelete'))
  if (confirmed) {
    emit('delete', props.vpk.path)
  }
}

const workshopId = computed(() => {
  const match = props.vpk.name.match(/\d{8,}/)
  return match ? match[0] : null
})

async function openWorkshop() {
  if (workshopId.value) {
    const url = `https://steamcommunity.com/sharedfiles/filedetails/?id=${workshopId.value}`
    await callApi(api => api.open_external_url(url), undefined)
  }
}

function formatSize(bytes: number): string {
  const mb = bytes / (1024 * 1024)
  if (mb < 1) {
    const kb = bytes / 1024
    return kb < 1 ? `${bytes} B` : `${kb.toFixed(2)} KB`
  }
  return `${mb.toFixed(2)} MB`
}

// ==================== 地图内部名称与悬浮气泡逻辑 ====================
const isMapPopoverOpen = ref(false)
const copiedMap = ref<string | null>(null)
let popoverTimeout: ReturnType<typeof setTimeout> | null = null

function openMapPopover() {
  if (popoverTimeout) {
    clearTimeout(popoverTimeout)
    popoverTimeout = null
  }
  isMapPopoverOpen.value = true
}

function closeMapPopover() {
  popoverTimeout = setTimeout(() => {
    isMapPopoverOpen.value = false
  }, 160)
}

function formatMapBadgeLabel(maps: MapInfo[]): string {
  if (maps.length === 1) {
    return maps[0].name
  }
  return $t('mapChapters', { count: maps.length }, `${maps.length}个章节`)
}

async function copyMapCommand(mapName: string) {
  const cmd = `map ${mapName}`
  try {
    await navigator.clipboard.writeText(cmd)
    copiedMap.value = mapName
    store.showNotification(`${$t('mapCommandCopied', '已复制')}: ${cmd}`, 'success')
    setTimeout(() => {
      if (copiedMap.value === mapName) {
        copiedMap.value = null
      }
    }, 1500)
  } catch (err) {
    console.error('Clipboard copy failed', err)
  }
}

async function copyAllMapCommands() {
  if (!props.vpk.maps || props.vpk.maps.length === 0) return
  const text = props.vpk.maps.map(m => `map ${m.name}`).join('\n')
  try {
    await navigator.clipboard.writeText(text)
    store.showNotification($t('allCommandsCopied', '已复制全部建图指令'), 'success')
  } catch (err) {
    console.error('Clipboard copy failed', err)
  }
}

const cardHoverTitle = computed(() => {
  if (!props.vpk.maps || props.vpk.maps.length === 0) {
    return props.vpk.display_name
  }
  const lines = [
    `${props.vpk.display_name}`,
    `【${$t('mapInternalNames', '地图内部名称')}】:`,
    ...props.vpk.maps.map((m, idx) => {
      const display = m.display_name ? ` (${m.display_name})` : ''
      return `  ${idx + 1}. ${m.name}${display}  ->  map ${m.name}`
    }),
    `(${$t('mapConsoleHint', '可在游戏控制台输入 map <代码> 启动')})`
  ]
  return lines.join('\n')
})
</script>

<style scoped>
.vpk-file-item {
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
  cursor: pointer;
  position: relative;
}

.vpk-file-item:hover {
  background: var(--bg-card-hover);
  border-color: var(--border-medium);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
  z-index: 20;
}

.vpk-file-item.is-selected {
  background: rgba(56, 189, 248, 0.08);
  border-color: var(--border-active);
  box-shadow: 0 0 0 1px var(--border-active), var(--shadow-sm);
}

.vpk-file-item.is-disabled {
  opacity: 0.65;
  filter: grayscale(0.25);
}

.vpk-file-item.is-disabled:hover {
  opacity: 0.9;
}

/* 自定义复选框 */
.custom-checkbox {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 22px;
  height: 22px;
  cursor: pointer;
  flex-shrink: 0;
}

.custom-checkbox input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.checkbox-mark {
  width: 18px;
  height: 18px;
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

/* ==================== 标准卡片模式 ==================== */
.card-content {
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 14px;
}

.thumbnail-wrapper {
  width: 80px;
  height: 54px;
  border-radius: var(--radius-sm);
  overflow: hidden;
  position: relative;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid var(--border-subtle);
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.thumbnail-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform var(--transition-normal);
}

.thumbnail-wrapper:hover .thumbnail-image {
  transform: scale(1.08);
}

.thumbnail-placeholder {
  color: var(--text-dim);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.05em;
}

.vpk-tag {
  background: rgba(255, 255, 255, 0.06);
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.thumbnail-zoom-hint {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  opacity: 0;
  transition: opacity var(--transition-fast);
}

.thumbnail-wrapper:hover .thumbnail-zoom-hint {
  opacity: 1;
}

.file-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.title-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-main);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.3;
}

.status-pill {
  font-size: 10px;
  font-weight: 600;
  padding: 2px 7px;
  border-radius: var(--radius-full);
  line-height: 1;
  flex-shrink: 0;
}

.status-pill.active {
  background: rgba(16, 185, 129, 0.15);
  color: var(--accent-success);
  border: 1px solid rgba(16, 185, 129, 0.3);
}

.status-pill.disabled {
  background: rgba(245, 158, 11, 0.15);
  color: var(--accent-warning);
  border: 1px solid rgba(245, 158, 11, 0.3);
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: var(--text-dim);
}

.file-filename {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 320px;
}

.dot-divider {
  opacity: 0.4;
}

.file-size-badge {
  color: var(--text-muted);
  font-size: 11px;
}

.tag-row {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 2px;
}

.workshop-link-btn {
  appearance: none;
  -webkit-appearance: none;
  font-size: 11px;
  font-weight: 600;
  color: #38bdf8;
  background: rgba(56, 189, 248, 0.12);
  border: 1px solid rgba(56, 189, 248, 0.35);
  padding: 3px 10px;
  border-radius: var(--radius-sm);
  display: inline-flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.workshop-link-btn:hover {
  background: rgba(56, 189, 248, 0.25);
  border-color: #38bdf8;
  color: #ffffff;
}

/* 操作按钮 */
.item-actions {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.action-btn-circle {
  appearance: none;
  -webkit-appearance: none;
  width: 32px;
  height: 32px;
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  background: #1c273c;
  border: 1px solid #2d3e58;
  color: #cbd5e1;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.action-btn-circle:hover {
  transform: scale(1.05);
}

.action-btn-circle.enable:hover {
  background: rgba(16, 185, 129, 0.25);
  color: var(--accent-success);
  border-color: rgba(16, 185, 129, 0.5);
}

.action-btn-circle.disable:hover {
  background: rgba(245, 158, 11, 0.25);
  color: var(--accent-warning);
  border-color: rgba(245, 158, 11, 0.5);
}

.action-btn-circle.delete:hover {
  background: rgba(244, 63, 94, 0.25);
  color: var(--accent-danger);
  border-color: rgba(244, 63, 94, 0.5);
}

.action-btn-circle.workshop:hover {
  background: rgba(56, 189, 248, 0.25);
  color: var(--accent-primary);
  border-color: rgba(56, 189, 248, 0.5);
}

/* ==================== 紧凑列表模式 ==================== */
.compact-row {
  padding: 8px 14px;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
}

.mini-thumb {
  width: 34px;
  height: 24px;
  border-radius: 4px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
}

.mini-thumb-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.mini-thumb-fallback {
  font-size: 9px;
  font-weight: 700;
  color: var(--text-dim);
}

.compact-info {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: baseline;
  gap: 10px;
}

.compact-title {
  font-weight: 600;
  color: var(--text-main);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 360px;
}

.compact-filename {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
  color: var(--text-dim);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.compact-status {
  flex-shrink: 0;
}

.compact-size {
  width: 75px;
  text-align: right;
  font-size: 12px;
  color: var(--text-muted);
  flex-shrink: 0;
}

.compact-actions .action-btn-circle {
  width: 26px;
  height: 26px;
}

.compact-title-wrap {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  max-width: 480px;
}

/* 附加备注药丸徽章 */
.note-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 1px 7px;
  border-radius: var(--radius-sm);
  background: rgba(168, 85, 247, 0.14);
  border: 1px solid rgba(168, 85, 247, 0.35);
  color: #c084fc;
  font-size: 11px;
  cursor: pointer;
  transition: all var(--transition-fast);
  max-width: 240px;
  user-select: none;
  flex-shrink: 0;
}

.note-pill:hover {
  background: rgba(168, 85, 247, 0.28);
  border-color: #c084fc;
  color: #f3e8ff;
  transform: translateY(-1px);
}

.note-pill .note-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ==================== 地图内部名称标签与浮动气泡 ==================== */
.map-badge-wrap {
  position: relative;
  display: inline-flex;
  align-items: center;
  flex-shrink: 0;
}

.map-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  background: rgba(14, 165, 233, 0.12);
  border: 1px solid rgba(14, 165, 233, 0.35);
  color: #38bdf8;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
  user-select: none;
}

.map-badge:hover,
.map-badge.is-active {
  background: rgba(14, 165, 233, 0.25);
  border-color: #38bdf8;
  color: #ffffff;
  box-shadow: 0 0 10px rgba(56, 189, 248, 0.25);
  transform: translateY(-1px);
}

.map-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 1px 7px;
  border-radius: var(--radius-sm);
  background: rgba(14, 165, 233, 0.12);
  border: 1px solid rgba(14, 165, 233, 0.35);
  color: #38bdf8;
  font-size: 11px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
  max-width: 220px;
  user-select: none;
  flex-shrink: 0;
}

.map-pill:hover,
.map-pill.is-active {
  background: rgba(14, 165, 233, 0.25);
  border-color: #38bdf8;
  color: #ffffff;
  box-shadow: 0 0 8px rgba(56, 189, 248, 0.25);
  transform: translateY(-1px);
}

.map-text {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.map-popover {
  position: absolute;
  top: calc(100% + 6px);
  left: 0;
  z-index: 100;
  width: 320px;
  background: #0f172a;
  border: 1px solid rgba(56, 189, 248, 0.4);
  box-shadow: 0 16px 36px rgba(0, 0, 0, 0.75), 0 4px 12px rgba(0, 0, 0, 0.4);
  border-radius: 8px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  cursor: default;
  pointer-events: auto;
}

.map-popover-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
  padding-bottom: 6px;
}

.map-popover-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 700;
  color: #38bdf8;
}

.map-header-icon {
  color: #38bdf8;
}

.map-count-tag {
  font-size: 11px;
  color: var(--text-dim);
  font-weight: 500;
}

.copy-all-btn {
  appearance: none;
  background: rgba(56, 189, 248, 0.12);
  border: 1px solid rgba(56, 189, 248, 0.3);
  color: #38bdf8;
  font-size: 11px;
  padding: 2px 7px;
  border-radius: 4px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  transition: all var(--transition-fast);
}

.copy-all-btn:hover {
  background: rgba(56, 189, 248, 0.25);
  border-color: #38bdf8;
  color: #ffffff;
}

.map-popover-hint {
  font-size: 11px;
  color: var(--text-muted);
  line-height: 1.3;
}

.map-popover-list {
  max-height: 180px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-right: 2px;
}

.map-popover-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 4px 6px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  transition: all var(--transition-fast);
  cursor: pointer;
}

.map-popover-item:hover {
  background: rgba(56, 189, 248, 0.12);
  border-color: rgba(56, 189, 248, 0.35);
}

.map-popover-item.is-copied {
  background: rgba(16, 185, 129, 0.18);
  border-color: rgba(16, 185, 129, 0.5);
}

.map-item-main {
  display: flex;
  align-items: center;
  gap: 6px;
  min-width: 0;
  overflow: hidden;
}

.map-item-idx {
  font-size: 11px;
  color: var(--text-dim);
  width: 14px;
  flex-shrink: 0;
}

.map-item-code {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
  color: #f1f5f9;
  font-weight: 600;
  white-space: nowrap;
}

.map-item-name {
  font-size: 11px;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.map-item-action {
  flex-shrink: 0;
  margin-left: 6px;
}

.copy-icon-btn {
  appearance: none;
  background: transparent;
  border: none;
  color: var(--text-dim);
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 2px;
  border-radius: 3px;
  transition: all var(--transition-fast);
}

.map-popover-item:hover .copy-icon-btn {
  color: #38bdf8;
}

.copied-badge {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  font-size: 10px;
  font-weight: 600;
  color: #34d399;
}

/* 气泡平滑动画 */
.popover-fade-enter-active,
.popover-fade-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.popover-fade-enter-from,
.popover-fade-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* 备注操作按钮 */
.action-btn-circle.note {
  color: #c084fc;
  border-color: rgba(168, 85, 247, 0.35);
}

.action-btn-circle.note:hover {
  background: rgba(168, 85, 247, 0.25);
  color: #f3e8ff;
  border-color: #c084fc;
}

/* 备注编辑模态框 */
.note-modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(10, 15, 26, 0.75);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  padding: 16px;
  animation: fadeIn 0.15s ease-out;
}

.note-modal {
  width: 100%;
  max-width: 480px;
  background: #141c2c;
  border: 1px solid #2d3e58;
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 40px -15px rgba(0, 0, 0, 0.6), 0 0 0 1px rgba(255, 255, 255, 0.05);
  overflow: hidden;
  animation: modalScale 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.note-modal-header {
  padding: 14px 18px;
  background: #192336;
  border-bottom: 1px solid #24334a;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.note-modal-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #e2e8f0;
}

.note-modal-close {
  appearance: none;
  background: transparent;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.note-modal-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #ffffff;
}

.note-modal-body {
  padding: 18px;
}

.note-modal-target {
  margin-bottom: 12px;
  padding: 8px 12px;
  background: #0f1624;
  border: 1px solid #1e2c42;
  border-radius: var(--radius-sm);
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.note-modal-target .target-name {
  font-size: 13px;
  font-weight: 600;
  color: #f1f5f9;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.note-modal-target .target-file {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 11px;
  color: #64748b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.note-modal-input {
  width: 100%;
  padding: 10px 14px;
  background: #0b111c;
  border: 1px solid #334155;
  border-radius: var(--radius-md);
  color: #f8fafc;
  font-size: 14px;
  outline: none;
  transition: all var(--transition-fast);
  box-sizing: border-box;
}

.note-modal-input:focus {
  border-color: #a855f7;
  box-shadow: 0 0 0 3px rgba(168, 85, 247, 0.2);
}

.note-modal-hint {
  margin: 10px 0 0 0;
  font-size: 12px;
  color: #94a3b8;
  line-height: 1.5;
}

.note-modal-footer {
  padding: 12px 18px;
  background: #111827;
  border-top: 1px solid #1e2c42;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.note-modal-footer .footer-right {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-left: auto;
}

.note-btn {
  padding: 6px 14px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all var(--transition-fast);
  border: 1px solid transparent;
}

.note-btn.clear {
  background: transparent;
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.3);
}

.note-btn.clear:hover {
  background: rgba(239, 68, 68, 0.15);
  border-color: #ef4444;
}

.note-btn.cancel {
  background: #1e293b;
  color: #94a3b8;
  border-color: #334155;
}

.note-btn.cancel:hover {
  background: #334155;
  color: #f1f5f9;
}

.note-btn.confirm {
  background: #9333ea;
  color: #ffffff;
}

.note-btn.confirm:hover {
  background: #a855f7;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes modalScale {
  from { transform: scale(0.96); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

/* ==================== 分组样式与高亮 ==================== */
.vpk-file-item.is-grouped {
  background: rgba(var(--group-rgb), 0.08) !important;
  border-color: rgba(var(--group-rgb), 0.35) !important;
}

.vpk-file-item.is-grouped:hover {
  background: rgba(var(--group-rgb), 0.15) !important;
  border-color: rgba(var(--group-rgb), 0.6) !important;
  box-shadow: 0 0 16px rgba(var(--group-rgb), 0.22);
}

.vpk-file-item.is-grouped.is-selected {
  background: rgba(var(--group-rgb), 0.22) !important;
  border-color: var(--group-color) !important;
  box-shadow: 0 0 0 1px var(--group-color), 0 0 16px rgba(var(--group-rgb), 0.35);
}

.vpk-file-item.is-grouped.is-compact {
  border-left: 3.5px solid var(--group-color) !important;
  background: rgba(var(--group-rgb), 0.07) !important;
}

.vpk-file-item.is-grouped.is-compact:hover {
  background: rgba(var(--group-rgb), 0.14) !important;
}

.vpk-file-item.is-grouped.is-compact.is-selected {
  background: rgba(var(--group-rgb), 0.22) !important;
}

/* 分组徽章 (Card Mode) */
.group-badge {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  background: rgba(var(--group-rgb), 0.15);
  border: 1px solid rgba(var(--group-rgb), 0.4);
  color: var(--group-color);
  font-size: 11px;
  font-weight: 600;
  padding: 2px 7px 2px 8px;
  border-radius: var(--radius-full);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.group-badge:hover {
  background: rgba(var(--group-rgb), 0.28);
  border-color: var(--group-color);
}

.group-badge-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
}

.group-badge-text {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-badge-remove {
  width: 14px;
  height: 14px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: inherit;
  opacity: 0.65;
  background: transparent;
  transition: all var(--transition-fast);
  margin-left: 2px;
}

.group-badge-remove:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

/* 分组胶囊 (Compact Mode) */
.group-pill {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: rgba(var(--group-rgb), 0.15);
  border: 1px solid rgba(var(--group-rgb), 0.35);
  color: var(--group-color);
  font-size: 10px;
  font-weight: 600;
  padding: 1px 6px 1px 6px;
  border-radius: var(--radius-full);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.group-pill:hover {
  background: rgba(var(--group-rgb), 0.26);
}

.group-pill-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.group-pill-text {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.group-pill-remove {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: inherit;
  opacity: 0.6;
  background: transparent;
  transition: all var(--transition-fast);
}

.group-pill-remove:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
}

/* 按钮悬停状态 */
.action-btn-circle.ungroup:hover {
  background: rgba(244, 63, 94, 0.2);
  color: var(--accent-danger);
  border-color: rgba(244, 63, 94, 0.4);
}

.action-btn-circle.group-add:hover {
  background: rgba(56, 189, 248, 0.2);
  color: var(--accent-primary);
  border-color: rgba(56, 189, 248, 0.4);
}
</style>
