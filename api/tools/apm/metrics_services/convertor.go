package metrics_services

import (
	"beedance-mcp/api/tools/apm"
	"beedance-mcp/api/tools/apm/list_services"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/table"
	"beedance-mcp/pkg/timeutils"
	"bytes"
	"errors"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func convert2Variables(request mcp.CallToolRequest) (ServiceMetricsVariables, error) {
	workspaceId := request.Header.Get(apm.WorkspaceIdHeaderName)
	if workspaceId == "" {
		loggers.Error("parse workspaceId from header failed", zap.Any("headers", request.Header))
		return ServiceMetricsVariables{}, errors.New("请求头未携带工作空间ID")
	}
	serviceNames, err := request.RequireStringSlice(apm.ServiceNamesParamName)
	if err != nil {
		loggers.Error("parse serviceNames failed", zap.Error(err))
		return ServiceMetricsVariables{}, fmt.Errorf("服务名称列表参数错误：%w", err)
	}
	start := request.GetString(apm.StartParamName, "")
	duration, err := timeutils.BuildDuration(start)
	if err != nil {
		loggers.Error("build duration failed", zap.String("start", start), zap.Error(err))
		return ServiceMetricsVariables{}, fmt.Errorf("构建duration参数错误：%w", err)
	}

	variables := ServiceMetricsVariables{}
	variables.WorkspaceID = workspaceId
	variables.Duration = duration
	variables.M0 = metricsM0Name
	variables.M1 = metricsM1Name
	variables.M2 = metricsM2Name
	variables.IDs = list_services.ServiceIDs(workspaceId, serviceNames)
	return variables, nil
}

func convert2Table(serviceMetricsResp ServiceMetricsResponse) *table.Table[string, string, int64] {
	metricsRegister := table.NewTable[string, string, int64]()

	// 0. 服务负载
	cpms := serviceMetricsResp.ServiceCPM.Values
	if len(cpms) > 0 {
		for _, cpm := range cpms {
			metricsRegister.Put(cpm.ID, metricsM0Name, cpm.Value)
		}
	}

	// 1. 服务成功率
	slas := serviceMetricsResp.ServiceSLA.Values
	if len(slas) > 0 {
		for _, sla := range slas {
			metricsRegister.Put(sla.ID, metricsM1Name, sla.Value)
		}
	}

	// 2. 服务的平均响应时间
	rts := serviceMetricsResp.ServiceRespTime.Values
	if len(rts) > 0 {
		for _, rt := range rts {
			metricsRegister.Put(rt.ID, metricsM2Name, rt.Value)
		}
	}

	return metricsRegister
}

func convert2Message(workspaceId string, serviceMetricsResp ServiceMetricsResponse) string {
	metricsRegister := convert2Table(serviceMetricsResp)
	ids := metricsRegister.Rows()

	var toolInvokeMessageBuffer bytes.Buffer
	if len(ids) > 0 {
		toolInvokeMessageBuffer.WriteString("服务的应用性能指标信息如下：\n")
		for _, id := range ids {
			metrics := metricsRegister.Row(id)
			svcName := list_services.ServiceName(workspaceId, id)
			cpm, sla, rt := metrics[metricsM0Name], metrics[metricsM1Name], metrics[metricsM2Name]
			toolInvokeMessageBuffer.WriteString(fmt.Sprintf(serviceMetricsInfoPattern, svcName, cpm, float64(sla)/float64(100), rt))
		}
	} else {
		toolInvokeMessageBuffer.WriteString("未查询到任何服务应用性能指标")
	}

	return toolInvokeMessageBuffer.String()
}
