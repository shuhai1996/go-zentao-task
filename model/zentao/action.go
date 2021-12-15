package zentao

import (
	"go-zentao-task/pkg/db"
	"time"
)

type Action struct {
	ID       int    `json:"id"`
	ObjectID  int `json:"objectID"`
	ObjectType   string `json:"objectType"`
	Product   string `json:"product"`
	Project   int `json:"project"`
	Execution   int `json:"execution"`
	Actor     string `json:"actor"`
	Action     string `json:"action"`
	Extra float32 `json:"extra"`
	Date time.Time `json:"date"`
}

func NewAction() *Action {
	return &Action{}
}

func (*Action) Create(objectID int, objectType string, product string, project int, execution int, actor string, action string) (int, error) {
	data := &Action{
		ObjectID: objectID,
		ObjectType: objectType,
		Product: product,
		Project: project,
		Execution:   execution,
		Actor:  actor,
		Action: action,
		Date: time.Now(),
	}
	op := db.Orm.Create(data)
	if op.Error != nil {
		return 0, op.Error
	}
	return data.ID, nil
}
