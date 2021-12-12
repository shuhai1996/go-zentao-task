package admin

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task/core"
	"go-zentao-task/pkg/rbac"
	"go-zentao-task/pkg/util"
	"strings"
)

type Interface struct {
}

type InterfaceIndexRequest struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Title    string `json:"title" binding:"-"`
	MenuId   int    `json:"menu_id" binding:"-"`
}

type InterfaceViewRequest struct {
	ID int `json:"id"`
}

type InterfaceCreateRequest struct {
	ParentMenuId int    `json:"parent_menu_id" binding:"required"`
	MenuId       int    `json:"menu_id"`
	Title        string `json:"title" binding:"required"`
	Mark         string `json:"mark" binding:"required"`
	Path         string `json:"path" binding:"-"`
}

type InterfaceUpdateStatusRequest struct {
	ID     int    `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type InterfaceUpdateRequest struct {
	ID           int    `json:"id" binding:"required"`
	ParentMenuId int    `json:"parent_menu_id" binding:"required"`
	MenuId       int    `json:"menu_id" binding:"required"`
	Title        string `json:"title" binding:"required"`
	Mark         string `json:"mark" binding:"required"`
}

func (*Interface) Index(c *core.Context) {
	var r InterfaceIndexRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	res, err := service.Intf.FindAll(r.Title, r.MenuId, r.Page, r.PageSize)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	list := []map[string]interface{}{}
	for _, v := range res {
		list = append(list, map[string]interface{}{
			"id":          v.ID,
			"title":       v.Title,
			"menu_title":  v.MenuTitle,
			"mark":        v.Mark,
			"path":        v.Path,
			"operator":    v.Operator,
			"create_time": v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	//获取父级菜单
	parents, err := service.Menu.FindAllByParentId(-1)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	menuList := []map[string]interface{}{}
	for _, v := range parents {

		//获取子菜单
		children := []map[string]interface{}{}

		childMenus, _ := service.Menu.FindAllByParentId(v.ID)
		if len(childMenus) > 0 {
			for _, c := range childMenus {
				children = append(children, map[string]interface{}{
					"id":    c.ID,
					"title": c.Title,
				})
			}
		}

		menuList = append(menuList, map[string]interface{}{
			"id":       v.ID,
			"title":    v.Title,
			"children": children,
		})
	}

	c.Success(gin.H{
		"list":     list,
		"menuList": menuList,
	})
}

func (*Interface) View(c *core.Context) {
	var r MenuViewRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	//获取父级菜单
	parents, err := service.Menu.FindAllByParentId(-1)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	menuList := []map[string]interface{}{}
	for _, v := range parents {

		//获取子菜单
		children := []map[string]interface{}{}

		childMenus, _ := service.Menu.FindAllByParentId(v.ID)
		if len(childMenus) > 0 {
			for _, c := range childMenus {
				children = append(children, map[string]interface{}{
					"id":    c.ID,
					"title": c.Title,
				})
			}
		}

		if len(children) > 0 {
			menuList = append(menuList, map[string]interface{}{
				"id":       v.ID,
				"title":    v.Title,
				"children": children,
			})
		}
	}

	if r.ID == 0 {
		c.Success(gin.H{
			"menu_list": menuList,
		})
		return
	}

	info, err := service.Intf.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "接口不存在", nil)
		return
	}

	c.Success(gin.H{
		"menu_list": menuList,
		"info": gin.H{
			"id":                info.ID,
			"title":             info.Title,
			"parent_menu_id":    info.ParentMenuId,
			"parent_menu_title": info.ParentMenuTitle,
			"menu_id":           info.MenuId,
			"menu_title":        info.MenuTitle,
			"mark":              info.Mark,
			"path":              info.Path,
		},
	})
}

func (*Interface) Create(c *core.Context) {
	var r InterfaceCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)
	r.Path = strings.ToLower(r.Path)

	info, err := service.Intf.FindOneByPath(r.Path)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID > 0 {
		c.Fail(400, "接口地址重复，请调整后重新提交", nil)
		return
	}
	//获取父级菜单名称
	parentMenuInfo, err := service.Menu.FindOne(r.ParentMenuId)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if parentMenuInfo.ID == 0 {
		c.Fail(400, "父级菜单获取失败", nil)
		return
	}

	//获取二级菜单名称
	menuInfo, err := service.Menu.FindOne(r.MenuId)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if menuInfo.ID == 0 {
		c.Fail(400, "二级菜单获取失败", nil)
		return
	}

	res, err := service.Intf.Create(r.Title, r.Path, r.Mark, r.ParentMenuId, r.MenuId, parentMenuInfo.Title, menuInfo.Title, 1, operator)
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

func (*Interface) UpdateStatus(c *core.Context) {
	var r InterfaceUpdateStatusRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	if r.Status != "off" && r.Status != "on" {
		c.Fail(400, "操作异常", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)
	status := 2
	if r.Status == "on" {
		status = 1
	}

	info, err := service.Intf.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "数据不存在", nil)
		return
	}

	//下架前检查该接口是否已被授权
	if r.Status != "on" {
		allInterface := rbac.Enforcer.GetAllNamedObjects("p") //获取已分配接口权限
		if in, _ := util.Contain(info.Path, allInterface); in {
			c.Fail(400, "接口已被分配权限，请撤销相关权限后再删除", nil)
			return
		}
	}

	res, err := service.Intf.UpdateStatus(r.ID, status, operator)
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

func (*Interface) Update(c *core.Context) {
	var r InterfaceUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)

	info, err := service.Intf.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "数据不存在", nil)
		return
	}

	//获取父级菜单名称
	parentMenuInfo, err := service.Menu.FindOne(r.ParentMenuId)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if parentMenuInfo.ID == 0 {
		c.Fail(400, "父级菜单获取失败", nil)
		return
	}

	//获取二级菜单名称
	menuInfo, err := service.Menu.FindOne(r.MenuId)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if menuInfo.ID == 0 {
		c.Fail(400, "二级菜单获取失败", nil)
		return
	}

	res, err := service.Intf.UpdateInfo(r.ID, r.Title, r.Mark, r.ParentMenuId, r.MenuId, parentMenuInfo.Title, menuInfo.Title, operator)
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
