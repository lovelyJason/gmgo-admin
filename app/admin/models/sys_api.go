package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"github.com/go-admin-team/go-admin-core/storage"

	"gmgo-admin/common/models"
)

type SysApi struct {
	Id               primitive.ObjectID `bson:"_id" json:"id"`
	ApiId            int                `bson:"apiId" json:"apiId"`
	Handle           string             `bson:"handle" json:"handle"`
	Title            string             `bson:"title" json:"title"`
	Path             string             `bson:"path" json:"path"`
	Action           string             `bson:"action" json:"action"`
	Type             string             `bson:"type" json:"type"`             // 接口类型,BUS:业务接口（可以用来分配权限的），SYS:系统接口
	Modifiable       bool               `bson:"modifiable" json:"modifiable"` // 是否可修改
	models.ModelTime `bson:"inline"`
	models.ControlBy `bson:"inline"`
}

func (*SysApi) TableName() string {
	return "sysApi"
}

func (e *SysApi) Generate() models.ActiveRecord {
	o := *e
	return &o
}

func (e *SysApi) GetId() interface{} {
	return e.Id
}

func SaveSysApi(message storage.Messager) (err error) {
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		err = fmt.Errorf("json Marshal error, %v", err.Error())
		return err
	}

	var l runtime.Routers
	err = json.Unmarshal(rb, &l)
	if err != nil {
		err = fmt.Errorf("json Unmarshal error, %s", err.Error())
		return err
	}
	dbList := sdk.Runtime.GetDb()
	for _, d := range dbList {
		for _, v := range l.List {
			if v.HttpMethod != "HEAD" ||
				strings.Contains(v.RelativePath, "/swagger/") ||
				strings.Contains(v.RelativePath, "/static/") ||
				strings.Contains(v.RelativePath, "/form-generator/") ||
				strings.Contains(v.RelativePath, "/sys/tables") {

				// 根据接口方法注释里的@Summary填充接口名称，适用于代码生成器
				// 可在此处增加配置路径前缀的if判断，只对代码生成的自建应用进行定向的接口名称填充
				jsonFile, _ := ioutil.ReadFile("docs/swagger.json")
				jsonData, _ := simplejson.NewFromReader(bytes.NewReader(jsonFile))
				urlPath := v.RelativePath
				idPatten := "(.*)/:(\\w+)" // 正则替换，把:id换成{id}
				reg, _ := regexp.Compile(idPatten)
				if reg.MatchString(urlPath) {
					urlPath = reg.ReplaceAllString(v.RelativePath, "${1}/{${2}}") // 把:id换成{id}
				}
				apiTitle, _ := jsonData.Get("paths").Get(urlPath).Get(strings.ToLower(v.HttpMethod)).Get("summary").String()

				err := d.Debug().Where(SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
					Attrs(SysApi{Handle: v.Handler, Title: apiTitle}).
					FirstOrCreate(&SysApi{}).
					//UpdateByUserId("handle", v.Handler).
					Error
				if err != nil {
					err := fmt.Errorf("Models SaveSysApi error: %s \r\n ", err.Error())
					return err
				}
			}
		}
	}
	return nil
}

func (e *SysApi) List(db *mongo.Database, filter SysApi, pageIndex, pageSize int) ([]*SysApi, int64, error) {
	var list []SysApi
	var total int64
	var ctx = context.Background()
	query := bson.M{}
	if filter.Type != "" {
		query["type"] = filter.Type
	}
	if filter.Title != "" {
		query["title"] = filter.Title
	}
	if filter.Path != "" {
		query["path"] = filter.Path
	}
	if filter.Action != "" {
		query["action"] = filter.Action
	}
	if pageIndex <= 0 {
		pageIndex = 1 // 默认值
	}
	if pageSize <= 0 {
		pageSize = 2000 // 默认值
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
	findOptions.SetSort(bson.M{"createdAt": -1})

	cursor, err := db.Collection(e.TableName()).Find(ctx, query, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &list); err != nil {
		return nil, 0, err
	}

	var result []*SysApi
	for i := range list {
		result = append(result, &list[i])
	}

	return result, total, nil
}

func (e *SysApi) GetOne(ctx context.Context, db *mongo.Database, apiId int, data *SysApi) error {
	filter := bson.M{"apiId": apiId}
	err := db.Collection(e.TableName()).FindOne(ctx, filter).Decode(&data)
	if err != nil {
		return err
	}
	return nil
}

func (e *SysApi) Count(ctx context.Context, db *mongo.Database, filter bson.M, count *int64) error {
	num, err := db.Collection(e.TableName()).CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	*count = num
	return nil
}

func (e *SysApi) InsertOne(ctx context.Context, db *mongo.Database, filter bson.M) (int, error) {
	collection := db.Collection(e.TableName())
	pipeline := mongo.Pipeline{
		{{"$group", bson.D{
			{"_id", nil},                             // 不分组
			{"maxApiId", bson.D{{"$max", "$apiId"}}}, // 获取最大值
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var result struct {
		MaxApiId int `bson:"maxApiId"`
	}
	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Max roleId: %d\n", result.MaxApiId)
	} else {
		fmt.Println("No documents found")
	}
	newApiId := result.MaxApiId + 1
	filter["apiId"] = newApiId
	_, err = collection.InsertOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return newApiId, nil
}

func (e *SysApi) UpdateByApiId(ctx context.Context, db *mongo.Database, apiId int, filter bson.M) error {
	_, err := db.Collection(e.TableName()).UpdateOne(ctx, bson.M{
		"apiId": apiId,
	}, bson.M{
		"$set": filter,
	})
	if err != nil {
		return err
	}

	return nil
}

func (e *SysApi) SoftDelByApiIds(ctx context.Context, db *mongo.Database, apiIds []int) error {
	now := time.Now()
	_, err := db.Collection(e.TableName()).UpdateMany(ctx, bson.M{"apiId": bson.M{"$in": apiIds}}, bson.M{
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
