package metrics_service_relations

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/timeutils"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func convert2Variables(request mcp.CallToolRequest) (ServiceRelationMetricsVariables, error) {
	workspaceId := request.Header.Get(apm.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from header failed", zap.Any("headers", request.Header))
		return ServiceRelationMetricsVariables{}, errors.New("请求头未携带工作空间ID")
	}
	serviceNames, err := request.RequireStringSlice(apm.ServiceNamesParamName)
	if err != nil {
		loggers.Error("parse serviceNames failed", zap.Error(err))
		return ServiceRelationMetricsVariables{}, fmt.Errorf("服务名称列表参数错误：%w", err)
	}
	start := request.GetString(apm.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("build duration failed", zap.String("start", start), zap.Error(err))
		return ServiceRelationMetricsVariables{}, fmt.Errorf("构建duration参数错误：%w", err)
	}

	variables := ServiceRelationMetricsVariables{}
	variables.WorkspaceID = workspaceId
	variables.Duration = duration
	variables.IDs = list_services.ServiceIDs(workspaceId, serviceNames)
	return variables, nil
}

func convert2ClientVariables(request mcp.CallToolRequest) (ServiceRelationMetricsVariables, error) {
	variables, err := convert2Variables(request)
	if err != nil {
		return variables, fmt.Errorf("转换服务调用关系客户端指标查询参数失败：%w", err)
	}
	variables.M0 = metricsClientM0Name
	variables.M1 = metricsClientM1Name
	return variables, nil
}

func convert2ServerVariables(request mcp.CallToolRequest) (ServiceRelationMetricsVariables, error) {
	variables, err := convert2Variables(request)
	if err != nil {
		return variables, fmt.Errorf("转换服务调用关系服务端查询参数失败：%w", err)
	}
	variables.M0 = metricsServerM0Name
	variables.M1 = metricsServerM1Name
	return variables, nil
}
