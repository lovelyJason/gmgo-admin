package models

import (
	"context"
	"gmgo-admin/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SysDept struct {
	Id               primitive.ObjectID `bson:"_id"`
	DeptId           int                `bson:"deptId" json:"deptId" gorm:"primaryKey;autoIncrement;"` //部门编码
	ParentId         int                `bson:"parentId" json:"parentId" gorm:""`                      //上级部门
	DeptPath         string             `bson:"deptPath" json:"deptPath" gorm:"size:255;"`             //
	DeptName         string             `bson:"deptName" json:"deptName"  gorm:"size:128;"`            //部门名称
	Sort             int                `bson:"sort" json:"sort" gorm:"size:4;"`                       //排序
	Leader           string             `bson:"leader" json:"leader" gorm:"size:128;"`                 //负责人
	Phone            string             `bson:"phone" json:"phone" gorm:"size:11;"`                    //手机
	Email            string             `bson:"email" json:"email" gorm:"size:64;"`                    //邮箱
	Status           int                `bson:"status" json:"status" gorm:"size:4;"`                   //状态
	models.ControlBy `bson:"inline"`
	models.ModelTime `bson:"inline"`
	DataScope        string    `bson:"-" json:"dataScope" gorm:"-"`
	Params           string    `bson:"-" json:"params" gorm:"-"`
	Children         []SysDept `bson:"-" json:"children" gorm:"-"`
}

func (*SysDept) TableName() string {
	return "sysDept"
}

func (e *SysDept) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysDept) GetId() interface{} {
	return e.DeptId
}

func (e *SysDept) List(db *mongo.Database, filter SysDept, pageIndex, pageSize int, list *[]SysDept) error {
	var ctx = context.Background()
	query := bson.M{}
	if filter.DeptName != "" {
		query["deptName"] = filter.DeptName
	}
	if filter.Status > 0 {
		query["status"] = filter.Status
	}
	if pageIndex <= 0 {
		pageIndex = 1 // 默认值
	}
	if pageSize <= 0 {
		pageSize = 20 // 默认值
	}
	skip := (pageIndex - 1) * pageSize

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))

	cursor, err := db.Collection(e.TableName()).Find(ctx, query, findOptions)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, list); err != nil {
		return err
	}
	return nil
}
