package zentao

// service.go:
import (
	"fmt"
	"go-zentao-task/model/zentao"
	"strconv"
	"time"
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
		Task:           task,
		User:           user,
		Estimate:       estimate,
		ProjectProduct: projectProduct,
		Action:         action,
	}
}

type Service struct {
	Task           *zentao.Task
	User           *zentao.User
	Estimate       *zentao.TaskEstimate
	ProjectProduct *zentao.ProjectProduct
	Action         *zentao.Action
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

func (service *Service) GetEstimateToday() (float64, error) {
	estimate, err := service.Estimate.GetToday(service.User.Account)
	if err != nil {
		fmt.Println("获取工时异常")
		return 0, nil
	}
	var consumed float64
	for _, v := range estimate {
		consumed += v.Consumed
	}
	fmt.Println("今日工时填写", consumed)
	return consumed, err
}

func (service *Service) UpdateTask(task int, estimate float64, action string) float64 {
	taskInfo, err := service.Task.FindOne(task)
	if err != nil {
		fmt.Println("获取任务异常")
		return 0
	}
	//productInfo, _:= service.ProjectProduct.FindOneByProject(taskInfo.Project)
	if action == "finished" {
		if taskInfo.FromBug != 0 {
			fmt.Println("任务从bug创建，不能直接完成")
			return 0
		}
		estimate = taskInfo.Left
		taskInfo.FinishedDate = time.Now()
	}
	if estimate > taskInfo.Left {
		fmt.Println("消耗工时不能大于剩余工时")
		return 0
	}
	//创建操作记录
	//service.Action.Create(task,"task", ","+ strconv.Itoa(productInfo.Product)+",", taskInfo.Project, taskInfo.Execution, taskInfo.AssignedTo, action, estimate)

	service.Task.UpdateOne(task, estimate+taskInfo.Consumed, taskInfo.Left-estimate, taskInfo.AssignedTo, taskInfo.FinishedDate)
	return estimate
}

func (service *Service) GetOptimumTasks() []int {
	var tasks []zentao.Task
	var res []int
	return service.OptimumTask(tasks, res, 0)
}

// OptimumTask 竞选出最优任务， 用于更新
func (service *Service) OptimumTask(tasks []zentao.Task, result []int, round int) []int {
	var tmpTasks []zentao.Task
	if len(tasks) == 0 {
		tasks = service.GetAllTaskNotDone()
		result = []int{}
		round = 0
	}
	if round > 3 {
		return result
	}
	for _, task := range tasks {
		switch task.Type {
		case zentao.TypeDiscuss: //优先讨论类型的任务
			result = append(result, task.ID)
		case zentao.TypeDev:
			if round == 1 && task.Left <= 8 && task.Status == zentao.StatusDo {
				result = append(result, task.ID)
			} else if round == 2 && task.Left > 8 && task.Status == zentao.StatusDo { // 第三轮竞选，剩余时间大于8天并且在doing状态的
				result = append(result, task.ID)
			} else if round == 3 && task.Status == zentao.StatusWait && task.Left > 0 { // 第四轮竞选，等待状态的case,并且剩余时间不能为0
				result = append(result, task.ID)
			} else {
				tmpTasks = append(tmpTasks, task)
			}
		}
	}
	if len(tmpTasks) > 0 {
		result = service.OptimumTask(tmpTasks, result, round+1)
	}
	return result
}

// ConsumeRecord 记录工时
func (service *Service) ConsumeRecord() float64 {
	estimate := 0.5
	var current float64
	tasks := service.GetOptimumTasks()
	for _, task := range tasks {
		current += service.UpdateTask(task, estimate, "recordestimate")
		if current >= 1 {
			break
		}
	}
	current, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", current), 64)
	return current
}
