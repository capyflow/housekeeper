import request from './request'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
}

export const userApi = {
  // 用户登录
  login(data: LoginRequest): Promise<LoginResponse> {
    return request.post('/user/login', data)
  }
}
