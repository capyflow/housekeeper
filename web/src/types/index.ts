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

// Record types (deprecated, use CheckInTarget instead)
export interface Record {
  id: string
  title: string
  create_ts: number
  record_ts: number[]
}

export interface CreateRecordReq {
  title: string
  record_ts: number[]
}

export interface UpdateRecordReq {
  id: string
  title: string
  record_ts: number[]
}

export interface ListRecordReq {
  page: number
  page_size: number
}

export interface ListRecordResp {
  total: number
  page: number
  list: Record[]
}

// CheckIn Target types
export interface TimeRange {
  enabled: boolean
  start: number // minutes from 0-1439
  end: number   // minutes from 0-1439
}

export interface CheckInTarget {
  id: string
  title: string
  description: string
  create_time: number
  update_time: number
  owner: string
  time_range: TimeRange
  total_days: number
  check_in_counts: { [key: string]: number } // YYYY-MM-DD -> count
}

export interface CreateCheckInTargetReq {
  title: string
  description: string
  time_range: TimeRange
}

export interface UpdateCheckInTargetReq {
  id: string
  title: string
  description: string
  time_range: TimeRange
}

export interface DeleteCheckInTargetReq {
  id: string
}

export interface CheckInTargetInfoReq {
  id: string
}

export interface ListCheckInTargetReq {
  page: number
  page_size: number
  owner?: string
}

export interface ListCheckInTargetResp {
  total: number
  page: number
  list: CheckInTarget[]
}

export interface CheckInRequest {
  id: string
  target_ts?: number // 0 means current time
}

export interface CheckInResponse {
  success: boolean
  message: string
  target_id: string
  check_in_time: number
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
