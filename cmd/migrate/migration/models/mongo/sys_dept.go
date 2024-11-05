package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SysDept struct {
	DeptId   primitive.ObjectID `bson:"_id" json:"deptId"`        // 部门编码
	ParentId int                `bson:"parentId" json:"parentId"` // 上级部门
	DeptPath string             `bson:"deptPath" json:"deptPath"` // 部门路径
	DeptName string             `bson:"deptName" json:"deptName"` // 部门名称
	Sort     int                `bson:"sort" json:"sort"`         // 排序
	Leader   string             `bson:"leader" json:"leader"`     // 负责人
	Phone    string             `bson:"phone" json:"phone"`       // 手机
	Email    string             `bson:"email" json:"email"`       // 邮箱
	Status   int                `bson:"status" json:"status"`     // 状态
	ControlBy
	ModelTime
}

func (SysDept) TableName() string {
	return "sys_dept"
}
