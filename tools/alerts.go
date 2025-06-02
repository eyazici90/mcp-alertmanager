package tools

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type contextKey string

func RegisterToolAlerts(s *server.MCPServer, url string) {
	s.AddTool(toolAlerts, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ctx = context.WithValue(ctx, contextKey("url"), url)
		return toolAlertsHandler(ctx, req)
	})
}

var (
	toolAlerts = mcp.NewTool("list_alerts",
		mcp.WithDescription("List of firing and pending alerts of the alertmanager instance. This tool uses `/api/v2/alerts` endpoint of Alertmanager API."),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:           "List of alerts",
			ReadOnlyHint:    ptr(true),
			DestructiveHint: ptr(false),
			OpenWorldHint:   ptr(true),
		}),
		mcp.WithBoolean("active",
			mcp.Title("Show active alerts"),
			mcp.Description("If true, the query will include alerts that have state as active."),
			mcp.DefaultBool(true),
		),
		mcp.WithBoolean("silenced",
			mcp.Title("Show silenced alerts"),
			mcp.Description("If true, the query will include alerts that have state as silenced."),
			mcp.DefaultBool(true),
		),
		mcp.WithBoolean("inhibited",
			mcp.Title("Show inhibited alerts"),
			mcp.Description("If true, the query will include alerts that have state as inhibited."),
			mcp.DefaultBool(true),
		),
		mcp.WithBoolean("unprocessed",
			mcp.Title("Show unprocessed alerts"),
			mcp.Description("If true, the query will include alerts that have state as unprocessed."),
			mcp.DefaultBool(true),
		),
	)
)

func toolAlertsHandler(ctx context.Context, tr mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	active, err := getToolReqParam[string](tr, "active", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	silenced, err := getToolReqParam[string](tr, "silenced", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	inhibited, err := getToolReqParam[string](tr, "inhibited", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	unprocessed, err := getToolReqParam[string](tr, "unprocessed", false)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	url := ctx.Value(contextKey("url")).(string)
	u := fmt.Sprintf("%s/api/v2/alerts", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to create request: %v", err)), nil
	}
	q := req.URL.Query()
	q.Add("active", active)
	q.Add("silenced", silenced)
	q.Add("inhibited", inhibited)
	q.Add("unprocessed", unprocessed)
	return doRequest(req), nil
}
