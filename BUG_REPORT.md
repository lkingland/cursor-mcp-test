# Bug Report: Error Invoking MCP Tools with Optional Parameters

## Summary


Cursor Desktop's MCP Client yields an error when attempting to invoke an MCP
server tool with optional paramters.  Example:

```
Parameter 'optionalStr' must be of type null,string, got string
```

## Reproduction

1. Clone Test MCP Server: `git clone https://github.com/lkingland/cursor-mcp-test`
2. Build: `cd cursor-mcp-test && go build`
3. Configure Cursor MCP settings (`~/.cursor/mcp.json`):
   ```json
   {
     "mcpServers": {
       "cursor-mcp-test": {
         "command": "/absolute/path/to/cursor-mcp-test"
       }
     }
   }
   ```
4. Restart Cursor
5. Ask Cursor to invoke "Example Tool" with:
   ```json
   {
     "requiredParam": "test-value",
     "optionalStr": "optional-value"
   }
   ```

## Analysis

A log of communication with the server is written to `/tmp/cursor-mcp-test.log`
which shows the tool invocation never exits the client.  This is therefore
happening client-side in the MCP Client and is likely a problem of not
supporting union set validation of the tool's schema.

The test server uses the following to define an input:

```go
type MyToolInput struct {
    RequiredParam string  `json:"requiredParam" jsonschema:"required,A required string parameter"`
    OptionalStr   *string `json:"optionalStr,omitempty" jsonschema:"An optional string parameter"`
}
```

Which generates the tool schema:
```json
{
  "type": "object",
  "required": ["requiredParam"],
  "properties": {
    "requiredParam": {
      "type": "string",
      "description": "A required string parameter"
    },
    "optionalStr": {
      "type": ["null", "string"],
      "description": "An optional string parameter"
    }
  },
  "additionalProperties": false
}
```

The pattern `"type": ["null", "string"]` is the standard way to represent optional/nullable fields:
- Go MCP servers: pointer types (`*string`, `*int`, `*bool`) generate union types
- Python MCP servers: `Optional[str]` generates union types
- TypeScript MCP servers: `string | null` generates union types

## Specification References

According to the [MCP specification (2025-06-18)](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/main/schema/2025-06-18/schema.json):

> A Tool's `inputSchema` is "A JSON Schema object defining the expected parameters for the tool"

The MCP specification requires `inputSchema` to be a valid JSON Schema object. **JSON Schema explicitly supports union types** via array notation for the `type` field (e.g., `"type": ["null", "string"]`), which is the standard pattern for optional/nullable parameters.

Per [JSON Schema specification](https://json-schema.org/understanding-json-schema/reference/type), a value validates successfully if it matches **any** of the types in a union array.

## Verification

The same parameters work correctly with other MCP clients. Unit tests using the official Go MCP SDK all pass:

```bash
$ go test -v ./pkg/server
=== RUN   TestMyTool/with_optional_string
--- PASS: TestMyTool/with_optional_string (0.00s)
=== RUN   TestMyTool/without_optional_string
--- PASS: TestMyTool/without_optional_string (0.00s)
=== RUN   TestMyTool/with_explicit_null
--- PASS: TestMyTool/with_explicit_null (0.00s)
=== RUN   TestMyTool/missing_required_param
--- PASS: TestMyTool/missing_required_param (0.00s)
=== RUN   TestMyTool/wrong_type_for_optional
--- PASS: TestMyTool/wrong_type_for_optional (0.00s)
```

## References

- [MCP Specification 2025-06-18](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/main/schema/2025-06-18/schema.json)
- [JSON Schema Specification - Type](https://json-schema.org/understanding-json-schema/reference/type)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Test Case Repository](https://github.com/lkingland/cursor-mcp-test)

## Environment

- **Cursor Version**: 2.0.75
- **VSCode Version**: 1.99.3
- **Commit**: 9e7a27b76730ca7fe4aecaeafc58bac1e2c82120
- **Date**: 2025-11-12T17:34:21.472Z
- **Platform**: Darwin arm64 25.0.0
- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk v1.1.0`

