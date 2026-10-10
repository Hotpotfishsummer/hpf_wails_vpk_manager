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

            <!-- 文件列表 -->
            <div class="items-container">
              <VpkFileItem
                v-for="vpk in files"
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
    </transition>
  </section>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { VpkFile, useVpkManagerStore } from '@/stores/vpk-manager'
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
const expanded = ref(true)

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
</style>
