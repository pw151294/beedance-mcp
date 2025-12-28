package create_rag_segment

const ragSegmentCreateToolName = "create_rag_segment"
const ragSegmentCreateToolDesc = "创建知识库文档分段"

const ragDatasetIdParamName = "datasetId"
const ragDatasetIdParamDesc = "知识库ID"
const ragDocumentIdParamName = "documentId"
const ragDocumentIdParamDesc = "文档ID"
const contentParamName = "content"
const contentParamDesc = "分段文本"
const ragCreateSegmentUrl = "/api/v1/datasets/documents/segment/create"
