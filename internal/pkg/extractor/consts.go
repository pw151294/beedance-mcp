package extractor

const semicolonSplitter = "；"

const endpointIdPrefix = "接口ID："
const endpointSlaPrefix = "成功率："
const endpointRtPrefix = "响应时间："

const endpointTraceDetailSuffix = "的链路详情如下：\n"
const endpointNamePrefix = "接口"
const traceIdPrefix = "链路ID：["
const rightBracket = "]"

const durationLabel = "总持续时长："
const durationUnit = "毫秒"

const slowTraceThreshold = 500 // 慢链路阈值：500毫秒
