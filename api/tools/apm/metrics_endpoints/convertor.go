package metrics_endpoints

import (
	"beedance-mcp/api/tools"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/internal/pkg/convertor"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/timeutils"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// 为服务简称serviceName增加"|{token}|"后缀，示例：auth --> auth|tok_3665d688b3284fa39eace7175cb0de14|
func convertServiceName2Code(request mcp.CallToolRequest, serviceName string) string {
	workspaceId := request.Header.Get(tools.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("workspaceId is empty", zap.Any("request", request))
		return ""
	}
	serviceIDs := list_services.ConvertServiceNames2IDs(request, workspaceId, []string{serviceName})
	return convertor.ConvertServiceID2Code(serviceIDs[0])
}

func convert2Variable(request mcp.CallToolRequest) (MetricsEndpointVariable, error) {
	metricsName, err := request.RequireString(metricsNameParamName)
	if err != nil {
		loggers.Error("parse metrics name param failed", zap.Any("request", request), zap.Error(err))
		return MetricsEndpointVariable{}, fmt.Errorf("获取指标名称失败：%w", err)
	}
	var query, queryOrder string
	switch metricsName {
	case endpointCpmMetricsName:
		query = endpointCpmGraphqlQuery
		queryOrder = "DES"
	case endpointSlaMetricsName:
		query = endpointSlaGraphqlQuery
		queryOrder = "ASC"
	case endpointRespTimeMetricsName:
		query = endpointRespTimeGraphqlQuery
		queryOrder = "DES"
	default:
		return MetricsEndpointVariable{}, fmt.Errorf("指标名称错误：%s", metricsName)
	}

	serviceName, err := request.RequireString(tools.ServiceNameParamName)
	if err != nil {
		loggers.Error("parse service name param failed", zap.Any("request", request), zap.Error(err))
		return MetricsEndpointVariable{}, fmt.Errorf("获取服务名称失败：%w", err)
	}
	serviceCode := convertServiceName2Code(request, serviceName)
	if serviceCode == "" {
		return MetricsEndpointVariable{}, fmt.Errorf("获取服务全称失败：serviceName:%s", serviceName)
	}

	start := request.GetString(tools.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("build duration from start failed", zap.Any("start", start), zap.Error(err))
		return MetricsEndpointVariable{}, fmt.Errorf("获取查询间隔参数失败：%w", err)
	}
	topN := request.GetInt(topNParamName, 10)

	condition := Condition{}
	condition.Name = metricsName
	condition.ParentService = serviceCode
	condition.Normal = true
	condition.Scope = "Endpoint"
	condition.TopN = topN
	condition.Order = queryOrder
	variable := MetricsEndpointVariable{}
	variable.Condition0 = condition
	variable.Duration = duration
	variable.GraphqlQuery = query
	return variable, nil
}

func convertEndpointMetric2Message(metricName string, metricVal MetricValue) string {
	pairs := strings.Split(metricVal.Id, "_")
	if len(pairs) != 2 {
		loggers.Warn("invalid endpoint id format", zap.String("endpointId", metricVal.Id))
		return ""
	}
	serviceName := convertor.ConvertServiceID2Name(pairs[0])
	if serviceName == "" {
		loggers.Warn("invalid service id format", zap.String("serviceId", pairs[0]))
	}

	val, err := strconv.ParseInt(metricVal.Value, 10, 64)
	if err != nil {
		loggers.Error("parse metric value failed", zap.String("metric", metricVal.Value), zap.Error(err))
		return ""
	}

	switch metricName {
	case endpointCpmMetricsName:
		return fmt.Sprintf(endpointCpmMetricsInfoPattern, serviceName, metricVal.Name, val)
	case endpointSlaMetricsName:
		return fmt.Sprintf(endpointSlaMetricsInfoPattern, serviceName, metricVal.Name, convertor.ConvertSlaVal2Rate(val))
	case endpointRespTimeMetricsName:
		return fmt.Sprintf(endpointRespTimeMetricsInfoPattern, serviceName, metricVal.Name, val)
	default:
		return ""
	}
}

func convert2Message(resp MetricsEndpointResponse) string {
	var metricsName string
	var metricValues []MetricValue
	var toolInvokeMessageBuffer bytes.Buffer

	switch {
	case len(resp.MetricsEndpointCpm) > 0:
		metricsName = endpointCpmMetricsName
		metricValues = resp.MetricsEndpointCpm
		toolInvokeMessageBuffer.WriteString(fmt.Sprintf("该服务负载最高的%d个接口信息如下：\n", len(metricValues)))
	case len(resp.MetricsEndpointRespTime) > 0:
		metricsName = endpointRespTimeMetricsName
		metricValues = resp.MetricsEndpointRespTime
		toolInvokeMessageBuffer.WriteString(fmt.Sprintf("该服务响应时间最长的%d个接口信息如下：\n", len(metricValues)))
	case len(resp.MetricsEndpointSla) > 0:
		metricsName = endpointSlaMetricsName
		metricValues = resp.MetricsEndpointSla
		toolInvokeMessageBuffer.WriteString(fmt.Sprintf("该服务成功率最低的%d个接口信息如下：\n", len(metricValues)))
	default:
		toolInvokeMessageBuffer.WriteString("未查询到该服务的任何接口")
	}

	for _, metricVal := range metricValues {
		metricMessage := convertEndpointMetric2Message(metricsName, metricVal)
		toolInvokeMessageBuffer.WriteString(metricMessage)
	}
	return toolInvokeMessageBuffer.String()
}
