# MCP Server for PipeCD

This project provides an MCP (Model Context Protocol) server for PipeCD, enabling integration and management of applications and deployments.

## Usage
Configure Claude or some MCP Clients with the environment variables below.

- PIPECD_HOST : host of the PipeCD control plane. for example, `demo.pipecd.dev:443`
- PIPECD_API_KEY_FILE : full path to the file which contains PipeCD API Key
- PIPECD_INSECURE : set this as `true` if you want to connect to control plane without ssl/tls

### Example Configuration
In this example, the built binary is placed at `/Users/sawada/go/bin/mcp-server-pipecd`.
This is the default install path when installing this server with `go install github.com/Warashi/mcp-server-pipecd@latest` with the home directory as `/Users/sawada`.

```json
{
  "mcpServers": {
    "pipecd": {
      "command": "/Users/sawada/go/bin/mcp-server-pipecd",
      "args": [],
      "env": {
        "PIPECD_HOST": "demo.pipecd.dev:443",
        "PIPECD_API_KEY_FILE": "/Users/sawada/.config/mcp-server-pipecd/api_key",
        "PIPECD_INSECURE": "false"
      }
    }
  }
}
```

## Example Screenshot
<img src="https://github.com/user-attachments/assets/1dd3da91-c9b0-460f-9c1d-b63f3ad584a3" alt="Claude Desktop Screenshot" />


## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

## Note
mcp-server-pipecd is indexed and certified by [MCP Review](https://mcpreview.com/mcp-servers/warashi/mcp-server-pipecd)
