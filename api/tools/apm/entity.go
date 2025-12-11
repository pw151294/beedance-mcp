package apm

// MetricValue 指标值
type MetricValue struct {
	ID    string `json:"id"`
	Value int64  `json:"value"`
}

// MetricValues 指标值列表
type MetricValues struct {
	Values []MetricValue `json:"values"`
}
