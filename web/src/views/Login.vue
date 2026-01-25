<template>
  <div class="login-container">
    <div class="theme-toggle-wrapper">
      <button class="theme-toggle" @click="themeStore.toggleTheme()" :title="themeStore.theme === 'light' ? '切换到深色模式' : '切换到浅色模式'">
        <svg v-if="themeStore.theme === 'light'" class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path>
        </svg>
        <svg v-else class="icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path>
        </svg>
      </button>
    </div>
    <div class="login-box">
      <div class="login-header">
        <div class="logo-wrapper">
          <svg class="logo-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"></path>
          </svg>
        </div>
        <h1 class="login-title">Housekeeper</h1>
        <p class="login-subtitle">欢迎回来，请登录您的账户</p>
      </div>

      <form class="login-form" @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="username">用户名</label>
          <input
            id="username"
            v-model="formData.username"
            type="text"
            class="input"
            placeholder="请输入用户名"
            required
          />
        </div>

        <div class="form-group">
          <label for="password">密码</label>
          <input
            id="password"
            v-model="formData.password"
            type="password"
            class="input"
            placeholder="请输入密码"
            required
          />
        </div>

        <div class="form-group">
          <label class="checkbox-label">
            <input v-model="formData.remember" type="checkbox" />
            <span>记住我</span>
          </label>
        </div>

        <button type="submit" class="btn btn-primary btn-block" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useThemeStore } from '@/stores/theme'
import { useToast } from '@/composables/useToast'
import { userApi } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()
const themeStore = useThemeStore()
const toast = useToast()

const formData = reactive({
  username: '',
  password: '',
  remember: false
})

const loading = ref(false)

const handleLogin = async () => {
  loading.value = true

  try {
    // 调用真实的登录API
    const response = await userApi.login({
      username: formData.username,
      password: formData.password
    })

    console.log('登录响应:', response)

    // 检查响应格式
    if (!response || !response.token) {
      console.error('响应格式错误:', response)
      toast.error('登录响应格式错误')
      return
    }

    // 保存token
    userStore.setToken(response.token)
    userStore.setUser({
      id: '1',
      username: formData.username,
      nickname: formData.username
    })

    console.log('Token已保存，准备跳转')

    toast.success('登录成功')

    // 跳转到首页
    await router.push('/')
    console.log('跳转完成')
  } catch (err: any) {
    console.error('登录失败:', err)
    toast.error(err.response?.data?.msg || '登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  width: 100%;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--accent-purple) 100%);
  position: relative;
  overflow: hidden;
}

.login-container::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 80%;
  height: 80%;
  background: radial-gradient(circle, rgba(255,255,255,0.1) 0%, transparent 70%);
  border-radius: 50%;
}

.login-container::after {
  content: '';
  position: absolute;
  bottom: -50%;
  left: -20%;
  width: 80%;
  height: 80%;
  background: radial-gradient(circle, rgba(255,255,255,0.05) 0%, transparent 70%);
  border-radius: 50%;
}

.theme-toggle-wrapper {
  position: absolute;
  top: 24px;
  right: 24px;
  z-index: 10;
}

.theme-toggle {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background-color: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s;
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.theme-toggle:hover {
  background-color: rgba(255, 255, 255, 0.3);
  transform: rotate(15deg) scale(1.05);
}

.theme-toggle .icon {
  width: 24px;
  height: 24px;
}

.login-box {
  width: 100%;
  max-width: 440px;
  padding: 48px;
  background-color: var(--bg-primary);
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  position: relative;
  z-index: 1;
  backdrop-filter: blur(10px);
  border: 1px solid var(--border-color);
}

.login-header {
  text-align: center;
  margin-bottom: 36px;
}

.logo-wrapper {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--accent-purple) 100%);
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8px 16px rgba(59, 130, 246, 0.3);
}

.logo-icon {
  width: 44px;
  height: 44px;
  color: white;
}

.login-title {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--accent-purple) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.login-subtitle {
  font-size: 15px;
  color: var(--text-secondary);
  line-height: 1.5;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.form-group label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  letter-spacing: 0.3px;
}

.checkbox-label {
  flex-direction: row;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: var(--primary-color);
}

.checkbox-label span {
  font-size: 14px;
  color: var(--text-secondary);
  font-weight: 400;
}

.btn-block {
  width: 100%;
  padding: 14px;
  font-size: 16px;
  margin-top: 12px;
  font-weight: 600;
  border-radius: 10px;
}

.btn-block:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
