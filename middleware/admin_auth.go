package middleware

import (
	"go-zentao-task/core"
	"go-zentao-task/service/admin"
	"go-zentao-task/service/errcode"
)

var service = admin.InitializeService()

type AuthHeader struct {
	Token string `header:"Authorization-Token" binding:"required"`
}

// 后台鉴权中间件
func AdminAuth(c *core.Context) {
	var h AuthHeader
	if err := c.ShouldBindHeader(&h); err != nil {
		c.FailWithErrCode(errcode.ErrAdminLoginExpired, nil)
		return
	}

	res, err := service.GetUserInfoByToken(h.Token)
	if err != nil || res.Code != 200 {
		c.FailWithErrCode(errcode.ErrAdminLoginExpired, nil)
		return
	}

	userInfo := res.Data
	c.Set("user_info", userInfo)
	c.Next()
}
