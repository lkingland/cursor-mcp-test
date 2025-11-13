package server

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type MyToolInput struct {
	RequiredParam string  `json:"requiredParam" jsonschema:"required,A required string parameter"`
	OptionalStr   *string `json:"optionalStr,omitempty" jsonschema:"An optional string parameter"`
}

type MyToolOutput struct {
	Message string `json:"message" jsonschema:"Output message"`
}

func New() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "cursor-mcp-test",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			HasTools:     true,
			HasResources: true,
		},
	)

	// A tool with optional parameters that generate union type schemas
	myTool := &mcp.Tool{
		Name:        "mytool",
		Title:       "Example Tool",
		Description: "A test tool with optional parameters.",
	}

	mcp.AddTool(server, myTool, myToolHandler)

	// A a simple resource (Cursor seems to require this)
	server.AddResource(
		&mcp.Resource{
			URI:         "cursor-mcp-test://readme",
			Name:        "README",
			Description: "Information about this test server",
			MIMEType:    "text/plain",
		},
		readmeHandler,
	)

	return server
}

func readmeHandler(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	content := `Cursor MCP Test Server
======================
This is a minimal MCP server for testing optional parameters`

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "cursor-mcp-test://readme",
				MIMEType: "text/plain",
				Text:     content,
			},
		},
	}, nil
}

func myToolHandler(ctx context.Context, req *mcp.CallToolRequest, input MyToolInput) (*mcp.CallToolResult, MyToolOutput, error) {
	msg := fmt.Sprintf("requiredParam=%s", input.RequiredParam)
	if input.OptionalStr != nil {
		msg += fmt.Sprintf(", optionalStr=%s", *input.OptionalStr)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: msg},
		},
	}, MyToolOutput{Message: msg}, nil
}
