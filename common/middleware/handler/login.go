package handler

import (
	"context"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Login struct {
	Username string `form:"UserName" json:"username" binding:"required"`
	Password string `form:"Password" json:"password" binding:"required"`
	Code     string `form:"Code" json:"code" binding:"required"`
	UUID     string `form:"UUID" json:"uuid" binding:"required"`
}

func (u *Login) GetUser(tx *gorm.DB) (user SysUser, role SysRole, err error) {
	err = tx.Table("sys_user").Where("username = ?  and status = '2'", u.Username).First(&user).Error
	if err != nil {
		log.Errorf("get user error, %s", err.Error())
		return
	}
	_, err = pkg.CompareHashAndPassword(user.Password, u.Password)
	if err != nil {
		log.Errorf("user login error, %s", err.Error())
		return
	}
	err = tx.Table("sys_role").Where("role_id = ? ", user.RoleId).First(&role).Error
	if err != nil {
		log.Errorf("get role error, %s", err.Error())
		return
	}
	return
}

func (u *Login) GetUserWithMongo(db *mongo.Database) (user SysUser, role SysRole, err error) {
	userCollection := db.Collection("sysUser")

	// 查找用户
	err = userCollection.FindOne(context.TODO(), bson.M{
		"username": u.Username,
		"status":   "2", // 1禁用 2启用
		"$or": []bson.M{
			{"deletedAt": bson.M{"$exists": false}},
			{"deletedAt": nil},
		},
	}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Errorf("user not found, %s", err.Error())
		} else {
			log.Errorf("get user error, %s", err.Error())
		}
		return
	}

	// 校验密码
	_, err = pkg.CompareHashAndPassword(user.Password, u.Password)
	if err != nil {
		log.Errorf("user login error, %s", err.Error())
		return
	}

	// 查找角色
	roleCollection := db.Collection("sysRole")
	err = roleCollection.FindOne(context.TODO(), bson.M{
		"roleId": user.RoleId,
	}).Decode(&role)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Errorf("role not found, %s", err.Error())
		} else {
			log.Errorf("get role error, %s", err.Error())
		}
		return
	}
	return
}
