package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
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
func (notification *Notification) SendNotification(url string, name string, estimate string) []byte {
	// 构造POST请求
	postBody := &PostBody{
		Msgtype: "markdown",
		Markdown: Marks{
			Content: "禅道<font color=\"warning\">工时</font>，请相关同事注意。\n>昵称<font color=\"comment\">" + name + "</font>\n>用时:<font color=\"comment\">" + estimate + "</font>\n>",
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
