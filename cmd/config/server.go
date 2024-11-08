package config

import (
	"encoding/json"
	"fmt"
	"github.com/go-admin-team/go-admin-core/config/source/file"
	"github.com/spf13/cobra"

	sdkConfig "github.com/go-admin-team/go-admin-core/sdk/config"
)

var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:     "config",
		Short:   "Get Application config info",
		Example: "gmgo-admin config -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
}

func run() {
	sdkConfig.Setup(file.NewSource(file.WithPath(configYml)))

	application, errs := json.MarshalIndent(sdkConfig.ApplicationConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("application:", string(application))

	jwt, errs := json.MarshalIndent(sdkConfig.JwtConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("jwt:", string(jwt))

	// todo 需要兼容
	database, errs := json.MarshalIndent(sdkConfig.DatabasesConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("database:", string(database))

	gen, errs := json.MarshalIndent(sdkConfig.GenConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("gen:", string(gen))

	loggerConfig, errs := json.MarshalIndent(sdkConfig.LoggerConfig, "", "   ") //转换成JSON返回的是byte[]
	if errs != nil {
		fmt.Println(errs.Error())
	}
	fmt.Println("logger:", string(loggerConfig))

}
