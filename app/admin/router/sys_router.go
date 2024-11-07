package router

import (
	"go-admin/app/admin/apis"
	"mime"

	"github.com/go-admin-team/go-admin-core/sdk/config"

	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/ws"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"

	"go-admin/common/middleware"
	"go-admin/common/middleware/handler"
	_ "go-admin/docs/admin"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")
	sysBaseRouter(g)
	// 静态文件
	sysStaticFileRouter(g)
	// swagger；注意：生产环境可以注释掉
	if config.ApplicationConfig.Mode != "prod" {
		sysSwaggerRouter(g)
	}
	// 需要认证
	sysCheckRoleRouterInit(g, authMiddleware)
	return g
}

func sysBaseRouter(r *gin.RouterGroup) {
	//后台任务，独立运行的服务，负责 WebSocket 连接的管理、消息的发送、客户端的注册和注销等
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendAllService()

	if config.ApplicationConfig.Mode != "prod" {
		r.GET("/", apis.GoAdmin)
	}
	r.GET("/info", handler.Ping)
}

func sysStaticFileRouter(r *gin.RouterGroup) {
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		return
	}
	r.Static("/static", "./static")
	if config.ApplicationConfig.Mode != "prod" {
		r.Static("/form-generator", "./static/form-generator")
	}
}

func sysSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/admin/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler(), ginSwagger.InstanceName("admin")))
}

func sysCheckRoleRouterInit(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	// 前端js的websocket无法设置请求头https://stackoverflow.com/questions/4361173/http-headers-in-websockets-client-api， 只有服务端请求服务端的环境才能自由携带头部
	wss := r.Group("").Use(authMiddleware.MiddlewareFunc())
	{
		// 创建websocket链接
		wss.GET("/ws/:id/:channel", func(c *gin.Context) {
			ws.WebsocketManager.WsClient(c, authMiddleware, 1)
		})
		wss.GET("/wslogout/:id/:channel", func(c *gin.Context) {
			ws.WebsocketManager.UnWsClient(c, authMiddleware, 1)
		})
	}
	wss2 := r.Group("/ws2")
	// 只能传递jwt.New的鉴权中间件实例进去，因为内部依赖mw的鉴权token
	{
		wss2.GET("/ws/:id/:channel", func(c *gin.Context) {
			ws.WebsocketManager.WsClient(c, authMiddleware, 2)
		})
		wss2.GET("/wslogout/:id/:channel", func(c *gin.Context) {
			ws.WebsocketManager.UnWsClient(c, authMiddleware, 2)
		})
	}

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		// Refresh time can be longer than token timeout
		v1.GET("/refresh_token", authMiddleware.RefreshHandler)
	}
	registerBaseRouter(v1, authMiddleware)
}

func registerBaseRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysMenu{}
	api2 := apis.SysDept{}
	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		v1auth.GET("/roleMenuTreeselect/:roleId", api.GetMenuTreeSelect)
		//v1.GET("/menuTreeselect", api.GetMenuTreeSelect)
		v1auth.GET("/roleDeptTreeselect/:roleId", api2.GetDeptTreeRoleSelect)
		v1auth.POST("/logout", handler.LogOut)
	}
}
