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
	"go.mongodb.org/mongo-driver/mongo"
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
		Type: c.Type,
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
	//var err error
	//var data models.SysApi
	//var apiModel = &models.SysApi{}
	//var filter = bson.M{}
	return 0, nil
}

// Update 修改SysApi对象
func (e *SysApi) Update(c *dto.SysApiUpdateReq, p *actions.DataPermission) error {
	var model = models.SysApi{}
	db := e.Orm.Debug().First(&model, c.GetId())
	if db.RowsAffected == 0 {
		return errors.New("无权更新该数据")
	}
	c.Generate(&model)
	db = e.Orm.Save(&model)
	if err := db.Error; err != nil {
		e.Log.Errorf("Service UpdateSysApi error:%s", err)
		return err
	}

	return nil
}

// Remove 删除SysApi
func (e *SysApi) Remove(d *dto.SysApiDeleteReq, p *actions.DataPermission) error {
	var data models.SysApi

	db := e.Orm.Model(&data).
		Scopes(
			actions.Permission(data.TableName(), p),
		).Delete(&data, d.GetId())
	if err := db.Error; err != nil {
		e.Log.Errorf("Service RemoveSysApi error:%s", err)
		return err
	}
	if db.RowsAffected == 0 {
		return errors.New("无权删除该数据")
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
