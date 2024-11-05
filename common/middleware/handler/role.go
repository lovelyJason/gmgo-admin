package handler

import "go-admin/common/models"

type SysRole struct {
	RoleId    int    `bson:"roleId" json:"roleId" gorm:"primaryKey;autoIncrement"` // 角色编码
	RoleName  string `bson:"roleName" json:"roleName" gorm:"size:128;"`            // 角色名称
	Status    string `bson:"status" json:"status" gorm:"size:4;"`                  //
	RoleKey   string `bson:"roleKey" json:"roleKey" gorm:"size:128;"`              //角色代码
	RoleSort  int    `bson:"roleSort" json:"roleSort" gorm:""`                     //角色排序
	Flag      string `bson:"flag" json:"flag" gorm:"size:128;"`                    //
	Remark    string `bson:"remark" json:"remark" gorm:"size:255;"`                //备注
	Admin     bool   `bson:"admin" json:"admin" gorm:"size:4;"`
	DataScope string `bson:"dataScope" json:"dataScope" gorm:"size:128;"`
	Params    string `bson:"params" json:"params" gorm:"-"`
	MenuIds   []int  `bson:"menuIds" json:"menuIds" gorm:"-"`
	DeptIds   []int  `bson:"DeptIds" json:"deptIds" gorm:"-"`
	models.ControlBy
	models.ModelTime
}

func (SysRole) TableName() string {
	return "sys_role"
}
