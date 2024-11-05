package models

type SysMenu struct {
	MenuId     int      `bson:"menuId" json:"menuId"`         // 菜单ID
	MenuName   string   `bson:"menuName" json:"menuName"`     // 菜单名称
	Title      string   `bson:"title" json:"title"`           // 标题
	Icon       string   `bson:"icon" json:"icon"`             // 图标
	Path       string   `bson:"path" json:"path"`             // 路径
	Paths      string   `bson:"paths" json:"paths"`           // 路径集合
	MenuType   string   `bson:"menuType" json:"menuType"`     // 菜单类型
	Action     string   `bson:"action" json:"action"`         // 操作
	Permission string   `bson:"permission" json:"permission"` // 权限
	ParentId   int      `bson:"parentId" json:"parentId"`     // 上级菜单ID
	NoCache    bool     `bson:"noCache" json:"noCache"`       // 是否不缓存
	Breadcrumb string   `bson:"breadcrumb" json:"breadcrumb"` // 面包屑
	Component  string   `bson:"component" json:"component"`   // 组件
	Sort       int      `bson:"sort" json:"sort"`             // 排序
	Visible    string   `bson:"visible" json:"visible"`       // 可见性
	IsFrame    string   `bson:"isFrame" json:"isFrame"`       // 是否为框架
	SysApi     []SysApi `bson:"sysApi" json:"sysApi"`         // 关联的API
	ModelTime
	ControlBy
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
