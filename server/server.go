package server

import (
	"github.com/eyazici90/mcp-alertmanager/tools"
	"github.com/mark3labs/mcp-go/server"
)

func New(amURL string) *server.MCPServer {
	srv := server.NewMCPServer(
		"mcp-alertmanager",
		"0.0.1",
		server.WithRecovery(),
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithInstructions(`
You are Virtual Assistant, a tool for interacting with Alertmanager API for different tasks related to monitoring and observability.
When investigating an alert through list_alerts, check if there is a matched routing by using get_status tool.
Try not to second guess information - if you don't know something or lack information, it's better to ask.
`),
	)
	tools.RegisterToolStatus(srv, amURL)
	tools.RegisterToolAlerts(srv, amURL)
	tools.RegisterToolSilences(srv, amURL)

	return srv
}
