package models

import (
	"time"
)

type SysOperaLog struct {
	Model
	Title         string    `bson:"title" json:"title"`                 // 操作模块
	BusinessType  string    `bson:"businessType" json:"businessType"`   // 操作类型
	BusinessTypes string    `bson:"businessTypes" json:"businessTypes"` // 业务类型
	Method        string    `bson:"method" json:"method"`               // 函数
	RequestMethod string    `bson:"requestMethod" json:"requestMethod"` // 请求方式: GET POST PUT DELETE
	OperatorType  string    `bson:"operatorType" json:"operatorType"`   // 操作类型
	OperName      string    `bson:"operName" json:"operName"`           // 操作者
	DeptName      string    `bson:"deptName" json:"deptName"`           // 部门名称
	OperUrl       string    `bson:"operUrl" json:"operUrl"`             // 访问地址
	OperIp        string    `bson:"operIp" json:"operIp"`               // 客户端ip
	OperLocation  string    `bson:"operLocation" json:"operLocation"`   // 访问位置
	OperParam     string    `bson:"operParam" json:"operParam"`         // 请求参数
	Status        string    `bson:"status" json:"status"`               // 操作状态 1:正常 2:关闭
	OperTime      time.Time `bson:"operTime" json:"operTime"`           // 操作时间
	JsonResult    string    `bson:"jsonResult" json:"jsonResult"`       // 返回数据
	Remark        string    `bson:"remark" json:"remark"`               // 备注
	LatencyTime   string    `bson:"latencyTime" json:"latencyTime"`     // 耗时
	UserAgent     string    `bson:"userAgent" json:"userAgent"`         // UA
	CreatedAt     time.Time `bson:"createdAt" json:"createdAt"`         // 创建时间
	UpdatedAt     time.Time `bson:"updatedAt" json:"updatedAt"`         // 最后更新时间
	ControlBy
}

func (SysOperaLog) TableName() string {
	return "sys_opera_log"
}
