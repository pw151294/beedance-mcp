package convertor

import (
	"beedance-mcp/pkg/loggers"
	"encoding/base64"
	"strings"

	"go.uber.org/zap"
)

const splitter = "|"

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
