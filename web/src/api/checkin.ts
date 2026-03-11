import request from './request'
import type {
  CreateCheckInTargetReq,
  ListCheckInTargetReq,
  ListCheckInTargetResp,
  CheckInTarget,
  UpdateCheckInTargetReq,
  CheckInRequest,
  CheckInResponse,
  TimeRange
} from '@/types'

export const checkinApi = {
  // 创建打卡任务
  create(data: CreateCheckInTargetReq) {
    return request.post<any, { info: CheckInTarget }>('/checkin/target/create', data)
  },

  // 更新打卡任务
  update(data: UpdateCheckInTargetReq) {
    return request.put<any, { Message: string }>('/checkin/target/update', data)
  },

  // 删除打卡任务
  delete(id: string) {
    return request.delete<any, { Message: string }>('/checkin/target/delete', { data: { id } })
  },

  // 查询打卡任务详情
  info(id: string) {
    return request.post<any, { info: CheckInTarget }>('/checkin/target/info', { id })
  },

  // 分页查询打卡任务
  list(params: ListCheckInTargetReq) {
    return request.post<any, ListCheckInTargetResp>('/checkin/target/list', params)
  },

  // 打卡
  checkIn(data: CheckInRequest) {
    return request.post<any, CheckInResponse>('/checkin/check_in', data)
  },

  // 获取打卡统计（热力图数据）
  stats(targetId: string, days: number = 365) {
    return request.post<any, { [key: string]: number }>('/checkin/stats', { id: targetId, days })
  },

  // 格式化时间显示
  formatTimeRange(timeRange: TimeRange): string {
    if (!timeRange.enabled) return '全天'
    const startHour = Math.floor(timeRange.start / 60)
    const startMin = timeRange.start % 60
    const endHour = Math.floor(timeRange.end / 60)
    const endMin = timeRange.end % 60
    
    const format = (h: number, m: number) => `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`
    
    if (timeRange.start > timeRange.end) {
      // 跨天情况，如 22:00-06:00
      return `${format(startHour, startMin)} ~ ${format(endHour, endMin)} (次日)`
    }
    return `${format(startHour, startMin)} ~ ${format(endHour, endMin)}`
  }
}
