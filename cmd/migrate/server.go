package migrate

import (
	"bytes"
	"fmt"
	"github.com/go-admin-team/go-admin-core/config/source/file"
	"github.com/go-admin-team/go-admin-core/sdk"
	sdkConfig "github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/spf13/cobra"
	"gmgo-admin/cmd/migrate/migration"
	_ "gmgo-admin/cmd/migrate/migration/version"
	_ "gmgo-admin/cmd/migrate/migration/version-local"
	"gmgo-admin/common/database"
	"gmgo-admin/common/models"
	"strconv"
	"text/template"
	"time"
)

var (
	configYml string
	generate  bool
	goAdmin   bool
	host      string
	StartCmd  = &cobra.Command{
		Use:     "migrate",
		Short:   "Initialize the database",
		Example: "gmgo-admin migrate -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

type Settings struct {
	Application ApplicationConfig `yaml:"application"`
}

type ApplicationConfig struct {
	Mode string `yaml:"mode"`
	Host string `yaml:"host"`
}

// fixme 在您看不见代码的时候运行迁移，我觉得是不安全的，所以编译后最好不要去执行迁移
func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().BoolVarP(&generate, "generate", "g", false, "generate migration file")
	StartCmd.PersistentFlags().BoolVarP(&goAdmin, "goAdmin", "a", false, "generate go-admin migration file")
	StartCmd.PersistentFlags().StringVarP(&host, "domain", "d", "*", "select tenant host")
}

func run() {
	fmt.Println("当前生效的configYml为：", configYml)
	LoggerConfig := sdkConfig.LoggerConfig
	fmt.Println(LoggerConfig) // 注意：需要调用sdkConfig.setup后才能拿到数据

	if !generate {
		fmt.Println(`start init`)
		//1. 读取配置,Setup为config注入数据
		sdkConfig.Setup(
			file.NewSource(file.WithPath(configYml)),
			initDB,
		)
	} else {
		fmt.Println(`generate migration file`)
		_ = genFile()
	}
	//fmt.Println(sdkConfig.DatabasesConfig)
	//var cfg appConfig.AppConfig
	//config.Scan(&cfg)
	fmt.Println("==========")
}

func migrateModel() error {
	if host == "" {
		host = "*"
	}
	db := sdk.Runtime.GetDbByKey(host)
	if db == nil {
		if len(sdk.Runtime.GetDb()) == 1 && host == "*" {
			for k, v := range sdk.Runtime.GetDb() {
				db = v
				host = k
				break
			}
		}
	}
	if db == nil {
		return fmt.Errorf("未找到数据库配置")
	}
	if sdkConfig.DatabasesConfig[host].Driver == "mysql" {
		//初始化数据库时候用
		db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	// 自动迁移数据库表与模型一致
	err := db.Debug().AutoMigrate(&models.Migration{})
	if err != nil {
		return err
	}
	// 设置或更新该迁移实例使用的数据库连接
	migration.Migrate.SetDb(db.Debug())
	//收集并排序所有版本，然后检查数据库表 sys_migration 中是否存在这些版本。
	//如果版本不存在，则会调用对应的迁移函数来执行迁移
	migration.Migrate.Migrate()
	return err
}
func initDB() {
	//3. 初始化数据库链接, 已扩展mongo连接
	database.Setup()
	//4. 数据库迁移
	fmt.Println("数据库迁移开始")
	_ = migrateModel()
	fmt.Println(`数据库基础数据初始化成功`)
}

func genFile() error {
	t1, err := template.ParseFiles("template/migrate.template")
	if err != nil {
		return err
	}
	m := map[string]string{}
	m["GenerateTime"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	m["Package"] = "version_local"
	if goAdmin {
		m["Package"] = "version"
	}
	var b1 bytes.Buffer
	err = t1.Execute(&b1, m)
	if goAdmin {
		pkg.FileCreate(b1, "./cmd/migrate/migration/version/"+m["GenerateTime"]+"_migrate.go")
	} else {
		pkg.FileCreate(b1, "./cmd/migrate/migration/version-local/"+m["GenerateTime"]+"_migrate.go")
	}
	return nil
}
