package zentao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
)

type Task struct {
	ID       int    `json:"id"`
	Project  int `json:"project"`
	Parent   int `json:"parent"`
	Module   int `json:"module"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Estimate string `json:"estimate"`
	Left     float32 `json:"left"`
	Execution int `json:"execution"`
	AssignedTo string `json:"assignedTo"`
	FromBug int `json:"fromBug"`
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

	if err := odb.Order("status desc").Order("id desc").Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}
