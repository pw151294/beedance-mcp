package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const url = "http://172.31.186.217:%d/mock/topology/send/deepflow"
const (
	periodA = 1
	periodB = 2
	periodC = 3
	periodD = 6
)
const (
	portA = 9021
	portB = 9006
	portC = 9011
	portD = 9016
)

var client *http.Client

func sendRequest(port, period int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(period)) * time.Second)

	params := make(map[string]interface{})
	reqBodyBytes, _ := json.Marshal(params)
	for {
		request, _ := http.NewRequest(http.MethodPost, fmt.Sprintf(url, port), bytes.NewBuffer(reqBodyBytes))
		log.Printf("send request, time: %s, port: %d", time.Now().Format(time.DateTime), port)
		resp, err := client.Do(request)
		if err != nil {
			log.Printf("send request failed, error: %v", err)
		} else {
			respBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("read response failed, error: %v", err)
			} else {
				log.Printf("send request success, response: %s", string(respBytes))
			}
			resp.Body.Close()
		}

		time.Sleep(time.Duration(period) * time.Second)
	}
}

func main() {
	client = &http.Client{}
	go sendRequest(portA, periodA)
	go sendRequest(portB, periodB)
	go sendRequest(portC, periodC)
	go sendRequest(portD, periodD)

	log.Println("begin send requests to nodes")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

}
