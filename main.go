package main

import (
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/api/tools/apm/metrics_endpoints"
	"beedance-mcp/api/tools/apm/metrics_service_relations"
	"beedance-mcp/api/tools/apm/metrics_services"
	"beedance-mcp/api/tools/apm/services_topology"
	"beedance-mcp/api/tools/trace/detail_trace"
	"beedance-mcp/api/tools/trace/list_traces"
	"beedance-mcp/configs"
	"beedance-mcp/internal/pkg/cache"
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/loggers"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
)

var configPath = flag.String("cfg", "./configs/config.toml", "config file path")

func main() {
	flag.Parse()
	// 加载配置
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading env config")
	}
	if err := configs.InitConfig(*configPath); err != nil {
		log.Fatalf("init config err: %v", err)
	}

	// 初始化
	if err := loggers.InitGlobalLogger(&configs.GlobalConfig.Log); err != nil {
		log.Fatalf("init logger err: %v", err)
	}
	graphql.InitClient()
	cache.InitCacheManager()

	// 创建 MCP 服务器
	s := server.NewMCPServer(
		"apm-collector",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	// 添加apm工具
	s.AddTool(list_services.ListServicesToolSchema(), list_services.InvokeListServicesTool)
	s.AddTool(metrics_services.MetricsServiceToolSchema(), metrics_services.InvokeMetricsServicesTool)
	s.AddTool(services_topology.ServicesTopologyToolSchema(), services_topology.InvokeServicesTopologyTool)
	s.AddTool(metrics_service_relations.MetricsServiceRelationToolSchema(), metrics_service_relations.InvokeMetricsServiceRelationTool)
	s.AddTool(metrics_endpoints.MetricsEndpointsToolSchema(), metrics_endpoints.InvokeMetricsEndpointsTool)

	// 添加trace工具
	s.AddTool(list_traces.ListTracesToolSchema(), list_traces.InvokeListTracesTool)
	s.AddTool(list_traces.EndpointsTracesToolSchema(), list_traces.InvokeEndpointsTracesTool)
	s.AddTool(detail_trace.DetailTraceToolSchema(), detail_trace.InvokeDetailTraceTool)
	s.AddTool(detail_trace.DetailTracesToolSchema(), detail_trace.InvokeDetailTracesTool)

	// 创建并启动 HTTP 服务器
	httpServer := server.NewStreamableHTTPServer(s)
	addr := os.Getenv("APM_SERVER_ADDR")
	if addr == "" {
		addr = ":9601" // 默认端口
	}

	log.Printf("Starting APM MCP server on %s", addr)
	if err := httpServer.Start(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
