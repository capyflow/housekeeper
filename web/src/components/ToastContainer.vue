<template>
  <Teleport to="body">
    <div class="toast-container">
      <TransitionGroup name="toast-list">
        <div
          v-for="toast in toasts"
          :key="toast.id"
          :class="['toast', `toast-${toast.type}`]"
        >
          <div class="toast-icon">
            <svg v-if="toast.type === 'success'" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
            <svg v-else-if="toast.type === 'error'" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
            <svg v-else-if="toast.type === 'warning'" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
            </svg>
            <svg v-else class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          <div class="toast-message">{{ toast.message }}</div>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { useToast } from '@/composables/useToast'

const { toasts } = useToast()
</script>

<style scoped>
.toast-container {
  position: fixed;
  top: 24px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 12px;
  pointer-events: none;
}

.toast {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background-color: var(--bg-primary);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  min-width: 280px;
  max-width: 480px;
  border: 1px solid var(--border-color);
  backdrop-filter: blur(10px);
  pointer-events: auto;
}

.toast-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.toast-icon .icon {
  width: 18px;
  height: 18px;
}

.toast-success {
  border-left: 4px solid var(--success-color);
}

.toast-success .toast-icon {
  background-color: rgba(16, 185, 129, 0.1);
  color: var(--success-color);
}

.toast-error {
  border-left: 4px solid var(--danger-color);
}

.toast-error .toast-icon {
  background-color: rgba(239, 68, 68, 0.1);
  color: var(--danger-color);
}

.toast-warning {
  border-left: 4px solid var(--warning-color);
}

.toast-warning .toast-icon {
  background-color: rgba(245, 158, 11, 0.1);
  color: var(--warning-color);
}

.toast-info {
  border-left: 4px solid var(--info-color);
}

.toast-info .toast-icon {
  background-color: rgba(59, 130, 246, 0.1);
  color: var(--info-color);
}

.toast-message {
  flex: 1;
  font-size: 15px;
  font-weight: 500;
  color: var(--text-primary);
  line-height: 1.5;
}

/* 列表过渡动画 */
.toast-list-move,
.toast-list-enter-active,
.toast-list-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toast-list-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.toast-list-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-list-leave-active {
  position: absolute;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .toast-container {
    top: 16px;
    left: 16px;
    right: 16px;
    transform: none;
    width: auto;
  }

  .toast {
    min-width: auto;
    max-width: 100%;
    padding: 12px 16px;
  }

  .toast-message {
    font-size: 14px;
  }

  .toast-icon {
    width: 20px;
    height: 20px;
  }

  .toast-icon .icon {
    width: 16px;
    height: 16px;
  }
}

@media (max-width: 480px) {
  .toast-container {
    top: 12px;
    left: 12px;
    right: 12px;
  }

  .toast {
    padding: 10px 14px;
    gap: 10px;
  }

  .toast-message {
    font-size: 13px;
  }
}
</style>
