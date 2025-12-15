package metrics_endpoints

const metricsEndpointsToolName = "metrics_endpoints"
const metricsEndpointsToolDesc = "查看服务所在接口的负载、响应时间、成功率指标的topN个元素"

const topNParamName = "topN"
const topNParamDesc = "查询的元素数量"
const metricsNameParamName = "metricsName"
const metricsNameParamDesc = "查询的指标：成功率endpoint_sla/负载endpoint_cpm/响应时间endpoint_resp_time"

const endpointCpmMetricsName = "endpoint_cpm"
const endpointSlaMetricsName = "endpoint_sla"
const endpointRespTimeMetricsName = "endpoint_resp_time"

const endpointCpmMetricsInfoPattern = "服务：%s；接口：%s；负载：%d次/分\n"
const endpointRespTimeMetricsInfoPattern = "服务：%s；接口：%s；响应时间：%d毫秒\n"
const endpointSlaMetricsInfoPattern = "服务：%s；接口：%s；成功率：%.2f\n"

const endpointCpmGraphqlQuery = `query queryData($duration: Duration!,$condition0: TopNCondition!) {endpoint_cpm0: sortMetrics(condition: $condition0, duration: $duration){
    name
    id
    value
    refId
  }}`
const endpointRespTimeGraphqlQuery = ` query queryData($duration: Duration!,$condition0: TopNCondition!) {endpoint_resp_time0: sortMetrics(condition: $condition0, duration: $duration){
    name
    id
    value
    refId
  }}`
const endpointSlaGraphqlQuery = `query queryData($duration: Duration!,$condition0: TopNCondition!) {endpoint_sla0: sortMetrics(condition: $condition0, duration: $duration){
    name
    id
    value
    refId
  }}
`
