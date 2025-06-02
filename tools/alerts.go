package tools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

var (
	toolAlerts = mcp.NewTool("alerts",
		mcp.WithDescription("List of firing and pending alerts of the alertmanager instance. This tool uses `/api/v1/alerts` endpoint of Alertmanager API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of alerts",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
	)
)

func toolAlertsHandler(ctx context.Context, tr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, err := getToolReqParam[string](tr, "param", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "target-url", nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	return doRequest(req), nil
}
