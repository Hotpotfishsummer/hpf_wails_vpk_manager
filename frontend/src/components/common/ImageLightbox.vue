<template>
  <teleport to="body">
    <transition name="lightbox-fade">
      <div v-if="isOpen" class="lightbox-overlay" @click="close">
        <div class="lightbox-dialog" @click.stop>
          <div class="lightbox-header">
            <h3 class="lightbox-title">{{ title || '模组缩略图' }}</h3>
            <button class="lightbox-close-btn" @click="close" title="关闭 (Esc)">
              <SvgIcon name="x" :size="18" />
            </button>
          </div>
          <div class="lightbox-image-wrap">
            <img :src="src" :alt="title" class="lightbox-img" />
          </div>
        </div>
      </div>
    </transition>
  </teleport>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import SvgIcon from '@/components/common/SvgIcon.vue'

const props = defineProps<{
  isOpen: boolean
  src: string
  title?: string
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

function close() {
  emit('close')
}

function handleKeyDown(e: KeyboardEvent) {
  if (props.isOpen && e.key === 'Escape') {
    close()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
.lightbox-overlay {
  position: fixed;
  inset: 0;
  z-index: 9999;
  background: rgba(4, 7, 13, 0.85);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px;
}

.lightbox-dialog {
  max-width: 90vw;
  max-height: 90vh;
  background: var(--bg-card);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-modal);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  animation: zoom-in 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.lightbox-header {
  padding: 12px 18px;
  background: rgba(0, 0, 0, 0.3);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.lightbox-title {
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.lightbox-close-btn {
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

.lightbox-close-btn:hover {
  background: #ef4444;
  border-color: #ef4444;
  color: #ffffff;
}

.lightbox-image-wrap {
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: auto;
  max-height: calc(90vh - 60px);
}

.lightbox-img {
  max-width: 100%;
  max-height: calc(85vh - 80px);
  object-fit: contain;
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-lg);
}

@keyframes zoom-in {
  from {
    opacity: 0;
    transform: scale(0.92);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.lightbox-fade-enter-active,
.lightbox-fade-leave-active {
  transition: opacity 0.2s ease;
}

.lightbox-fade-enter-from,
.lightbox-fade-leave-to {
  opacity: 0;
}
</style>
