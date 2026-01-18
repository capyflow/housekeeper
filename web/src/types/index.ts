// ShareBoard types
export interface ShareBoard {
  id: string
  device_name: string
  board_name: string
  content: string
  owner: string
  create_time: number
  update_time: number
}

export interface CreateShareBoardReq {
  device_name: string
  board_name: string
  content: string
  owner: string
}

export interface UpdateShareBoardReq {
  id: string
  board_name: string
  content: string
  device_name: string
}

export interface ListShareBoardReq {
  page: number
  page_size: number
  owner?: string
}

export interface ListShareBoardResp {
  total: number
  page: number
  list: ShareBoard[]
}

// User types
export interface User {
  id: string
  username: string
  nickname?: string
}

// API Response
export interface ApiResponse<T = any> {
  code?: number
  message?: string
  data?: T
}
