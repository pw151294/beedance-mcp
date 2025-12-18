package endpoints_traces

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/api/tools/trace/list_traces"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/timeutils"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func convert2Variables(request mcp.CallToolRequest) ([]list_traces.ListTracesVariable, error) {
	workspaceId := request.Header.Get(tools.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from request headers failed,", zap.Any("headers", request.Header))
		return nil, errors.New("获取工作空间ID参数失败")
	}
	// 1. 服务名称 + 接口名称列表
	serviceName, err := request.RequireString(tools.ServiceNameParamName)
	if err != nil {
		loggers.Error("parse serviceName from request failed,", zap.Any("request", request), zap.Error(err))
		return nil, fmt.Errorf("获取服务名称参数失败：%w", err)
	}
	serviceIds := list_services.ConvertServiceNames2IDs(request, workspaceId, []string{serviceName})
	serviceId := serviceIds[0]
	endpointNames, err := request.RequireStringSlice(endpointNamesParamName)
	if err != nil || len(endpointNames) == 0 {
		loggers.Error("parse endpointNames from request failed,", zap.Any("request", request), zap.Error(err))
		return nil, fmt.Errorf("获取接口名称列表失败：%w", err)
	}

	// 2.开始时间 链路状态
	start := request.GetString(tools.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("parse duration from request failed", zap.Any("start", start), zap.Error(err))
		return nil, fmt.Errorf("获取查询时间参数失败：%v", err)
	}
	traceState := request.GetString(traceStateParamName, "ALL")

	// 3.查询条件
	variables := make([]list_traces.ListTracesVariable, 0, len(endpointNames))
	for _, endpointName := range endpointNames {
		condition := list_traces.Condition{}
		condition.TraceState = traceState
		condition.ServiceId = serviceId
		condition.EndpointId = convertor.ConvertServiceIDAndEndpointName2EndpointID(serviceId, endpointName)
		condition.QueryOrder = queryOrder
		condition.QueryDuration = duration
		condition.Paging = list_traces.Paging{PageNum: pageNum, PageSize: pageSize}

		variable := list_traces.ListTracesVariable{}
		variable.Condition = condition
		variables = append(variables, variable)
	}
	return variables, nil
}
