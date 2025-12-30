package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func sendRequest() {
	url := "http://172.30.34.73/alarm-manager/event/page/agg"
	method := "POST"

	payload := strings.NewReader(`{
    "aggId": "1999813425740251136",
    "pageNo": 1,
    "pageSize": 100000,
    "total": 0
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Printf("创建请求失败: %v\n", err)
		return
	}
	req.Header.Add("Token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY3MDA5ODg3LCJpYXQiOjE3NjcwMDYyODd9.IhFEculgOY5JIpsRl5WYCD4Rs2HCERZna1YoGKgh1po")
	req.Header.Add("workspace-id", "1")
	req.Header.Add("Content-Type", "application/json")

	_, err = client.Do(req)
	if err != nil {
		log.Printf("发送请求失败: %v\n", err)
		return
	}
}

func main() {
	log.Println("定时任务启动，每3秒发送一次请求...")

	// 创建一个每3秒触发一次的定时器
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// 持续监听定时器信号
	for range ticker.C {
		log.Printf("[%s] 开始执行定时任务...\n", time.Now().Format("2006-01-02 15:04:05"))
		sendRequest()
	}
}
