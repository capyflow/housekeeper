import { ref } from 'vue'

export type ToastType = 'success' | 'error' | 'warning' | 'info'

interface ToastMessage {
  id: number
  message: string
  type: ToastType
  duration: number
}

const toasts = ref<ToastMessage[]>([])
let idCounter = 0

export function useToast() {
  const show = (message: string, type: ToastType = 'info', duration = 3000) => {
    const id = ++idCounter
    const toast: ToastMessage = {
      id,
      message,
      type,
      duration
    }

    toasts.value.push(toast)

    setTimeout(() => {
      const index = toasts.value.findIndex(t => t.id === id)
      if (index > -1) {
        toasts.value.splice(index, 1)
      }
    }, duration)
  }

  const success = (message: string, duration = 3000) => {
    show(message, 'success', duration)
  }

  const error = (message: string, duration = 3000) => {
    show(message, 'error', duration)
  }

  const warning = (message: string, duration = 3000) => {
    show(message, 'warning', duration)
  }

  const info = (message: string, duration = 3000) => {
    show(message, 'info', duration)
  }

  return {
    toasts,
    show,
    success,
    error,
    warning,
    info
  }
}
