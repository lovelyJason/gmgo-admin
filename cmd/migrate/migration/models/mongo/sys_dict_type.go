package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DictType struct {
	DictId   primitive.ObjectID `bson:"_id" json:"dictId"`        // 字典编码
	DictName string             `bson:"dictName" json:"dictName"` // 字典名称
	DictType string             `bson:"dictType" json:"dictType"` // 字典类型
	Status   int                `bson:"status" json:"status"`     // 状态
	Remark   string             `bson:"remark" json:"remark"`     // 备注
	ControlBy
	ModelTime
}

func (DictType) TableName() string {
	return "sys_dict_type"
}
