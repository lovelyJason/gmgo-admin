package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"github.com/go-admin-team/go-admin-core/sdk/service"
	"gmgo-admin/app/admin/models"
	"gmgo-admin/app/admin/service/dto"
	"gmgo-admin/common/actions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type SysApi struct {
	service.Service
}

// GetPage 获取SysApi列表
func (e *SysApi) GetPage(c *dto.SysApiGetPageReq, p *actions.DataPermission) ([]*models.SysApi, int64, error) {
	var err error
	//var data models.SysApi
	model := &models.SysApi{}
	filter := models.SysApi{
		Title:  c.Title,
		Path:   c.Path,
		Action: c.Action,
		Type:   c.Type,
	}
	list, total, err := model.List(e.Mongo, filter, c.PageIndex, c.PageSize)
	if err != nil {
		e.Log.Errorf("Service GetSysApiPage error:%s", err)
		return nil, 0, err
	}
	return list, total, err
}

// Get 获取SysApi对象with id
func (e *SysApi) Get(ctx context.Context, d *dto.SysApiGetReq, p *actions.DataPermission, data *models.SysApi) *SysApi {
	model := &models.SysApi{}
	var err error
	//var data models.SysApi
	//err := e.Orm.Model(&data).
	//	Scopes(
	//		actions.Permission(data.TableName(), p),
	//	).
	//	FirstOrInit(model, d.GetId()).Error
	//if err != nil {
	//	e.Log.Errorf("db error:%s", err)
	//	_ = e.AddError(err)
	//	return e
	//}
	//if model.Id == 0 {
	//	err = errors.New("查看对象不存在或无权查看")
	//	e.Log.Errorf("Service GetSysApi error: %s", err)
	//	_ = e.AddError(err)
	//	return e
	//}
	//return e
	if p != nil {
		// Add any permission-related filtering logic here
	}
	apiId := d.GetId()
	err = model.GetOne(ctx, e.Mongo, apiId, data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			err = errors.New("查看对象不存在或无权查看")
			e.Log.Errorf("Service GetSysApi error: %s", err)
			_ = e.AddError(err)
			return e
		}
		e.Log.Errorf("db error: %s", err)
		_ = e.AddError(err)
		return e
	}

	return e
}

// Insert 创建SysApi对象
func (e *SysApi) Insert(c *dto.SysApiInsertReq) (int, error) {
	var err error
	var apiModel = &models.SysApi{}
	var filter = bson.M{
		"$or": []bson.M{
			// 查询 title 相同的记录数量
			{
				"title": c.Title,
			},
			// 查询 path 和 action 都相同的记录数量
			{
				"path":   c.Path,
				"action": c.Action,
			},
		},
	}
	var ctx = context.Background()
	var i int64
	err = apiModel.Count(ctx, e.Mongo, filter, &i)
	if err != nil {
		e.Log.Errorf("db error: %s", err)
		return 0, err
	}
	if i > 0 {
		err := errors.New("api标题或者api地址已存在！")
		e.Log.Errorf("db error: %s", err)
		return 0, err
	}
	insertFilter := bson.M{
		"title":      c.Title,
		"action":     c.Action,
		"handle":     c.Handle,
		"path":       c.Path,
		"type":       c.Type,
		"modifiable": 1,
		"createdAt":  time.Now(),
		"createBy":   c.CreateBy,
	}
	apiId, err := apiModel.InsertOne(ctx, e.Mongo, insertFilter)
	if err != nil {
		e.Log.Errorf("db insert api error: %s", err)
		return 0, err
	}
	return apiId, nil
}

// Update 修改SysApi对象
func (e *SysApi) Update(c *dto.SysApiUpdateReq, p *actions.DataPermission) error {
	var err error
	var apiModel = &models.SysApi{}
	ctx := context.Background()
	filter := bson.M{
		"handle":    c.Handle,
		"title":     c.Title,
		"type":      c.Type,
		"action":    c.Action,
		"path":      c.Path,
		"updatedAt": time.Now(),
		"updateBy":  c.UpdateBy,
	}
	err = apiModel.UpdateByApiId(ctx, e.Mongo, c.ApiId, filter)
	if err != nil {
		e.Log.Errorf("Service UpdateSysApi error: %s", err)
		return err
	}
	return nil
}

// Remove 删除SysApi
func (e *SysApi) Remove(d *dto.SysApiDeleteReq, p *actions.DataPermission) error {
	var err error
	apiModel := models.SysApi{}
	ctx := context.Background()

	err = apiModel.SoftDelByApiIds(ctx, e.Mongo, d.ApiIds)

	if err != nil {
		e.Log.Errorf("Error found in  RemoveSysUser : %s", err)
		return err
	}
	return nil
}

// CheckStorageSysApi 创建SysApi对象
func (e *SysApi) CheckStorageSysApi(c *[]runtime.Router) error {
	for _, v := range *c {
		err := e.Orm.Debug().Where(models.SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
			Attrs(models.SysApi{Handle: v.Handler}).
			FirstOrCreate(&models.SysApi{}).Error
		if err != nil {
			err := fmt.Errorf("Service CheckStorageSysApi error: %s \r\n ", err.Error())
			return err
		}
	}
	return nil
}
