package endpoint_analyzer

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/configs"
	"beedance-mcp/internal/pkg/cache"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/loggers"
	"context"
	"log"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	configPath  = "/Users/panwei/Downloads/working/2025.12/beedance-mcp/configs/config.toml"
	workspaceId = "63"
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY2OTk4NzI0LCJpYXQiOjE3NjY5OTUxMjR9.QP7f1wta3KzClzr4r0awA2mbKMDoVoT0z7ZayB7Jyx0"
)

func TestInvokeEndpointErrorAnalyzerTool(t *testing.T) {
	if err := configs.InitConfig(configPath); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := loggers.InitGlobalLogger(&configs.GlobalConfig.Log); err != nil {
		log.Fatalf("Error initializing logger, %s", err)
	}
	cache.InitCacheManager()

	graphql.InitClient()
	request := mcp.CallToolRequest{
		Request: mcp.Request{},
		Header:  make(map[string][]string),
		Params:  mcp.CallToolParams{},
	}
	request.Header.Set(tools.WorkspaceIdHeaderName, workspaceId)
	request.Header.Set(tools.TokenHeaderName, token)
	request.Header.Set("Content-Type", "application/json")

	arguments := make(map[string]any)
	arguments[tools.StartParamName] = "2025-12-29 13:00:00"
	// 这里需要填入真实的 endpointId，可以通过 list_endpoints 工具获取
	arguments[endpointIdsParamName] = []string{"YXV0aHx0b2tfOGQ3M2Y1ZDEyNDhlNGY3YzgyMzljYjc1NTExYWZkNTV8.1_UE9TVDovYXV0aC9hdXRoZW50aWNhdGlvbi9hdXRob3JpemF0aW9uL2FwcGxpY2F0aW9u"}
	request.Params.Arguments = arguments

	result, err := InvokeEndpointErrorAnalyzerTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking endpoint error analyzer tool, %s", err)
	}
	log.Println(convertor.ConvertToolCallResult2Text(result))
}

func TestInvokeEndpointSlowAnalyzerTool(t *testing.T) {
	if err := configs.InitConfig(configPath); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := loggers.InitGlobalLogger(&configs.GlobalConfig.Log); err != nil {
		log.Fatalf("Error initializing logger, %s", err)
	}
	cache.InitCacheManager()

	graphql.InitClient()
	request := mcp.CallToolRequest{
		Request: mcp.Request{},
		Header:  make(map[string][]string),
		Params:  mcp.CallToolParams{},
	}
	request.Header.Set(tools.WorkspaceIdHeaderName, workspaceId)
	request.Header.Set(tools.TokenHeaderName, token)
	request.Header.Set("Content-Type", "application/json")

	arguments := make(map[string]any)
	arguments[tools.StartParamName] = "2025-12-28 16:00:00"
	// 这里需要填入真实的 endpointId，可以通过 list_endpoints 工具获取
	arguments[endpointIdsParamName] = []string{"bjllfHRva184ZDczZjVkMTI0OGU0ZjdjODIzOWNiNzU1MTFhZmQ1NXw=.1_UE9TVDovYXBpL245ZS9ydW0vY291bnQ="}
	request.Params.Arguments = arguments

	result, err := InvokeEndpointSlowAnalyzerTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking endpoint slow analyzer tool, %s", err)
	}
	log.Println(convertor.ConvertToolCallResult2Text(result))
}
