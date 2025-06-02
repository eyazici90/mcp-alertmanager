package tools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterToolStatus(s *server.MCPServer, url string) {
	s.AddTool(toolStatus, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx = context.WithValue(ctx, contextKey("url"), url)
		return toolStatusHandler(ctx, req)
	})
}

var toolStatus = mcp.NewTool("get_status",
	mcp.WithDescription("Get current status of an Alertmanager instance and its cluster. This tool uses `/api/v2/status` endpoint of Alertmanager API."),
	mcp.WithToolAnnotation(mcp.ToolAnnotation{
		Title:           "Get status",
		ReadOnlyHint:    ptr(true),
		DestructiveHint: ptr(false),
		OpenWorldHint:   ptr(true),
	}),
)

func toolStatusHandler(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := ctx.Value(contextKey("url")).(string)
	u := fmt.Sprintf("%s/api/v2/status", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("create request: %v", err)), nil
	}
	return doRequest(req), nil
}
