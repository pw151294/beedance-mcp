package list_services

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
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY2NTc1MTg1LCJpYXQiOjE3NjY1NzE1ODV9.Il-ynClYOdM4VnO2p7mqJV559L4tqAmmYRP1kWq0jvE"
)

func TestListServicesTool(t *testing.T) {
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
	request.Params.Arguments = arguments

	_, err := InvokeListServicesTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking list services tool, %s", err)
	}
}
