package handler

import (
	"gmgo-admin/common/models"
	"gorm.io/gorm"
)

type SysUser struct {
	UserId   int    `bson:"userId" gorm:"primaryKey;autoIncrement;comment:编码"  json:"userId"`
	Username string `bson:"username" json:"username" gorm:"size:64;comment:用户名"`
	Password string `bson:"password" json:"-" gorm:"size:128;comment:密码"`
	NickName string `bson:"nickName" json:"nickName" gorm:"size:128;comment:昵称"`
	Phone    string `bson:"phone" json:"phone" gorm:"size:11;comment:手机号"`
	RoleId   int    `bson:"roleId" json:"roleId" gorm:"size:20;comment:角色ID"`
	Salt     string `bson:"salt" json:"-" gorm:"size:255;comment:加盐"`
	Avatar   string `bson:"avatar" json:"avatar" gorm:"size:255;comment:头像"`
	Sex      string `bson:"ses" json:"sex" gorm:"size:255;comment:性别"`
	Email    string `bson:"email" json:"email" gorm:"size:128;comment:邮箱"`
	DeptId   int    `bson:"deptId" json:"deptId" gorm:"size:20;comment:部门"`
	PostId   int    `bson:"postId" json:"postId" gorm:"size:20;comment:岗位"`
	Remark   string `bson:"remark" json:"remark" gorm:"size:255;comment:备注"`
	Status   string `bson:"status" json:"status" gorm:"size:4;comment:状态"`
	DeptIds  []int  `bson:"deptIds" json:"deptIds" gorm:"-"`
	PostIds  []int  `bson:"postIds" json:"postIds" gorm:"-"`
	RoleIds  []int  `bson:"roleIds" json:"roleIds" gorm:"-"`
	//Dept     *SysDept `json:"dept"`
	models.ControlBy
	models.ModelTime
}

func (*SysUser) TableName() string {
	return "sys_user"
}

func (e *SysUser) AfterFind(_ *gorm.DB) error {
	e.DeptIds = []int{e.DeptId}
	e.PostIds = []int{e.PostId}
	e.RoleIds = []int{e.RoleId}
	return nil
}
