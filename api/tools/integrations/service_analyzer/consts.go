package service_analyzer

const serviceErrorAnalyzerToolName = "service_error_analyzer"
const serviceErrorAnalyzerToolDesc = "对服务进行错误链路根因分析"
const serviceSlowAnalyzerToolName = "service_slow_analyzer"
const serviceSlowAnalyzerToolDesc = "对错误进行慢链路根因分析"

const semicolonSplitter = "；"
const lineBreaker = "\n"

const endpointIdPrefix = "接口ID："
const endpointSlaPrefix = "成功率："
const endpointRtPrefix = "响应时间："
const metricsEndpointMessagePrefix = "该服务"
const endpointSlaMetricsName = "endpoint_sla"
const endpointRtMetricsName = "endpoint_resp_time"

const endpointTraceDetailSuffix = "的链路详情如下：\n"
const endpointNamePrefix = "接口"
const traceIdPrefix = "链路ID：["
const rightBracket = "]"

const metricsNameParamName = "metricsName"
const topNParamName = "topN"
const traceStateParamName = "state"
const endpointIdsParamName = "endpointIds"
const traceIDsParamName = "traceIds"
