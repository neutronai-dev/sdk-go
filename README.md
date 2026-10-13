# Neutron AI Go SDK

Go 1.25+ SDK for Neutron AI memory and MCP operations.

## Install

```sh
go get github.com/neutronai-dev/sdk-go
```

## Usage

```go
client, err := neutron.NewClient(neutron.Options{
    BaseURL: "https://neutronai.dev/api",
    Token: "your_token",
    NucleusID: "nucleus_demo",
})
if err != nil {
    return err
}

_, err = client.Remember(ctx, map[string]any{
    "scopeId": "user:derek",
    "type": "preference",
    "text": "Prefers concise implementation notes",
})
```

The SDK covers the core memory lifecycle, audit and cache operations, and MCP tool discovery. `TenantID` remains a compatibility alias and must match `NucleusID` when both are provided.
