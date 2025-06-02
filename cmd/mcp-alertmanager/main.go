package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/eyazici90/mcp-alertmanager/tools"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(
		&transport,
		"transport",
		"stdio",
		"Transport type (stdio or sse)",
	)
	addr := flag.String("sse-address", "localhost:8000", "The host and port to start the sse server on")
	basePath := flag.String("base-path", "", "Base path for the sse server")
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	flag.Parse()

	if err := run(transport, *addr, *basePath, parseLevel(*logLevel)); err != nil {
		panic(err)
	}
}

func run(transport, addr, basePath string, logLevel slog.Level) error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})))

	s := server.NewMCPServer(
		"mcp-alertmanager",
		"0.0.1",
		server.WithRecovery(),
		server.WithLogging(),
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
	)

	tools.RegisterToolAlerts(s)

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
