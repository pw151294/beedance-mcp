package detail_trace

const detailTraceToolName = "detail_trace"
const detailTraceToolDesc = "根据链路ID查询链路详情"

const traceIDParamName = "traceId"
const traceIDParamDesc = "链路ID"

const eventPropertyName = "event"
const errorKindPropertyName = "error.kind"
const messagePropertyName = "message"
const stackPropertyName = "stack"
const errorLogInfoPattern = "错误%d：错误类型：%s；错误信息：%s；错误堆栈：%s；"
const stackLengthThreshold = 100

const spanInfoPattern = "应用：%s；接口：%s；耗时：%d毫秒；组件：%s；状态：%s；层级：%s；属性：%s；日志：%s\n"
const detailTraceGraphqlQuery = `query queryTrace($traceId: ID!) {
  trace: queryTrace(traceId: $traceId) {
    spans {
      traceId
      segmentId
      spanId
      parentSpanId
      refs {
        traceId
        parentSegmentId
        parentSpanId
        type
      }
      serviceCode
      serviceInstanceName
      startTime
      endTime
      endpointName
      type
      peer
      component
      isError
      layer
      tags {
        key
        value
      }
      logs {
        time
        data {
          key
          value
        }
      }
      attachedEvents {
        startTime {
          seconds
          nanos
        }
        event
        endTime {
          seconds
          nanos
        }
        tags {
          key
          value
        }
        summary {
          key
          value
        }
      }
    }
  }
  }`
