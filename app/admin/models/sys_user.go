package models

import (
	"context"
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	"go-admin/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type SysUser struct {
	Id               primitive.ObjectID `bson:"_id"`
	UserId           int                `bson:"userId"  json:"userId"`
	Username         string             `bson:"username" json:"username"` // 用户名
	Password         string             `bson:"password" json:"-"`
	NickName         string             `bson:"nickName" json:"nickName"` // 昵称
	Phone            string             `bson:"phone" json:"phone"`       // 手机号
	RoleId           int                `bson:"roleId" json:"roleId"`     // 角色id
	Salt             string             `bson:"salt" json:"-"`            // 加盐
	Avatar           string             `bson:"avatar" json:"avatar"`
	Sex              string             `bson:"sex" json:"sex"`
	Email            string             `bson:"email" json:"email"`
	DeptId           int                `bson:"deptId" json:"deptId"`
	PostId           int                `bson:"postId" json:"postId"` // 岗位
	Remark           string             `bson:"remark" json:"remark"`
	Status           string             `bson:"status" json:"status"`
	DeptIds          []int              `bson:"-" json:"deptIds"`
	PostIds          []int              `bson:"-" json:"postIds"`
	RoleIds          []int              `bson:"-" json:"roleIds"`
	Dept             *SysDept           `json:"dept"`
	models.ControlBy `bson:"inline"`
	models.ModelTime `bson:"inline"`
}

func (*SysUser) TableName() string {
	return "sysUser"
}

func (e *SysUser) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysUser) GetId() interface{} {
	return e.UserId
}

// Encrypt 加密
func (e *SysUser) Encrypt() (err error) {
	if e.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		e.Password = string(hash)
		return
	}
}

func (e *SysUser) BeforeCreate(_ *gorm.DB) error {
	return e.Encrypt()
}

func (e *SysUser) BeforeUpdate(_ *gorm.DB) error {
	var err error
	if e.Password != "" {
		err = e.Encrypt()
	}
	return err
}

func (e *SysUser) AfterFind(_ *gorm.DB) error {
	e.DeptIds = []int{e.DeptId}
	e.PostIds = []int{e.PostId}
	e.RoleIds = []int{e.RoleId}
	return nil
}

func (e *SysUser) Count(ctx context.Context, db *mongo.Database, filter bson.M, count *int64) error {
	num, err := db.Collection(e.TableName()).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	*count = num
	return nil
}

func (e *SysUser) List(db *mongo.Database, filter SysUser, pageIndex, pageSize int) ([]*SysUser, int64, error) {
	var list []SysUser
	var total int64
	var ctx = context.Background()
	query := bson.M{}
	if filter.Username != "" {
		query["username"] = filter.Username
	}
	if filter.Phone != "" {
		query["phone"] = filter.Phone
	}
	if filter.Status != "" {
		query["status"] = filter.Status
	}
	if pageIndex <= 0 {
		pageIndex = 1 // 默认值
	}
	if pageSize <= 0 {
		pageSize = 20 // 默认值
	}
	query["$or"] = []bson.M{
		{"deletedAt": bson.M{"$exists": false}},
		{"deletedAt": nil},
	}
	skip := (pageIndex - 1) * pageSize

	total, _ = db.Collection(e.TableName()).CountDocuments(ctx, query)

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))

	cursor, err := db.Collection(e.TableName()).Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &list); err != nil {
		return nil, 0, err
	}

	var result []*SysUser
	for i := range list {
		result = append(result, &list[i])
	}

	return result, total, nil
}
func (e *SysUser) UpdateByUserId(ctx context.Context, db *mongo.Database, userId int, filter bson.M) error {
	_, err := db.Collection(e.TableName()).UpdateOne(ctx, bson.M{
		"userId": userId,
	}, bson.M{
		"$set": filter,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *SysUser) InsertOne(ctx context.Context, db *mongo.Database, filter bson.M) (int, error) {
	collection := db.Collection(e.TableName())
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", nil}, // 不分组
			{"maxUserId", bson.D{{"$max", "$userId"}}}, // 获取最大值
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result struct {
		MaxUserId int `bson:"maxUserId"`
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Max roleId: %d\n", result.MaxUserId)
	} else {
		fmt.Println("No documents found")
	}
	newUserId := result.MaxUserId + 1
	filter["userId"] = newUserId
	_, err = collection.InsertOne(ctx, filter)
	return newUserId, nil
}

func (e *SysUser) SoftDelByUserIds(ctx context.Context, db *mongo.Database, userIds []int, updateBy int) error {
	// 获取当前时间
	now := time.Now()

	// 使用 UpdateMany 进行软删除
	_, err := db.Collection(e.TableName()).UpdateMany(ctx, bson.M{"userId": bson.M{"$in": userIds}}, bson.M{
		"$set": bson.M{
			"deletedAt": now,      // 设置 deletedAt 为当前时间
			"updateBy":  updateBy, // 设置更新人
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
