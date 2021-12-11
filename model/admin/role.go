package admin

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type Role struct {
	ID         int       `json:"id"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Mark       string    `json:"mark"`
	Status     int       `json:"status"`
	Operator   string    `json:"operator"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (Role) TableName() string {
	return "admin_role"
}

func NewRole() *Role {
	return &Role{}
}

func (*Role) FindAll(title, otype string, pageindex, pagesize int) ([]Role, error) {
	var result []Role
	odb := db.Orm.Where("status = ?", 1)

	if title != "" {
		odb = odb.Where("title like ?", "%"+title+"%")
	}
	if otype != "" {
		odb = odb.Where("type = ?", otype)
	}

	if err := odb.Order("status").Order("id desc").Offset((pageindex - 1) * pagesize).Limit(pagesize).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func (*Role) FindOne(id int) (*Role, error) {
	var result Role
	if err := db.Orm.Where(&Role{
		ID:     id,
		Status: 1,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*Role) UpdateInfo(id int, title, mark, operator string) (int64, error) {
	op := db.Orm.Model(&Role{}).Where(&Role{ID: id}).Updates(map[string]interface{}{
		"title":       title,
		"mark":        mark,
		"operator":    operator,
		"update_time": time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*Role) UpdateStatus(id, status int, operator string) (int64, error) {
	op := db.Orm.Model(&Role{}).Where(&Role{ID: id}).Updates(map[string]interface{}{
		"status":      status,
		"operator":    operator,
		"update_time": time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*Role) Create(roleType, title, mark string, status int, operator string) (int, error) {
	data := &Role{
		Type:       roleType,
		Title:      title,
		Mark:       mark,
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
