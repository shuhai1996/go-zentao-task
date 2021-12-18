package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-zentao-task/pkg/util"
	"io/ioutil"
	"net/http"
)

type Notification struct {
}

type Marks struct {
	Content string `json:"content"`
}

// PostBody post 请求
type PostBody struct {
	Msgtype  string `json:"msgtype"`
	Markdown Marks  `json:"markdown"`
}

func NewNotification() *Notification {
	return &Notification{}
}

// SendNotification 发送机器人报警
func (notification *Notification) SendNotification(name string, estimate string, count string, ids []int) []byte {
	url := util.GetRobotUrl()
	idSt, err := json.Marshal(ids)
	// 构造POST请求
	postBody := &PostBody{
		Msgtype: "markdown",
		Markdown: Marks{
			Content: "禅道<font color=\"warning\">工时</font>，请相关同事注意。\n>昵称<font color=\"comment\">" + name + "</font>\n>手动填写用时:<font color=\"comment\">" + estimate + "</font>\n>自动填写用时:<font color=\"comment\">" + count + "</font>\n>任务id:<font color=\"comment\">" + string(idSt) + "</font>",
		},
	}
	// struct 转json
	body, _ := json.Marshal(postBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	statuscode := resp.StatusCode
	hea := resp.Header
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println(hea)
	fmt.Println(statuscode)
	return body
}
