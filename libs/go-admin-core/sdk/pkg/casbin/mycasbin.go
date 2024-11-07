package mycasbin

import (
	"errors"
	mongoAdapter "github.com/casbin/mongodb-adapter/v3"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/log"
	"github.com/casbin/casbin/v2/model"
	"github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	redisWatcher "github.com/go-admin-team/redis-watcher/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	gormAdapter "github.com/go-admin-team/gorm-adapter/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Initialize the model from a string.
var text = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.sub == p.sub && (keyMatch2(r.obj, p.obj) || keyMatch(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

var (
	enforcer *casbin.SyncedEnforcer // 内置读写锁，保证并发安全,用于加载策略，判断权限，实时更新策略
	once     sync.Once
)

func SetupWithMongo(options *options.ClientOptions, dbName string) *casbin.SyncedEnforcer {
	if dbName == "" {
		panic(errors.New("MongoDB database name cannot be empty"))
	}
	once.Do(func() {
		// 选择适配器
		apter, err := mongoAdapter.NewAdapterWithClientOption(options, dbName) // 好像mongo的这个adapter不支持指定集合名称，只能是casbin_rule
		if err != nil {
			panic(err)
		}

		// 加载模型
		m, err := model.NewModelFromString(text)
		if err != nil {
			panic(err)
		}

		// 初始化 Enforcer
		enforcer, err = casbin.NewSyncedEnforcer(m, apter)
		if err != nil {
			panic(err)
		}

		// 加载策略
		err = enforcer.LoadPolicy()
		if err != nil {
			panic(err)
		}

		// 启用日志
		log.SetLogger(&Logger{})
		enforcer.EnableLog(true)

		// ================测试数据 添加一条策略（例如：角色 admin 对 /api/v1/sys-user 有 GET 权限）
		//policy := []string{"admin", "/api/v1/sys-tet", "GET"}
		//_, err = enforcer.AddNamedPolicy("p", policy)
		//if err != nil {
		//	fmt.Println("Error adding policy:", err)
		//} else {
		//	fmt.Println("Policy added successfully")
		//}
		//
		//// 检查当前策略是否正确添加
		//policies := enforcer.GetPolicy()
		//fmt.Println("Current policies:", policies)
		//
		//// 保存策略到数据库
		//err = enforcer.SavePolicy()
		//if err != nil {
		//	fmt.Println("Error saving policy:", err)
		//} else {
		//	fmt.Println("Policies saved successfully")
		//}
		// ============================
	})

	return enforcer
}

func Setup(db *gorm.DB, _ string) *casbin.SyncedEnforcer {
	once.Do(func() {
		// 1. 初始化适配器，并存储在数据库表中
		Apter, err := gormAdapter.NewAdapterByDBUseTableName(db, "sys", "casbin_rule")
		if err != nil && err.Error() != "invalid DDL" {
			panic(err)
		}
		// 2. 创建casbin模型，包含请求，策略，匹配条件
		m, err := model.NewModelFromString(text)
		if err != nil {
			panic(err)
		}
		// 3. 创建enforcer
		enforcer, err = casbin.NewSyncedEnforcer(m, Apter)
		if err != nil {
			panic(err)
		}
		// 4. 从适配器（数据库）中加载权限策略存储于内存
		err = enforcer.LoadPolicy()
		if err != nil {
			panic(err)
		}
		// 5. 从redis动态观察，TODO:
		// set redis watcher if redis config is not nil
		if config.CacheConfig.Redis != nil {
			w, err := redisWatcher.NewWatcher(config.CacheConfig.Redis.Addr, redisWatcher.WatcherOptions{
				Options: redis.Options{
					Network:  "tcp",
					Password: config.CacheConfig.Redis.Password,
				},
				Channel:    "/casbin",
				IgnoreSelf: false,
			})
			if err != nil {
				panic(err)
			}

			err = w.SetUpdateCallback(updateCallback)
			if err != nil {
				panic(err)
			}
			err = enforcer.SetWatcher(w)
			if err != nil {
				panic(err)
			}
		}

		log.SetLogger(&Logger{})
		enforcer.EnableLog(true)
	})

	return enforcer
}

func updateCallback(msg string) {
	l := logger.NewHelper(sdk.Runtime.GetLogger())
	l.Infof("casbin updateCallback msg: %v", msg)
	err := enforcer.LoadPolicy()
	if err != nil {
		l.Errorf("casbin LoadPolicy err: %v", err)
	}
}
