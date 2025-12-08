package list_services

import (
	"beedance-mcp/internal/pkg/graphql"
	"beedance-mcp/pkg/httputils"
	"beedance-mcp/pkg/loggers"
	"beedance-mcp/pkg/table"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

var serviceRegister *ServicesRegister

type ServicesRegister struct {
	name2Id *table.Table[string, string, string] // workspaceId -> serviceName -> serviceID
	id2Name *table.Table[string, string, string] // workspaceId -> serviceID -> serviceName
}

func RefreshJustForTest(workspaceId string, listServicesResp ListServicesResponse) {
	if serviceRegister == nil {
		serviceRegister = newServiceRegister()
	}
	serviceRegister.refresh(workspaceId, listServicesResp)
}

func ClearServicesRegister() {
	serviceRegister = nil // 置空 等待下次初始化
}

func InitServicesRegister(request mcp.CallToolRequest) {
	if serviceRegister != nil {
		return
	}

	serviceRegister = newServiceRegister()
	variables, err := convert2Variables(request)
	if err != nil {
		loggers.Error("get list services graph request variables failed", zap.Error(err))
		return
	}
	headers, err := httputils.BuildHeaders(request)
	if err != nil {
		loggers.Error("get list services graph request headers failed", zap.Error(err))
		return
	}
	graphqlResp, err := graphql.DoGraphqlRequest[ListServicesVariables, ListServicesResponse](graphQuery, headers, variables)
	if err != nil {
		loggers.Error("get list services graph request failed", zap.Error(err))
		return
	}

	servicesResp := graphqlResp.Data
	loggers.Info("list services response: ", zap.Any("servicesResp", servicesResp))
	serviceRegister.refresh(variables.WorkspaceID, servicesResp)
}

func ServiceIDs(workspaceId string, serviceNames []string) []string {
	return serviceRegister.getServiceIDs(workspaceId, serviceNames)
}

func ServiceName(workspaceId string, serviceID string) string {
	return serviceRegister.getServiceName(workspaceId, serviceID)
}

func newServiceRegister() *ServicesRegister {
	return &ServicesRegister{
		name2Id: table.NewTable[string, string, string](),
		id2Name: table.NewTable[string, string, string](),
	}
}
func (sr *ServicesRegister) refresh(workspaceId string, listServicesResp ListServicesResponse) {
	services := listServicesResp.Services
	if workspaceId == "" || len(services) == 0 {
		return
	}
	if sr.name2Id == nil {
		sr.name2Id = table.NewTable[string, string, string]()
	}
	if sr.id2Name == nil {
		sr.id2Name = table.NewTable[string, string, string]()
	}

	sr.name2Id.Clear()
	sr.id2Name.Clear()
	for _, svc := range services {
		sr.name2Id.Put(workspaceId, svc.Label, svc.ID)
		sr.id2Name.Put(workspaceId, svc.ID, svc.Label)
	}
}
func (sr *ServicesRegister) getServiceIDs(workspaceId string, serviceNames []string) []string {
	serviceIDs := make([]string, 0, 0)
	if len(serviceNames) == 0 {
		return serviceIDs
	}

	for _, svcName := range serviceNames {
		svcID, ok := sr.name2Id.Get(workspaceId, svcName)
		if ok && svcID != "" {
			serviceIDs = append(serviceIDs, svcID)
		}
	}

	return serviceIDs
}
func (sr *ServicesRegister) getServiceName(workspaceId string, serviceID string) string {
	svcName, ok := sr.id2Name.Get(workspaceId, serviceID)
	if ok && svcName != "" {
		return svcName
	}
	return ""
}
