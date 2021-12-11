package admin

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task/core"
)

type Role struct {
}

type RoleIndexRequest struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

type RoleViewRequest struct {
	ID int `json:"id"`
}

type RoleCreateRequest struct {
	Title string `json:"title" binding:"required"`
	Mark  string `json:"mark" binding:"required"`
}

type RoleUpdateStatusRequest struct {
	ID     int    `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type RoleUpdateRequest struct {
	ID    int    `json:"id" binding:"required"`
	Title string `json:"title" binding:"required"`
	Mark  string `json:"mark" binding:"required"`
}

func (*Role) Index(c *core.Context) {
	var r RoleIndexRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	res, err := service.Role.FindAll(r.Title, r.Type, r.Page, r.PageSize)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}
	list := []map[string]interface{}{}
	for _, v := range res {
		list = append(list, map[string]interface{}{
			"id":          v.ID,
			"title":       v.Title,
			"type":        v.Type,
			"mark":        v.Mark,
			"operator":    v.Operator,
			"create_time": v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	c.Success(gin.H{
		"list": list,
	})
}

func (*Role) View(c *core.Context) {
	var r RoleViewRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	typeList := []string{"role"}
	if r.ID == 0 {
		c.Success(gin.H{
			"role_type_list": typeList,
		})
		return
	}

	info, err := service.Role.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "角色不存在", nil)
		return
	}

	c.Success(gin.H{
		"role_type_list": typeList,
		"info": gin.H{
			"id":    info.ID,
			"title": info.Title,
			"type":  info.Type,
			"mark":  info.Mark,
		},
	})
}

func (*Role) Create(c *core.Context) {
	var r RoleCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)
	roleType := "role"

	res, err := service.Role.Create(roleType, r.Title, r.Mark, 1, operator)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if res <= 0 {
		c.Fail(400, "新增失败", nil)
		return
	}

	c.Success(nil)
}

func (*Role) UpdateStatus(c *core.Context) {
	var r RoleUpdateStatusRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	if r.Status != "off" && r.Status != "on" {
		c.Fail(40100, "操作异常", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)

	status := 2
	if r.Status == "on" {
		status = 1
	}

	info, err := service.Role.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "数据不存在", nil)
		return
	}

	if info.Type == "super" {
		c.Fail(400, "超管不支持修改状态", nil)
		return
	}

	res, err := service.Role.UpdateStatus(r.ID, status, operator)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if res <= 0 {
		c.Fail(400, "修改失败", nil)
		return
	}

	c.Success(nil)
}

func (*Role) Update(c *core.Context) {
	var r RoleUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)

	info, err := service.Role.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "数据不存在", nil)
		return
	}

	res, err := service.Role.UpdateInfo(r.ID, r.Title, r.Mark, operator)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if res <= 0 {
		c.Fail(400, "修改失败", nil)
		return
	}

	c.Success(nil)
}
