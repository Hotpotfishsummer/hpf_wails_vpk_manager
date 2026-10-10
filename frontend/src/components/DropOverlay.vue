<template>
  <transition name="drop-overlay-fade">
    <div v-if="props.show" class="drop-overlay-container">
      <div class="drop-overlay-backdrop"></div>
      <div class="drop-overlay-box" :class="{ 'no-directory': !props.hasDirectory }">
        <div class="drop-icon-wrapper">
          <SvgIcon
            :name="props.hasDirectory ? 'import' : 'alert-triangle'"
            :size="52"
            class="drop-icon"
          />
        </div>

        <h2 class="drop-title">
          {{ props.hasDirectory ? $t('dropFilesToInstall') : $t('dropFilesNoDirectory') }}
        </h2>

        <p class="drop-desc">
          {{ props.hasDirectory ? $t('dropFilesHint') : $t('needDirectoryFirst') }}
        </p>

        <div v-if="props.hasDirectory" class="drop-badges">
          <span class="badge-item">.ZIP</span>
          <span class="badge-item">.7Z</span>
          <span class="badge-item">.TAR.ZST</span>
          <span class="badge-item">.VPK</span>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import SvgIcon from '@/components/common/SvgIcon.vue'

const { t: $t } = useI18n()

const props = defineProps<{
  show: boolean
  hasDirectory: boolean
}>()
</script>

<style scoped>
.drop-overlay-container {
  position: fixed;
  inset: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none; /* 穿透鼠标事件，保证底层 drop 正常冒泡到 window */
}

.drop-overlay-backdrop {
  position: absolute;
  inset: 0;
  background: rgba(15, 23, 42, 0.78);
  backdrop-filter: blur(10px);
}

.drop-overlay-box {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 40px 60px;
  max-width: 580px;
  background: rgba(30, 41, 59, 0.9);
  border: 2px dashed #38bdf8;
  border-radius: 20px;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6), 0 0 35px rgba(56, 189, 248, 0.35);
  animation: drop-pulse 1.8s ease-in-out infinite alternate;
}

.drop-overlay-box.no-directory {
  border-color: #f59e0b;
  box-shadow: 0 20px 50px rgba(0, 0, 0, 0.6), 0 0 35px rgba(245, 158, 11, 0.35);
}

.drop-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 90px;
  height: 90px;
  border-radius: 50%;
  background: rgba(56, 189, 248, 0.12);
  color: #38bdf8;
  margin-bottom: 20px;
  box-shadow: 0 0 20px rgba(56, 189, 248, 0.25);
  animation: drop-bounce 1.4s ease-in-out infinite alternate;
}

.drop-overlay-box.no-directory .drop-icon-wrapper {
  background: rgba(245, 158, 11, 0.12);
  color: #f59e0b;
  box-shadow: 0 0 20px rgba(245, 158, 11, 0.25);
}

.drop-title {
  margin: 0 0 10px 0;
  font-size: 24px;
  font-weight: 700;
  color: #f8fafc;
  letter-spacing: 0.5px;
}

.drop-desc {
  margin: 0 0 24px 0;
  font-size: 14px;
  color: #94a3b8;
  line-height: 1.6;
}

.drop-badges {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
}

.badge-item {
  display: inline-block;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  background: rgba(56, 189, 248, 0.15);
  color: #38bdf8;
  border: 1px solid rgba(56, 189, 248, 0.3);
  letter-spacing: 0.5px;
}

@keyframes drop-pulse {
  0% {
    transform: scale(0.98);
  }
  100% {
    transform: scale(1.02);
  }
}

@keyframes drop-bounce {
  0% {
    transform: translateY(0);
  }
  100% {
    transform: translateY(-8px);
  }
}

/* 过渡动画 */
.drop-overlay-fade-enter-active,
.drop-overlay-fade-leave-active {
  transition: opacity 0.25s cubic-bezier(0.16, 1, 0.3, 1), transform 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.drop-overlay-fade-enter-from,
.drop-overlay-fade-leave-to {
  opacity: 0;
  transform: scale(0.96);
}
</style>
