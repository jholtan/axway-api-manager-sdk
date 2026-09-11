# Axway API Manager SDK

A generated Go client for the Axway API Manager REST API v1.4. The client is generated from
`specs/api-manager-V_1_4-oas3.json` and exposes every operation and schema in that specification.

## Installation
```bash
go get github.com/jholtan/axway-api-manager-sdk
```

## Quick Start
```go
package main

import (
    "context"
    "encoding/json"
    "log"

    "github.com/jholtan/axway-api-manager-sdk/apimanager/openapi"
)

func main() {
    client, err := openapi.NewBasicAuthClient(
        "https://localhost:8075/api/portal/v1.4",
        "username",
        "password",
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }

    response, err := client.GetAllApisLight1(context.Background(), nil)
    if err != nil {
        log.Fatal(err)
    }
    defer response.Body.Close()

    var apis []openapi.APIDefinitionLight
    if err := json.NewDecoder(response.Body).Decode(&apis); err != nil {
        log.Fatal(err)
    }
}
```
