package list_traces

const listEndpointsToolName = "list_endpoints"
const listEndpointsToolDesc = "查询服务关联的接口名称和ID"

const endpointsTracesToolName = "endpoints_traces"
const endpointsTracesToolDesc = "查询接口列表中所有接口的链路信息"

const endpointIdsParamName = "endpointIds"
const endpointIdsParamDesc = "接口ID列表"

const traceStateParamName = "state"
const traceStateParamDesc = "查询链路状态：成功SUCCESS/失败ERROR/所有ALL"

const queryOrder = "BY_DURATION"
const pageNum = 1
const pageSize = 10

const endpointInfoPattern = "接口名称：%s，接口ID：%s，服务名称：%s，服务ID：%s\n"
const traceInfoPattern = "链路ID：%s；接口：%s；总持续时长：%d毫秒；链路状态：%s\n"

const listEndpointsGraphqlQuery = `query queryEndpoints($serviceId: ID!, $keyword: String!) {
  pods: findEndpoint(serviceId: $serviceId, keyword: $keyword, limit: 100) {
    id
    value: name
     label: name
    }
}`
const listTracesGraphqlQuery = `query queryTraces($condition: TraceQueryCondition) {
  data: queryBasicTraces(condition: $condition) {
    traces {
      key: segmentId
      endpointNames
      duration
      start
      isError
      traceIds
    }
  }}`
