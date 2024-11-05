package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type DictData struct {
	DictCode  primitive.ObjectID `bson:"_id" json:"dictCode"`        // 字典编码
	DictSort  int                `bson:"dictSort" json:"dictSort"`   // 显示顺序
	DictLabel string             `bson:"dictLabel" json:"dictLabel"` // 数据标签
	DictValue string             `bson:"dictValue" json:"dictValue"` // 数据键值
	DictType  string             `bson:"dictType" json:"dictType"`   // 字典类型
	CssClass  string             `bson:"cssClass" json:"cssClass"`   // CSS 类
	ListClass string             `bson:"listClass" json:"listClass"` // 列表 CSS 类
	IsDefault string             `bson:"isDefault" json:"isDefault"` // 是否默认
	Status    int                `bson:"status" json:"status"`       // 状态
	Default   string             `bson:"default" json:"default"`     // 默认值
	Remark    string             `bson:"remark" json:"remark"`       // 备注
	ControlBy
	ModelTime
}

func (DictData) TableName() string {
	return "sys_dict_data"
}
