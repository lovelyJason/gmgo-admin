package models

type SysTables struct {
	TableId             int    `bson:"tableId" json:"tableId"`                         // 表编码
	TBName              string `bson:"tableName" json:"tableName"`                     // 表名称
	TableComment        string `bson:"tableComment" json:"tableComment"`               // 表备注
	ClassName           string `bson:"className" json:"className"`                     // 类名
	TplCategory         string `bson:"tplCategory" json:"tplCategory"`                 // 模板类别
	PackageName         string `bson:"packageName" json:"packageName"`                 // 包名
	ModuleName          string `bson:"moduleName" json:"moduleName"`                   // go文件名
	ModuleFrontName     string `bson:"moduleFrontName" json:"moduleFrontName"`         // 前端文件名
	BusinessName        string `bson:"businessName" json:"businessName"`               // 业务名称
	FunctionName        string `bson:"functionName" json:"functionName"`               // 功能名称
	FunctionAuthor      string `bson:"functionAuthor" json:"functionAuthor"`           // 功能作者
	PkColumn            string `bson:"pkColumn" json:"pkColumn"`                       // 主键列
	PkGoField           string `bson:"pkGoField" json:"pkGoField"`                     // Go字段
	PkJsonField         string `bson:"pkJsonField" json:"pkJsonField"`                 // JSON字段
	Options             string `bson:"options" json:"options"`                         // 选项
	TreeCode            string `bson:"treeCode" json:"treeCode"`                       // 树形编码
	TreeParentCode      string `bson:"treeParentCode" json:"treeParentCode"`           // 树形父编码
	TreeName            string `bson:"treeName" json:"treeName"`                       // 树形名称
	Tree                bool   `bson:"tree" json:"tree"`                               // 是否树形
	Crud                bool   `bson:"crud" json:"crud"`                               // 是否 CRUD
	Remark              string `bson:"remark" json:"remark"`                           // 备注
	IsDataScope         int    `bson:"isDataScope" json:"isDataScope"`                 // 是否数据范围
	IsActions           int    `bson:"isActions" json:"isActions"`                     // 是否操作
	IsAuth              int    `bson:"isAuth" json:"isAuth"`                           // 是否认证
	IsLogicalDelete     string `bson:"isLogicalDelete" json:"isLogicalDelete"`         // 逻辑删除标志
	LogicalDelete       bool   `bson:"logicalDelete" json:"logicalDelete"`             // 逻辑删除
	LogicalDeleteColumn string `bson:"logicalDeleteColumn" json:"logicalDeleteColumn"` // 逻辑删除列
	ModelTime
	ControlBy
}

func (SysTables) TableName() string {
	return "sys_tables"
}
