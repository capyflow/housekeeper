package model

// 收藏，通常是一个地址
type Collection struct {
	Id             string `json:"id" bson:"id"`                           // id
	Title          string `json:"title" bson:"title"`                     // 标题
	Description    string `json:"description" bson:"description"`         // 描述
	Url            string `json:"url" bson:"url"`                         // url
	CollectionType string `json:"collection_type" bson:"collection_type"` // 收藏类型
}
