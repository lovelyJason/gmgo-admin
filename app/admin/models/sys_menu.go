package models

import (
	"context"
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	"go-admin/cmd/migrate/migration/models"
	commonModels "go-admin/common/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type SysMenu struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	MenuId     int                `bson:"menuId" json:"menuId"`          // autoIncrement的菜单id
	MenuName   string             `bson:"menuName" json:"menuName"`      // 菜单英文名称
	Title      string             `bson:"title" json:"title"`            // 菜单中文名称
	Icon       string             `bson:"icon" json:"icon"`              // 菜单图标
	Path       string             `bson:"path" json:"path"`              // 菜单路径
	Paths      string             `bson:"paths" json:"paths" `           // 菜单的菜单id路径
	MenuType   string             `bson:"menuType" json:"menuType" `     // 菜单类型
	Action     string             `bson:"action" json:"action" `         // api请求方式
	Permission string             `bson:"permission" json:"permission" ` // 权限标识
	ParentId   int                `bson:"parentId" json:"parentId" `     // 父菜单id
	NoCache    bool               `bson:"noCache" json:"noCache" `
	Breadcrumb string             `bson:"breadcrumb" json:"breadcrumb" `
	Component  string             `bson:"component" json:"component" `
	Sort       int                `bson:"sort" json:"sort"`
	Visible    string             `bson:"visible" json:"visible"`
	IsFrame    string             `bson:"isFrame" json:"isFrame"` // 是否外链
	// 注意，即使表里没有sysApi这个字段，但是通过聚合查询加进来的，也要注明bson
	SysApi                 []SysApi             `bson:"sysApi" json:"sysApi"  // gorm:"many2many:sys_menu_api_rule"`
	Apis                   []primitive.ObjectID `bson:"apis" json:"apis"`
	DataScope              string               `bson:"-" json:"dataScope"`
	Params                 string               `bson:"-" json:"params"`
	RoleId                 int                  `bson:"-"`
	Children               []SysMenu            `bson:"-" json:"children,omitempty"`
	IsSelect               bool                 `bson:"-" json:"is_select"`
	commonModels.ControlBy `bson:"inline"`
	commonModels.ModelTime `bson:"inline"`
}

type SysMenuSlice []SysMenu

func (x SysMenuSlice) Len() int           { return len(x) }
func (x SysMenuSlice) Less(i, j int) bool { return x[i].Sort < x[j].Sort }
func (x SysMenuSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func (*SysMenu) TableName() string {
	return "sysMenu"
}

func (e *SysMenu) Generate() commonModels.ActiveRecord {
	o := *e
	return &o
}

func (e *SysMenu) GetId() interface{} {
	return e.MenuId
}

func (e *SysMenu) List(ctx context.Context, db *mongo.Database, filter bson.M, pageIndex, pageSize int, list *[]SysMenu) error {
	//query := bson.M{}
	//if filter.RoleId > 0 { // TODO: 目前菜单表中没有roleId
	//	query["roleId"] = filter.RoleId
	//}
	//if filter["menuName"] != "" {
	//	query["menuName"] = filter["menuName"]
	//}
	//if filter["visible"] != "" {
	//	query["visible"] = filter["visible"]
	//}
	apiModel := &models.SysApi{}
	if pageIndex <= 0 {
		pageIndex = 1 // 默认值
	}
	if pageSize <= 0 {
		pageSize = 2000 // 默认值
	}
	filter["$or"] = []bson.M{
		{"deletedAt": bson.M{"$exists": false}},
		{"deletedAt": nil},
	}
	skip := (pageIndex - 1) * pageSize

	//findOptions := options.Find()
	//findOptions.SetSkip(int64(skip))
	//findOptions.SetLimit(int64(pageSize))
	//findOptions.SetSort(bson.M{"sort": 1})

	//cursor, err := db.Collection(e.TableName()).Find(ctx, filter, findOptions)
	pipeline := mongo.Pipeline{
		{
			{"$match", filter}, // 匹配菜单ID
		},
		{
			{"$lookup", bson.M{
				"from":         apiModel.TableName(), // 关联的集合
				"localField":   "apis",               // sysMenu 中的字段
				"foreignField": "_id",                // sysApi 中的字段
				"as":           "sysApi",             // 输出字段名称
			}},
		},
		{
			{"$addFields", bson.M{
				"sysApi": bson.M{"$ifNull": []interface{}{"$sysApi", []interface{}{}}}, // 如果没有找到 sysApi，则设置为空数组
			}},
		},
		// 分页操作：跳过指定数量的文档
		{{"$skip", int64(skip)}},

		// 限制查询结果的数量
		{{"$limit", int64(pageSize)}},

		// 排序操作
		{{"$sort", bson.M{"sort": 1}}},
	}

	// 执行聚合查询
	cursor, err := db.Collection(e.TableName()).Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, list); err != nil {
		return err
	}

	return nil
}

// 关联表版
func (e *SysMenu) RoleMenuList(ctx context.Context, db *mongo.Database, filter bson.M, list *[]SysMenu) error {
	cursor, err := db.Collection("sysRoleMenu").Find(ctx, filter)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, list); err != nil {
		return err
	}

	return nil
}

//	func (e *SysMenu) GetOneByMenuId(ctx context.Context, db *mongo.Database, menuId int, data *SysMenu) error {
//		filter := bson.M{
//			"menuId": menuId,
//		}
//		err := db.Collection(e.TableName()).FindOne(ctx, filter).Decode(&data)
//		if err != nil {
//			return err
//		}
//		return nil
//	}

func (e *SysMenu) Get(ctx context.Context, db *mongo.Database, filter bson.M, data *[]SysMenu) error {
	collection := db.Collection(e.TableName())

	// 执行查询
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute find: %w", err)
	}
	defer cursor.Close(ctx)

	// 检查光标是否为空
	if !cursor.Next(ctx) {
		return mongo.ErrNoDocuments // 如果没有找到文档，返回相应的错误
	}

	// 解析结果到 data
	if err := cursor.Decode(data); err != nil {
		return fmt.Errorf("failed to decode result: %w", err)
	}

	return nil // 查询成功
}

func (e *SysMenu) GetOneByMenuId(ctx context.Context, db *mongo.Database, menuId int, data *SysMenu) error {
	apiModel := models.SysApi{}
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{"menuId": menuId}}, // 匹配菜单ID
		},
		{
			{"$lookup", bson.M{
				"from":         apiModel.TableName(), // 关联的集合
				"localField":   "apis",               // sysMenu 中的字段
				"foreignField": "_id",                // sysApi 中的字段
				"as":           "sysApi",             // 输出字段名称
			}},
		},
		{
			{"$addFields", bson.M{
				"sysApi": bson.M{"$ifNull": []interface{}{"$sysApi", []interface{}{}}}, // 如果没有找到 sysApi，则设置为空数组
			}},
		},
	}

	// 执行聚合查询
	collection := db.Collection(e.TableName())
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err := cursor.Decode(data); err != nil {
			return err
		}

		// 调试输出
		fmt.Printf("Decoded data: %+v\n", data)
	} else {
		return mongo.ErrNoDocuments // 如果没有找到文档
	}

	return nil
}

func (e *SysMenu) UpdateByMenuId(ctx context.Context, db *mongo.Database, menuId int, filter bson.M) error {
	_, err := db.Collection(e.TableName()).UpdateOne(ctx, bson.M{
		"menuId": menuId,
	}, bson.M{
		"$set": filter,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *SysMenu) InsertOne(ctx context.Context, db *mongo.Database, filter bson.M) (int, error) {
	collection := db.Collection(e.TableName())
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", nil}, // 不分组
			{"maxMenuId", bson.D{{"$max", "$menuId"}}}, // 获取最大值
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result struct {
		MaxMenuId int `bson:"maxMenuId"`
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Max menuId: %d\n", result.MaxMenuId)
	} else {
		fmt.Println("No documents found")
	}
	newMenuId := result.MaxMenuId + 1
	filter["menuId"] = newMenuId
	_, err = collection.InsertOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return newMenuId, nil
}

func (e *SysMenu) SoftDelByUserIds(ctx context.Context, db *mongo.Database, menuIds []int) error {
	now := time.Now()
	_, err := db.Collection(e.TableName()).UpdateMany(ctx, bson.M{"menuId": bson.M{"$in": menuIds}}, bson.M{
		"$set": bson.M{
			"deletedAt": now,
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
