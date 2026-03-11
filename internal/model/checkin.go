package model

import "time"

// CheckInTarget 打卡任务
type CheckInTarget struct {
	Id            string   `json:"id" bson:"id"`                // 任务 ID
	Title         string   `json:"title" bson:"title"`          // 任务名称
	Description   string   `json:"description" bson:"description"` // 描述
	CreateTime    int64    `json:"create_time" bson:"create_time"`
	UpdateTime    int64    `json:"update_time" bson:"update_time"`
	Owner         string   `json:"owner" bson:"owner"`                    // 所有者
	TimeRange     TimeRange `json:"time_range" bson:"time_range"`         // 打卡时间段 (空表示全天)
	TotalDays     int      `json:"total_days" bson:"total_days"`         // 总打卡天数
	CheckInCounts map[string]int `json:"check_in_counts" bson:"check_in_counts"` // 每日打卡次数 key: YYYY-MM-DD, value: 次数
}

// TimeRange 时间段
type TimeRange struct {
	Enabled bool `json:"enabled" bson:"enabled"` // 是否启用时间段限制
	Start   int  `json:"start" bson:"start"`     // 开始时间 (分钟从 0-1439)
	End     int  `json:"end" bson:"end"`         // 结束时间 (分钟从 0-1439)
}

// CreateCheckInTargetReq 创建打卡任务请求
type CreateCheckInTargetReq struct {
	Title       string    `json:"title"`        // 任务名称
	Description string    `json:"description"`  // 描述
	TimeRange   TimeRange `json:"time_range"`   // 打卡时间段
}

// UpdateCheckInTargetReq 更新打卡任务请求
type UpdateCheckInTargetReq struct {
	Id          string    `json:"id"`           // 任务 ID
	Title       string    `json:"title"`        // 任务名称
	Description string    `json:"description"`  // 描述
	TimeRange   TimeRange `json:"time_range"`   // 打卡时间段
}

// DeleteCheckInTargetReq 删除打卡任务请求
type DeleteCheckInTargetReq struct {
	Id string `json:"id"` // 任务 ID
}

// CheckInTargetInfoReq 查询打卡任务详情请求
type CheckInTargetInfoReq struct {
	Id string `json:"id"` // 任务 ID
}

// ListCheckInTargetReq 分页查询打卡任务请求
type ListCheckInTargetReq struct {
	Page     int    `json:"page"`      // 页码
	PageSize int    `json:"page_size"` // 每页数量
	Owner    string `json:"owner"`     // 所有者 (可选)
}

// ListCheckInTargetResp 分页查询响应
type ListCheckInTargetResp struct {
	Total int64               `json:"total"` // 总数
	Page  int                 `json:"page"`  // 当前页
	List  []CheckInTarget     `json:"list"`  // 列表
}

// CheckInRequest 打卡请求
type CheckInRequest struct {
	Id       string `json:"id"`       // 任务 ID
	TargetTs int64  `json:"target_ts"` // 目标打卡时间戳 (0 表示当前时间)
}

// CheckInResponse 打卡响应
type CheckInResponse struct {
	Success     bool   `json:"success"`      // 是否成功
	Message     string `json:"message"`      // 消息
	TargetId    string `json:"target_id"`    // 任务 ID
	CheckInTime int64  `json:"check_in_time"` // 打卡时间戳
	TotalDays   int    `json:"total_days"`   // 总打卡天数
}

// GetDateKey 获取日期键 YYYY-MM-DD
func GetDateKey(ts int64) string {
	t := time.UnixMilli(ts)
	return t.Format("2006-01-02")
}

// NormalizeDayStart 归一化到当天开始
func NormalizeDayStart(ts int64) int64 {
	t := time.UnixMilli(ts)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).UnixMilli()
}
