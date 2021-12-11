package zendao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
)

type Task struct {
	ID         int      `json:"id"`
	Project   string    `json:"project"`
	Parent    string    `json:"parent"`
	Module    string    `json:"module"`
	Status    string    `json:"status"`
	Type   	  string    `json:"type"`
	Name      string    `json:"name"`
	Deleted    int    	`json:"deleted"`
}

func (Task) TableName() string {
	return "zt_task"
}


func NewTask() *Task {
	return &Task{}
}

func (*Task) FindOne(id int) (*Task, error) {
	var result Task
	if err := db.Orm.Where(&Task{
		ID:     id,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}