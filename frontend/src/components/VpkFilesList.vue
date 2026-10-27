<template>
  <div class="vpk-files-list">
    <!-- 搜索与筛选工具栏 -->
    <div class="filter-toolbar">
      <div class="search-box">
        <SvgIcon name="search" :size="15" class="search-icon" />
        <input
          v-model="store.searchQuery"
          type="text"
          class="search-input"
          :placeholder="$t('searchFilesPlaceholder')"
          :aria-label="$t('searchFilesPlaceholder')"
        />
        <button
          v-if="store.searchQuery"
          class="search-clear-btn"
          :title="$t('clearSearch')"
          @click="store.searchQuery = ''"
        >
          <SvgIcon name="x" :size="14" />
        </button>
      </div>

      <!-- 快速筛选状态 Chips -->
      <div class="filter-chips">
        <button
          class="chip-btn"
          :class="{ active: store.filterStatus === 'all' }"
          @click="store.filterStatus = 'all'"
        >
          <span>{{ $t('filterAll', '全部') }}</span>
          <span class="chip-count">{{ store.totalModStats.total }}</span>
        </button>
        <button
          class="chip-btn active-state"
          :class="{ active: store.filterStatus === 'enabled' }"
          @click="store.filterStatus = 'enabled'"
        >
          <span class="chip-dot dot-green"></span>
          <span>{{ $t('filterEnabled', '已启用') }}</span>
          <span class="chip-count">{{ store.totalModStats.enabled }}</span>
        </button>
        <button
          class="chip-btn disabled-state"
          :class="{ active: store.filterStatus === 'disabled' }"
          @click="store.filterStatus = 'disabled'"
        >
          <span class="chip-dot dot-amber"></span>
          <span>{{ $t('filterDisabled', '已禁用') }}</span>
          <span class="chip-count">{{ store.totalModStats.disabled }}</span>
        </button>
      </div>

      <!-- 排序与视图切换 -->
      <div class="controls-right">
        <!-- 排序下拉 / 切换按钮 -->
        <div class="sort-group">
          <button
            class="sort-tab-btn"
            :class="{ active: store.sortKey === 'name' }"
            @click="toggleSort('name')"
            :title="$t('sortByName', '按名称排序')"
          >
            {{ $t('sortName', '名称') }}
            <span v-if="store.sortKey === 'name'" class="sort-indicator">{{ store.sortAsc ? '▲' : '▼' }}</span>
          </button>
          <button
            class="sort-tab-btn"
            :class="{ active: store.sortKey === 'size' }"
            @click="toggleSort('size')"
            :title="$t('sortBySize', '按文件大小排序')"
          >
            {{ $t('sortSize', '大小') }}
            <span v-if="store.sortKey === 'size'" class="sort-indicator">{{ store.sortAsc ? '▲' : '▼' }}</span>
          </button>
          <button
            class="sort-tab-btn"
            :class="{ active: store.sortKey === 'date' }"
            @click="toggleSort('date')"
            :title="$t('sortByDate', '按修改时间排序')"
          >
            {{ $t('sortDate', '时间') }}
            <span v-if="store.sortKey === 'date'" class="sort-indicator">{{ store.sortAsc ? '▲' : '▼' }}</span>
          </button>
        </div>

        <!-- 视图模式切换 (卡片 vs 紧凑) -->
        <div class="view-mode-toggle">
          <button
            class="view-btn"
            :class="{ active: store.viewMode === 'card' }"
            @click="store.setViewMode('card')"
            :title="$t('cardView', '标准卡片模式')"
          >
            <SvgIcon name="grid" :size="15" />
          </button>
          <button
            class="view-btn"
            :class="{ active: store.viewMode === 'compact' }"
            @click="store.setViewMode('compact')"
            :title="$t('compactView', '紧凑列表模式')"
          >
            <SvgIcon name="list" :size="15" />
          </button>
        </div>

        <!-- 刷新模组列表按钮 -->
        <button
          class="list-refresh-btn"
          @click="store.refreshVpks()"
          :disabled="store.isScanning"
          :title="$t('refreshTooltip')"
        >
          <SvgIcon name="refresh" :size="15" :spin="store.isScanning" />
        </button>
      </div>
    </div>

    <!-- 本地 VPK 模组区块 -->
    <VpkSection
      :title="$t('localVpkFiles')"
      :subtitle="$t('localVpkFilesSubtitle')"
      :empty-text="$t('noLocalVpkFiles')"
      :files="store.filteredLocalVpkFiles"
      :total-count="store.localVpkFiles.length"
      :selected-set="store.selectedLocalVpks"
      :is-scanning="store.isScanning"
      @toggle-select="store.toggleLocalVpkSelection"
      @select-all="store.selectAllLocal"
      @deselect-all="store.deselectAllLocal"
      @enable="(path: string) => store.deleteVpks(true, [path])"
      @disable="(path: string) => store.deleteVpks(true, [path])"
      @delete="(path: string) => store.deleteVpks(false, [path])"
      @open-folder="store.openAddonsFolder"
    />

    <!-- 创意工坊 VPK 模组区块 -->
    <VpkSection
      :title="$t('workshopFiles')"
      :subtitle="$t('workshopFilesSubtitle')"
      :empty-text="$t('noWorkshopFiles')"
      :files="store.filteredWorkshopVpkFiles"
      :total-count="store.workshopVpkFiles.length"
      :selected-set="store.selectedWorkshopVpks"
      :is-scanning="store.isScanning"
      @toggle-select="store.toggleWorkshopVpkSelection"
      @select-all="store.selectAllWorkshop"
      @deselect-all="store.deselectAllWorkshop"
      @enable="(path: string) => store.deleteVpks(true, [path])"
      @disable="(path: string) => store.deleteVpks(true, [path])"
      @delete="(path: string) => store.deleteVpks(false, [path])"
      @open-folder="store.openWorkshopFolder"
    />
  </div>
</template>

<script setup lang="ts">
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import VpkSection from '@/components/VpkSection.vue'
import SvgIcon from '@/components/common/SvgIcon.vue'

const store = useVpkManagerStore()
const { t: $t } = useI18n()

function toggleSort(key: 'name' | 'size' | 'date') {
  if (store.sortKey === key) {
    store.sortAsc = !store.sortAsc
  } else {
    store.sortKey = key
    store.sortAsc = true
  }
}
</script>

<style scoped>
.vpk-files-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.filter-toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  padding: 8px 12px;
  box-shadow: var(--shadow-sm);
}

/* 搜索框 */
.search-box {
  flex: 1;
  min-width: 200px;
  display: flex;
  align-items: center;
  background: rgba(0, 0, 0, 0.35);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 0 10px;
  gap: 8px;
  height: 34px;
  transition: all var(--transition-fast);
}

.search-box:focus-within {
  border-color: var(--border-focus);
  box-shadow: 0 0 0 2px var(--accent-primary-glow);
  background: rgba(0, 0, 0, 0.5);
}

.search-icon {
  color: var(--text-dim);
}

.search-input {
  flex: 1;
  background: transparent !important;
  border: none !important;
  padding: 0 !important;
  color: var(--text-main);
  font-size: 13px;
}

.search-input:focus {
  box-shadow: none !important;
}

.search-clear-btn {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  color: var(--text-dim);
  display: flex;
  align-items: center;
  justify-content: center;
}

.search-clear-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

/* 状态筛选 Chips */
.filter-chips {
  display: flex;
  align-items: center;
  gap: 6px;
}

.chip-btn {
  appearance: none;
  -webkit-appearance: none;
  height: 28px;
  padding: 0 10px;
  border-radius: var(--radius-full);
  background: #162032;
  border: 1px solid var(--border-medium);
  font-size: 12px;
  font-weight: 500;
  color: #cbd5e1;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.chip-btn:hover {
  background: #1e2c44;
  color: #fff;
  border-color: var(--border-strong);
}

.chip-btn.active {
  background: #0284c7 !important;
  border-color: #38bdf8 !important;
  color: #ffffff !important;
  font-weight: 700;
  box-shadow: 0 1px 6px rgba(2, 132, 199, 0.4);
}

.chip-count {
  font-size: 10px;
  font-weight: 700;
  background: rgba(0, 0, 0, 0.3);
  padding: 1px 6px;
  border-radius: var(--radius-full);
  color: inherit;
}

.chip-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.dot-green {
  background: var(--accent-success);
}

.dot-amber {
  background: var(--accent-warning);
}

.chip-btn.active.active-state {
  background: #059669 !important;
  border-color: #10b981 !important;
  color: #ffffff !important;
  box-shadow: 0 1px 6px rgba(16, 185, 129, 0.4);
}

.chip-btn.active.disabled-state {
  background: #d97706 !important;
  border-color: #f59e0b !important;
  color: #ffffff !important;
  box-shadow: 0 1px 6px rgba(245, 158, 11, 0.4);
}

/* 控制按钮组 */
.controls-right {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-left: auto;
}

/* 排序组 */
.sort-group {
  display: flex;
  background: #101726;
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  padding: 2px;
  gap: 2px;
}

.sort-tab-btn {
  appearance: none;
  -webkit-appearance: none;
  height: 26px;
  padding: 0 10px;
  font-size: 12px;
  font-weight: 500;
  color: #cbd5e1;
  background: transparent;
  border: none;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.sort-tab-btn:hover {
  background: #1c273e;
  color: #ffffff;
}

.sort-tab-btn.active {
  background: #0284c7 !important;
  color: #ffffff !important;
  font-weight: 700;
  box-shadow: 0 1px 4px rgba(2, 132, 199, 0.4);
}

.sort-indicator {
  font-size: 9px;
  color: inherit;
}

/* 视图模式切换 */
.view-mode-toggle {
  display: flex;
  background: #101726;
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  padding: 2px;
  gap: 2px;
}

.view-btn {
  appearance: none;
  -webkit-appearance: none;
  width: 28px;
  height: 26px;
  border-radius: 4px;
  background: transparent;
  border: none;
  color: #cbd5e1;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.view-btn:hover {
  background: #1c273e;
  color: #ffffff;
}

.view-btn.active {
  background: #0284c7 !important;
  color: #ffffff !important;
  box-shadow: 0 1px 4px rgba(2, 132, 199, 0.4);
}

/* 模组列表工具栏刷新按钮 */
.list-refresh-btn {
  appearance: none;
  -webkit-appearance: none;
  width: 32px;
  height: 32px;
  border-radius: var(--radius-md);
  background: #101726;
  border: 1px solid var(--border-medium);
  color: #cbd5e1;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.list-refresh-btn:hover:not(:disabled) {
  background: #1c273e;
  border-color: rgba(6, 182, 212, 0.5);
  color: #22d3ee;
  box-shadow: 0 0 8px rgba(6, 182, 212, 0.2);
}

.list-refresh-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
</style>
