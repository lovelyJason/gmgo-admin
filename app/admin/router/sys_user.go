package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"go-admin/app/admin/apis"
	"go-admin/common/actions"
	"go-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysUserRouter)
}

// 需认证的路由代码
func registerSysUserRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	api := apis.SysUser{}
	r := v1.Group("/sys-user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionAction())
	{
		r.GET("/get", api.GetPage)
		r.POST("/get/:id", api.Get)
		r.POST("/add", api.Insert)
		r.POST("/edit", api.Update)
		r.POST("/del", api.Delete)
	}

	user := v1.Group("/user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole()).Use(actions.PermissionAction())
	{
		user.GET("/profile", api.GetProfile)
		user.POST("/avatar", api.InsetAvatar)
		user.POST("/pwd/set", api.UpdatePwd)
		user.POST("/pwd/reset", api.ResetPwd)
		user.PUT("/status", api.UpdateStatus)
	}
	v1auth := v1.Group("").Use(authMiddleware.MiddlewareFunc())
	{
		v1auth.GET("/getinfo", api.GetInfo)
	}
}
