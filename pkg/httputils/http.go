package httputils

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/pkg/loggers"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func BuildHeaders(request mcp.CallToolRequest) (map[string]string, error) {
	workspaceId := request.Header.Get(tools.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from header failed", zap.Any("headers", request.Header))
		return nil, errors.New("请求头未携带工作空间ID")
	}
	token := request.Header.Get(tools.TokenHeaderName)
	if token == "" {
		loggers.Error("parse token from header failed", zap.Any("headers", request.Header))
		return nil, errors.New("请求头未携带Token认证令牌")
	}
	return map[string]string{
		"workspace-id": workspaceId,
		"token":        token,
		"Content-Type": "application/json",
	}, nil
}
