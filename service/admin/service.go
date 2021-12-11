package admin

import (
	"crypto/md5"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"go-zentao-task/pkg/gredis"
	"go-zentao-task/pkg/util"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

const (
	LoginErrTimesPrefix = "demo-gin:admin:auth:err_times:"
	UserInfoPrefix      = "demo-gin:admin:userinfo:"
)

func (service *Service) Login(account, password, sessionid, nctoken, sig, scene, ip string) (*util.Response, error) {
	conn := gredis.RedisPool.Get()
	defer conn.Close()

	times, err := redis.Int(conn.Do("get", LoginErrTimesPrefix+account))
	if err != nil && err != redis.ErrNil {
		return util.ReturnResponse(400, err.Error(), nil), err
	}

	if times >= 5 {
		return util.ReturnResponse(400, "频繁登录，请5分钟后重试", nil), err
	}

	userinfo, err := service.User.FindOneByAccount(account)
	if err != nil {
		return util.ReturnResponse(400, err.Error(), nil), err
	}

	if userinfo.ID == 0 {
		conn.Do("incr", LoginErrTimesPrefix+account)        //nolint
		conn.Do("expire", LoginErrTimesPrefix+account, 300) //nolint
		return util.ReturnResponse(400, "用户不存在", nil), nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userinfo.Password), []byte(password)); err != nil {
		conn.Do("incr", LoginErrTimesPrefix+account)        //nolint
		conn.Do("expire", LoginErrTimesPrefix+account, 300) //nolint
		return util.ReturnResponse(400, "密码错误", nil), nil
	}

	token := fmt.Sprintf("%x", md5.Sum(uuid.NewV4().Bytes()))
	if _, err := conn.Do("hmset", redis.Args{}.Add(UserInfoPrefix+token).AddFlat(userinfo)...); err != nil {
		return nil, err
	}

	conn.Do("expire", UserInfoPrefix+token, 1800) //nolint 无操作 30分钟后过期

	return util.ReturnResponse(200, "success", map[string]string{
		"access_token": token,
	}), err
}

func (*Service) Logout(token string) (*util.Response, error) {
	conn := gredis.RedisPool.Get()
	defer conn.Close()

	_, err := conn.Do("del", UserInfoPrefix+token) //删除token
	if err != nil && err != redis.ErrNil {
		return util.ReturnResponse(400, err.Error(), nil), err
	}

	return util.ReturnResponse(200, "success", nil), nil
}

func (service *Service) GetUserInfoByToken(token string) (*util.Response, error) {
	conn := gredis.RedisPool.Get()
	defer conn.Close()

	cache, err := redis.Values(conn.Do("HGETALL", UserInfoPrefix+token))
	if err != nil && err != redis.ErrNil {
		return util.ReturnResponse(400, err.Error(), nil), err
	}

	data := service.User

	if err := redis.ScanStruct(cache, data); err != nil {
		return util.ReturnResponse(400, err.Error(), nil), err
	}

	if data.Account == "" {
		return util.ReturnResponse(400, "登录已失效", nil), nil
	}

	conn.Do("expire", UserInfoPrefix+token, 1800) //nolint 鉴权成功后，token重置有效期为30分钟

	return util.ReturnResponse(200, "success", map[string]interface{}{
		"user_id":    data.ID,
		"user_name":  data.UserName,
		"role_id":    data.RoleId,
		"role_title": data.RoleTitle,
		"role_type":  data.RoleType,
		"account":    data.Account,
	}), nil
}

//获取admin菜单
func (service *Service) SuperMenus() ([]map[string]interface{}, error) {
	//获取父级菜单
	parents, err := service.Menu.FindAllByParentId(-1)
	if err != nil {
		return nil, err
	}

	menuList := []map[string]interface{}{}
	for _, v := range parents {
		//获取子菜单
		children := []map[string]interface{}{}

		childMenus, _ := service.Menu.FindAllByParentId(v.ID)
		if len(childMenus) > 0 {
			for _, c := range childMenus {
				children = append(children, map[string]interface{}{
					"title": c.Title,
					"link":  c.Link,
				})
			}
		}

		menuList = append(menuList, map[string]interface{}{
			"title":    v.Title,
			"children": children,
		})
	}
	return menuList, err
}

//获取普通角色菜单
func (service *Service) RoleMenus(roleId int) ([]map[string]interface{}, error) {
	children, err := service.RoleMenu.FindAllByRole(roleId)
	if err != nil {
		return nil, err
	}

	list := []map[string]interface{}{}
	parentMenus := map[int]map[string]interface{}{}
	for _, v := range children {
		//获取菜单信息
		menuInfo, err := service.Menu.FindOne(v.MenuId)
		if err != nil {
			return nil, err
		}

		if menuInfo.ID == 0 {
			continue
		}
		parentTitle := ""
		parentChild := []map[string]interface{}{}
		if _, ok := parentMenus[menuInfo.ParentId]; ok {
			parentTitle = parentMenus[menuInfo.ParentId]["title"].(string)
			parentChild = parentMenus[menuInfo.ParentId]["child"].([]map[string]interface{})
		} else {
			parentMenuInfo, err := service.Menu.FindOne(menuInfo.ParentId)
			if err != nil {
				return nil, err
			}
			parentTitle = parentMenuInfo.Title

		}

		parentChild = append(parentChild, map[string]interface{}{
			"title": menuInfo.Title,
			"link":  menuInfo.Link,
		})

		parentMenus[menuInfo.ParentId] = map[string]interface{}{
			"title": parentTitle,
			"child": parentChild,
		}
	}

	for _, v := range parentMenus {
		list = append(list, v)
	}

	return list, nil
}

func (service *Service) AddRoleMenu(roleId int, menuidsRaw, operator string) (*util.Response, error) {
	//获取角色信息
	roleInfo, err := service.Role.FindOne(roleId)
	if err != nil {
		return util.ReturnResponse(400, "角色信息获取失败", nil), err
	}

	if roleInfo.ID == 0 {
		return util.ReturnResponse(400, "角色不存在", nil), err
	}

	//获取已勾选菜单
	res, err := service.RoleMenu.FindAllByRole(roleId)
	if err != nil {
		return util.ReturnResponse(400, "获取已勾选菜单失败", nil), err

	}

	hadList := map[int]int{}
	for _, v := range res {
		hadList[v.MenuId] = v.ID
	}

	//查找需要新增菜单
	menuIds := strings.Split(menuidsRaw, ",") //menuids为空时，做删除操作
	for _, menuId := range menuIds {
		id, _ := strconv.Atoi(menuId)
		if menuId == "" || id <= 0 {
			continue
		}

		if _, ok := hadList[id]; ok {
			delete(hadList, id)
		} else {
			service.RoleMenu.Create(roleId, roleInfo.Type, id, 1, operator) //nolint
		}
	}

	updateIds := []int{}
	for _, v := range hadList {
		updateIds = append(updateIds, v)
	}

	//批量下架
	if len(updateIds) > 0 {
		_, err := service.RoleMenu.UpdateStatusByIds(updateIds, 2, operator)
		if err != nil {
			return util.ReturnResponse(400, "批量下架授权菜单失败", nil), err
		}
	}

	return util.ReturnResponse(200, "success", nil), nil
}

func (service *Service) ViewRoleMenu(roleId int) (*util.Response, error) {
	//获取已勾选菜单
	res, err := service.RoleMenu.FindAllByRole(roleId)
	if err != nil {
		return util.ReturnResponse(400, "获取已勾选菜单失败", nil), err

	}

	hadList := map[int]int{}
	for _, v := range res {
		hadList[v.MenuId] = v.ID
	}

	//获取父级菜单
	parents, err := service.Menu.FindAllByParentId(-1)
	if err != nil {
		return util.ReturnResponse(400, "获取父级菜单列表失败", nil), err
	}

	menuList := []map[string]interface{}{}
	for _, v := range parents {
		//获取子菜单
		children := []map[string]interface{}{}

		childMenus, _ := service.Menu.FindAllByParentId(v.ID)
		if len(childMenus) > 0 {
			for _, c := range childMenus {
				check := "N"

				//查看该菜单是否已被勾选
				if _, ok := hadList[c.ID]; ok {
					check = "Y"
				}

				children = append(children, map[string]interface{}{
					"id":    c.ID,
					"title": c.Title,
					"check": check,
				})
			}
		}

		if len(children) == 0 {
			continue
		}

		menuList = append(menuList, map[string]interface{}{
			"id":       v.ID,
			"title":    v.Title,
			"children": children,
		})
	}

	return util.ReturnResponse(200, "success", menuList), nil
}
