package metrics_services

const metricsServiceToolName = "metrics_services"
const metricsServiceToolDescription = "查询服务的接口负载、成功率还有平均响应时间"

const startParamName = "start"
const startParamDesc = "查询时间范围的起始时间，遵循'YYYY-MM-DD HH::mm:ss'的格式"
const serviceNamesParamName = "serviceNames"
const serviceNamesParamDesc = "服务名称列表"

const metricsM0Name = "service_cpm"
const metricsM1Name = "service_sla"
const metricsM2Name = "service_resp_time"

const serviceMetricsInfoPattern = "服务名称：%s；负载：%d次/分；成功率：%.2f；响应时间：%d毫秒\n"

const graphqlQuery = `query queryData($duration: Duration!,$ids: [ID!]!,$m0: String!,$m1: String!,$m2: String!) {
 service_cpm: getValues(metric: {
     name: $m0
     ids: $ids
   }, duration: $duration) {
     values {
       id
       value
     }
   } 
 service_sla: getValues(metric: {
     name: $m1
     ids: $ids
   }, duration: $duration) {
     values {
       id
       value
     }
   }
 service_resp_time: getValues(metric: {
     name: $m2
     ids: $ids
   }, duration: $duration) {
     values {
       id
       value
     }
   }
}`
