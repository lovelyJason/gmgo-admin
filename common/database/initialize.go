package database

import (
	"context"
	"fmt"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	toolsConfig "github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	mycasbin "github.com/go-admin-team/go-admin-core/sdk/pkg/casbin"
	toolsDB "github.com/go-admin-team/go-admin-core/tools/database"
	. "github.com/go-admin-team/go-admin-core/tools/gorm/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"

	"go-admin/common/global"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Setup 配置数据库
func Setup() {
	fmt.Println("toolsConfig.DatabasesConfig", len(toolsConfig.DatabasesConfig))
	for k := range toolsConfig.DatabasesConfig {
		setupSimpleDatabase(k, toolsConfig.DatabasesConfig[k])
	}
}

func setupSimpleDatabase(host string, c *toolsConfig.Database) {
	if global.Driver == "" {
		global.Driver = c.Driver
	}
	log.Infof("%s => %s", host, pkg.Green(c.Source))
	if c.Driver == "mongo" {
		mongoClientOptions := options.Client().ApplyURI(c.Source).SetConnectTimeout(10 * time.Second)

		// 创建 MongoDB 客户端
		client, err := mongo.Connect(context.TODO(), mongoClientOptions)
		if err != nil {
			log.Fatal(pkg.Red("MongoDB connect error: "), err)
		}

		// Ping the MongoDB server
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(pkg.Red("MongoDB ping error: "), err)
		}
		log.Info(pkg.Green("MongoDB connect success !"))

		// 这里可以设置 MongoDB 相关的操作，如：
		// - 创建集合
		// - 设置相关的中间件等
		// - 其他数据库操作

		// 在这里你可能想要将 MongoDB 客户端存储到 sdk.Runtime 中
		sdk.Runtime.SetMongoClient(host, client)
		db := client.Database(c.Db)
		sdk.Runtime.SetMongoDb(host, db)

		e := mycasbin.SetupWithMongo(mongoClientOptions, "sysCasbinRule")
		sdk.Runtime.SetCasbin(host, e)
	} else {
		registers := make([]toolsDB.ResolverConfigure, len(c.Registers))
		for i := range c.Registers {
			registers[i] = toolsDB.NewResolverConfigure(
				c.Registers[i].Sources,
				c.Registers[i].Replicas,
				c.Registers[i].Policy,
				c.Registers[i].Tables)
		}
		resolverConfig := toolsDB.NewConfigure(c.Source, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxIdleTime, c.ConnMaxLifeTime, registers)
		db, err := resolverConfig.Init(&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: New(
				logger.Config{
					SlowThreshold: time.Second,
					Colorful:      true,
					LogLevel: logger.LogLevel(
						log.DefaultLogger.Options().Level.LevelForGorm()),
				},
			),
		}, opens[c.Driver])

		if err != nil {
			log.Fatal(pkg.Red(c.Driver+" connect error :"), err)
		} else {
			log.Info(pkg.Green(c.Driver + " connect success !"))
		}

		e := mycasbin.Setup(db, "")

		sdk.Runtime.SetDb(host, db)
		sdk.Runtime.SetCasbin(host, e)
	}
}
