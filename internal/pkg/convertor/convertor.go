package convertor

import (
	"beedance-mcp/pkg/loggers"
	"encoding/base64"
	"strings"

	"go.uber.org/zap"
)

const splitter = "|"
const callSplitter = "-"
const endpointSplitter = "_"

func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func ConvertServiceIDAndEndpointName2EndpointID(serviceID, endpointName string) string {
	if endpointName == "" {
		return ""
	}
	return serviceID + endpointSplitter + encodeBase64([]byte(endpointName))
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
