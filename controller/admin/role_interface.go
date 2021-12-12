package admin

import (
	"github.com/gin-gonic/gin"
	"go-zentao-task/core"
	"go-zentao-task/pkg/rbac"
	"strconv"
	"strings"
)

type RoleInterface struct {
}

type RoleInterfaceAddRequest struct {
	RoleId int    `json:"role_id" binding:"required"`
	Paths  string `json:"paths" binding:"-"`
}

type RoleInterfaceViewRequest struct {
	RoleId int `json:"role_id" binding:"required"`
}

func (*RoleInterface) Add(c *core.Context) {
	var r RoleInterfaceAddRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	hadPoliceList := rbac.Enforcer.GetFilteredNamedPolicy("p", 0, "role_id_"+strconv.Itoa(r.RoleId))
	for _, values := range hadPoliceList {
		rbac.Enforcer.RemoveNamedPolicy("p", "role_id_"+strconv.Itoa(r.RoleId), values[1], values[2]) //nolint
	}

	paths := strings.Split(r.Paths, ",") //paths为空时，做删除操作
	for _, path := range paths {
		if path == "" {
			continue
		}
		_, err := rbac.Enforcer.AddNamedPolicy("p", "role_id_"+strconv.Itoa(r.RoleId), path, "POST")
		if err != nil {
			c.Fail(400, err.Error(), nil)
			return
		}
	}

	// Save the policy back to DB.
	if err := rbac.Enforcer.SavePolicy(); err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	c.Success(nil)
}

/**
提供页面初始化数据
*/
func (*RoleInterface) View(c *core.Context) {
	var r RoleInterfaceViewRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.Fail(40100, "缺少参数，请重试", nil)
		return
	}

	//获取已勾选接口列表
	checkList := map[string]string{}
	hadPoliceList := rbac.Enforcer.GetFilteredNamedPolicy("p", 0, "role_id_"+strconv.Itoa(r.RoleId))
	for _, values := range hadPoliceList {
		checkList[values[1]] = values[2]
	}

	//获取所有interface
	all, err := service.Intf.FindAll("", 0, 1, 100)
	if err != nil {
		c.Fail(400, err.Error(), nil)
		return
	}

	list := []interface{}{}
	parent_key := map[int]map[string]interface{}{}
	child_key := map[int]map[string]interface{}{}
	path_key := map[int][]interface{}{}

	for _, v := range all {
		check := "N"

		//查看该path是否已被勾选
		if _, ok := checkList[v.Path]; ok {
			check = "Y"
		}

		path_key[v.MenuId] = append(path_key[v.MenuId], map[string]string{
			"path":  v.Path,
			"title": v.Title,
			"check": check,
		})

		child_key[v.MenuId] = map[string]interface{}{
			"title": v.MenuTitle,
			"paths": path_key[v.MenuId],
		}

		parent_key[v.ParentMenuId] = map[string]interface{}{
			"title": v.ParentMenuTitle,
			"child": child_key,
		}
	}

	for _, v := range parent_key {

		child := []map[string]interface{}{}
		temp := v["child"].(map[int]map[string]interface{})

		for _, c := range temp {
			child = append(child, c)
		}

		list = append(list, map[string]interface{}{
			"title":    v["title"].(string),
			"children": child,
		})
	}

	c.Success(gin.H{
		"list": list,
	})
}
