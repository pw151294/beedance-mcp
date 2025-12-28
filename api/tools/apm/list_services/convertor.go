package list_services

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/internal/pkg/cache"
	"beedance-mcp/pkg/loggers"
	"bytes"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// ConvertServiceNames2IDs 将服务简称转换成服务ID（服务简称不携带`|token{token}|`后缀，不能通过base64编码获取服务ID）
func ConvertServiceNames2IDs(request mcp.CallToolRequest, workspaceId string, serviceNames []string) []string {
	name2Id := cache.GetByKey[map[string]string](cache.Name2Id, workspaceId, func() any {
		name2Id := make(map[string]string)
		listServicesResp, err := ListServices(request)
		if err != nil {
			loggers.Error("list services failed", zap.Any("request", request), zap.Error(err))
			return name2Id
		}

		services := listServicesResp.Services
		if len(services) > 0 {
			for _, svc := range services {
				name2Id[svc.Label] = svc.ID
			}
		}
		return name2Id
	})

	svcIds := make([]string, 0, len(serviceNames))
	for _, svcName := range serviceNames {
		svcIds = append(svcIds, name2Id[svcName])
	}

	return svcIds
}

func convert2Variables(request mcp.CallToolRequest) (ListServicesVariables, error) {
	workspaceId := request.Header.Get(tools.WorkspaceIdHeaderName)
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
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf(serviceInfoPattern, svc.Label, svc.ID))
		}
	} else {
		toolInvokeMessageBuffer.WriteString("未查询到任务服务")
	}
	return toolInvokeMessageBuffer.String()
}
