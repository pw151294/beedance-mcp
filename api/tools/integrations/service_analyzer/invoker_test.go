package service_analyzer

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/configs"
	"beedance-mcp/internal/pkg/cache"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/internal/pkg/extractor"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/loggers"
	"context"
	"encoding/json"
	"log"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

const (
	configPath  = "/Users/panwei/Downloads/working/2025.12/beedance-mcp/configs/config.toml"
	workspaceId = "63"
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY2OTk1MDUwLCJpYXQiOjE3NjY5OTE0NTB9.qh4968Lq__VullAEQSvdf1vB-xC2fLDCrZDZ_OoPDmE"
)

func TestInvokeServiceErrorAnalyzerTool(t *testing.T) {
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
	arguments[tools.ServiceNameParamName] = "auth"
	arguments[tools.StartParamName] = "2025-12-28 13:00:00"
	request.Params.Arguments = arguments

	result, err := InvokeServiceErrorAnalyzerTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking service error analyzer tool, %s", err)
	}
	log.Println(convertor.ConvertToolCallResult2Text(result))
}

func TestInvokeServiceSlowAnalyzerTool(t *testing.T) {
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
	arguments[tools.ServiceNameParamName] = "st-logplatform-service"
	arguments[tools.StartParamName] = "2025-12-28 13:00:00"
	request.Params.Arguments = arguments

	result, err := InvokeServiceSlowAnalyzerTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking service slow analyzer tool, %s", err)
	}
	log.Println(convertor.ConvertToolCallResult2Text(result))
}

func Test_extractEndpointTraces(t *testing.T) {
	content := "接口GET:/permission/userRoleCountAndLatestTime的链路详情如下：\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.72.17665654800730123]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：38毫秒；链路状态：成功\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.66.17665925400175243]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：37毫秒；链路状态：失败\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.66.17665752000302205]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：15毫秒；链路状态：成功\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.66.17665735200251865]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：14毫秒；链路状态：成功\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.70.17666286000421177]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：14毫秒；链路状态：成功\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.68.17665663200420313]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：13毫秒；链路状态：成功\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.74.17665664400740339]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：11毫秒；链路状态：成功\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.74.17665667400350403]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：10毫秒；链路状态：成功\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.69.17665956000525651]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：9毫秒；链路状态：成功\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.75.17665652401150015]；接口：[GET:/permission/userRoleCountAndLatestTime]；总持续时长：9毫秒；链路状态：成功\n接口GET:/permission/dataGroupTree的链路详情如下：\n链路ID：[be82591448aa43de9951c24e4ddb9fbf.71.17665653091080053]；接口：[GET:/permission/dataGroupTree]；总持续时长：40毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.77.17667339108241687]；接口：[GET:/permission/dataGroupTree]；总持续时长：6毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.74.17667295210121583]；接口：[GET:/permission/dataGroupTree]；总持续时长：5毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.75.17667241631891489]；接口：[GET:/permission/dataGroupTree]；总持续时长：4毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.79.17666284034410701]；接口：[GET:/permission/dataGroupTree]；总持续时长：4毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.75.17666284071380697]；接口：[GET:/permission/dataGroupTree]；总持续时长：3毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.70.17667338336161677]；接口：[GET:/permission/dataGroupTree]；总持续时长：3毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.72.17666285183470709]；接口：[GET:/permission/dataGroupTree]；总持续时长：3毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.74.17667242488211487]；接口：[GET:/permission/dataGroupTree]；总持续时长：3毫秒；链路状态：成功\n链路ID：[b8f689f8a98049dfaff9e569a93bd0b9.75.17667246039071521]；接口：[GET:/permission/dataGroupTree]；总持续时长：3毫秒；链路状态：成功\n"
	endpointTraces := extractor.ExtractEndpointTraces(content)
	bytes, _ := json.Marshal(endpointTraces)
	log.Println(string(bytes))
}
