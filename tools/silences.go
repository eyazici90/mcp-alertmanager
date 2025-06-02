package tools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterToolSilences(s *server.MCPServer, url string) {
	s.AddTool(toolSilences, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx = context.WithValue(ctx, contextKey("url"), url)
		return toolSilencesHandler(ctx, req)
	})
}

var toolSilences = mcp.NewTool("list_silences",
	mcp.WithDescription("List of current silences of the Alertmanager instance. This tool uses `/api/v2/silences` endpoint of Alertmanager API."),
	mcp.WithToolAnnotation(mcp.ToolAnnotation{
		Title:           "List of silences",
		ReadOnlyHint:    ptr(true),
		DestructiveHint: ptr(false),
		OpenWorldHint:   ptr(true),
	}),
)

func toolSilencesHandler(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url := ctx.Value(contextKey("url")).(string)
	u := fmt.Sprintf("%s/api/v2/silences", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("create request: %v", err)), nil
	}
	return doRequest(req), nil
}
