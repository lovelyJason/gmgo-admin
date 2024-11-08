package cmd

import (
	"errors"
	"fmt"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"gmgo-admin/cmd/app"
	"gmgo-admin/common/global"
	"gmgo-admin/initialize"
	"os"

	"github.com/spf13/cobra"

	"gmgo-admin/cmd/api"
	"gmgo-admin/cmd/config"
	"gmgo-admin/cmd/migrate"
	"gmgo-admin/cmd/version"
)

var rootCmd = &cobra.Command{
	Use:          "gmgo-admin",
	Short:        "gmgo-admin",
	SilenceUsage: true,
	Long:         `gmgo-admin`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(pkg.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	PreRun: func(cmd *cobra.Command, args []string) {
		setup()
	},
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func setup() {
	//initialize.OtherInit()
}

func tip() {
	usageStr := `欢迎使用 ` + pkg.Green(`go-admin `+global.Version) + ` 可以使用 ` + pkg.Red(`-h`) + ` 查看命令`
	usageStr1 := `也可以参考 https://doc.go-admin.dev/guide/ksks 的相关内容`
	fmt.Printf("%s\n", usageStr)
	fmt.Printf("%s\n", usageStr1)
}

func init() {
	initialize.OtherInit()
	rootCmd.AddCommand(api.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
	rootCmd.AddCommand(version.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
	rootCmd.AddCommand(app.StartCmd)
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
