package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WxMessage struct {
	MsgType string        `json:"msgtype"`
	Text    WxMessageText `json:"text"`
}

type WxMessageText struct {
	Content string `json:"content"`
}

func PushWX(url, messageText string) {

	msg := WxMessage{
		MsgType: "text",
		Text: WxMessageText{
			Content: messageText,
		},
	}
	jsonData, _ := json.Marshal(msg)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体
}
