package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestMockAPMMetrics(t *testing.T) {
	params := make(map[string]interface{})
	requestBodyBytes, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("marshal request body failed: %v", err)
	}

	request, err := http.NewRequest(http.MethodPost, fmt.Sprintf(url, 9016), bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		log.Fatalf("create request failed: %v", err)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("send request failed: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		log.Fatalf("send request failed: %v", response.Status)
	}
	defer response.Body.Close()
	respBytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("read response body failed: %v", err)
	}
	log.Printf("response body: %s", string(respBytes))
}
