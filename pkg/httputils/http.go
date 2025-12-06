package httputils

import "github.com/mark3labs/mcp-go/mcp"

func BuildHeaders(request mcp.CallToolRequest) (map[string]string, error) {
	workspaceId, err := request.RequireString("workspaceId")
	if err != nil {
		return nil, err
	}
	token, err := request.RequireString("token")
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"workspace-id": workspaceId,
		"token":        token,
		"Content-Type": "application/json",
	}, nil
}
