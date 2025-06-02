package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func ptr[T any](v T) *T {
	return &v
}

func RegisterToolAlerts(s *server.MCPServer) {
	s.AddTool(toolAlerts, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolAlertsHandler(ctx, req)
	})
}

type toolReqParamType interface {
	string | float64 | bool | []string | []any
}

func getToolReqParam[T toolReqParamType](tcr mcp.CallToolRequest, param string, required bool) (T, error) {
	var val T
	arg, ok := tcr.GetArguments()[param]
	if ok {
		val, ok = arg.(T)
		if !ok {
			return val, fmt.Errorf("%s has wrong type: %T", param, arg)
		}
	} else if required {
		return val, fmt.Errorf("%s param is required", param)
	}
	return val, nil
}

func doRequest(req *http.Request) *mcp.CallToolResult {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to do request: %v", err))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to read response body: %v", err))
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("unexpected response status code %v: %s", resp.StatusCode, string(body)))
	}
	return mcp.NewToolResultText(string(body))
}
