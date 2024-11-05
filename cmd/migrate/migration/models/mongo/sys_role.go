package models

type SysRole struct {
	RoleId    int       `bson:"roleId" json:"roleId"`       // 角色编码
	RoleName  string    `bson:"roleName" json:"roleName"`   // 角色名称
	Status    string    `bson:"status" json:"status"`       // 状态
	RoleKey   string    `bson:"roleKey" json:"roleKey"`     // 角色代码
	RoleSort  int       `bson:"roleSort" json:"roleSort"`   // 角色排序
	Flag      string    `bson:"flag" json:"flag"`           // 标志
	Remark    string    `bson:"remark" json:"remark"`       // 备注
	Admin     bool      `bson:"admin" json:"admin"`         // 是否为管理员
	DataScope string    `bson:"dataScope" json:"dataScope"` // 数据范围
	SysMenu   []SysMenu `bson:"sysMenu" json:"sysMenu"`     // 关联菜单
	ControlBy
	ModelTime
}

func (SysRole) TableName() string {
	return "sys_role"
}
