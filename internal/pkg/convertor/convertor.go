package convertor

import (
	"beedance-mcp/pkg/loggers"
	"encoding/base64"
	"strings"

	"go.uber.org/zap"
)

const splitter = "|"
const callSplitter = "-"

func decodeBase64(s string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(s)
}

func ConvertID2Name(serviceID string) string {
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
