package admin

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type User struct {
	ID         int       `json:"id" redis:"id"`
	RoleId     int       `json:"role_id" redis:"role_id"`
	RoleType   string    `json:"role_type" redis:"role_type"`
	RoleTitle  string    `json:"role_title" redis:"role_title"`
	UserName   string    `json:"user_name" redis:"user_name"`
	Account    string    `json:"account" redis:"account"`
	Password   string    `json:"password" redis:"-"`
	Status     int       `json:"status" redis:"-"`
	Operator   string    `json:"operator" redis:"-"`
	CreateTime time.Time `json:"create_time" redis:"-"`
	UpdateTime time.Time `json:"update_time" redis:"-"`
}

func (User) TableName() string {
	return "admin_user"
}

func NewUser() *User {
	return &User{}
}

func (*User) FindAll(userName, roleTitle string, pageindex, pagesize int) ([]User, error) {
	var result []User
	odb := db.Orm.Where("status = ?", 1)

	if userName != "" {
		odb = odb.Where("user_name like ?", "%"+userName+"%")
	}

	if roleTitle != "" {
		odb = odb.Where("role_type like ?", "%"+roleTitle+"%")
	}

	if err := odb.Order("status").Order("id desc").Offset((pageindex - 1) * pagesize).Limit(pagesize).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func (*User) FindOne(id int) (*User, error) {
	var result User
	if err := db.Orm.Where(&User{
		ID:     id,
		Status: 1,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*User) FindOneByAccount(Acount string) (*User, error) {
	var result User
	if err := db.Orm.Where(&User{
		Account: Acount,
		Status:  1,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*User) UpdateInfo(id, roleId int, roleType, roleTitle, username, password, operator string) (int64, error) {
	op := db.Orm.Model(&User{}).Where(&User{ID: id}).Updates(map[string]interface{}{
		"role_id":     roleId,
		"role_type":   roleType,
		"role_title":  roleTitle,
		"user_name":   username,
		"password":    password,
		"operator":    operator,
		"update_time": time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*User) UpdateStatus(id, status int, operator string) (int64, error) {
	op := db.Orm.Model(&User{}).Where(&User{ID: id}).Updates(map[string]interface{}{
		"status":      status,
		"operator":    operator,
		"update_time": time.Now(),
	})
	if op.Error != nil {
		return 0, op.Error
	}
	return op.RowsAffected, nil
}

func (*User) Create(roleId int, roleType, roleTitle, userName, account, password string, status int, operator string) (int, error) {
	data := &User{
		RoleId:     roleId,
		RoleType:   roleType,
		RoleTitle:  roleTitle,
		UserName:   userName,
		Account:    account,
		Password:   password,
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
