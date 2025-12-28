package create_rag_segment

import (
	"beedance-mcp/pkg/loggers"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

func convert2Variable(request mcp.CallToolRequest) (RagSegmentCreateVariable, error) {
	datasetId, err := request.RequireString(ragDatasetIdParamName)
	if err != nil {
		loggers.Error("parse datasetId from request failed,", zap.Any("request", request))
		return RagSegmentCreateVariable{}, errors.New("获取知识库ID参数失败")
	}
	documentId, err := request.RequireString(ragDocumentIdParamName)
	if err != nil {
		loggers.Error("parse documentId from request failed,", zap.Any("request", request))
		return RagSegmentCreateVariable{}, errors.New("获取知识库文档ID参数失败")
	}
	content, err := request.RequireString(contentParamName)
	if err != nil {
		loggers.Error("parse content from request failed,", zap.Any("request", request))
		return RagSegmentCreateVariable{}, errors.New("获取分段文本参数失败")
	}

	variable := RagSegmentCreateVariable{}
	variable.DatasetId = datasetId
	variable.DocumentId = documentId
	variable.Content = content
	return variable, nil
}
