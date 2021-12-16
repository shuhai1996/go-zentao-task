package zentao

import (
	"github.com/jinzhu/gorm"
	"go-zentao-task/pkg/db"
	"time"
)

type TaskEstimate struct {
	ID       int     `json:"id"`
	Task     string  `json:"task"`
	Date     string  `json:"date"`
	Left     string  `json:"left"`
	Consumed float64 `json:"consumed"`
	Account  string  `json:"account"`
	Work     string  `json:"work"`
}

func (TaskEstimate) TableName() string {
	return "zt_taskestimate"
}

func (TaskEstimate) GetToday(account string) ([]TaskEstimate, error) {
	var result []TaskEstimate
	date := time.Now().Format("2006-01-02")
	if err := db.Orm.Where(&TaskEstimate{
		Date:    date,
		Account: account,
	}).Find(&result).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	return result, nil
}

func NewTaskEstimate() *TaskEstimate {
	return &TaskEstimate{}
}
