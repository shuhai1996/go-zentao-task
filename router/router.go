package router

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"go-zentao-task/controller"
	"go-zentao-task/controller/admin"
	"go-zentao-task/controller/docs"
	"go-zentao-task/core"
	"go-zentao-task/middleware"
	"net/http"
)

func Register(env string) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(static.Serve("/g/static", static.LocalFile("static", false)))

	r.LoadHTMLGlob("template/**/*")
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "error/404", nil)
	})
	r.Any("/g/", func(c *gin.Context) {
		c.JSON(http.StatusOK, nil)
	})
	r.GET("/g/healthcheck", controller.HealthCheck)

	if env != "production" {
		r.GET("/g/docs/swagger", core.Handle(middleware.CORS), docs.Swagger)
	}

	r.Use(middleware.SetV())
	r.Use(core.Handle(middleware.Logging))
	r.Use(core.Handle(middleware.I18n))
	r.Use(core.Handle(middleware.Session))
	r.Use(middleware.Secure())

	g := r.Group("/g")

	g.POST("/admin/auth/login", core.Handle(new(admin.Auth).Login))
	ad := g.Group(
		"/admin",
		core.Handle(middleware.AdminAuth),
		core.Handle(middleware.RBAC),
	)
	{
		//菜单
		ad.POST("menu/add", core.Handle(new(admin.Menu).Create))
		ad.POST("menu/list", core.Handle(new(admin.Menu).Index))
		ad.POST("menu/edit", core.Handle(new(admin.Menu).Update))
		ad.POST("menu/set_status", core.Handle(new(admin.Menu).UpdateStatus))
		ad.POST("menu/view", core.Handle(new(admin.Menu).View))

		//角色
		ad.POST("role/add", core.Handle(new(admin.Role).Create))
		ad.POST("role/list", core.Handle(new(admin.Role).Index))
		ad.POST("role/edit", core.Handle(new(admin.Role).Update))
		ad.POST("role/set_status", core.Handle(new(admin.Role).UpdateStatus))
		ad.POST("role/view", core.Handle(new(admin.Role).View))

		//接口
		ad.POST("interface/add", core.Handle(new(admin.Interface).Create))
		ad.POST("interface/list", core.Handle(new(admin.Interface).Index))
		ad.POST("interface/edit", core.Handle(new(admin.Interface).Update))
		ad.POST("interface/set_status", core.Handle(new(admin.Interface).UpdateStatus))
		ad.POST("interface/view", core.Handle(new(admin.Interface).View))

		//用户
		ad.POST("user/add", core.Handle(new(admin.User).Create))
		ad.POST("user/list", core.Handle(new(admin.User).Index))
		ad.POST("user/edit", core.Handle(new(admin.User).Update))
		ad.POST("user/set_status", core.Handle(new(admin.User).UpdateStatus))
		ad.POST("user/view", core.Handle(new(admin.User).View))

		//接口权限分配
		ad.POST("role_interface/add", core.Handle(new(admin.RoleInterface).Add))
		ad.POST("role_interface/view", core.Handle(new(admin.RoleInterface).View))

		//角色权限分配
		ad.POST("role_menu/add", core.Handle(new(admin.RoleMenu).Add))
		ad.POST("role_menu/view", core.Handle(new(admin.RoleMenu).View))

		ad.POST("auth/logout", core.Handle(new(admin.Auth).Logout))
		ad.POST("auth/info", core.Handle(new(admin.Auth).GetAuthInfo)) //菜单列表+用户信息
	}

	g.GET("/v1/v2/i18n", core.Handle(func(c *core.Context) {
		c.String(http.StatusOK, c.Tr("author.info", 18)+" "+c.Tr("section.language"))
	}))

	g.GET("/v1/v2/session", core.Handle(func(c *core.Context) {
		var count int
		v := c.GetSession("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		c.SetSession("count", count)
		c.SaveSession() //nolint
		c.JSON(http.StatusOK, gin.H{"count": count})
	}))

	return r
}
