package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SysRoleDept struct {
	RoleId primitive.ObjectID `bson:"roleId" json:"roleId"` // 角色ID
	DeptId primitive.ObjectID `bson:"deptId" json:"deptId"` // 部门ID
}

func (SysRoleDept) TableName() string {
	return "sys_role_dept"
}
