package initialize

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"go-admin/common/global"
	"go-admin/utils"
)

func OtherInit() {
	dr, err := utils.ParseDuration("7d")
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration("1d")
	if err != nil {
		panic(err)
	}
	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
}