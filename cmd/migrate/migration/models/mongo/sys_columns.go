package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SysColumns struct {
	ColumnId           primitive.ObjectID `bson:"_id" json:"columnId"`                          // 主键ID
	TableId            int                `bson:"tableId" json:"tableId"`                       // 表ID
	ColumnName         string             `bson:"columnName" json:"columnName"`                 // 列名
	ColumnComment      string             `bson:"column_comment" json:"columnComment"`          // 列注释
	ColumnType         string             `bson:"column_type" json:"columnType"`                // 列类型
	GoType             string             `bson:"go_type" json:"goType"`                        // Go类型
	GoField            string             `bson:"go_field" json:"goField"`                      // Go字段名
	JsonField          string             `bson:"json_field" json:"jsonField"`                  // JSON字段名
	IsPk               string             `bson:"is_pk" json:"isPk"`                            // 是否主键
	IsIncrement        string             `bson:"is_increment" json:"isIncrement"`              // 是否自增
	IsRequired         string             `bson:"is_required" json:"isRequired"`                // 是否必填
	IsInsert           string             `bson:"is_insert" json:"isInsert"`                    // 是否可插入
	IsEdit             string             `bson:"is_edit" json:"isEdit"`                        // 是否可编辑
	IsList             string             `bson:"is_list" json:"isList"`                        // 是否在列表中显示
	IsQuery            string             `bson:"is_query" json:"isQuery"`                      // 是否可查询
	QueryType          string             `bson:"query_type" json:"queryType"`                  // 查询类型
	HtmlType           string             `bson:"html_type" json:"htmlType"`                    // HTML类型
	DictType           string             `bson:"dict_type" json:"dictType"`                    // 字典类型
	Sort               int                `bson:"sort" json:"sort"`                             // 排序
	List               string             `bson:"list" json:"list"`                             // 列表标志
	Pk                 bool               `bson:"pk" json:"pk"`                                 // 是否主键标志
	Required           bool               `bson:"required" json:"required"`                     // 是否必填标志
	SuperColumn        bool               `bson:"super_column" json:"superColumn"`              // 是否超级列
	UsableColumn       bool               `bson:"usable_column" json:"usableColumn"`            // 是否可用列
	Increment          bool               `bson:"increment" json:"increment"`                   // 是否自增标志
	Insert             bool               `bson:"insert" json:"insert"`                         // 是否可插入标志
	Edit               bool               `bson:"edit" json:"edit"`                             // 是否可编辑标志
	Query              bool               `bson:"query" json:"query"`                           // 是否可查询标志
	Remark             string             `bson:"remark" json:"remark"`                         // 备注
	FkTableName        string             `bson:"fkTableName" json:"fkTableName"`               // 外键表名
	FkTableNameClass   string             `bson:"fkTableNameClass" json:"fkTableNameClass"`     // 外键表类名
	FkTableNamePackage string             `bson:"fkTableNamePackage" json:"fkTableNamePackage"` // 外键表包名
	FkLabelId          string             `bson:"fkLabelId" json:"fkLabelId"`                   // 外键标签ID
	FkLabelName        string             `bson:"fkLabelName" json:"fkLabelName"`               // 外键标签名
	ModelTime
	ControlBy
}

func (SysColumns) TableName() string {
	return "sys_columns"
}
