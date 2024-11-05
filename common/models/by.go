package models

import (
	"time"
)

type ControlBy struct {
	CreateBy int `bson:"createBy" json:"createBy"`
	UpdateBy int `bson:"updateBy" json:"updateBy"`
}

// SetCreateBy 设置创建人id
func (e *ControlBy) SetCreateBy(createBy int) {
	e.CreateBy = createBy
}

// SetUpdateBy 设置修改人id
func (e *ControlBy) SetUpdateBy(updateBy int) {
	e.UpdateBy = updateBy
}

type Model struct {
	Id int `json:"id"`
}

type ModelTime struct {
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
	DeletedAt time.Time `bson:"deletedAt" json:"deletedAt"`
}
