package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var mcpClient *client.Client

func TestMain(m *testing.M) {
	var (
		err  error
		code int
	)
	defer func() {
		if err != nil {
			slog.Error(err.Error())
			os.Exit(code)
		}
	}()
	ctx := context.Background()

	mcpClient, err = client.NewSSEMCPClient("http://localhost:8000/sse")
	if err != nil {
		return
	}
	err = mcpClient.Start(ctx)
	if err != nil {
		return
	}
	mcpClient.OnNotification(func(notification mcp.JSONRPCNotification) {
		fmt.Printf("Received notification: %s\n", notification.Method)
	})

	_, err = mcpClient.Initialize(ctx, mcp.InitializeRequest{})
	if err != nil {
		return
	}

	code = m.Run()
}

func TestMCPServer(t *testing.T) {
	ctx := context.Background()

	t.Run("tools", func(t *testing.T) {
		resp, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
		require.NoError(t, err)
		assert.Len(t, resp.Tools, 3)
	})
	t.Run("get_status", func(t *testing.T) {
		var req mcp.CallToolRequest
		req.Params.Name = "get_status"
		resp, err := mcpClient.CallTool(ctx, req)
		assert.NoError(t, err)

		printToolResult(resp)
	})
	t.Run("list_alerts", func(t *testing.T) {
		var req mcp.CallToolRequest
		req.Params.Name = "list_alerts"
		req.Params.Arguments = map[string]any{
			"active":    "true",
			"silenced":  "true",
			"inhibited": "false",
		}
		resp, err := mcpClient.CallTool(ctx, req)
		assert.NoError(t, err)

		printToolResult(resp)
	})
	t.Run("list_silences", func(t *testing.T) {
		var req mcp.CallToolRequest
		req.Params.Name = "list_alerts"
		resp, err := mcpClient.CallTool(ctx, req)
		assert.NoError(t, err)

		printToolResult(resp)
	})
}

// Helper function to print tool results
func printToolResult(result *mcp.CallToolResult) {
	for _, content := range result.Content {
		if textContent, ok := content.(mcp.TextContent); ok {
			fmt.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	}
}
