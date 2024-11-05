package service

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"gorm.io/gorm/clause"

	"github.com/casbin/casbin/v2"

	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
)

type SysRole struct {
	service.Service
}

// GetPage 获取SysRole列表
func (e *SysRole) GetPage(c *dto.SysRoleGetPageReq) ([]*models.SysRole, int64, error) {
	var err error
	//var data models.SysRole
	model := &models.SysRole{}
	filter := models.SysRole{}

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
	err = menuModel.List(ctx, e.Mongo, menuFilter, 0, 0, &dataMenu)
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

// Update 修改SysRole对象
func (e *SysRole) Update(c *dto.SysRoleUpdateReq, cb *casbin.SyncedEnforcer) error {
	var err error
	var model = models.SysRole{}
	var mlist = make([]models.SysMenu, 0)
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
	// TODO:POLICY重构
	// 清除 sys_casbin_rule 权限表里 当前角色的所有记录
	_, err = cb.RemoveFilteredPolicy(0, model.RoleKey)
	if err != nil {
		e.Log.Errorf("delete policy error:%s", err)
		return err
	}
	mp := make(map[string]interface{}, 0)
	polices := make([][]string, 0) // 类似于["角色名称", "/api/v1/sys-user", "GET"]
	for _, menu := range mlist {
		for _, api := range menu.SysApi {
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
	var model = models.SysRole{}
	tx.Preload("SysMenu").Preload("SysDept").First(&model, c.GetId())
	//删除 SysRole 时，同时删除角色所有 关联其它表 记录 (SysMenu 和 SysMenu)
	db := tx.Select(clause.Associations).Delete(&model)

	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}

	// 清除 sys_casbin_rule 权限表里 当前角色的所有记录
	_, _ = cb.RemoveFilteredPolicy(0, model.RoleKey)

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
	var model = models.SysRole{}
	tx.First(&model, c.GetId())
	c.Generate(&model)
	// 更新关联的数据，使用 FullSaveAssociations 模式
	db := tx.Session(&gorm.Session{FullSaveAssociations: true}).Debug().Save(&model)
	if err = db.Error; err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	return nil
}

// GetById 获取SysRole对象
func (e *SysRole) GetById(roleId int) ([]string, error) {
	permissions := make([]string, 0)
	model := models.SysRole{}
	model.RoleId = roleId
	if err := e.Orm.Model(&model).Preload("SysMenu").First(&model).Error; err != nil {
		return nil, err
	}
	l := *model.SysMenu
	for i := 0; i < len(l); i++ {
		if l[i].Permission != "" {
			permissions = append(permissions, l[i].Permission)
		}
	}
	return permissions, nil
}
