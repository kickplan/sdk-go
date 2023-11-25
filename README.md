# sdk-go

[![main](https://github.com/kickplan/sdk-go/actions/workflows/main.yml/badge.svg)](https://github.com/kickplan/sdk-go/actions/workflows/main.yml)

## Installation

```bash
go get github.com/kickplan/sdk-go
```

Set environment variables:

```bash
export KICKPLAN_ENDPOINT=https://api...
export KICKPLAN_ACCESS_TOKEN=...
```

## Usage

```go
package main

import (
    "context"
    "log"

    kickplan "github.com/kickplan/sdk-go"
    "github.com/kickplan/sdk-go/adapter"
    "github.com/kickplan/sdk-go/eval"
)

func main() {
    ctx := context.Background()
    client := kickplan.NewClient()

    b, err := client.GetBool(ctx, "my-flag", false, eval.Context{
        "account_id": "123",
    })
    if err != nil {
        if errors.Is(err, adapter.ErrFlagNotFound) {
            log.Printf("flag %q not found", flag)
            return
        }
        log.Fatalf("failed to get flag: %v", err)
    }

    log.Printf("my-flag: %v", b)
}
```

See [examples](examples) for more.
