package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SysUser struct {
	UserId   int    `bson:"userId" json:"userId" comment:"编码"`      // 编码
	Username string `bson:"username" json:"username" comment:"用户名"` // 用户名
	Password string `bson:"password" json:"-" comment:"密码"`         // 密码 (不返回)
	NickName string `bson:"nickName" json:"nickName" comment:"昵称"`  // 昵称
	Phone    string `bson:"phone" json:"phone" comment:"手机号"`       // 手机号
	RoleId   int    `bson:"roleId" json:"roleId" comment:"角色ID"`    // 角色ID
	Salt     string `bson:"-" json:"-" comment:"加盐"`                // 加盐 (不返回)
	Avatar   string `bson:"avatar" json:"avatar" comment:"头像"`      // 头像
	Sex      string `bson:"sex" json:"sex" comment:"性别"`            // 性别
	Email    string `bson:"email" json:"email" comment:"邮箱"`        // 邮箱
	DeptId   int    `bson:"deptId" json:"deptId" comment:"部门"`      // 部门
	PostId   int    `bson:"postId" json:"postId" comment:"岗位"`      // 岗位
	Remark   string `bson:"remark" json:"remark" comment:"备注"`      // 备注
	Status   string `bson:"status" json:"status" comment:"状态"`      // 状态
	ModelTime
	ControlBy
}

func (*SysUser) TableName() string {
	return "sysUser"
}

// Encrypt 加密
func (e *SysUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

func (e *SysUser) BeforeCreate(_ *gorm.DB) error {
	return e.Encrypt()
}
