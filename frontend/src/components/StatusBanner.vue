<template>
  <div :class="['status-banner', `status-${store.notification.type}`]">
    <div class="status-content">
      <div class="status-icon-box">
        <SvgIcon
          :name="store.notification.type === 'success' ? 'check-circle' : (store.notification.type === 'error' ? 'alert-triangle' : 'info')"
          :size="16"
        />
      </div>
      <span class="status-text">{{ store.notification.message }}</span>
      <button @click="store.notification.message = ''" class="status-close-btn" title="关闭通知">
        <SvgIcon name="x" :size="14" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useVpkManagerStore } from '@/stores/vpk-manager'
import SvgIcon from '@/components/common/SvgIcon.vue'

const store = useVpkManagerStore()
</script>

<style scoped>
.status-banner {
  border-bottom: 1px solid var(--border-subtle);
  animation: slide-down 0.2s ease;
  z-index: 90;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
}

.status-success {
  background: rgba(16, 185, 129, 0.15);
  border-color: rgba(16, 185, 129, 0.3);
  color: var(--accent-success);
}

.status-error {
  background: rgba(244, 63, 94, 0.15);
  border-color: rgba(244, 63, 94, 0.3);
  color: var(--accent-danger);
}

.status-info {
  background: rgba(56, 189, 248, 0.15);
  border-color: rgba(56, 189, 248, 0.3);
  color: var(--accent-primary);
}

.status-content {
  max-width: 1360px;
  margin: 0 auto;
  padding: 10px 24px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-icon-box {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.status-text {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
  color: #fff;
  line-height: 1.4;
}

.status-close-btn {
  width: 24px;
  height: 24px;
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--transition-fast);
}

.status-close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

@keyframes slide-down {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
