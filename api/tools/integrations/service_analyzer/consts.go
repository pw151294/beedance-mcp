package service_analyzer

const serviceErrorAnalyzerToolName = "service_error_analyzer"
const serviceErrorAnalyzerToolDesc = "对服务进行错误链路根因分析"
const serviceSlowAnalyzerToolName = "service_slow_analyzer"
const serviceSlowAnalyzerToolDesc = "对错误进行慢链路根因分析"

const lineBreaker = "\n"

const metricsEndpointMessagePrefix = "该服务"
const endpointSlaMetricsName = "endpoint_sla"
const endpointRtMetricsName = "endpoint_resp_time"

const metricsNameParamName = "metricsName"
const topNParamName = "topN"
const traceStateParamName = "state"
const endpointIdsParamName = "endpointIds"
const traceIDsParamName = "traceIds"
