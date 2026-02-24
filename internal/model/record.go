package model

// 创建打卡记录请求
type CreateRecordReq struct {
	Title    string  `json:"title"`     // 标题
	RecordTs []int64 `json:"record_ts"` // 打卡时间戳
}

// 创建打卡记录响应
type CreateRecordResp struct {
	Record
}

// 更新打卡记录请求（仅支持修改标题和打卡时间）
type UpdateRecordReq struct {
	Id       string  `json:"id"`        // 打卡记录ID
	Title    string  `json:"title"`     // 标题
	RecordTs []int64 `json:"record_ts"` // 打卡时间戳
}

// 删除打卡记录请求
type DeleteRecordReq struct {
	Id string `json:"id"` // 打卡记录ID
}

// 查询打卡记录详情请求
type RecordInfoReq struct {
	Id string `json:"id"` // 打卡记录ID
}

// 分页查询打卡记录请求
type ListRecordReq struct {
	Page     int `json:"page"`      // 页码，从1开始
	PageSize int `json:"page_size"` // 每页数量
}

// 分页查询打卡记录响应
type ListRecordResp struct {
	Total int64    `json:"total"` // 总数
	Page  int      `json:"page"`  // 当前页码
	List  []Record `json:"list"`  // 打卡记录列表
}

// 打卡记录
type Record struct {
	Id       string  `json:"id" bson:"id"`
	Title    string  `json:"title" bson:"title"`
	CreateTs int64   `json:"create_ts" bson:"create_ts"`
	RecordTs []int64 `json:"record_ts" bson:"record_ts"` // 打卡时间戳
}
