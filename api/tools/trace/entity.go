package trace

import "beedance-mcp/api/tools"

type EndpointQueryRequest struct {
	ServiceId string         `json:"serviceId"`
	Duration  tools.Duration `json:"duration"`
}

type Pod struct {
	Id    string `json:"id"`
	Value string `json:"value"`
	Label string `json:"label"`
}

type EndpointQueryRespData struct {
	Pods []Pod `json:"pods"`
}

type EndpointQueryResp struct {
	Data EndpointQueryRespData `json:"data"`
}
