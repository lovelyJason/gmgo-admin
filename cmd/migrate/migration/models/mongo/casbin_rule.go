package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CasbinRule struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"` // 使用 ObjectID 类型作为 ID，`omitempty` 表示省略零值
	Ptype string             `bson:"ptype" json:"ptype"`                // Ptype 字段
	V0    string             `bson:"v0" json:"v0"`                      // V0 字段
	V1    string             `bson:"v1" json:"v1"`                      // V1 字段
	V2    string             `bson:"v2" json:"v2"`                      // V2 字段
	V3    string             `bson:"v3" json:"v3"`                      // V3 字段
	V4    string             `bson:"v4" json:"v4"`                      // V4 字段
	V5    string             `bson:"v5" json:"v5"`                      // V5 字段
}

// TableName 返回 CasbinRule 对应的 MongoDB 集合名称
func (CasbinRule) TableName() string {
	return "sys_casbin_rule"
}
