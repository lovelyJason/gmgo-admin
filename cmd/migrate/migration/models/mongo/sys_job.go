package models

type SysJob struct {
	JobId          int    `bson:"jobId" json:"jobId"`                   // 编码
	JobName        string `bson:"jobName" json:"jobName"`               // 名称
	JobGroup       string `bson:"jobGroup" json:"jobGroup"`             // 任务分组
	JobType        int    `bson:"jobType" json:"jobType"`               // 任务类型
	CronExpression string `bson:"cronExpression" json:"cronExpression"` // cron表达式
	InvokeTarget   string `bson:"invokeTarget" json:"invokeTarget"`     // 调用目标
	Args           string `bson:"args" json:"args"`                     // 目标参数
	MisfirePolicy  int    `bson:"misfirePolicy" json:"misfirePolicy"`   // 执行策略
	Concurrent     int    `bson:"concurrent" json:"concurrent"`         // 是否并发
	Status         int    `bson:"status" json:"status"`                 // 状态
	EntryId        int    `bson:"entry_id" json:"entry_id"`             // job启动时返回的id
	ModelTime
	ControlBy
}

func (SysJob) TableName() string {
	return "sys_job"
}
