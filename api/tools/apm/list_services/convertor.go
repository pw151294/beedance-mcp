package list_services

import (
	"beedance-mcp/api/tools/apm"
	"bytes"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

func convert2Variables(request mcp.CallToolRequest) (ListServicesVariables, error) {
	workspaceId, err := request.RequireString(apm.WorkspaceIdParamName)
	if err != nil {
		return ListServicesVariables{}, fmt.Errorf("工具空间ID参数错误：%w", err)
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
