package models

type TbDemo struct {
	Model
	Name string `bson:"name" json:"name" comment:"名称"`
	ModelTime
	ControlBy
}

func (TbDemo) TableName() string {
	return "tb_demo"
}
