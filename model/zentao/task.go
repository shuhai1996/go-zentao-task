package zentao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
)

type Task struct {
	ID       int    `json:"id"`
	Project  string `json:"project"`
	Parent   string `json:"parent"`
	Module   string `json:"module"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Estimate string `json:"estimate"`
	Left     string `json:"left"`
	Deleted  int    `json:"deleted"`
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
		ID: id,
	}).First(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return &result, nil
}

func (*Task) FindAll(assignedTo, status string) ([]Task, error) {
	var result []Task
	odb := db.Orm.Where("assignedTo = ?", assignedTo)

	if status != "" {
		odb = odb.Where("status = ?", "status")
	} else {
		odb = odb.Where("status = 'doing' or status = 'wait'")
	}

	if err := odb.Order("status").Order("id desc").Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}
