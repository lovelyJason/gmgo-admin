package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	sdkConfig "github.com/go-admin-team/go-admin-core/sdk/config"
)

func WithContextDb(c *gin.Context) {
	databaseConfig := sdkConfig.DatabaseConfig
	fmt.Println(databaseConfig)
	if databaseConfig.Driver == "mongo" {
		db := sdk.Runtime.GetMongoDbByKey(c.Request.Host)
		client := sdk.Runtime.GetMongoClientByKey(c.Request.Host)
		c.Set("db", db)
		c.Set("client", client)
	} else {
		db := sdk.Runtime.GetDbByKey(c.Request.Host)
		c.Set("db", db.WithContext(c))
	}

	c.Next()
}
