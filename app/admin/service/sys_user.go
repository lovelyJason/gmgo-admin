package service

import (
	"context"
	"errors"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	"go-admin/common/actions"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type SysUser struct {
	service.Service
}

// GetPage 获取SysUser列表
func (e *SysUser) GetPage(c *dto.SysUserGetPageReq, p *actions.DataPermission) ([]*models.SysUser, int64, error) {
	var err error
	//var data models.SysUser
	model := &models.SysUser{}
	filter := models.SysUser{}
	//err = e.Orm.Debug().Preload("Dept").
	//	Scopes(
	//		cDto.MakeCondition(c.GetNeedSearch()),
	//		cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
	//		actions.Permission(data.TableName(), p),
	//	).
	//	Find(list).Limit(-1).Offset(-1).
	//	Count(count).Error
	//if err != nil {
	//	e.Log.Errorf("db error: %s", err)
	//	return err
	//}
	//return nil
	list, total, err := model.List(e.Mongo, filter, c.PageIndex, c.PageSize)
	if err != nil {
		e.Log.Errorf("Service GetSysApiPage error:%s", err)
		return nil, 0, err
	}
	return list, total, err

}

// Get 获取SysUser对象
func (e *SysUser) Get(ctx context.Context, d *dto.SysUserById, p *actions.DataPermission, model *models.SysUser) error {
	var data models.SysUser
	collection := e.Mongo.Collection("sysUser")

	// TODO: dataScope的处理，依赖于开关EnableDP
	//err := e.Orm.Model(&data).Debug().
	//	Scopes(
	//		actions.Permission(data.TableName(), p),
	//	).
	//	First(model, d.GetId()).Error
	filter := bson.M{"userId": d.GetId()} // 使用 BSON 过滤器查找用户

	err := collection.FindOne(ctx, filter).Decode(&data) // 查询并解码
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	*model = data
	return nil
}

// Insert 创建SysUser对象
func (e *SysUser) Insert(c *dto.SysUserInsertReq) (int, error) {
	var err error
	var data models.SysUser
	var userModel = &models.SysUser{}
	var filter = bson.M{
		"username": c.Username,
	}
	var ctx = context.Background()
	var i int64
	err = userModel.Count(ctx, e.Mongo, filter, &i)
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return 0, err
	}
	if i > 0 {
		err := errors.New("用户名已存在！")
		e.Log.Errorf("db error: %s", err)
		return 0, err
	}
	c.Generate(&data)

	insertFilter := bson.M{
		"username":  data.Username,
		"password":  "$2a$10$/Glr4g9Svr6O0kvjsRJCXu3f0W8/dsP3XZyVNi1019ratWpSPMyw.", // 123456
		"nickName":  data.NickName,
		"phone":     data.Phone,
		"roleId":    data.RoleId,
		"avatar":    data.Avatar,
		"sex":       data.Sex,
		"email":     data.Email,
		"deptId":    data.DeptId,
		"postId":    data.PostId,
		"remark":    data.Remark,
		"status":    data.Status,
		"createdAt": time.Now(),
		"createBy":  c.CreateBy,
	}
	userId, err := userModel.InsertOne(ctx, e.Mongo, insertFilter)
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return 0, err
	}
	return userId, nil
}

// Update 修改SysUser对象
func (e *SysUser) Update(c *dto.SysUserUpdateReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	c.Generate(&model)
	ctx := context.Background()
	userModel := &models.SysUser{}
	filter := bson.M{
		"username":  model.Username,
		"nickName":  model.NickName,
		"phone":     model.Phone,
		"roleId":    model.RoleId,
		"avatar":    model.Avatar,
		"sex":       model.Sex,
		"email":     model.Email,
		"deptId":    model.DeptId,
		"postId":    model.PostId,
		"remark":    model.Remark,
		"status":    model.Status,
		"updatedAt": time.Now(),
		"updateBy":  c.UpdateBy,
	}
	err = userModel.UpdateByUserId(ctx, e.Mongo, c.UserId, filter)
	if err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	return nil
}

// UpdateAvatar 更新用户头像
func (e *SysUser) UpdateAvatar(c *dto.UpdateSysUserAvatarReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	err = e.Orm.Table(model.TableName()).Where("user_id =? ", c.UserId).Updates(c).Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	return nil
}

// UpdateStatus 更新用户状态
func (e *SysUser) UpdateStatus(c *dto.UpdateSysUserStatusReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	db := e.Orm.Scopes(
		actions.Permission(model.TableName(), p),
	).First(&model, c.GetId())
	if err = db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")

	}
	err = e.Orm.Table(model.TableName()).Where("user_id =? ", c.UserId).Updates(c).Error
	if err != nil {
		e.Log.Errorf("Service UpdateSysUser error: %s", err)
		return err
	}
	return nil
}

// ResetPwd 重置用户密码
func (e *SysUser) ResetPwd(c *dto.ResetSysUserPwdReq, p *actions.DataPermission) error {
	var err error
	var model models.SysUser
	c.Generate(&model)
	ctx := context.Background()
	userModel := &models.SysUser{}
	//err = e.Orm.Omit("username", "nick_name", "phone", "role_id", "avatar", "sex").Save(&model).Error
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)
	if err != nil {
		e.Log.Errorf("At Service ResetSysUserPwd error: %s", err)
		return err
	}
	filter := bson.M{
		"password": string(encryptedPassword),
	}
	userModel.UpdateByUserId(ctx, e.Mongo, c.GetId(), filter)
	if err != nil {
		e.Log.Errorf("At Service ResetSysUserPwd error: %s", err)
		return err
	}
	return nil
}

// Remove 删除SysUser
func (e *SysUser) Remove(c *dto.SysUserById, p *actions.DataPermission) error {
	var err error
	userModel := models.SysUser{}
	ctx := context.Background()

	err = userModel.SoftDelByUserIds(ctx, e.Mongo, c.UserIds, c.UpdateBy)

	if err != nil {
		e.Log.Errorf("Error found in  RemoveSysUser : %s", err)
		return err
	}
	return nil
}

// UpdatePwd 修改SysUser对象密码
func (e *SysUser) UpdatePwd(id int, oldPassword, newPassword string, p *actions.DataPermission) error {
	var err error

	if newPassword == "" {
		return nil
	}
	c := &models.SysUser{}

	err = e.Orm.Model(c).
		Scopes(
			actions.Permission(c.TableName(), p),
		).Select("UserId", "Password", "Salt").
		First(c, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("无权更新该数据")
		}
		e.Log.Errorf("db error: %s", err)
		return err
	}
	var ok bool
	ok, err = pkg.CompareHashAndPassword(c.Password, oldPassword)
	if err != nil {
		e.Log.Errorf("CompareHashAndPassword error, %s", err.Error())
		return err
	}
	if !ok {
		err = errors.New("incorrect Password")
		e.Log.Warnf("user[%d] %s", id, err.Error())
		return err
	}
	c.Password = newPassword
	db := e.Orm.Model(c).Where("user_id = ?", id).
		Select("Password", "Salt").
		Updates(c)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error: %s", err)
		return err
	}
	if db.RowsAffected == 0 {
		err = errors.New("set password error")
		log.Warnf("db update error")
		return err
	}
	return nil
}

func (e *SysUser) GetProfile(c *dto.SysUserById, user *models.SysUser, roles *[]models.SysRole, posts *[]models.SysPost) error {
	err := e.Orm.Preload("Dept").First(user, c.GetId()).Error
	if err != nil {
		return err
	}
	err = e.Orm.Find(roles, user.RoleId).Error
	if err != nil {
		return err
	}
	err = e.Orm.Find(posts, user.PostIds).Error
	if err != nil {
		return err
	}

	return nil
}
