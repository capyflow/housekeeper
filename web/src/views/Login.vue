<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h1 class="login-title">Housekeeper</h1>
        <p class="login-subtitle">管理后台登录</p>
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

      <div v-if="error" class="error-message">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userApi } from '@/api/user'

const router = useRouter()
const userStore = useUserStore()

const formData = reactive({
  username: '',
  password: '',
  remember: false
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  error.value = ''
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
      error.value = '登录响应格式错误'
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

    // 跳转到首页
    await router.push('/')
    console.log('跳转完成')
  } catch (err: any) {
    console.error('登录失败:', err)
    error.value = err.response?.data?.msg || '登录失败，请检查用户名和密码'
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-box {
  width: 100%;
  max-width: 400px;
  padding: 40px;
  background-color: var(--bg-primary);
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-title {
  font-size: 28px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 8px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  font-size: 14px;
  font-weight: 500;
  color: var(--text-primary);
}

.checkbox-label {
  flex-direction: row;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.checkbox-label input[type="checkbox"] {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.checkbox-label span {
  font-size: 14px;
  color: var(--text-secondary);
}

.btn-block {
  width: 100%;
  padding: 12px;
  font-size: 16px;
  margin-top: 8px;
}

.btn-block:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  margin-top: 16px;
  padding: 12px;
  background-color: #fee;
  border: 1px solid #fcc;
  border-radius: 6px;
  color: var(--danger-color);
  font-size: 14px;
  text-align: center;
}
</style>
