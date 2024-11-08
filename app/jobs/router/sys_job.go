package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"gmgo-admin/app/jobs/apis"
	models2 "gmgo-admin/app/jobs/models"
	dto2 "gmgo-admin/app/jobs/service/dto"
	"gmgo-admin/common/actions"
	"gmgo-admin/common/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerSysJobRouter)
}

// 需认证的路由代码
func registerSysJobRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r := v1.Group("/sysjob").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
	{
		sysJob := &models2.SysJob{}
		r.GET("", actions.PermissionAction(), actions.IndexAction(sysJob, new(dto2.SysJobSearch), func() interface{} {
			list := make([]models2.SysJob, 0)
			return &list
		}))
		r.GET("/:id", actions.PermissionAction(), actions.ViewAction(new(dto2.SysJobById), func() interface{} {
			return &dto2.SysJobItem{}
		}))
		r.POST("/add", actions.CreateAction(new(dto2.SysJobControl)))
		r.POST("/edit", actions.PermissionAction(), actions.UpdateAction(new(dto2.SysJobControl)))
		r.POST("/del", actions.PermissionAction(), actions.DeleteAction(new(dto2.SysJobById)))
	}
	sysJob := apis.SysJob{}

	v1.GET("/job/remove/:id", sysJob.RemoveJobForService)
	v1.GET("/job/start/:id", sysJob.StartJobForService)
}
