package models

type SysApi struct {
	ApiId  int    `json:"id"`
	Handle string `json:"handle"`
	Title  string `json:"title"`
	Path   string `json:"path"`
	Type   string `json:"type"`
	Action string `json:"action"`
	ModelTime
	ControlBy
}

func (SysApi) TableName() string {
	return "sysApi"
}
