package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/go-admin-team/go-admin-core/sdk/config"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"gmgo-admin/app/admin/models"
	"gmgo-admin/app/admin/service/dto"
)

type SysRole struct {
	service.Service
}

type RoleMenuResult struct {
	Title       string          `bson:"title"`
	Permission  string          `bson:"permission"`
	ApisDetails []models.SysApi `bson:"apisDetails"`
}

// GetPage 获取SysRole列表
func (e *SysRole) GetPage(c *dto.SysRoleGetPageReq) ([]*models.SysRole, int64, error) {
	var err error
	//var data models.SysRole
	model := &models.SysRole{}
	filter := models.SysRole{
		RoleName: c.RoleName,
		RoleKey:  c.RoleKey,
		Status:   c.Status,
	}

	//err = e.Orm.Model(&data).Preload("SysMenu").
	//	Scopes(
	//		cDto.MakeCondition(c.GetNeedSearch()),
	//		cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
	//	).
	//	Find(list).Limit(-1).Offset(-1).
	//	Count(count).Error
	//if err != nil {
	//	e.Log.Errorf("db error:%s", err)
	//	return err
	//}
	//return nil
	list, total, err := model.List(e.Mongo, filter, c.PageIndex, c.PageSize)
	if err != nil {
		e.Log.Errorf("Service GetSysRolePage error:%s", err)
		return nil, 0, err
	}
	return list, total, err
}

// Get 获取SysRole对象
func (e *SysRole) Get(ctx context.Context, d *dto.SysRoleGetReq, model *models.SysRole) error {
	//var err error
	//db := e.Orm.First(model, d.GetId())
	//err = db.Error
	//if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
	//	err = errors.New("查看对象不存在或无权查看")
	//	e.Log.Errorf("db error:%s", err)
	//	return err
	//}
	//if err != nil {
	//	e.Log.Errorf("db error:%s", err)
	//	return err
	//}
	//model.MenuIds, err = e.GetRoleMenuId(model.RoleId)
	//if err != nil {
	//	e.Log.Errorf("get menuIds error, %s", err.Error())
	//	return err
	//}
	//return nil
	// 注意这里是roleID
	err := model.GetOneByRoleId(ctx, e.Mongo, d.Id, model)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	//model.MenuIds, err = e.GetRoleMenuId(model.RoleId) // 废弃，直接通过聚合查询
	return nil
}

// Insert 创建SysRole对象
func (e *SysRole) Insert(c *dto.SysRoleInsertReq, cb *casbin.SyncedEnforcer) error {
	var err error
	var data models.SysRole
	var dataMenu []models.SysMenu
	// 预先加载与SysMenu实例相关联的SysApi实例,然后查询那些menu_id在c.MenuIds列表中的SysMenu实例
	// 并将查询结果填充到dataMenu这个切片中。
	//err = e.Orm.Preload("SysApi").Where("menu_id in ?", c.MenuIds).Find(&dataMenu).Error
	// 1. 先拉取角色所有菜单
	roleModel := &models.SysRole{}
	menuModel := &models.SysMenu{}
	ctx := context.Background()
	menuFilter := bson.M{
		"_id": bson.M{"$in": c.Menus},
	}
	err = menuModel.List(ctx, e.Mongo, menuFilter, 0, 0, &dataMenu) // 拉取请求参数菜单id对应的菜单表数据列表，最重要的是要拉取到SysApi
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	c.SysMenu = dataMenu
	c.Generate(&data) // 将请求参数bind到data上,方便后续事务操作menu表和roleMenu关联表

	var count int64
	//err = tx.Model(&data).Where("role_key = ?", c.RoleKey).Count(&count).Error
	roleFilter := bson.M{
		"$or": []bson.M{
			{"roleName": c.RoleName},
			{"roleKey": c.RoleKey},
			{"roleId": c.RoleId},
		},
	}
	err = roleModel.Count(ctx, e.Mongo, roleFilter, &count)
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}

	if count > 0 {
		err = errors.New("角色id或者角色名或者权限字符已存在，需更换在提交！")
		e.Log.Errorf("db error:%s", err)
		return err
	}

	insertFilter := c.GenerateBson()
	newRoleId, err := roleModel.InsertOne(ctx, e.Mongo, insertFilter)
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	fmt.Println(newRoleId)

	// TODO: 好像没有生成数据库表啊， TODO:这里menuitem Apis有值，但是SysApi没查询上来
	mp := make(map[string]interface{}, 0)
	polices := make([][]string, 0) // // 类似于["角色名称", "/api/v1/sys-user", "GET"]
	for _, menu := range dataMenu {
		for _, api := range menu.SysApi {
			if mp[data.RoleKey+"-"+api.Path+"-"+api.Action] != "" {
				mp[data.RoleKey+"-"+api.Path+"-"+api.Action] = ""
				polices = append(polices, []string{data.RoleKey, api.Path, api.Action})
			}
		}
	}

	if len(polices) <= 0 {
		return nil
	}

	// 写入 sys_casbin_rule 权限表里 当前角色数据的记录
	_, err = cb.AddNamedPolicies("p", polices)
	if err != nil {
		return err
	}

	return nil
}

func getRoleMenuResults(c *mongo.Collection, menus []primitive.ObjectID) ([]RoleMenuResult, error) {
	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"menus", bson.D{
					{"$in", menus},
				}},
			}},
		},
		{
			{"$lookup", bson.D{
				{"from", "sysMenu"},
				{"localField", "menus"},
				{"foreignField", "_id"},
				{"as", "menu_details"},
			}},
		},
		{
			{"$unwind", "$menu_details"},
		},
		{
			{"$match", bson.D{
				{"menu_details.permission", bson.D{{"$ne", nil}}},
			}},
		},
		{
			{"$project", bson.D{
				{"_id", 0},
				{"title", "$menu_details.title"},
				{"permission", "$menu_details.permission"},
				{"apis", "$menu_details.apis"},
			}},
		},
		{
			{"$lookup", bson.D{
				{"from", "sysApi"},
				{"localField", "apis"},
				{"foreignField", "_id"},
				{"as", "apisDetails"},
			}},
		},
		{
			{"$project", bson.D{
				{"title", 1},
				{"permission", 1},
				{"apisDetails", 1},
			}},
		},
	}

	// 执行聚合查询
	cursor, err := c.Aggregate(context.Background(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// 存储查询结果
	var results []RoleMenuResult
	for cursor.Next(context.Background()) {
		var result RoleMenuResult
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// Update 修改SysRole对象
func (e *SysRole) Update(c *dto.SysRoleUpdateReq, cb *casbin.SyncedEnforcer) error {
	var err error
	var model = models.SysRole{}
	//tx.Preload("SysMenu").First(&model, c.GetId())                      // 预加载了SysRole模型中的SysMenu关联。这意味着通过roleId查询一个SysRole实例时，GORM会自动查询并填充与该角色相关联的所有菜单（SysMenu）
	//tx.Preload("SysApi").Where("menu_id in ?", c.MenuIds).Find(&mlist)  // SysMenu模型中的SysApi关联.many2many:sys_menu_api_rule。查询SysMenu实例时，自动查询并填充与该菜单相关联的所有API（SysApi）,存到mlist。结果：sysMenu.SysApi
	//err = tx.Model(&model).Association("SysMenu").Delete(model.SysMenu) // Association: 这个方法用于访问角色模型的关联SysMenu:从sys_role_menu这样的多对多关联表中删除了相应的记录model.SysMenu
	// 思路就是一开始获取到前端传递的menuIds查到菜单列表存到mlist，然后删除这个角色所有菜单，再重新添加。
	ctx := context.Background()
	roleModel := &models.SysRole{}

	c.Generate(&model)
	//model.SysMenu = &mlist // mlist是前端传递过来的为角色赋予的菜单列表，在后文查询
	// 更新关联的数据
	roleFilter := c.GenerateBson()
	//err = roleModel.UpdateRoleMenuWithRoleId(ctx, e.Mongo, c.RoleId, roleFilter, c.MenuIds)
	err = roleModel.UpdateByRoleId(ctx, e.Mongo, c.RoleId, roleFilter)
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	// 清除 sys_casbin_rule 权限表里 当前角色的所有记录
	_, err = cb.RemoveFilteredPolicy(0, model.RoleKey)
	if err != nil {
		e.Log.Errorf("delete policy error:%s", err)
		return err
	}
	// 生成apis策略表
	menusWithApis, err := getRoleMenuResults(e.Mongo.Collection(roleModel.TableName()), c.Menus)

	mp := make(map[string]interface{}, 0)
	polices := make([][]string, 0) // // 遍历menus，并遍历每一项menu的SysApi，生成类似于["角色名称", "/api/v1/sys-user", "GET"]
	for _, menu := range menusWithApis {
		for _, api := range menu.ApisDetails {
			if mp[model.RoleKey+"-"+api.Path+"-"+api.Action] != "" {
				mp[model.RoleKey+"-"+api.Path+"-"+api.Action] = ""
				//_, err = cb.AddNamedPolicy("p", model.RoleKey, api.Path, api.Action)
				polices = append(polices, []string{model.RoleKey, api.Path, api.Action})
			}
		}
	}
	if len(polices) <= 0 {
		return nil
	}

	// 写入 sys_casbin_rule 权限表里 当前角色数据的记录
	_, err = cb.AddNamedPolicies("p", polices)
	if err != nil {
		return err
	}
	return nil
}

// Remove 删除SysRole
func (e *SysRole) Remove(c *dto.SysRoleDeleteReq, cb *casbin.SyncedEnforcer) error {
	var err error
	roleModel := &models.SysRole{}
	ctx := context.Background()
	filter := bson.M{
		"roleId": c.RoleId,
	}
	err = roleModel.GetOne(ctx, e.Mongo, filter, roleModel)
	if err != nil {
		return err
	}
	err = roleModel.SoftDelByRoleId(ctx, e.Mongo, c.RoleId)
	if err != nil {
		return err
	}
	// 清除 sys_casbin_rule 权限表里 当前角色的所有记录
	_, _ = cb.RemoveFilteredPolicy(0, roleModel.RoleKey)

	return nil
}

// GetRoleMenuId 获取角色对应的菜单ids
//
//	func (e *SysRole) GetRoleMenuId(roleId int) ([]int, error) {
//		menuIds := make([]int, 0)
//		model := models.SysRole{}
//		model.RoleId = roleId
//		if err := e.Orm.Model(&model).Preload("SysMenu").First(&model).Error; err != nil {
//			return nil, err
//		}
//		l := *model.SysMenu
//		for i := 0; i < len(l); i++ {
//			menuIds = append(menuIds, l[i].MenuId)
//		}
//		return menuIds, nil
//	}
func (e *SysRole) GetRoleMenuId(roleId int) ([]int, error) {
	menuModel := &models.SysMenu{}
	filter := bson.M{
		"roleId": roleId,
	}
	var menuList = make([]models.SysMenu, 0)
	menuIds := make([]int, 0)
	ctx := context.Background()

	menuModel.RoleMenuList(ctx, e.Mongo, filter, &menuList)

	for i := 0; i < len(menuList); i++ {
		menuIds = append(menuIds, menuList[i].MenuId)
	}
	return menuIds, nil
}

func (e *SysRole) UpdateDataScope(c *dto.RoleDataScopeReq) *SysRole {
	var err error
	tx := e.Orm
	if config.DatabaseConfig.Driver != "sqlite3" {
		tx := e.Orm.Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
	}
	var dlist = make([]models.SysDept, 0)
	var model = models.SysRole{}
	tx.Preload("SysDept").First(&model, c.RoleId)
	tx.Where("dept_id in ?", c.DeptIds).Find(&dlist)
	// 删除SysRole 和 SysDept 的关联关系
	err = tx.Model(&model).Association("SysDept").Delete(model.SysDept)
	if err != nil {
		e.Log.Errorf("delete SysDept error:%s", err)
		_ = e.AddError(err)
		return e
	}
	c.Generate(&model)
	model.SysDept = dlist
	// 更新关联的数据，使用 FullSaveAssociations 模式
	db := tx.Model(&model).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	if db.RowsAffected == 0 {
		_ = e.AddError(errors.New("无权更新该数据"))
		return e
	}
	return e
}

// UpdateStatus 修改SysRole对象status
func (e *SysRole) UpdateStatus(c *dto.UpdateStatusReq) error {
	var err error
	ctx := context.Background()
	roleModel := &models.SysRole{}
	filter := bson.M{
		"status":    c.Status,
		"updateBy":  c.UpdateBy,
		"updatedAt": time.Now(),
	}
	err = roleModel.UpdateByRoleId(ctx, e.Mongo, c.RoleId, filter)
	if err != nil {
		e.Log.Errorf("Service UpdateSysRole error: %s", err)
		return err
	}
	return nil
}

// GetById 获取SysRole对象
func (e *SysRole) GetById(roleId int) ([]string, error) {
	ctx := context.Background()
	roleModel := &models.SysRole{}
	menuModel := &models.SysMenu{}

	// 构建聚合管道
	pipeline := []bson.M{
		{"$match": bson.M{"roleId": roleId}}, // 匹配特定角色ID
		{"$lookup": bson.M{
			"from":         menuModel.TableName(), // 菜单表的名称
			"localField":   "menus",               // 角色集合中的菜单 ObjectId 数组
			"foreignField": "_id",                 // 菜单集合中的 ObjectId 字段
			"as":           "menu_details",        // 将匹配到的菜单数据放到新的字段 menu_details 中
		}},
		{"$unwind": "$menu_details"}, // 展开 menu_details 数组
		{"$match": bson.M{
			"menu_details.permission": bson.M{"$ne": ""}, // 只选择具有权限字段的菜单项
		}},
		{"$project": bson.M{
			"permissions": "$menu_details.permission", // 只保留 permissions 字段
		}},
		{"$group": bson.M{
			"_id":             nil,                             // 将所有结果聚合在一起
			"permissionsList": bson.M{"$push": "$permissions"}, // 聚集所有权限到数组
		}},
		{"$project": bson.M{
			"_id":             0, // 不需要 _id 字段
			"permissionsList": 1, // 只保留 permissionsList 数组
		}},
	}

	// 执行聚合查询
	cur, err := e.Mongo.Collection(roleModel.TableName()).Aggregate(ctx, pipeline, options.Aggregate().SetAllowDiskUse(true))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var result struct {
		PermissionsList []string `bson:"permissionsList"`
	}

	if cur.Next(ctx) {
		if err := cur.Decode(&result); err != nil {
			return nil, err
		}
	} else if err := cur.Err(); err != nil {
		return nil, err
	}

	return result.PermissionsList, nil
}
