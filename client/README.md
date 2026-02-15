# SWE Client

This package provides a Go client for the SWEService gRPC server.

## Installation

To use this client in another Go repository, you can import it:

```go
import "github.com/jenujari/go-swe-api/client"
```

## Usage

Here is a basic example of how to use the client:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/jenujari/go-swe-api/client"
)

func main() {
    // Initialize the client
    c, err := client.NewSWEClient("localhost:5678")
    if err != nil {
        log.Fatalf("could not create client: %v", err)
    }
    defer c.Close()

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    // Ping the server
    pingResp, err := c.Ping(ctx)
    if err != nil {
        log.Fatalf("ping failed: %v", err)
    }
    fmt.Printf("Server Status: %s, Version: %s\n", pingResp.Status, pingResp.Version)

    // Get position for Sun
    posResp, err := c.GetPos(ctx, "2026-01-26T00:00:00Z", "Sun")
    if err != nil {
        log.Fatalf("get pos failed: %v", err)
    }
    
    if sun, ok := posResp.Results["Sun"]; ok {
        fmt.Printf("Sun Longitude: %f\n", sun.Longitude)
    }
}
```

## API Documentation

### `NewSWEClient(address string) (*SWEClient, error)`
Creates a new client instance. The address should be the host and port of the gRPC server.

### `Close() error`
Closes the underlying gRPC connection.

### `Ping(ctx context.Context) (*pb.PingResponse, error)`
Checks if the server is alive and returns version information.

### `GetPos(ctx context.Context, timeStr string, planetName string) (*pb.PosResponse, error)`
Returns the position(s) of planets. If `planetName` is empty, it returns positions for all supported planets.

### `FindConjunction(ctx context.Context, start, end, planet1, planet2 string, orb int32, step float64) (*pb.ConjunctionResponse, error)`
Searches for conjunctions between two planets in a given time range.
