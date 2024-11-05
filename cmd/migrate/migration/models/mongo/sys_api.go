package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SysApi struct {
	Id     primitive.ObjectID `bson:"_id" json:"id"`        // 主键编码
	Handle string             `bson:"handle" json:"handle"` // handle
	Title  string             `bson:"title" json:"title"`   // 标题
	Path   string             `bson:"path" json:"path"`     // 地址
	Type   string             `bson:"type" json:"type"`     // 接口类型
	Action string             `bson:"action" json:"action"` // 请求类型
	ModelTime
	ControlBy
}

func (SysApi) TableName() string {
	return "sys_api"
}
