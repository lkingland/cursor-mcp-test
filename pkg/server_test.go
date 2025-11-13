package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMyTool(t *testing.T) {
	session := setupSession(t)
	ctx := context.Background()

	tests := []struct {
		name    string
		args    any
		want    string
		wantErr bool
	}{
		{
			name: "with_optional_string",
			args: map[string]any{
				"requiredParam": "foo",
				"optionalStr":   "bar",
			},
			want: "requiredParam=foo, optionalStr=bar",
		},
		{
			name: "without_optional_string",
			args: map[string]any{
				"requiredParam": "foo",
			},
			want: "requiredParam=foo",
		},
		{
			name: "with_explicit_null",
			args: json.RawMessage(`{"requiredParam": "foo", "optionalStr": null}`),
			want: "requiredParam=foo",
		},
		{
			name:    "missing_required_param",
			args:    map[string]any{"optionalStr": "bar"},
			wantErr: true,
		},
		{
			name:    "wrong_type_for_optional",
			args:    json.RawMessage(`{"requiredParam": "foo", "optionalStr": 123}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var argsJSON json.RawMessage
			switch v := tt.args.(type) {
			case json.RawMessage:
				argsJSON = v
			default:
				argsJSON, _ = json.Marshal(v)
			}

			result, err := session.CallTool(ctx, &mcp.CallToolParams{
				Name:      "mytool",
				Arguments: argsJSON,
			})

			if tt.wantErr {
				if err == nil && !result.IsError {
					t.Fatal("expected error, got success")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.IsError {
				t.Fatalf("tool returned error: %+v", result)
			}

			var output MyToolOutput
			if err := unmarshalOutput(result.StructuredContent, &output); err != nil {
				t.Fatalf("failed to unmarshal output: %v", err)
			}

			if output.Message != tt.want {
				t.Errorf("got %q, want %q", output.Message, tt.want)
			}
		})
	}
}

func TestListTools(t *testing.T) {
	session := setupSession(t)

	result, err := session.ListTools(context.Background(), &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	if len(result.Tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(result.Tools))
	}

	tool := result.Tools[0]
	if tool.Name != "mytool" {
		t.Errorf("expected tool name 'mytool', got %q", tool.Name)
	}

	schema, ok := tool.InputSchema.(map[string]any)
	if !ok {
		t.Fatal("schema is not a map")
	}

	// Verify union type for optionalStr
	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("schema has no properties")
	}

	optionalStr, ok := props["optionalStr"].(map[string]any)
	if !ok {
		t.Fatal("optionalStr property not found")
	}

	typeVal := optionalStr["type"]
	t.Logf("optionalStr type: %v", typeVal)

	// Verify it's a union type ["null", "string"]
	typeArr, ok := typeVal.([]any)
	if !ok || len(typeArr) != 2 {
		t.Errorf("expected union type array, got %T: %v", typeVal, typeVal)
	}
}

// Helper functions

func setupSession(t *testing.T) *mcp.ClientSession {
	t.Helper()
	server := New()
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	t1, t2 := mcp.NewInMemoryTransports()

	ctx := context.Background()
	if _, err := server.Connect(ctx, t1, nil); err != nil {
		t.Fatalf("server connect failed: %v", err)
	}

	session, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatalf("client connect failed: %v", err)
	}

	return session
}

func unmarshalOutput(structuredContent any, output *MyToolOutput) error {
	if rawMsg, ok := structuredContent.(json.RawMessage); ok {
		return json.Unmarshal(rawMsg, output)
	}
	bytes, _ := json.Marshal(structuredContent)
	return json.Unmarshal(bytes, output)
}
