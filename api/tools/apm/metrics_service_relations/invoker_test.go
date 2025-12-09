package metrics_service_relations

import (
	"beedance-mcp/api/tools/apm"
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
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY1Mjg5MTUzLCJpYXQiOjE3NjUyODU1NTN9.n-0Q-MPAeN2O5CXGn2z1nrS42lAWQvQBqooCp04rNH0"
)

func TestInvokeMetricsServiceRelationTool(t *testing.T) {
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
	request.Header.Set(apm.WorkspaceIdHeaderName, workspaceId)
	request.Header.Set(apm.TokenHeaderName, token)
	request.Header.Set("Content-Type", "application/json")
	arguments := make(map[string]any)
	arguments[apm.ServiceNamesParamName] = []string{"st-logplatform-service", "auth"}
	request.Params.Arguments = arguments

	_, err := InvokeMetricsServiceRelationTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking metrics service relation tool, %s", err)
	}
}
