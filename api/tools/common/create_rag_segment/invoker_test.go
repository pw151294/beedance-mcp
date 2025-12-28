package create_rag_segment

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
	workspaceId = "1"
	token       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxIiwiYWNjb3VudCI6ImFkbWluIiwiZXhwIjoxNzY2OTEyMTc0LCJpYXQiOjE3NjY5MDg1NzR9.r2rTGggal_kMI927cNvQQ94ZcyQG8Q5ZLwvOT5EHd6g"
)

func TestInvokeCreateRagSegmentTool(t *testing.T) {
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
	arguments[ragDatasetIdParamName] = "cff13dcb1bc04950859c8efeff0773c5"
	arguments[ragDocumentIdParamName] = "a127c605302146fb982ddac8096ea239"
	arguments[contentParamName] = "```\\n我现在需要构建一个可观测性分析智能体，为这个智能体嵌入了5个工具，工具及其参数描述和响应示例如下：\\n```\\n1. list_services：查询当前工作空间下的所有服务的名称还有ID\\n2. metrics_services：查询服务的接口负载、成功率还有平均响应时间\\n   - serviceNames：服务名称列表（必填）\\n   - start：查询时间范围的起始时间，遵循'YYYY-MM-DD HH::mm:ss'的格式（可选，默认半小时内）\\n   - 响应示例：\\n3. metrics_endpoints：查看服务所在接口的负载、响应时间、成功率指标的topN个元素\\n   - serviceName：服务名称（必填）\\n   - metricsName：查询的指标：成功率endpoint_sla/负载endpoint_cpm/响应时间endpoint_resp_time（必填）\\n   - start：查询时间范围的起始时间，遵循'YYYY-MM-DD HH::mm:ss'的格式（可选，默认半小时内）\\n   - topN：查询的元素数量（可选，默认查询前5）\\n4. endpoints_traces：查询接口列表中所有接口的链路信息\\n   - serviceName：服务名称（必填）\\n   - endpointNames：接口名称列表（必填）\\n   - start：查询时间范围的起始时间，遵循'YYYY-MM-DD HH::mm:ss'的格式（可选，默认半小时内）\\n   - state：查询链路状态：成功SUCCESS/失败ERROR/所有ALL（可选，默认ALL）\\n5. detail_traces：根据链路ID数组批量查询链路详情\\n   - traceIds：链路ID数组（必填）\\n```\\n这些工具具备了应用性能分析或者链路分析的功能，我现在需要和智能体针对不同类型问题的处理流程和解决方案预设一份规约，请结合你对可观测领域的理解，结合这些工具信息，编写一份5000字以内的Spec规约，可以让智能体通过遵循这份规约来帮助用户进行服务的应用性能指标查询和分析、链路分析、故障定位、错误根因分析等一系列可观测功能。\\n```"
	request.Params.Arguments = arguments

	_, err := InvokeCreateRagSegmentTool(context.Background(), request)
	if err != nil {
		log.Fatalf("Error invoking create rag segment tool, %s", err)
	}
}
