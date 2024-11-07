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
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type SysRole struct {
	Id                     primitive.ObjectID   `bson:"_id" json:"id"`
	RoleId                 int                  `bson:"roleId" json:"roleId"`     // 角色编码
	RoleName               string               `bson:"roleName" json:"roleName"` // 角色名称
	Status                 string               `bson:"status" json:"status"`     // 状态 1禁用 2正常
	RoleKey                string               `bson:"roleKey" json:"roleKey"`   //角色代码
	RoleSort               int                  `bson:"roleSort" json:"roleSort"` //角色排序
	Flag                   string               `bson:"flag" json:"flag"`         //
	Remark                 string               `bson:"remark" json:"remark"`     //备注
	Admin                  bool                 `bson:"admin" json:"admin"`
	DataScope              string               `bson:"dataScope" json:"dataScope"`
	Params                 string               `bson:"-" json:"params"`
	Menus                  []primitive.ObjectID `bson:"menus" json:"menus"`
	SysMenu                *[]SysMenu           `bson:"sysMenu" json:"sysMenu"` // gorm:"many2many:sys_role_menu;foreignKey:RoleId;joinForeignKey:role_id;references:MenuId;joinReferences:menu_id;"`
	DeptIds                []int                `bson:"-" json:"deptIds"`
	SysDept                []SysDept            `bson:"-" json:"sysDept"` // gorm:"many2many:sys_role_dept;foreignKey:RoleId;joinForeignKey:role_id;references:DeptId;joinReferences:dept_id;"
	commonModels.ControlBy `bson:"inline"`
	commonModels.ModelTime `bson:"inline"`
}

func (*SysRole) TableName() string {
	return "sysRole"
}

func (e *SysRole) Generate() commonModels.ActiveRecord {
	o := *e
	return &o
}

func (e *SysRole) GetId() interface{} {
	return e.RoleId
}
func (e *SysRole) List(db *mongo.Database, filter SysRole, pageIndex, pageSize int) ([]*SysRole, int64, error) {
	var list []SysRole
	var total int64
	var ctx = context.Background()
	query := bson.M{}
	if filter.RoleName != "" {
		query["roleName"] = filter.RoleName
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

	var result []*SysRole
	for i := range list {
		result = append(result, &list[i])
	}

	return result, total, nil
}

func (e *SysRole) Count(ctx context.Context, db *mongo.Database, filter bson.M, count *int64) error {
	num, err := db.Collection(e.TableName()).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	*count = num
	return nil
}

func (e *SysRole) GetOne(ctx context.Context, db *mongo.Database, filter bson.M, data *SysRole) error {
	menuMoldel := models.SysMenu{}
	pipeline := mongo.Pipeline{
		{
			{"$match", filter}, // 匹配角色ID
		},
		{
			{"$lookup", bson.M{
				"from":         menuMoldel.TableName(), // 关联的集合
				"localField":   "menus",                // sysRole 中的字段
				"foreignField": "_id",                  // sysMenu 中的字段
				"as":           "sysMenu",              // 输出字段名称
			}},
		},
		{
			{"$addFields", bson.M{
				"sysMenu": bson.M{"$ifNull": []interface{}{"$sysMenu", []interface{}{}}}, // 如果没有找到 sysApi，则设置为空数组
			}},
		},
	}
	collection := db.Collection(e.TableName())
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err != nil {
		return err
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(data); err != nil {
			fmt.Printf("Decoded sysMenu: %+v\n", data.SysMenu)

			return err
		}

		// 调试输出
		fmt.Printf("Decoded data: %+v\n", data)
	} else {
		return mongo.ErrNoDocuments // 如果没有找到文档
	}
	return nil
}

func (e *SysRole) GetOneByRoleId(ctx context.Context, db *mongo.Database, roleId int, data *SysRole) error {
	menuMoldel := models.SysMenu{}
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.M{"roleId": roleId}}, // 匹配角色ID
		},
		{
			{"$lookup", bson.M{
				"from":         menuMoldel.TableName(), // 关联的集合
				"localField":   "menus",                // sysRole 中的字段
				"foreignField": "_id",                  // sysMenu 中的字段
				"as":           "sysMenu",              // 输出字段名称
			}},
		},
		{
			{"$addFields", bson.M{
				"sysMenu": bson.M{"$ifNull": []interface{}{"$sysMenu", []interface{}{}}}, // 如果没有找到 sysApi，则设置为空数组
			}},
		},
	}
	collection := db.Collection(e.TableName())
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	if err != nil {
		return err
	}
	if cursor.Next(ctx) {
		if err := cursor.Decode(data); err != nil {
			fmt.Printf("Decoded sysMenu: %+v\n", data.SysMenu)

			return err
		}

		// 调试输出
		fmt.Printf("Decoded data: %+v\n", data)
	} else {
		return mongo.ErrNoDocuments // 如果没有找到文档
	}
	return nil
}

// 事务操作，插入角色的时候再插入角色菜单表
func (e *SysRole) InserOneWithRoleMenu(ctx context.Context, db *mongo.Database, filter SysRole) error {
	// 开启事务
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	// 开始事务
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	txnOpts := options.Transaction().
		SetReadConcern(readconcern.Local()).                   // 使用Local读取关注
		SetWriteConcern(writeconcern.Majority()).              // 使用WMajority写入关注
		SetMaxCommitTime(&[]time.Duration{5 * time.Second}[0]) // 注意这里取地址

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 在这里执行数据库操作
		collection := db.Collection(e.TableName())
		// 直接传递结构体进去会吧bson._id也带上，从而默认0值覆盖了自增id
		filter.Id = primitive.NewObjectID()
		_, err := collection.InsertOne(sessCtx, filter)
		if err != nil {
			return nil, err
		}
		roleMenuCollection := db.Collection("sysRoleMenu")
		sysMenuSlice := *filter.SysMenu // 解引用以获得 []SysMenu
		// 构建bson.A, 因为insertMany不支持传入SysMenu这样的结构体
		var menuSlice bson.A
		for _, menu := range sysMenuSlice {
			menuSlice = append(menuSlice, bson.M{
				"roleId": filter.RoleId,
				"menuId": menu.MenuId,
			}) // 将每个 SysMenu 添加到接口切片中
		}
		_, err = roleMenuCollection.InsertMany(sessCtx, menuSlice)
		if err != nil {
			return nil, err
		}

		return nil, nil // 如果没有错误，返回nil表示提交事务
	}, txnOpts)

	if err != nil {
		// 处理事务错误（这里将自动回滚事务）
		log.Fatal(err)
		return err
	}
	return nil
}

// 文档引用版
func (e *SysRole) InsertOne(ctx context.Context, db *mongo.Database, filter bson.M) (int, error) {
	collection := db.Collection(e.TableName())
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", nil}, // 不分组
			{"maxRoleId", bson.D{{"$max", "$roleId"}}}, // 获取最大值
		}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result struct {
		MaxRoleId int `bson:"maxRoleId"`
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Max roleId: %d\n", result.MaxRoleId)
	} else {
		fmt.Println("No documents found")
	}
	newRoleId := result.MaxRoleId + 1
	filter["roleId"] = newRoleId
	_, err = collection.InsertOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return newRoleId, nil
}

// 关联表版
func (e *SysRole) UpdateRoleMenuWithRoleId(ctx context.Context, db *mongo.Database, roleId int, filter bson.M, menuIds []int) error {
	session, err := db.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	// 开始事务
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	txnOpts := options.Transaction().
		SetReadConcern(readconcern.Local()).                   // 使用Local读取关注
		SetWriteConcern(writeconcern.Majority()).              // 使用WMajority写入关注
		SetMaxCommitTime(&[]time.Duration{5 * time.Second}[0]) // 注意这里取地址
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		_, err := db.Collection(e.TableName()).UpdateOne(ctx, bson.M{"roleId": roleId}, bson.M{"$set": filter})
		if err != nil {
			return nil, err
		}
		filter2 := bson.M{"roleId": roleId}
		_, err = db.Collection("sysRoleMenu").DeleteMany(ctx, filter2)
		if err != nil {
			return nil, err
		}
		// 插入菜单到关联表
		var documents bson.A
		for _, menuId := range menuIds {
			doc := bson.M{
				"roleId": roleId,
				"menuId": menuId,
			}
			documents = append(documents, doc)
		}
		_, err = db.Collection("sysRoleMenu").InsertMany(ctx, documents)
		if err != nil {
			// 处理错误
			log.Fatalf("Failed to insert documents: %v", err)
		}

		return nil, nil // 如果没有错误，返回nil表示提交事务
	}, txnOpts)

	return nil
}

// 文档引用版
func (e *SysRole) UpdateByRoleId(ctx context.Context, db *mongo.Database, roleId int, filter bson.M) error {
	_, err := db.Collection(e.TableName()).UpdateOne(ctx, bson.M{
		"roleId": roleId,
	}, bson.M{
		"$set": filter,
	})
	if err != nil {
		return err
	}
	return nil
}
