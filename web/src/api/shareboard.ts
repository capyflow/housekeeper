import request from './request'
import type {
  ShareBoard,
  CreateShareBoardReq,
  UpdateShareBoardReq,
  ListShareBoardReq,
  ListShareBoardResp
} from '@/types'

export const shareBoardApi = {
  // 创建共享看板
  create(data: CreateShareBoardReq) {
    return request.post<any, { info: ShareBoard }>('/notes/share_board/create', data)
  },

  // 更新共享看板
  update(data: UpdateShareBoardReq) {
    return request.put<any, { Message: string }>('/notes/share_board/update', data)
  },

  // 删除共享看板
  delete(id: string) {
    return request.delete<any, { Message: string }>('/notes/share_board/delete', {
      data: { id }
    })
  },

  // 分页查询共享看板
  list(params: ListShareBoardReq) {
    return request.post<any, ListShareBoardResp>('/notes/share_board/list', params)
  }
}
