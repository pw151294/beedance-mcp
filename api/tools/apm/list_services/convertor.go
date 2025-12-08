package list_services

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/pkg/loggers"
	"bytes"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func convert2Variables(request mcp.CallToolRequest) (ListServicesVariables, error) {
	workspaceId := request.Header.Get(apm.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from header failed", zap.Any("headers", request.Header))
		return ListServicesVariables{}, errors.New("请求头未携带工作空间ID")
	}

	variables := ListServicesVariables{}
	variables.WorkspaceID = workspaceId
	variables.Layer = layer
	return variables, nil
}

func convert2Message(response ListServicesResponse) string {
	var toolInvokeMessageBuffer bytes.Buffer
	services := response.Services
	if len(services) > 0 {
		toolInvokeMessageBuffer.WriteString("查询到以下服务信息：\n")
		for _, svc := range services {
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf(serviceInfoPattern, svc.Value, svc.Label, svc.ID))
		}
	} else {
		toolInvokeMessageBuffer.WriteString("未查询到任务服务")
	}
	return toolInvokeMessageBuffer.String()
}
