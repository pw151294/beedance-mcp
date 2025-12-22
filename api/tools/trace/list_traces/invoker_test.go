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
	workspaceId = "3"
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY2Mzk2OTE5LCJpYXQiOjE3NjYzOTMzMTl9.gbXaDknQYO3XOuUD8CcPiV7jWii1crGKcPvrIn61OP0"
)

func TestInvokeListTracesTool(t *testing.T) {
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
	//arguments[tools.StartParamName] = "2025-12-12 14:00:00"
	arguments[tools.ServiceNameParamName] = "auth"
	arguments[endpointNameParamName] = "POST:/authentication/authorization"
	//arguments[traceStateParamName] = "ERROR"
	request.Params.Arguments = arguments

	_, err := InvokeListTracesTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking list traces tool, %s", err)
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
	arguments[endpointIdsParamName] = []string{"YXV0aHx0b2tfMzY2NWQ2ODhiMzI4NGZhMzllYWNlNzE3NWNiMGRlMTR8.1_R0VUOi92YWxpZGF0ZUludGVncmFsaXR5"}
	arguments["state"] = "ERROR"
	request.Params.Arguments = arguments

	_, err := InvokeEndpointsTracesTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking endpoints traces tool, %s", err)
	}
}
