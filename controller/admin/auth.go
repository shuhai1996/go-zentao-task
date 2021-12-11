package admin

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task/core"
	"go-zentao-task/service/admin"
)

var service = admin.InitializeService()

type Auth struct {
}

type AuthLoginRequest struct {
	Account    string `json:"account" binding:"required"`
	Password   string `json:"password" binding:"required"`
	NcToken    string `json:"nc_token" binding:"required"`
	Sig        string `json:"sig" binding:"required"`
	Csessionid string `json:"csessionid" binding:"required"`
	Scene      string `json:"scene" binding:"required"`
}

//登录
func (*Auth) Login(c *core.Context) {
	var r AuthLoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}
	res, err := service.Login(r.Account, r.Password, r.Csessionid, r.NcToken, r.Sig, r.Scene, c.ClientIP())
	if err != nil {
		c.Fail(400, err.Error(), res)
		return
	}

	if res.Code != 200 {
		c.Fail(400, res.Msg, res)
		return
	}

	c.Success(res.Data)
}

func (*Auth) Logout(c *core.Context) {
	//清除用户token数据
	_, err := service.Logout(c.GetHeader("Authorization-Token"))
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	c.Success(nil)
}

func (*Auth) GetAuthInfo(c *core.Context) {
	userInfo := c.GetStringMap("user_info")
	var menuList []map[string]interface{}
	var err error

	if userInfo["role_type"].(string) == "super" && userInfo["account"].(string) == "admin" {
		menuList, err = service.SuperMenus()

	} else {
		menuList, err = service.RoleMenus(userInfo["role_id"].(int))
	}

	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	c.Success(gin.H{
		"menu_list": menuList,
		"user_info": gin.H{
			"user_name":  userInfo["user_name"],
			"role_title": userInfo["role_title"],
		},
	})
}
