package admin

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task/core"
)

type RoleMenu struct {
}

type RoleMenuAddRequest struct {
	RoleId  int    `json:"role_id" binding:"required"`
	MenuIds string `json:"menuids" binding:"-"`
}

type RoleMenuViewRequest struct {
	RoleId int `json:"role_id" binding:"required"`
}

func (*RoleMenu) Add(c *core.Context) {
	var r RoleMenuAddRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)
	res, err := service.AddRoleMenu(r.RoleId, r.MenuIds, operator)
	if err != nil {
		c.Fail(400, err.Error(), res)
		return
	}

	if res.Code != 200 {
		c.Fail(400, res.Msg, nil)
		return
	}

	c.Success(nil)
}

/**
提供页面初始化数据
*/
func (*RoleMenu) View(c *core.Context) {
	var r RoleMenuViewRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	res, err := service.ViewRoleMenu(r.RoleId)
	if err != nil {
		c.Fail(400, err.Error(), res)
		return
	}

	if res.Code != 200 {
		c.Fail(400, res.Msg, nil)
		return
	}

	c.Success(gin.H{
		"list": res.Data,
	})
}
