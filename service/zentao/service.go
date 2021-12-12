package zentao

// service.go:
import (
	"fmt"
	"go-zentao-task/model/zentao"
)

func InitializeService() *Service {
	task := zentao.NewTask()
	user := zentao.NewUser()
	estimate := zentao.NewTaskEstimate()
	service := NewService(task, user, estimate)
	return service
}

func NewService(
	task *zentao.Task, user *zentao.User, estimate *zentao.TaskEstimate) *Service {
	return &Service{
		Task:     task,
		User:     user,
		Estimate: estimate,
	}
}

type Service struct {
	Task     *zentao.Task
	User     *zentao.User
	Estimate *zentao.TaskEstimate
}

func (service *Service) TaskView(id int) string {
	//获取任务
	t, err := service.Task.FindOne(id)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	if t.ID == 0 {
		fmt.Println("任务不存在")
		return ""
	}
	return t.Name
}

func (service *Service) GetAllTaskNotDone() []zentao.Task {
	tasks, err := service.Task.FindAll(service.User.Account, "")
	if err != nil {
		fmt.Println("获取任务异常")
		return nil
	}
	return tasks
}

func (service *Service) GetEstimateToday() (float32, error) {
	estimate, err := service.Estimate.GetToday(service.User.Account)
	if err != nil {
		fmt.Println("获取工时异常")
		return 0, nil
	}
	var consumed float32
	for _, v := range estimate {
		consumed += v.Consumed
	}
	fmt.Println("今日工时填写", consumed)
	return consumed, err
}
