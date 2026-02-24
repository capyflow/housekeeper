package model

// 创建共享看板请求
type CreateShareBoardReq struct {
	DeviceName string `json:"device_name"` // 设备名称
	BoardName  string `json:"board_name"`  // 看板名称
	Content    string `json:"content"`     // 共享内容
	Owner      string `json:"owner"`       // 所有者
}

// 创建共享看板响应
type CreateShareBoardResp struct {
	ShareBoard
}

// 创建笔记请求
type CreateNoteReq struct {
	Title   string `json:"title"`   // 笔记标题
	Content string `json:"content"` // 笔记内容
	Cover   string `json:"cover"`   // 封面fid
	Owner   string `json:"owner"`   // 所有者
}

// 创建笔记响应
type CreateNoteResp struct {
	Note
}

type ShareBoard struct {
	Id         string `json:"id" bson:"id"`                   // 共享看板id
	DeviceName string `json:"device_name" bson:"device_name"` // 设备名称
	BoardName  string `json:"board_name" bson:"board_name"`   // 看板名称
	Content    string `json:"content" bson:"content"`         // 共享内容
	Owner      string `json:"owner" bson:"owner"`             // 所有者
	CreateTime int64  `json:"create_time" bson:"create_time"` // 创建时间
	UpdateTime int64  `json:"update_time" bson:"update_time"` // 更新时间
}

// 更新共享看板请求
type UpdateShareBoardReq struct {
	Id         string `json:"id"`          // 看板ID
	BoardName  string `json:"board_name"`  // 看板名称
	Content    string `json:"content"`     // 共享内容
	DeviceName string `json:"device_name"` // 设备名称
}

// 更新笔记请求
type UpdateNoteReq struct {
	Id      string `json:"id"`      // 笔记ID
	Title   string `json:"title"`   // 笔记标题
	Content string `json:"content"` // 笔记内容
	Cover   string `json:"cover"`   // 封面fid
}

// 删除共享看板请求
type DeleteShareBoardReq struct {
	Id string `json:"id"` // 看板ID
}

// 删除笔记请求
type DeleteNoteReq struct {
	Id string `json:"id"` // 笔记ID
}

// 查询笔记详情请求
type NoteInfoReq struct {
	Id string `json:"id"` // 笔记ID
}

// 分页查询共享看板请求
type ListShareBoardReq struct {
	Page     int    `json:"page"`      // 页码，从1开始
	PageSize int    `json:"page_size"` // 每页数量
	Owner    string `json:"owner"`     // 所有者（可选，用于过滤）
}

// 分页查询笔记请求
type ListNoteReq struct {
	Page     int    `json:"page"`      // 页码，从1开始
	PageSize int    `json:"page_size"` // 每页数量
	Owner    string `json:"owner"`     // 所有者（可选，用于过滤）
}

// 分页查询共享看板响应
type ListShareBoardResp struct {
	Total int64        `json:"total"` // 总数
	Page  int          `json:"page"`  // 当前页码
	List  []ShareBoard `json:"list"`  // 看板列表
}

// 分页查询笔记响应
type ListNoteResp struct {
	Total int64  `json:"total"` // 总数
	Page  int    `json:"page"`  // 当前页码
	List  []Note `json:"list"`  // 笔记列表
}

// 笔记
type Note struct {
	Id       string `json:"id" bson:"id"`               // 笔记ID
	Title    string `json:"title" bson:"title"`         // 笔记标题
	Content  string `json:"content" bson:"content"`     // 笔记内容
	Owner    string `json:"owner" bson:"owner"`         // 所有者
	CreateTs int64  `json:"create_ts" bson:"create_ts"` // 创建时间
	UpdateTs int64  `json:"update_ts" bson:"update_ts"` // 更新时间
	ModifyTs int64  `json:"modify_ts" bson:"modify_ts"` // 修改时间
	Cover    string `json:"cover" bson:"cover"`         // 封面fid
}
