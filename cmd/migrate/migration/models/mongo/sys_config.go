package models

type SysConfig struct {
	Model
	ConfigName  string `bson:"configName" json:"configName"`   // 配置名称
	ConfigKey   string `bson:"configKey" json:"configKey"`     // 配置键
	ConfigValue string `bson:"configValue" json:"configValue"` // 配置值
	ConfigType  string `bson:"configType" json:"configType"`   // 配置类型
	IsFrontend  int    `bson:"isFrontend" json:"isFrontend"`   // 是否前台
	Remark      string `bson:"remark" json:"remark"`           // 备注
	ControlBy
	ModelTime
}

func (SysConfig) TableName() string {
	return "sys_config"
}
