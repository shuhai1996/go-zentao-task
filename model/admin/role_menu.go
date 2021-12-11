package admin

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type RoleMenu struct {
	ID         int       `json:"id"`
	RoleId     int       `json:"role_id"`
	RoleType   string    `json:"role_type"`
	MenuId     int       `json:"menu_id"`
	Status     int       `json:"status"`
	Operator   string    `json:"operator"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (RoleMenu) TableName() string {
	return "admin_role_menu"
}

func NewRoleMenu() *RoleMenu {
	return &RoleMenu{}
}

func (*RoleMenu) FindAllByRole(role_id int) ([]RoleMenu, error) {
	var result []RoleMenu
	if err := db.Orm.Where(&RoleMenu{
		RoleId: role_id,
		Status: 1,
	}).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func (*RoleMenu) FindOne(id int) (*RoleMenu, error) {
	var result RoleMenu
	if err := db.Orm.Where(&RoleMenu{
		ID:     id,
		Status: 1,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*RoleMenu) UpdateStatusByIds(ids []int, status int, operator string) (int64, error) {
	op := db.Orm.Model(&RoleMenu{}).Where("id IN (?)", ids).Updates(map[string]interface{}{
		"status":      status,
		"operator":    operator,
		"update_time": time.Now(),
	})

	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*RoleMenu) UpdateStatus(id, status int, operator string) (int64, error) {
	op := db.Orm.Model(&RoleMenu{}).Where(&RoleMenu{ID: id}).Updates(map[string]interface{}{
		"status":      status,
		"operator":    operator,
		"update_time": time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*RoleMenu) Create(roleId int, roleType string, menuId, status int, operator string) (int, error) {
	data := &RoleMenu{
		RoleId:     roleId,
		RoleType:   roleType,
		MenuId:     menuId,
		Status:     status,
		Operator:   operator,
		CreateTime: time.Now(),
	}
	op := db.Orm.Create(data)
	if op.Error != nil {
		return 0, op.Error
	}
	return data.ID, nil
}

func (*RoleMenu) FindAllByMenu(menuId int) ([]RoleMenu, error) {
	var result []RoleMenu
	if err := db.Orm.Where(&RoleMenu{
		MenuId: menuId,
		Status: 1,
	}).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}
