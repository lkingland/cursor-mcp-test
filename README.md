# Cursor MCP Test Server

This is a minimal MCP server designed to reproduce the issue where Cursor Desktop incorrectly rejects valid values for optional parameters that use JSON Schema union types (e.g., `["null", "string"]`).

## Reproducing

1. build the MCP server with `go build`
2. add `cursor-mcp-test` to registered servers `~/.cursor/mcp.go`:
```json
{
  "mcpServers": {
    "cursor-mcp-test": {
      "command": "/absolute/path/to/cursor-mcp-test"
    }
  }
}
```
3. Restart Cursor and ask it to invoke the server's "Example Tool"
