package get_current_time

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetCurrentTimeToolSchema() mcp.Tool {
	return mcp.NewTool(
		getCurrentTimeToolName,
		mcp.WithDescription(getCurrentTimeToolDesc),
	)
}

func InvokeGetCurrentTimeTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultText(time.Now().Format(time.DateTime)), nil
}
