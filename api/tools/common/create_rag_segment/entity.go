package create_rag_segment

type RagSegmentCreateVariable struct {
	DatasetId  string        `json:"datasetId"`
	DocumentId string        `json:"documentId"`
	Content    string        `json:"content"`
	Keywords   []interface{} `json:"keywords"`
}

type RagSegmentCreateResponse struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	TraceId interface{} `json:"traceId"`
}
