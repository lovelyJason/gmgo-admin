package models

type SysPost struct {
	PostId   int    `bson:"postId" json:"postId"`     // 岗位编号
	PostName string `bson:"postName" json:"postName"` // 岗位名称
	PostCode string `bson:"postCode" json:"postCode"` // 岗位代码
	Sort     int    `bson:"sort" json:"sort"`         // 岗位排序
	Status   int    `bson:"status" json:"status"`     // 状态
	Remark   string `bson:"remark" json:"remark"`     // 描述
	ControlBy
	ModelTime
}

func (SysPost) TableName() string {
	return "sys_post"
}
