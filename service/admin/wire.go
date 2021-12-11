//+build wireinject

package admin

import (
	"go-zentao-task/model/admin"
	"github.com/google/wire"
)

type Service struct {
	Intf     *admin.Intf
	User     *admin.User
	Role     *admin.Role
	Menu     *admin.Menu
	RoleMenu *admin.RoleMenu
}

func NewService(
	intf *admin.Intf,
	user *admin.User,
	role *admin.Role,
	menu *admin.Menu,
	roleMenu *admin.RoleMenu) *Service {
	return &Service{
		Intf:     intf,
		User:     user,
		Role:     role,
		Menu:     menu,
		RoleMenu: roleMenu,
	}
}

func InitializeService() *Service {
	wire.Build(
		NewService,
		admin.NewInf,
		admin.NewUser,
		admin.NewRole,
		admin.NewMenu,
		admin.NewRoleMenu,
	)
	return &Service{}
}
