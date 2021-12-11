package admin

import (
	"go-zentao-task/core"
	"go-zentao-task/pkg/rbac"
	"github.com/gin-gonic/gin"
)

type Menu struct {
}

type MenuIndexRequest struct {
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Title    string `json:"title"`
}

type MenuViewRequest struct {
	ID int `json:"id"`
}

type MenuCreateRequest struct {
	ParentId int    `json:"parent_id"`
	Sort     int    `json:"sort" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Link     string `json:"link"`
}

type MenuUpdateStatusRequest struct {
	ID     int    `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type MenuUpdateRequest struct {
	ID       int    `json:"id" binding:"required"`
	ParentId int    `json:"parent_id" binding:"required"`
	Sort     int    `json:"sort" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Link     string `json:"link"`
}

func (*Menu) Index(c *core.Context) {
	var r MenuIndexRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	res, err := service.Menu.FindAll(r.Title, r.Page, r.PageSize)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	list := []map[string]interface{}{}
	for _, v := range res {
		parent_title := ""
		if v.ParentId > 0 { //获取父级菜单
			parent_info, err := service.Menu.FindOne(v.ParentId)
			if err == nil && parent_info.ID > 0 {
				parent_title = parent_info.Title
			}
		}

		list = append(list, map[string]interface{}{
			"id":           v.ID,
			"title":        v.Title,
			"parent_title": parent_title,
			"sort":         v.Sort,
			"link":         v.Link,
			"operator":     v.Operator,
			"create_time":  v.CreateTime.Format("2006-01-02 15:04:05"),
		})
	}

	c.Success(gin.H{
		"list": list,
	})
}

func (*Menu) View(c *core.Context) {
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

	parentList := []map[string]interface{}{}
	for _, v := range parents {
		parentList = append(parentList, map[string]interface{}{
			"id":    v.ID,
			"title": v.Title,
		})
	}

	if r.ID == 0 {
		c.Success(gin.H{
			"parent_list": parentList,
		})
		return
	}

	info, err := service.Menu.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "菜单不存在", nil)
		return
	}

	c.Success(gin.H{
		"parent_list": parentList,
		"info": gin.H{
			"id":        info.ID,
			"title":     info.Title,
			"parent_id": info.ParentId,
			"sort":      info.Sort,
			"link":      info.Link,
		},
	})
}

func (*Menu) Create(c *core.Context) {
	var r MenuCreateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)

	info, err := service.Menu.FindOneByTitle(r.Title)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID > 0 {
		c.Fail(400, "菜单名重复，请调整后重新提交", nil)
		return
	}

	if r.ParentId == 0 {
		r.ParentId = -1
	}

	if r.ParentId > 0 && r.Link == "" {
		c.Fail(400, "二级菜单链接不能为空", nil)
		return
	}

	res, err := service.Menu.Create(r.ParentId, r.Title, r.Link, r.Sort, 1, operator)
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

func (*Menu) UpdateStatus(c *core.Context) {
	var r MenuUpdateStatusRequest
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

	info, err := service.Menu.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "数据不存在", nil)
		return
	}

	//下架前检查该菜单是否已被授权
	if r.Status != "on" {
		//检查菜单授权、检查菜单下的接口授权
		if info.ParentId == -1 {
			//检查是否有子菜单
			child, err := service.Menu.FindAllByParentId(r.ID)
			if err != nil {
				c.Fail(400, err.Error(), nil)
				return
			}
			if len(child) > 0 {
				c.Fail(400, "请先删除子菜单", nil)
				return
			}
		} else {
			//检查已分配菜单权限
			roleMenu, err := service.RoleMenu.FindAllByMenu(r.ID)
			if err != nil {
				c.Fail(400, err.Error(), nil)
				return
			}

			if len(roleMenu) > 0 {
				c.Fail(400, "菜单已被分配权限，请撤销相关权限后再删除", nil)
				return
			}

			//检查已分配接口权限
			allInterface := rbac.Enforcer.GetAllNamedObjects("p")
			if len(allInterface) > 0 {
				records, err := service.Intf.FindAllByPaths(allInterface)
				if err != nil {
					c.Fail(400, err.Error(), nil)
					return
				}

				for _, v := range records {
					if v.MenuId == r.ID {
						c.Fail(400, "菜单下的接口已被分配权限，请撤销相关权限后再删除", nil)
						return
					}
				}
			}

		}
	}

	res, err := service.Menu.UpdateStatus(r.ID, status, operator)
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

func (*Menu) Update(c *core.Context) {
	var r MenuUpdateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	operator := c.GetStringMap("user_info")["account"].(string)

	info, err := service.Menu.FindOne(r.ID)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if info.ID == 0 {
		c.Fail(400, "数据不存在", nil)
		return
	}

	check, err := service.Menu.FindOneByTitle(r.Title)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	if check.ID > 0 && check.ID != r.ID {
		c.Fail(400, "菜单名重复，请调整后重新提交", nil)
		return
	}

	res, err := service.Menu.UpdateInfo(r.ID, r.ParentId, r.Sort, r.Title, r.Link, operator)
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
