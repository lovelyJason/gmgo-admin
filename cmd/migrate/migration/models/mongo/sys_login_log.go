package models

import (
	"time"
)

type SysLoginLog struct {
	Model
	Username      string           `bson:"username" json:"username"`           // 用户名
	Status        string           `bson:"status" json:"status"`               // 状态
	Ipaddr        string           `bson:"ipaddr" json:"ipaddr"`               // IP地址
	LoginLocation string           `bson:"loginLocation" json:"loginLocation"` // 归属地
	Browser       string           `bson:"browser" json:"browser"`             // 浏览器
	Os            string           `bson:"os" json:"os"`                       // 系统
	Platform      string           `bson:"platform" json:"platform"`           // 固件
	LoginTime     time.Time        `bson:"loginTime" json:"loginTime"`         // 登录时间
	Remark        string           `bson:"remark" json:"remark"`               // 备注
	Msg           string           `bson:"msg" json:"msg"`                     // 信息
	CreatedAt     time.Time        `bson:"createdAt" json:"createdAt"`         // 创建时间
	UpdatedAt     time.Time        `bson:"updatedAt" json:"updatedAt"`         // 最后更新时间
	ControlBy     `bson:",inline"` // 包含 ControlBy 字段
}

func (SysLoginLog) TableName() string {
	return "sys_login_log"
}
