import request from './request'
import type {
  CreateRecordReq,
  ListRecordReq,
  ListRecordResp,
  Record,
  UpdateRecordReq
} from '@/types'

export const recordApi = {
  create(data: CreateRecordReq) {
    return request.post<any, { info: Record }>('/record/create', data)
  },

  update(data: UpdateRecordReq) {
    return request.put<any, { Message: string }>('/record/update', data)
  },

  delete(id: string) {
    return request.delete<any, { Message: string }>('/record/delete', { data: { id } })
  },

  info(id: string) {
    return request.post<any, { info: Record }>('/record/info', { id })
  },

  list(params: ListRecordReq) {
    return request.post<any, ListRecordResp>('/record/list', params)
  }
}
