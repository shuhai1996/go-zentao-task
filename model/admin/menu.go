package admin

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type Menu struct {
	ID         int       `json:"id"`
	ParentId   int       `json:"parent_id"`
	Sort       int       `json:"sort"`
	Title      string    `json:"title"`
	Link       string    `json:"link"`
	Status     int       `json:"status"`
	Operator   string    `json:"operator"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}

func (Menu) TableName() string {
	return "admin_menu"
}

func NewMenu() *Menu {
	return &Menu{}
}

func (*Menu) FindAll(title string, pageindex, pagesize int) ([]Menu, error) {
	var result []Menu
	odb := db.Orm.Where("status = ?", 1)

	if title != "" {
		odb = odb.Where("title like ?", "%"+title+"%")
	}

	if err := odb.Order("status").Order("id desc").Offset((pageindex - 1) * pagesize).Limit(pagesize).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func (*Menu) FindAllByParentId(parentId int) ([]Menu, error) {
	var result []Menu
	if err := db.Orm.Where(&Menu{
		ParentId: parentId,
		Status:   1,
	}).Order("status").Order("sort desc").Order("id desc").Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func (*Menu) FindOne(id int) (*Menu, error) {
	var result Menu
	if err := db.Orm.Where(&Menu{
		ID:     id,
		Status: 1,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*Menu) FindOneByTitle(title string) (*Menu, error) {
	var result Menu
	if err := db.Orm.Where(&Menu{
		Title:  title,
		Status: 1,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*Menu) UpdateInfo(id, parentId, sort int, title, link, operator string) (int64, error) {
	op := db.Orm.Model(&Menu{}).Where(&Menu{ID: id}).Updates(map[string]interface{}{
		"parent_id":   parentId,
		"title":       title,
		"sort":        sort,
		"link":        link,
		"operator":    operator,
		"update_time": time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*Menu) UpdateStatus(id, status int, operator string) (int64, error) {
	op := db.Orm.Model(&Menu{}).Where(&Menu{ID: id}).Updates(map[string]interface{}{
		"status":      status,
		"operator":    operator,
		"update_time": time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*Menu) Create(parentId int, title, link string, sort, status int, operator string) (int, error) {
	data := &Menu{
		ParentId:   parentId,
		Title:      title,
		Sort:       sort,
		Link:       link,
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
