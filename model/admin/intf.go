package admin

import (
	"go-zentao-task/pkg/db"
	"errors"
	"github.com/jinzhu/gorm"
	"time"
)

type Intf struct {
	ID              int       `json:"id"`
	ParentMenuId    int       `json:"parent_menu_id"`
	ParentMenuTitle string    `json:"parent_menu_title"`
	MenuId          int       `json:"menu_id"`
	MenuTitle       string    `json:"menu_title"`
	Title           string    `json:"title"`
	Mark            string    `json:"mark"`
	Path            string    `json:"path"`
	Status          int       `json:"status"`
	Operator        string    `json:"operator"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
}

func (Intf) TableName() string {
	return "admin_interface"
}

func NewInf() *Intf {
	return &Intf{}
}

func (*Intf) FindAll(title string, menu_id, pageindex, pagesize int) ([]Intf, error) {
	var result []Intf

	odb := db.Orm.Where("status = ?", 1)

	if title != "" {
		odb = odb.Where("title like ?", "%"+title+"%")
	}

	if menu_id > 0 {
		odb = odb.Where("menu_id = ?", menu_id)
	}

	if err := odb.Order("status").Order("id desc").Offset((pageindex - 1) * pagesize).Limit(pagesize).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return result, nil
}

func (*Intf) FindOne(id int) (*Intf, error) {
	var result Intf
	if err := db.Orm.Where(&Intf{
		ID:     id,
		Status: 1,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*Intf) UpdateInfo(id int, title, mark string, parentMenuId, menuId int, ParentMenuTitle, menuTitle, operator string) (int64, error) {
	op := db.Orm.Model(&Intf{}).Where(&Intf{ID: id}).Updates(map[string]interface{}{
		"parent_menu_id":    parentMenuId,
		"parent_menu_title": ParentMenuTitle,
		"menu_id":           menuId,
		"menu_title":        menuTitle,
		"title":             title,
		"mark":              mark,
		"operator":          operator,
		"update_time":       time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*Intf) UpdateStatus(id, status int, operator string) (int64, error) {
	op := db.Orm.Model(&Intf{}).Where(&Intf{ID: id}).Updates(map[string]interface{}{
		"status":      status,
		"operator":    operator,
		"update_time": time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*Intf) Create(title, path, mark string, parentMenuId, menuId int, parentMenuTitle, menuTitle string, status int, operator string) (int, error) {
	data := &Intf{
		ParentMenuId:    parentMenuId,
		ParentMenuTitle: parentMenuTitle,
		MenuId:          menuId,
		MenuTitle:       menuTitle,
		Title:           title,
		Mark:            mark,
		Path:            path,
		Status:          status,
		Operator:        operator,
		CreateTime:      time.Now(),
	}
	op := db.Orm.Create(data)
	if op.Error != nil {
		return 0, op.Error
	}
	return data.ID, nil
}

func (*Intf) FindOneByPath(path string) (*Intf, error) {
	var result Intf
	if err := db.Orm.Where(&Intf{
		Path:   path,
		Status: 1,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*Intf) FindAllByPaths(paths []string) ([]Intf, error) {
	var result []Intf

	if len(paths) == 0 {
		return result, errors.New("paths is empty")
	}

	odb := db.Orm.Where("status = ?", 1)
	if err := odb.Where("path in (?)", paths).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return result, nil
}
