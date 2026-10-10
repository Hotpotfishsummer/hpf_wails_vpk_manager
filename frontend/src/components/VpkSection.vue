<template>
  <section class="vpk-section">
    <div class="section-header" @click="toggleExpanded">
      <div class="header-titles">
        <div class="title-with-badge">
          <h2 class="section-title">{{ title }}</h2>
          <span class="count-pill">
            {{ isFiltering ? `${files.length} / ${totalCount}` : files.length }}
          </span>
        </div>
        <p class="section-subtitle">{{ subtitle }}</p>
      </div>

      <div class="header-actions">
        <button
          class="folder-quick-btn"
          @click.stop="$emit('open-folder')"
          :title="$t('openInFileManager')"
        >
          <SvgIcon name="folder-open" :size="15" />
          <span class="folder-btn-text">{{ $t('openFolder', '打开目录') }}</span>
        </button>

        <span class="section-chevron" :class="{ 'is-expanded': expanded }">
          <SvgIcon name="chevron-down" :size="18" />
        </span>
      </div>
    </div>

    <transition name="expand">
      <div v-if="expanded" class="section-content-wrapper">
        <div class="section-content">
          <div v-if="files.length === 0 && isScanning" class="file-list">
            <SkeletonLoader :item-count="3" />
          </div>
          <div v-else-if="files.length === 0" class="empty-message-box">
            <SvgIcon name="info" :size="20" class="empty-info-icon" />
            <p>{{ isFiltering ? $t('noSearchResults') : emptyText }}</p>
          </div>
          <div v-else class="file-list" :class="{ 'compact-mode': store.viewMode === 'compact' }">
            <!-- 全选与当前区块统计 -->
            <div class="section-toolbar">
              <label class="custom-checkbox" @click.stop>
                <input
                  type="checkbox"
                  :checked="allSelected"
                  @change="toggleAll"
                />
                <span class="checkbox-mark" :class="{ 'is-checked': allSelected }">
                  <SvgIcon v-if="allSelected" name="check" :size="12" stroke-width="3.2" class="check-icon" />
                </span>
                <span class="checkbox-label-text">{{ allSelected ? $t('deselectAll') : $t('selectAll') }}</span>
              </label>

              <span v-if="selectedCount > 0" class="selection-pill">
                {{ $t('selectedFilesWithSize', { count: selectedCount, size: selectedTotalSize }) }}
              </span>
            </div>

            <!-- 分组卡片列表 (如果有分组) -->
            <div v-if="displayGroups.length > 0" class="groups-container">
              <div
                v-for="item in displayGroups"
                :key="item.group.id"
                class="vpk-group-card"
                :class="{ 'is-collapsed': item.group.collapsed }"
                :style="{
                  '--group-color': item.group.color,
                  '--group-rgb': hexToRgb(item.group.color),
                }"
              >
                <!-- 分组头部 -->
                <div class="group-header" @click="store.toggleGroupCollapse(item.group.id)">
                  <div class="group-header-left">
                    <span class="group-chevron" :class="{ 'is-collapsed': item.group.collapsed }">
                      <SvgIcon name="chevron-down" :size="15" />
                    </span>

                    <label class="custom-checkbox group-checkbox" @click.stop>
                      <input
                        type="checkbox"
                        :checked="item.allSelected"
                        @change="toggleGroupSelect(item)"
                      />
                      <span class="checkbox-mark" :class="{ 'is-checked': item.allSelected }">
                        <SvgIcon v-if="item.allSelected" name="check" :size="11" stroke-width="3.2" class="check-icon" />
                        <span v-else-if="item.someSelected" class="indeterminate-bar"></span>
                      </span>
                    </label>

                    <span class="group-dot" :style="{ backgroundColor: item.group.color }"></span>
                    <span class="group-title" :title="item.group.name">{{ item.group.name }}</span>
                    <span class="group-count-badge">{{ item.files.length }} {{ $t('modsLoaded') }}</span>
                    <span v-if="item.group.collapsed" class="group-collapsed-hint">
                      {{ $t('collapsedGroupCount', { count: item.files.length }) }}
                    </span>
                  </div>

                  <div class="group-header-right" @click.stop>
                    <button
                      class="group-action-btn edit"
                      :title="$t('editGroup')"
                      @click="store.openGroupModal('edit', item.group.id)"
                    >
                      <SvgIcon name="edit" :size="12" />
                      <span class="btn-text">{{ $t('editGroup') }}</span>
                    </button>
                    <button
                      class="group-action-btn dissolve"
                      :title="$t('dissolveGroup')"
                      @click="dissolveGroupConfirm(item.group)"
                    >
                      <SvgIcon name="trash" :size="12" />
                      <span class="btn-text">{{ $t('dissolveGroup') }}</span>
                    </button>
                  </div>
                </div>

                <!-- 分组文件内容 -->
                <transition name="group-body-collapse">
                  <div v-if="!item.group.collapsed" class="group-body">
                    <div class="items-container">
                      <VpkFileItem
                        v-for="vpk in item.files"
                        :key="vpk.path"
                        :vpk="vpk"
                        :is-selected="selectedSet.has(vpk.path)"
                        @toggle="$emit('toggle-select', vpk.path)"
                        @enable="$emit('enable', $event)"
                        @disable="$emit('disable', $event)"
                        @delete="$emit('delete', $event)"
                      />
                    </div>
                  </div>
                </transition>
              </div>
            </div>

            <!-- 未分组文件列表 -->
            <div v-if="ungroupedFiles.length > 0" class="ungrouped-section">
              <div v-if="displayGroups.length > 0" class="ungrouped-header">
                <span class="ungrouped-title">{{ $t('ungroupedMods') }}</span>
                <span class="ungrouped-count">({{ ungroupedFiles.length }})</span>
              </div>
              <div class="items-container">
                <VpkFileItem
                  v-for="vpk in ungroupedFiles"
                  :key="vpk.path"
                  :vpk="vpk"
                  :is-selected="selectedSet.has(vpk.path)"
                  @toggle="$emit('toggle-select', vpk.path)"
                  @enable="$emit('enable', $event)"
                  @disable="$emit('disable', $event)"
                  @delete="$emit('delete', $event)"
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </section>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { VpkFile, useVpkManagerStore, getVpkKey, hexToRgb, VpkGroup } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import VpkFileItem from '@/components/VpkFileItem.vue'
import SkeletonLoader from '@/components/SkeletonLoader.vue'
import SvgIcon from '@/components/common/SvgIcon.vue'

const props = defineProps<{
  title: string
  subtitle: string
  emptyText: string
  files: VpkFile[]
  totalCount?: number
  selectedSet: Set<string>
  isScanning?: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle-select', path: string): void
  (e: 'select-all', paths: string[]): void
  (e: 'deselect-all', paths: string[]): void
  (e: 'open-folder'): void
  (e: 'enable', path: string): void
  (e: 'disable', path: string): void
  (e: 'delete', path: string): void
}>()

const store = useVpkManagerStore()
const { t: $t } = useI18n()
const expanded = ref(true)

interface GroupDisplay {
  group: VpkGroup
  files: VpkFile[]
  allSelected: boolean
  someSelected: boolean
  selectedCount: number
}

// 收集在当前 Section 中有文件归属的所有分组（同一个组的文件强制聚集显示在一起）
const displayGroups = computed<GroupDisplay[]>(() => {
  const result: GroupDisplay[] = []
  for (const grp of store.groups) {
    const keySet = new Set(grp.file_keys || [])
    const grpFiles = props.files.filter(f => keySet.has(getVpkKey(f.path)))
    if (grpFiles.length > 0) {
      const selected = grpFiles.filter(f => props.selectedSet.has(f.path))
      result.push({
        group: grp,
        files: grpFiles,
        allSelected: selected.length === grpFiles.length,
        someSelected: selected.length > 0 && selected.length < grpFiles.length,
        selectedCount: selected.length,
      })
    }
  }
  return result
})

// 未分组文件
const ungroupedFiles = computed<VpkFile[]>(() => {
  return props.files.filter(f => store.getFileGroup(f.path) === undefined)
})

function toggleGroupSelect(item: GroupDisplay) {
  const paths = item.files.map(f => f.path)
  if (item.allSelected) {
    emit('deselect-all', paths)
  } else {
    emit('select-all', paths)
  }
}

async function dissolveGroupConfirm(group: VpkGroup) {
  const confirmed = await store.showConfirm(
    $t('confirmDissolveGroup', {
      name: group.name,
      count: group.file_keys?.length || 0,
    })
  )
  if (confirmed) {
    store.dissolveGroup(group.id)
  }
}

const isFiltering = computed(() => {
  return props.totalCount !== undefined && props.files.length < props.totalCount
})

const allSelected = computed(() => {
  return props.files.length > 0 && props.files.every((vpk) => props.selectedSet.has(vpk.path))
})

const selectedCount = computed(() => {
  return props.files.filter((vpk) => props.selectedSet.has(vpk.path)).length
})

const selectedTotalSize = computed(() => {
  const bytes = props.files
    .filter((vpk) => props.selectedSet.has(vpk.path))
    .reduce((acc, curr) => acc + curr.size, 0)
  const mb = bytes / (1024 * 1024)
  if (mb < 1) {
    const kb = bytes / 1024
    return kb < 1 ? `${bytes} B` : `${kb.toFixed(2)} KB`
  }
  return `${mb.toFixed(2)} MB`
})

function toggleExpanded() {
  expanded.value = !expanded.value
}

function toggleAll() {
  const currentPaths = props.files.map((v) => v.path)
  if (allSelected.value) {
    emit('deselect-all', currentPaths)
  } else {
    emit('select-all', currentPaths)
  }
}
</script>

<style scoped>
.vpk-section {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: border-color var(--transition-fast);
}

.vpk-section:hover {
  border-color: var(--border-medium);
}

.section-header {
  padding: 14px 18px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid var(--border-subtle);
  user-select: none;
  transition: background var(--transition-fast);
}

.section-header:hover {
  background: rgba(255, 255, 255, 0.04);
}

.header-titles {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.title-with-badge {
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-title {
  font-size: 15px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.01em;
}

.count-pill {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 8px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.08);
  color: var(--text-muted);
  border: 1px solid var(--border-subtle);
}

.section-subtitle {
  font-size: 12px;
  color: var(--text-dim);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.folder-quick-btn {
  appearance: none;
  -webkit-appearance: none;
  height: 28px;
  padding: 0 12px;
  border-radius: var(--radius-sm);
  background: #182337;
  border: 1px solid var(--border-medium);
  color: #cbd5e1;
  font-size: 12px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.folder-quick-btn:hover {
  background: #22314a;
  color: #fff;
  border-color: var(--border-strong);
}

.section-chevron {
  color: var(--text-dim);
  display: flex;
  align-items: center;
  transition: transform var(--transition-normal);
}

.section-chevron.is-expanded {
  transform: rotate(0deg);
}

.section-chevron:not(.is-expanded) {
  transform: rotate(-90deg);
}

.section-content {
  padding: 14px 18px;
}

.empty-message-box {
  padding: 24px;
  text-align: center;
  color: var(--text-dim);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  font-size: 13px;
}

.empty-info-icon {
  color: var(--text-dim);
  opacity: 0.6;
}

.section-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 12px;
  margin-bottom: 12px;
  border-bottom: 1px solid var(--border-subtle);
}

.checkbox-label-text {
  font-size: 12px;
  font-weight: 600;
  color: #cbd5e1;
  margin-left: 8px;
}

.selection-pill {
  font-size: 11px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: var(--radius-full);
  background: rgba(56, 189, 248, 0.15);
  color: var(--accent-primary);
  border: 1px solid rgba(56, 189, 248, 0.3);
}

.items-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

/* 紧凑模式下的间距更小 */
.compact-mode .items-container {
  gap: 4px;
}

/* 自定义复选框 */
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

/* ==================== 分组卡片与折叠容器 ==================== */
.groups-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 12px;
}

.vpk-group-card {
  background: rgba(var(--group-rgb), 0.04);
  border: 1px solid rgba(var(--group-rgb), 0.35);
  border-radius: var(--radius-md);
  overflow: hidden;
  transition: all var(--transition-fast);
}

.vpk-group-card:hover {
  border-color: rgba(var(--group-rgb), 0.55);
}

.group-header {
  padding: 10px 14px;
  background: rgba(var(--group-rgb), 0.12);
  border-bottom: 1px solid rgba(var(--group-rgb), 0.25);
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  user-select: none;
  transition: all var(--transition-fast);
}

.group-header:hover {
  background: rgba(var(--group-rgb), 0.18);
}

.vpk-group-card.is-collapsed .group-header {
  border-bottom: none;
}

.group-header-left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.group-chevron {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  color: var(--group-color);
  transition: transform var(--transition-normal);
}

.group-chevron.is-collapsed {
  transform: rotate(-90deg);
}

.group-checkbox {
  margin: 0;
}

.indeterminate-bar {
  width: 8px;
  height: 2px;
  background: #ffffff;
  border-radius: 1px;
}

.group-dot {
  width: 9px;
  height: 9px;
  border-radius: 50%;
  flex-shrink: 0;
  box-shadow: 0 0 6px var(--group-color);
}

.group-title {
  font-size: 13px;
  font-weight: 700;
  color: #ffffff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.group-count-badge {
  font-size: 11px;
  font-weight: 600;
  padding: 1px 7px;
  border-radius: var(--radius-full);
  background: rgba(var(--group-rgb), 0.2);
  color: var(--group-color);
  border: 1px solid rgba(var(--group-rgb), 0.35);
  flex-shrink: 0;
}

.group-collapsed-hint {
  font-size: 11px;
  color: var(--text-dim);
  margin-left: 4px;
}

.group-header-right {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-shrink: 0;
}

.group-action-btn {
  padding: 4px 8px;
  font-size: 11px;
  font-weight: 500;
  border-radius: var(--radius-sm);
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: rgba(0, 0, 0, 0.25);
  border: 1px solid rgba(255, 255, 255, 0.08);
  color: var(--text-muted);
  cursor: pointer;
  transition: all var(--transition-fast);
}

.group-action-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.group-action-btn.dissolve:hover {
  background: rgba(244, 63, 94, 0.2);
  color: var(--accent-danger);
  border-color: rgba(244, 63, 94, 0.4);
}

.group-body {
  padding: 10px;
}

/* 紧凑模式下分组内部的间距 */
.compact-mode .group-body {
  padding: 6px;
}

/* 未分组部分标题 */
.ungrouped-section {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ungrouped-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 4px 2px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-dim);
}

.ungrouped-count {
  font-size: 11px;
  color: var(--text-dim);
}

/* 折叠/展开过渡动画 */
.group-body-collapse-enter-active,
.group-body-collapse-leave-active {
  transition: all var(--transition-normal);
  overflow: hidden;
}

.group-body-collapse-enter-from,
.group-body-collapse-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
