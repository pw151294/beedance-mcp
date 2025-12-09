package list_services

const (
	listServicesToolName = "list_services"
	listServicesToolDesc = "查询当前工作空间下的所有服务的名称还有ID"
)

const layer = "GENERAL"
const serviceInfoPattern = "服务名称：%s；服务ID：%s\n"

const graphqlQuery = `query queryServices($layer: String!, $workspaceId:String) {
  services: listServicesNew(layer: $layer, workspaceId: $workspaceId) {
    id
    value: name
    label: shortName
    group
    layers
    normal
    groupName
  }
}`
