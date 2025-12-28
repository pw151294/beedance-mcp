package list_traces

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/configs"
	"beedance-mcp/internal/pkg/cache"
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
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY2OTIyMzQzLCJpYXQiOjE3NjY5MTg3NDN9.mHEFK9gjIFBt78S1ZkfTQKwrTfcNJkxki4nQFUxSMAU"
)

func TestInvokeListEndpointsTool(t *testing.T) {
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
	arguments[tools.StartParamName] = "2025-12-18 13:00:00"
	arguments[tools.ServiceNamesParamName] = []string{"auth"}
	request.Params.Arguments = arguments

	_, err := InvokeListEndpointTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking list endpoints tool, %s", err)
	}
}

func TestInvokeEndpointsTracesTool(t *testing.T) {
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
	arguments[tools.StartParamName] = "2025-12-18 13:00:00"
	arguments[tools.ServiceNameParamName] = "auth"
	arguments[endpointIdsParamName] = []string{"YXV0aHx0b2tfOGQ3M2Y1ZDEyNDhlNGY3YzgyMzljYjc1NTExYWZkNTV8.1_R0VUOi9wZXJtaXNzaW9uL2RhdGFHcm91cFRyZWU=", "YXV0aHx0b2tfOGQ3M2Y1ZDEyNDhlNGY3YzgyMzljYjc1NTExYWZkNTV8.1_R0VUOi9wZXJtaXNzaW9uL3VzZXJSb2xlQ291bnRBbmRMYXRlc3RUaW1l"}
	//arguments["state"] = "ERROR"
	request.Params.Arguments = arguments

	_, err := InvokeEndpointsTracesTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking endpoints traces tool, %s", err)
	}
}
