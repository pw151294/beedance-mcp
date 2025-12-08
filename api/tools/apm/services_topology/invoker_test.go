package services_topology

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/configs"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/loggers"
	"encoding/json"
	"log"
	"testing"
)

const (
	configPath       = "/Users/panwei/Downloads/working/2025.12/beedance-mcp/configs/config.toml"
	workspaceId      = "3"
	token            = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY1MTk2ODIxLCJpYXQiOjE3NjUxOTMyMjF9.K7PCj9dkyssU-8xepAKTljxjhW_btk0xmAGES-A0Xo8"
	listServiceQuery = `query queryServices($layer: String!, $workspaceId:String) {
  services: listServicesNew(layer: $layer, workspaceId: $workspaceId) {
    id
    value: name
    label: shortName
    group
    layers
    normal
    groupName
  }
}`
)

func TestGetServicesTopology(t *testing.T) {
	if err := configs.InitConfig(configPath); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	if err := loggers.InitGlobalLogger(&configs.GlobalConfig.Log); err != nil {
		log.Fatalf("Error initializing logger, %s", err)
	}
	graphql.InitClient()

	// 准备测试数据
	duration := apm.Duration{
		Start: "2025-12-05 0934",
		End:   "2025-12-05 1004",
		Step:  "MINUTE",
	}
	headers := map[string]string{
		"workspace-id": workspaceId,
		"Token":        token,
		"Content-Type": "application/json",
	}
	ids := []string{
		"YXV0aHx0b2tfMzY2NWQ2ODhiMzI4NGZhMzllYWNlNzE3NWNiMGRlMTR8.1",
		"c3QtbG9ncGxhdGZvcm0tc2VydmljZXx0b2tfMzY2NWQ2ODhiMzI4NGZhMzllYWNlNzE3NWNiMGRlMTR8.1",
	}
	variables := ServiceTopologyVariables{}
	variables.WorkspaceID = workspaceId
	variables.Duration = duration
	variables.IDs = ids

	// 初始化服务ID-服务名称之间的关系
	vars := list_services.ListServicesVariables{}
	vars.WorkspaceID = workspaceId
	vars.Layer = "GENERAL"
	graphqlResp, err := graphql.DoGraphqlRequest[list_services.ListServicesVariables, list_services.ListServicesResponse](listServiceQuery, headers, vars)
	if err != nil {
		log.Fatalf("list services graphql request failed: %v", err)
	}
	listServicesResp := graphqlResp.Data
	log.Printf("list services graphql response: %v", listServicesResp)
	list_services.RefreshJustForTest(workspaceId, listServicesResp)

	// 测试
	resp, err := graphql.DoGraphqlRequest[ServiceTopologyVariables, ServiceTopologyResponse](graphqlQuery, headers, variables)
	if err != nil {
		log.Fatalf("Error in GraphqlDoGraphqlRequest: %v", err)
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error in json.Marshal: %v", err)
	}
	log.Printf("resp: %s", string(bytes))

	topoRegister = newTopoRegister()
	topoRegister.refresh(workspaceId, resp.Data)
	message := convert2Message(workspaceId, resp.Data)
	log.Printf("tool invoke message: %s", message)
}
