import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import router from '@/router'

const instance: AxiosInstance = axios.create({
  baseURL: '/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
instance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
instance.interceptors.response.use(
  (response: AxiosResponse) => {
    const data = response.data
    // 后端返回格式: {body: {...}, head: {ec, em, ...}, path: "..."}
    if (data && typeof data === 'object' && 'body' in data) {
      return data.body
    }
    // 兼容其他格式
    if (data && typeof data === 'object' && 'data' in data) {
      return data.data
    }
    // 否则直接返回整个响应数据
    return data
  },
  (error) => {
    if (error.response?.status === 401) {
      // 清除 token
      localStorage.removeItem('token')
      // 使用 router 跳转到登录页面
      router.push('/login')
    }
    return Promise.reject(error)
  }
)

export default instance
