package convertor

import (
	"beedance-mcp/pkg/loggers"
	"encoding/base64"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

const splitter = "|"
const callSplitter = "-"
const endpointSplitter = "_"

func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func ConvertEndpointID2ServiceIDAndEndpointName(endpointId string) (string, string) {
	pairs := strings.Split(endpointId, endpointSplitter)
	if len(pairs) != 2 {
		loggers.Warn("invalid endpoint id", zap.String("endpointId", endpointId))
		return "", ""
	}
	endpointName, err := decodeBase64(pairs[1])
	if err != nil {
		loggers.Warn("invalid endpoint code", zap.String("endpointCode", pairs[1]), zap.Error(err))
		return "", ""
	}

	return ConvertServiceID2Name(pairs[0]), string(endpointName)
}

func ConvertEndpointID2Name(endpointId string) string {
	pairs := strings.Split(endpointId, endpointSplitter)
	if len(pairs) != 2 {
		loggers.Warn("invalid endpoint id", zap.String("endpointId", endpointId))
		return ""
	}
	endpointName, err := decodeBase64(pairs[1])
	if err != nil {
		loggers.Warn("invalid endpoint code", zap.String("endpointCode", pairs[1]), zap.Error(err))
	}
	return string(endpointName)
}

func ConvertServiceID2Name(serviceID string) string {
	// 示例："YXV0aHx0b2tfMzY2NWQ2ODhiMzI4NGZhMzllYWNlNzE3NWNiMGRlMTR8.1" 只需要解码.之前的内容即可
	parts := strings.SplitN(serviceID, ".", 2)
	base64Part := parts[0]

	serviceNameBytes, err := decodeBase64(base64Part)
	if err != nil {
		loggers.Warn("decode serviceID failed", zap.String("serviceID", serviceID), zap.String("base64Part", base64Part), zap.Error(err))
		return ""
	}
	fullServiceName := string(serviceNameBytes)
	if strings.HasSuffix(fullServiceName, splitter) {
		idx := strings.Index(fullServiceName, splitter)
		return fullServiceName[:idx]
	}
	return fullServiceName
}

func ConvertCallID2ServiceIDs(callID string) (string, string) {
	svcIds := strings.Split(callID, callSplitter)
	if len(svcIds) != 2 {
		loggers.Warn("invalid call id format", zap.String("callID", callID))
		return "", ""
	}
	return svcIds[0], svcIds[1]
}

func ConvertServiceID2Code(serviceId string) string {
	pairs := strings.Split(serviceId, ".")
	if len(pairs) != 2 {
		loggers.Warn("invalid service id format", zap.String("serviceId", serviceId))
		return ""
	}
	serviceNameBytes, err := decodeBase64(pairs[0])
	if err != nil {
		loggers.Warn("decode serviceId failed", zap.String("serviceId", serviceId))
		return ""
	}
	return string(serviceNameBytes)
}

func ConvertServiceCode2Name(serviceCode string) string {
	if strings.HasSuffix(serviceCode, splitter) {
		idx := strings.Index(serviceCode, splitter)
		return serviceCode[:idx]
	}
	return serviceCode
}

func ConvertBool2Desc(isError bool) string {
	if isError {
		return "失败"
	} else {
		return "成功"
	}
}

func ConvertSlaVal2Rate(sla int64) float64 {
	return float64(sla) / float64(100)
}

func ConvertToolCallResult2Text(result *mcp.CallToolResult) string {
	contents := result.Content
	if len(contents) == 0 {
		return ""
	} else {
		return contents[0].(mcp.TextContent).Text
	}
}
