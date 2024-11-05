package models

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InitDb 初始化数据库
func InitDb(uri string, filePath string) error {
	// 创建 MongoDB 客户端
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	// 连接到 MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	// 读取 JSON 文件
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// 解析 JSON 数据
	var collections map[string][]interface{}
	if err = json.Unmarshal(data, &collections); err != nil {
		return err
	}

	// 遍历每个集合及其文档
	for collectionName, docs := range collections {
		collection := client.Database("your_database").Collection(collectionName)

		// 使用 InsertMany 插入数据
		_, err = collection.InsertMany(ctx, docs)
		if err != nil {
			return err
		}

		log.Printf("集合 %s 数据插入成功\n", collectionName)
	}

	return nil
}
