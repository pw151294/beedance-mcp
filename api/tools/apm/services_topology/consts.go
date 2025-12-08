package services_topology

const (
	serviceTopologyToolName = "services_topology"
	serviceTopologyToolDesc = "查询服务之间的调用拓扑图"
)

const serviceCallInfoPattern = "调用方：%s；调用方类型：%s；被调用方：%s；被调用方类型：%s\n"

const graphqlQuery = `query queryData($duration: Duration!, $serviceIds: [ID!]!) {
  topology: getServicesTopology(duration: $duration, serviceIds: $serviceIds) {
    nodes {
      id
      name
      type
      isReal
    }
    calls {
      id
      source
      detectPoints
      target
    }
  }
}`
