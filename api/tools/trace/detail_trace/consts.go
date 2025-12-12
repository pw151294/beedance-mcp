package detail_trace

const detailTraceToolName = "detail_trace"
const detailTraceToolDesc = "根据链路ID查询链路详情"

const traceIDParamName = "traceId"
const traceIDParamDesc = "链路ID"

const spanInfoPattern = "应用：%s；接口：%s；耗时：%d毫秒；组件：%s；状态：%s；层级：%s；标记：%s\n"
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
