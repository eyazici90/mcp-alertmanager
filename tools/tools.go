package tools

import (
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

type contextKey string

type toolReqParamType interface {
	string | float64 | bool | []string | []any
}

func getToolReqParam[T toolReqParamType](ctr mcp.CallToolRequest, param string, required bool) (T, error) {
	var t T
	arg, ok := ctr.GetArguments()[param]
	if ok {
		t, ok = arg.(T)
		if !ok {
			return t, fmt.Errorf("%s has wrong type: %T", param, arg)
		}
	} else if required {
		return t, fmt.Errorf("%s param is required", param)
	}
	return t, nil
}

func doRequest(req *http.Request) *mcp.CallToolResult {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("do request: %v", err))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("read response body: %v", err))
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return mcp.NewToolResultError(fmt.Sprintf("unexpected response status code %v: %s", resp.StatusCode, string(body)))
	}
	return mcp.NewToolResultText(string(body))
}

func ptr[T any](v T) *T {
	return &v
}
