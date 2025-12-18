package detail_trace

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
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY2MDUyMzE5LCJpYXQiOjE3NjYwNDg3MTl9.7gi_NBbBgIFER99MutqfEu_juDyqz73BW0VeDpYpljg"
)

func TestInvokeDetailTraceTool(t *testing.T) {
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
	arguments[traceIDParamName] = "fc19534e3bbb43dfab67b77a1ba0cd30.96.17658928428474411"
	request.Params.Arguments = arguments

	_, err := InvokeDetailTraceTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking detail trace tool, %s", err)
	}
}

func TestInvokeDetailTracesTool(t *testing.T) {
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
	// 这里沿用 detail_trace 的 traceId 入参；如需覆盖更多场景可替换为另一条有效 traceId
	arguments[traceIDsParamName] = []string{"fc19534e3bbb43dfab67b77a1ba0cd30.87.17660494184401361"}
	request.Params.Arguments = arguments

	_, err := InvokeDetailTracesTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking detail traces tool, %s", err)
	}
}
