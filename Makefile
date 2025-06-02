NAME:=mcp-alertmanager

tidy:
	rm -f go.sum; go mod tidy

vet:
	go vet ./...

fmt:
	gofmt -l -s -w ./
	goimports -l -w ./

install-linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.57.2

lint: install-linter
	./bin/golangci-lint run

test: tidy
	go test ./...

build: tidy
	go build -o ./bin/mcp-alertmanager ./cmd/mcp-alertmanager/main.go

run: build ## Run the MCP server in stdio mode.
	./bin/mcp-alertmanager

run-sse: build ## Run the MCP server in SSE mode.
	./bin/mcp-alertmanager --transport sse --log-level debug
