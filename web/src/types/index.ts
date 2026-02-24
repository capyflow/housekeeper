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
  owner?: string // 可选，后端会从token中解析
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

// Note types
export interface Note {
  id: string
  title: string
  content: string
  owner: string
  create_ts: number
  update_ts: number
  modify_ts: number
  cover: string
}

export interface ListNoteReq {
  page: number
  page_size: number
  owner?: string
}

export interface ListNoteResp {
  total: number
  page: number
  list: Note[]
}

export interface CreateNoteReq {
  title: string
  content: string
  cover?: string
  owner?: string
}

export interface UpdateNoteReq {
  id: string
  title: string
  content: string
  cover?: string
}

export interface DeleteNoteReq {
  id: string
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
