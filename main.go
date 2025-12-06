package main

import (
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/configs"
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
