package jwtauth

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type JwtService struct{}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList model.JwtBlacklist
//@return: err error

func (jwtService *JwtService) JsonInBlacklist(token string, db *mongo.Database) error {
	ctx := context.Background()
	filter := bson.M{
		"jwt":       token,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	}
	_, err := db.Collection("jwtBlacklists").InsertOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func (jwtService *JwtService) IsBlacklist(c *gin.Context, token string, db *mongo.Database) bool {
	// 创建一个新的上下文
	ctx := context.Background()

	// 创建一个空的结果结构体
	blacklistEntry := struct {
		Jwt string `bson:"jwt"`
	}{}

	// 查询 MongoDB 中是否存在对应的记录
	err := db.Collection("jwtBlacklists").FindOne(ctx, bson.M{"jwt": token}).Decode(&blacklistEntry)

	// 判断是否找到了记录
	return err == nil // 如果 err 为 nil，表示找到了记录
}
