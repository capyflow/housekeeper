<template>
  <Teleport to="body">
    <Transition name="toast-slide">
      <div v-if="visible" :class="['toast', `toast-${type}`]">
        <div class="toast-icon">
          <svg v-if="type === 'success'" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
          <svg v-else-if="type === 'error'" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
          <svg v-else-if="type === 'warning'" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
          </svg>
          <svg v-else class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
        </div>
        <div class="toast-message">{{ message }}</div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

export interface ToastProps {
  message: string
  type?: 'success' | 'error' | 'warning' | 'info'
  duration?: number
}

const props = withDefaults(defineProps<ToastProps>(), {
  type: 'info',
  duration: 3000
})

const visible = ref(false)
let timer: ReturnType<typeof setTimeout> | null = null

const show = () => {
  visible.value = true
  if (timer) {
    clearTimeout(timer)
  }
  timer = setTimeout(() => {
    visible.value = false
  }, props.duration)
}

const hide = () => {
  visible.value = false
  if (timer) {
    clearTimeout(timer)
    timer = null
  }
}

// 监听 message 变化，自动显示
watch(() => props.message, (newVal) => {
  if (newVal) {
    show()
  }
}, { immediate: true })

defineExpose({
  show,
  hide
})
</script>

<style scoped>
.toast {
  position: fixed;
  top: 24px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 20px;
  background-color: var(--bg-primary);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  z-index: 9999;
  min-width: 280px;
  max-width: 480px;
  border: 1px solid var(--border-color);
  backdrop-filter: blur(10px);
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

/* 过渡动画 */
.toast-slide-enter-active,
.toast-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.toast-slide-enter-from {
  opacity: 0;
  transform: translateX(-50%) translateY(-20px);
}

.toast-slide-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-10px);
}
</style>
