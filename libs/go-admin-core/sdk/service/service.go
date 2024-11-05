package service

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/storage"
	"gorm.io/gorm"
)

type Service struct {
	Orm   *gorm.DB
	Mongo *mongo.Database
	Msg   string
	MsgID string
	Log   *logger.Helper
	Error error
	Cache storage.AdapterCache
}

func (db *Service) AddError(err error) error {
	if db.Error == nil {
		db.Error = err
	} else if err != nil {
		db.Error = fmt.Errorf("%v; %w", db.Error, err)
	}
	return db.Error
}
