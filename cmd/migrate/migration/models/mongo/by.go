package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ControlBy 记录创建和更新者的信息
type ControlBy struct {
	CreateBy int `bson:"createBy" json:"createBy" comment:"创建者" `
	UpdateBy int `bson:"updateBy" json:"updateBy" comment:"更新者"`
}

// Model 定义了一个包含主键的模型
type Model struct {
	Id primitive.ObjectID `bson:"_id,omitempty" json:"id" comment:"主键编码"` // 使用 ObjectID 作为 MongoDB 的主键
}

// ModelTime 包含创建、更新时间和删除时间
type ModelTime struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt" comment:"创建时间"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt" comment:"最后更新时间"`
}
