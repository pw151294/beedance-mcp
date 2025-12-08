package metrics_service_relations

const (
	metricsServiceRelationsToolName = "metrics_service_relations"
	metricsServiceRelationsToolDesc = "查询当前工作空间下服务调用的负载和响应时间"
)

const metricsClientM0Name = "service_relation_client_cpm"
const metricsClientM1Name = "service_relation_client_resp_time"
const metricsServerM0Name = "service_relation_server_cpm"
const metricsServerM1Name = "service_relation_server_resp_time"

const serviceRelationInfoPattern = "调用方：%s；被调用方：%s；调用方类型：%s；被调用方类型：%s；负载：%d次/分；响应时间：%d毫秒\n"

const serviceRelationClientGraphQuery = `query queryData($duration: Duration!,$ids: [ID!]!,$m0: String!,$m1: String!) {
  service_relation_client_cpm: getValues(metric: {
      name: $m0
      ids: $ids
    }, duration: $duration) {
      values {
        id
        value
      }
    } 
  service_relation_client_resp_time: getValues(metric: {
      name: $m1
      ids: $ids
    }, duration: $duration) {
      values {
        id
        value
      }
    }
}`

const serviceRelationServerGraphqlQuery = `query queryData($duration: Duration!,$ids: [ID!]!,$m0: String!,$m1: String!) {
  service_relation_server_resp_time: getValues(metric: {
      name: $m0
      ids: $ids
    }, duration: $duration) {
      values {
        id
        value
      }
    } 
  service_relation_server_cpm: getValues(metric: {
      name: $m1
      ids: $ids
    }, duration: $duration) {
      values {
        id
        value
      }
    }
}`
