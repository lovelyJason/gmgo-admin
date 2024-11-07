package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"strings"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"go-admin/app/admin/models"
	"go-admin/app/admin/service/dto"
	cModels "go-admin/common/models"

	"github.com/go-admin-team/go-admin-core/sdk/service"
)

type SysMenu struct {
	service.Service
}

// GetPage 获取SysMenu列表
func (e *SysMenu) GetPage(c *dto.SysMenuGetPageReq, menus *[]models.SysMenu) *SysMenu {
	var menu = make([]models.SysMenu, 0)
	err := e.getPage(c, &menu).Error
	if err != nil {
		_ = e.AddError(err)
		return e
	}
	for i := 0; i < len(menu); i++ {
		if menu[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&menu, menu[i])
		*menus = append(*menus, menusInfo)
	}
	return e
}

// getPage 菜单分页列表
func (e *SysMenu) getPage(c *dto.SysMenuGetPageReq, list *[]models.SysMenu) *SysMenu {
	var err error
	//var data models.SysMenu
	model := &models.SysMenu{}
	filter := bson.M{}
	ctx := context.Background()

	//err = e.Orm.Model(&data).
	//	Scopes(
	//		cDto.OrderDest("sort", false),
	//		cDto.MakeCondition(c.GetNeedSearch()),
	//	).Preload("SysApi").
	//	Find(list).Error
	//if err != nil {
	//	e.Log.Errorf("getSysMenuPage error:%s", err)
	//	_ = e.AddError(err)
	//	return e
	//}
	//return e
	err = model.List(ctx, e.Mongo, filter, 0, 0, list)
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}
	return e
}

// Get 获取SysMenu对象
func (e *SysMenu) Get(ctx context.Context, d *dto.SysMenuGetReq, model *models.SysMenu) *SysMenu {
	var err error

	err = model.GetOneByMenuId(ctx, e.Mongo, d.Id, model)
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		err = errors.New("查看对象不存在或无权查看")
		e.Log.Errorf("GetSysMenu error:%s", err)
		_ = e.AddError(err)
		return e
	}
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
		return e
	}

	return e
}

// Insert 创建SysMenu对象
func (e *SysMenu) Insert(c *dto.SysMenuInsertReq) *SysMenu {
	var err error
	var data models.SysMenu
	ctx := context.Background()
	c.Generate(&data)

	// 1. mysql操作：拉取sysMenu上的apis，然后查表sysApi，然后更新菜单表和菜单api关联表
	menuModel := &models.SysMenu{}
	insertFilter := c.GenerateBson()

	newMenuId, err := menuModel.InsertOne(ctx, e.Mongo, insertFilter)
	data.MenuId = newMenuId
	e.initPaths(ctx, &data)

	if err != nil {
		e.Log.Errorf("db error:%s", err)
		_ = e.AddError(err)
	}
	fmt.Println("newMenuId", newMenuId)
	return e
}

func (e *SysMenu) initPaths(ctx context.Context, menu *models.SysMenu) error {
	var err error
	menuModel := &models.SysMenu{}
	parentMenu := new(models.SysMenu)
	if menu.ParentId != 0 {
		// 查询到parentMenu对象

		err = menuModel.GetOneByMenuId(ctx, e.Mongo, menu.ParentId, parentMenu)
		if err != nil {
			return err
		}
		if parentMenu.Paths == "" {
			err = errors.New("父级paths异常，请尝试对当前节点父级菜单进行更新操作！")
			return err
		}
		// 拼接父菜单paths和自身菜单id，如/0/2/自身的自增id
		menu.Paths = parentMenu.Paths + "/" + pkg.IntToString(menu.MenuId)
	} else {
		menu.Paths = "/0/" + pkg.IntToString(menu.MenuId)
	}
	// 更新sysMenu表的paths字段, 放在外面一起更新吧
	menuModel.UpdateByMenuId(ctx, e.Mongo, menu.MenuId, bson.M{
		"paths": menu.Paths,
	})
	return err
}

// Update 修改SysMenu对象
func (e *SysMenu) Update(c *dto.SysMenuUpdateReq) *SysMenu {
	var err error

	var model = models.SysMenu{}
	// 1.先查询单条信息, 合并数据库现有字段与修改字段
	menuModel := &models.SysMenu{}
	ctx := context.Background()
	err = menuModel.GetOneByMenuId(ctx, e.Mongo, c.MenuId, &model)
	if err != nil {
		e.Log.Errorf("get menu error:%s", err)
		_ = e.AddError(err)
		return e
	}

	oldPath := model.Paths
	updateMenuFilter := c.GenerateBson()
	//jsonData, err := bson.MarshalExtJSON(updateMenuFilter, false, true)
	//if err == nil {
	//	fmt.Println("JSON Format:", string(jsonData))
	//}
	//model.SysApi = apiList // mongo聚合查询自带
	// 2.修改这条记录
	menuModel.UpdateByMenuId(ctx, e.Mongo, model.MenuId, updateMenuFilter)
	var menuList []models.SysMenu
	// 3.查询所有以oldPath为前缀的menu, 修改所有子menu 的path，如果不同的话
	filter := bson.M{
		"paths": bson.M{
			"$regex":   "^" + oldPath, // 以 oldPath 为前缀
			"$options": "i",           // 可选：不区分大小写
		},
	}
	menuModel.Get(ctx, e.Mongo, filter, &menuList)
	for _, v := range menuList {
		v.Paths = strings.Replace(v.Paths, oldPath, model.Paths, 1)
		updateMenuFilter = bson.M{"paths": strings.Replace(v.Paths, oldPath, model.Paths, 1)}
		menuModel.UpdateByMenuId(ctx, e.Mongo, v.MenuId, updateMenuFilter)
	}
	return e
}

// Remove 删除SysMenu
func (e *SysMenu) Remove(d *dto.SysMenuDeleteReq) *SysMenu {
	var err error
	menuModal := models.SysMenu{}
	ctx := context.Background()

	err = menuModal.SoftDelByUserIds(ctx, e.Mongo, d.MenuIds)

	if err != nil {
		e.Log.Errorf("Error found in  RemoveMenu : %s", err)
		_ = e.AddError(err)
	}
	return e
}

// GetList 获取菜单数据
func (e *SysMenu) GetList(c *dto.SysMenuGetPageReq, list *[]models.SysMenu) error {
	var err error
	//var data models.SysMenu

	//err = e.Orm.Model(&data).
	//	Scopes(
	//		cDto.MakeCondition(c.GetNeedSearch()),
	//	).
	//	Find(list).Error
	//if err != nil {
	//	e.Log.Errorf("db error:%s", err)
	//	return err
	//}
	model := &models.SysMenu{}
	filter := bson.M{}
	ctx := context.Background()
	err = model.List(ctx, e.Mongo, filter, 0, 0, list)
	if err != nil {
		e.Log.Errorf("db error:%s", err)
		return err
	}
	return nil
}

// SetLabel 修改角色中 设置菜单基础数据
func (e *SysMenu) SetLabel() (m []dto.MenuLabel, err error) {
	var list []models.SysMenu
	err = e.GetList(&dto.SysMenuGetPageReq{}, &list)

	m = make([]dto.MenuLabel, 0)
	for i := 0; i < len(list); i++ {
		if list[i].ParentId != 0 {
			continue
		}
		e := dto.MenuLabel{}
		e.Id = list[i].Id
		e.MenuId = list[i].MenuId
		e.Label = list[i].Title
		deptsInfo := menuLabelCall(&list, e)

		m = append(m, deptsInfo)
	}
	return
}

// GetSysMenuByRoleName 左侧菜单
func (e *SysMenu) GetSysMenuByRoleName(roleName ...string) ([]models.SysMenu, error) {
	var MenuList []models.SysMenu
	var role models.SysRole
	var err error
	admin := false
	for _, s := range roleName {
		if s == "admin" {
			admin = true
		}
	}

	if len(roleName) > 0 && admin {
		var data []models.SysMenu
		err = e.Orm.Where(" menu_type in ('M','C')").
			Order("sort").
			Find(&data).
			Error
		MenuList = data
	} else {
		err = e.Orm.Model(&role).Preload("SysMenu", func(db *gorm.DB) *gorm.DB {
			return db.Where(" menu_type in ('M','C')").Order("sort")
		}).Where("role_name in ?", roleName).Find(&role).
			Error
		MenuList = *role.SysMenu
	}

	if err != nil {
		e.Log.Errorf("db error:%s", err)
	}
	return MenuList, err
}

// menuLabelCall 递归构造组织数据
func menuLabelCall(eList *[]models.SysMenu, dept dto.MenuLabel) dto.MenuLabel {
	list := *eList

	min := make([]dto.MenuLabel, 0)
	for j := 0; j < len(list); j++ {

		if dept.MenuId != list[j].ParentId {
			continue
		}
		mi := dto.MenuLabel{}
		mi.Id = list[j].Id
		mi.MenuId = list[j].MenuId
		mi.Label = list[j].Title
		mi.Children = []dto.MenuLabel{}
		if list[j].MenuType != "F" {
			ms := menuLabelCall(eList, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	if len(min) > 0 {
		dept.Children = min
	} else {
		dept.Children = nil
	}
	return dept
}

// menuCall 构建菜单树
func menuCall(menuList *[]models.SysMenu, menu models.SysMenu) models.SysMenu {
	list := *menuList

	min := make([]models.SysMenu, 0)
	for i := 0; i < len(list); i++ {

		if menu.MenuId != list[i].ParentId {
			continue
		}
		mi := models.SysMenu{}
		mi.Id = list[i].Id
		mi.MenuId = list[i].MenuId
		mi.MenuName = list[i].MenuName
		mi.Title = list[i].Title
		mi.Icon = list[i].Icon
		mi.Path = list[i].Path
		mi.MenuType = list[i].MenuType
		mi.Action = list[i].Action
		mi.Permission = list[i].Permission
		mi.ParentId = list[i].ParentId
		mi.NoCache = list[i].NoCache
		mi.Breadcrumb = list[i].Breadcrumb
		mi.Component = list[i].Component
		mi.Sort = list[i].Sort
		mi.Visible = list[i].Visible
		mi.CreatedAt = list[i].CreatedAt
		mi.UpdatedAt = list[i].UpdatedAt
		mi.DeletedAt = list[i].DeletedAt
		mi.SysApi = list[i].SysApi
		mi.Children = []models.SysMenu{}

		if mi.MenuType != cModels.Button { // 菜单不是按钮类型，需要进一步处理子菜单
			ms := menuCall(menuList, mi)
			min = append(min, ms)
		} else {
			min = append(min, mi)
		}
	}
	menu.Children = min
	return menu
}

func menuDistinct(menuList []models.SysMenu) (result []models.SysMenu) {
	distinctMap := make(map[int]struct{}, len(menuList))
	for _, menu := range menuList {
		if _, ok := distinctMap[menu.MenuId]; !ok {
			distinctMap[menu.MenuId] = struct{}{}
			result = append(result, menu)
		}
	}
	return result
}

func recursiveSetMenu(db *mongo.Database, mIds []int, menus *[]models.SysMenu) error {
	if len(mIds) == 0 || menus == nil {
		return nil
	}
	menuModel := &models.SysMenu{}
	// 构建查询条件
	filter := bson.M{
		"menuType": bson.M{"$in": []string{cModels.Directory, cModels.Menu, cModels.Button}},
		"menuId":   bson.M{"$in": mIds}, // 查找 menu_id 在给定列表中的文档
	}

	// 查找菜单
	var subMenus []models.SysMenu
	cursor, err := db.Collection(menuModel.TableName()).Find(context.Background(), filter)
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &subMenus); err != nil {
		return err
	}
	subIds := make([]int, 0)
	for _, menu := range subMenus {
		if menu.ParentId != 0 {
			subIds = append(subIds, menu.ParentId)
		}
		if menu.MenuType != cModels.Button {
			*menus = append(*menus, menu)
		}
	}
	return recursiveSetMenu(db, subIds, menus)
}

// SetMenuRole 获取左侧菜单树使用
func (e *SysMenu) SetMenuRole(roleName string) (m []models.SysMenu, err error) {
	menus, err := e.getByRoleName(roleName)
	m = make([]models.SysMenu, 0)
	for i := 0; i < len(menus); i++ {
		if menus[i].ParentId != 0 {
			continue
		}
		menusInfo := menuCall(&menus, menus[i])
		m = append(m, menusInfo)
	}
	return
}

func (e *SysMenu) getByRoleName(roleName string) ([]models.SysMenu, error) {
	var data []models.SysMenu
	roleModel := &models.SysRole{}
	var err error

	ctx := context.Background() // 创建一个上下文
	// 根据sysRoleMenu表聚合查询查到普通角色的菜单返回
	if roleName == "admin" {
		filter := bson.M{
			"menuType":  bson.M{"$in": []string{"M", "C"}},
			"deletedAt": nil,
		}
		cursor, err := e.Mongo.Collection("sysMenu").Find(ctx, filter, options.Find().SetSort(bson.M{"sort": 1}))
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &data); err != nil {
			return nil, err
		}
	} else {
		filter := bson.M{
			"roleKey": roleName,
		}
		err = roleModel.GetOne(ctx, e.Mongo, filter, roleModel)
		if roleModel.SysMenu != nil {
			mIds := make([]int, 0)
			for _, menu := range *roleModel.SysMenu {
				mIds = append(mIds, menu.MenuId)
			}

			if err := recursiveSetMenu(e.Mongo, mIds, &data); err != nil {
				return nil, err
			}

			data = menuDistinct(data)
		}
	}

	sort.Sort(models.SysMenuSlice(data))
	return data, err
}
