package list_traces

const listTracesToolName = "list_traces"
const listTracesToolDesc = "查询服务所在链路信息"

const endpointNameParamName = "endpointName"
const endpointNameParamDesc = "接口名称"
const traceStateParamName = "state"
const traceStateParamDesc = "查询链路状态：成功SUCCESS/失败ERROR/所有ALL"

const queryOrder = "BY_DURATION"
const pageNum = 1
const pageSize = 10

const traceInfoPattern = "链路ID：%s；接口：%s；总持续时长：%d毫秒；链路状态：%s\n"
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
