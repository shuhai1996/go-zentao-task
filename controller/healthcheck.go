package controller

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-zentao-task/service"
	"net/http"
)

type HealthCheckRequest struct {
	Timestamp string `form:"timestamp" binding:"required"`
	Sign      string `form:"sign" binding:"required"`
}

func HealthCheck(c *gin.Context) {
	var r HealthCheckRequest
	if err := c.ShouldBindQuery(&r); err != nil {
		c.String(http.StatusOK, "sign error")
		return
	}

	key := "dj7lp4xkbd6udeuo67fopno75n2cmurf"
	if r.Sign != fmt.Sprintf("%x", md5.Sum([]byte(r.Timestamp+key))) {
		c.String(http.StatusOK, "sign error")
		return
	}

	if err := service.PingMysql(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err := service.PingRedis(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	if err := service.PingLogRedis(); err != nil {
		c.String(http.StatusOK, err.Error())
		return
	}
	c.String(http.StatusOK, "success")
}
