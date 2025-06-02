package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/eyazici90/mcp-alertmanager/tools"
	"github.com/mark3labs/mcp-go/server"
)

var (
	transport string
	addr      string
	basePath  string
	logLevel  string
	amURL     string
)

func main() {
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&addr, "sse-address", "localhost:8000", "The host and port to start the sse server on")
	flag.StringVar(&basePath, "base-path", "", "Base path for the sse server")
	flag.StringVar(&logLevel, "log-level", "info", "Log level (debug, info, warn, error)")
	flag.StringVar(&amURL, "alertmanager-url", "https://localhost:9093", "Alertmanager URL")
	flag.Parse()

	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: parseLevel(logLevel)})))

	s := server.NewMCPServer(
		"mcp-alertmanager",
		"0.0.1",
		server.WithRecovery(),
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)
	tools.RegisterToolStatus(s, amURL)
	tools.RegisterToolAlerts(s, amURL)
	tools.RegisterToolSilences(s, amURL)

	switch transport {
	case "stdio":
		slog.Info("Starting Alertmanager MCP server using stdio transport")
		if err := server.ServeStdio(s); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	case "sse":
		srv := server.NewSSEServer(s,
			server.WithStaticBasePath(basePath),
		)
		slog.Info("Starting Alertmanager MCP server using SSE transport", "address", addr, "basePath", basePath)
		if err := srv.Start(addr); err != nil {
			return err
		}
	default:
		return fmt.Errorf(
			"invalid transport type: %s. Must be 'stdio' or 'sse'",
			transport,
		)
	}
	return nil
}

func parseLevel(level string) slog.Level {
	var l slog.Level
	if err := l.UnmarshalText([]byte(level)); err != nil {
		return slog.LevelInfo
	}
	return l
}
