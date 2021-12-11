package admin

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task/core"
	"go-zentao-task/pkg/rbac"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type User struct {
}

type UserIndexRequest struct {
	Page      int    `json:"page" binding:"required"`
	PageSize  int    `json:"page_size" binding:"required"`
	UserName  string `json:"user_name"`
	RoleTitle string `json:"role_title"`
}

type UserViewRequest struct {
	ID int `json:"id"`
}

//TODO 缺少长度校验
type UserCreateRequest struct {
	RoleId   int    `json:"role_id" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	ID       int    `json:"id" binding:"required"`
	RoleId   int    `json:"role_id" binding:"required"`
	UserName string `json:"user_name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateStatusRequest struct {
	ID     int    `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

func (*User) Index(c *core.Context) {
	var r UserIndexRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	res, err := service.User.FindAll(r.UserName, r.RoleTitle, r.Page, r.PageSize)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}
	list := []map[string]interface{}{}
	for _, v := range res {
		list = append(list, map[string]interface{}{
			"id":          v.ID,
			"user_name":   v.UserName,
			"account":     v.Account,
			"role_type":   v.RoleType,
			"role_title":  v.RoleTitle,
			"operator":    v.Operator,
			"create_time": v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	c.Success(gin.H{
		"list": list,
	})
}

func (*User) View(c *core.Context) {
	var r UserViewRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	//获取角色列表 默认取100条
	roles, err := service.Role.FindAll("", "role", 1, 100)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	rolelist := []map[string]interface{}{}
	for _, v := range roles {
		rolelist = append(rolelist, map[string]interface{}{
			"id":    v.ID,
			"title": v.Title,
			"type":  v.Type,
		})
	}

	if r.ID == 0 {
		c.Success(gin.H{
			"role_list": rolelist,
		})
		return
	}

	info, err := service.User.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "用户不存在", nil)
		return
	}

	c.Success(gin.H{
		"role_list": rolelist,
		"info": gin.H{
			"id":         info.ID,
			"user_name":  info.UserName,
			"account":    info.Account,
			"role_id":    info.RoleId,
			"role_type":  info.RoleType,
			"role_title": info.RoleTitle,
		},
	})
}

func (*User) Create(c *core.Context) {
	var r UserCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)

	//超管只能有一个
	//获取角色名称
	roleInfo, err := service.Role.FindOne(r.RoleId)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if roleInfo.ID == 0 {
		c.Fail(400, "角色不存在", nil)
		return
	}

	if roleInfo.Type == "super" {
		c.Fail(400, "超管不支持添加", nil)
		return
	}

	//检查用户名是否已存在
	check, err := service.User.FindOneByAccount(r.Account)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if check.ID > 0 {
		c.Fail(400, "该账号已被占用", nil)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	id, err := service.User.Create(r.RoleId, roleInfo.Type, roleInfo.Title, r.UserName, r.Account, string(password), 1, operator)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if id <= 0 {
		c.Fail(400, "添加失败", nil)
		return
	}

	rbac.Enforcer.AddNamedGroupingPolicy("g", r.Account, "role_id_"+strconv.Itoa(r.RoleId)) //nolint
	rbac.Enforcer.SavePolicy()                                                              //nolint
	c.Success(nil)
}

func (*User) Update(c *core.Context) {
	var r UserUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)
	info, err := service.User.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "改记录不存在", nil)
		return
	}

	//获取角色名称
	roleInfo, err := service.Role.FindOne(r.RoleId)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if roleInfo.ID == 0 {
		c.Fail(400, "角色不存在", nil)
		return
	}

	if roleInfo.Type == "super" && info.RoleType != "super" {
		c.Fail(400, "不支持编辑角色为超管", nil)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	res, err := service.User.UpdateInfo(r.ID, r.RoleId, roleInfo.Type, roleInfo.Title, r.UserName, string(password), operator)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if res <= 0 {
		c.Fail(400, "编辑失败", nil)
		return
	}

	c.Success(nil)
}

func (*User) UpdateStatus(c *core.Context) {
	var r UserUpdateStatusRequest
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

	info, err := service.User.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "改记录不存在", nil)
		return
	}

	if info.RoleType == "super" {
		c.Fail(400, "超管不支持修改状态", nil)
		return
	}

	res, err := service.User.UpdateStatus(r.ID, status, operator)
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
