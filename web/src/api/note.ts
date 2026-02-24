import request from './request'
import type {
  CreateNoteReq,
  DeleteNoteReq,
  ListNoteReq,
  ListNoteResp,
  Note,
  UpdateNoteReq
} from '@/types'

export const noteApi = {
  // 创建笔记
  create(data: CreateNoteReq) {
    return request.post<any, { info: Note }>('/notes/note/create', data)
  },

  // 更新笔记
  update(data: UpdateNoteReq) {
    return request.put<any, { Message: string }>('/notes/note/update', data)
  },

  // 删除笔记
  delete(id: string) {
    const payload: DeleteNoteReq = { id }
    return request.delete<any, { Message: string }>('/notes/note/delete', { data: payload })
  },

  // 分页查询笔记
  list(params: ListNoteReq) {
    return request.post<any, ListNoteResp>('/notes/note/list', params)
  }
}
