<template>
  <transition name="fade">
    <div v-if="store.confirmState.show" class="modal-overlay" @click.self="resolveConfirm(false)">
      <div class="modal-content confirm-dialog">
        <div class="modal-body">
          <div class="confirm-icon-box">
            <SvgIcon name="alert-triangle" :size="24" class="warning-icon" />
          </div>
          <p class="confirm-message">{{ store.confirmState.message }}</p>
        </div>
        <div class="modal-footer">
          <button class="btn-cancel" @click="resolveConfirm(false)">{{ $t('cancel') }}</button>
          <button class="btn-confirm-danger" @click="resolveConfirm(true)">{{ $t('confirm') }}</button>
        </div>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { useVpkManagerStore } from '@/stores/vpk-manager'
import { useI18n } from 'vue-i18n'
import SvgIcon from '@/components/common/SvgIcon.vue'

const store = useVpkManagerStore()
const { t: $t } = useI18n()

function resolveConfirm(value: boolean) {
  store.resolveConfirm(value)
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(4, 7, 14, 0.85);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
  padding: 20px;
}

.confirm-dialog {
  background: #151e30 !important;
  border: 1.5px solid #384c6c;
  border-radius: var(--radius-lg);
  max-width: 440px;
  width: 100%;
  box-shadow: 0 25px 70px rgba(0, 0, 0, 0.9), 0 0 0 1px rgba(244, 63, 94, 0.25), 0 0 35px rgba(244, 63, 94, 0.15);
  overflow: hidden;
  animation: modal-scale 0.2s cubic-bezier(0.16, 1, 0.3, 1);
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

.modal-body {
  padding: 24px 24px 18px;
  display: flex;
  align-items: flex-start;
  gap: 16px;
  background: #101826;
}

.confirm-icon-box {
  width: 42px;
  height: 42px;
  border-radius: var(--radius-md);
  background: rgba(244, 63, 94, 0.18);
  border: 1px solid rgba(244, 63, 94, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-danger);
  flex-shrink: 0;
}

.confirm-message {
  color: #f8fafc;
  font-size: 14px;
  line-height: 1.6;
  white-space: pre-line;
  margin-top: 2px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 14px 24px;
  background: #0d131f;
  border-top: 1px solid #243248;
}

.btn-cancel {
  appearance: none;
  -webkit-appearance: none;
  height: 34px;
  padding: 0 18px;
  border-radius: var(--radius-sm);
  background: #1e293b;
  border: 1px solid #3b4d6b;
  color: #cbd5e1;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.btn-cancel:hover {
  background: #2b3952;
  color: #ffffff;
  border-color: #4f6488;
}

.btn-confirm-danger {
  appearance: none;
  -webkit-appearance: none;
  height: 34px;
  padding: 0 20px;
  border-radius: var(--radius-sm);
  background: #e11d48;
  border: 1px solid #f43f5e;
  color: #ffffff !important;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all var(--transition-fast);
  box-shadow: 0 2px 8px rgba(225, 29, 72, 0.4);
}

.btn-confirm-danger:hover {
  background: #be123c;
  border-color: #fb7185;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
