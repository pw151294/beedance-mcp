package graphql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}
type GraphqlRequest[V any] struct {
	Query     string `json:"query"`
	Variables V      `json:"variables"`
}
type GraphqlResponse[T any] struct {
	Data   T       `json:"data"`
	Errors []Error `json:"errors,omitempty"`
}

// DoGraphqlRequest 泛型函数 发送graphql请求
func DoGraphqlRequest[V, T any](query string, headers map[string]string, variables V) (*GraphqlResponse[T], error) {
	graphqlClient := GetClient()
	defer PutClient(graphqlClient)

	// 1. 构建http请求
	graphqlReq := GraphqlRequest[V]{}
	graphqlReq.Variables = variables
	graphqlReq.Query = query
	reqBody, err := json.Marshal(graphqlReq)
	if err != nil {
		return nil, fmt.Errorf("序列化graphql请求失败: %w", err)
	}
	httpReq, err := http.NewRequest("POST", graphqlClient.url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建http请求失败: %w", err)
	}
	if len(headers) > 0 {
		for k, v := range headers {
			httpReq.Header.Add(k, v)
		}
	}

	// 2. 发送请求并获取响应
	resp, err := graphqlClient.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送graphql请求失败：%w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取graphql响应失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("graphql请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 3. 解析响应体
	var graphqlResp GraphqlResponse[T]
	if err := json.Unmarshal(respBody, &graphqlResp); err != nil {
		return nil, fmt.Errorf("解析graphql响应体失败：%w", err)
	}
	if len(graphqlResp.Errors) > 0 {
		errs := make([]error, 0, len(graphqlResp.Errors))
		for _, e := range graphqlResp.Errors {
			errs = append(errs, fmt.Errorf("GraphQL 错误: %s", e.Message))
		}
		return nil, errors.Join(errs...)
	}

	return &graphqlResp, nil
}

// DoHttpRequest 泛型函数 发送HTTP请求
func DoHttpRequest[REQ, RES any](url string, headers map[string]string, variables REQ) (*RES, error) {
	graphqlClient := GetClient()
	defer PutClient(graphqlClient)

	// 1. 构建http请求
	reqBody, err := json.Marshal(variables)
	if err != nil {
		return nil, fmt.Errorf("序列化请求体失败: %w", err)
	}
	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建http请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if len(headers) > 0 {
		for k, v := range headers {
			httpReq.Header.Add(k, v)
		}
	}

	// 2. 发送请求并获取响应
	resp, err := graphqlClient.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送http请求失败：%w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取http响应失败: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
	}

	// 3. 解析响应体
	var res RES
	if err := json.Unmarshal(respBody, &res); err != nil {
		return nil, fmt.Errorf("解析http响应体失败：%w", err)
	}

	return &res, nil
}
