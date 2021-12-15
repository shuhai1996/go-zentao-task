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
	projectProduct := zentao.NewProjectProduct()
	action := zentao.NewAction()
	service := NewService(task, user, estimate, projectProduct, action)
	return service
}

func NewService(
	task *zentao.Task, user *zentao.User, estimate *zentao.TaskEstimate, projectProduct *zentao.ProjectProduct, action *zentao.Action) *Service {
	return &Service{
		Task:     task,
		User:     user,
		Estimate: estimate,
		ProjectProduct: projectProduct,
		Action: action,
	}
}

type Service struct {
	Task     *zentao.Task
	User     *zentao.User
	Estimate *zentao.TaskEstimate
	ProjectProduct *zentao.ProjectProduct
	Action *zentao.Action
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

func (service *Service) UpdateTask(task int, estimate float32, action string) float32{
	taskInfo,err :=  service.Task.FindOne(task)
	if err !=nil {
		fmt.Println("获取任务异常")
		return 0
	}
	productInfo, _:= service.ProjectProduct.FindOneByProject(taskInfo.Project)
	if action== "finished" {
		if taskInfo.FromBug !=0 {
			fmt.Println("任务从bug创建，不能直接完成")
			return 0
		}
		estimate = taskInfo.Left
	}
	if estimate > taskInfo.Left {
		fmt.Println("消耗工时不能大于剩余工时")
		return 0
	}
	//创建操作记录
	service.Action.Create(task,"task", ","+string(productInfo.Product)+",", taskInfo.Project, taskInfo.Execution, taskInfo.AssignedTo, action)

	return estimate
}
