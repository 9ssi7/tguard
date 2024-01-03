# tguard Package

The `tguard` package is designed to manage data with specific timeouts. It offers functionalities to start, cancel, and manage data based on their time-to-live (TTL) settings.

## Installation

```bash
go get github.com/9ssi7/tguard
```

## Documentation

[![Go Reference](https://pkg.go.dev/badge/github.com/9ssi7/tguard.svg)](https://pkg.go.dev/github.com/9ssi7/tguard)

## Usage

### Import the Package

```go
import "github.com/9ssi7/tguard"
```

### Define Your Data Structure

Define the structure of your data, which should include an identifier (`Id`) and any other relevant fields.

```go
type TestData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
```

### Implement the Identity Checker

Implement an identity checker function to verify the data's identity.

```go
identityChecker := func(id string, data TestData) bool {
	return id == data.Id
}
```

### Create Configuration

Create a configuration for the `tguard` service, specifying the fallback function, identity checker, default TTL, and interval.

```go
config := tguard.Config[TestData]{
	Fallback:        func(data TestData) {},
	IdentityChecker: identityChecker,
	DefaultTTL:      time.Minute * 5,
	Interval:        time.Second * 10,
}
```

### Create and Start the Service

Create a new `tguard` service instance using the configuration and start the service.

```go
g := tguard.New(config)
ctx := context.Background()
go g.Connect(ctx)
```

### Manage Data

Use the `Start` method to add data and the `Cancel` method to remove data if needed.

```go
data := TestData{
	Id:   "1",
	Name: "test",
}
_ = g.Start(ctx, data)
```

## Full Example

Here's a complete example demonstrating how to use the `tguard` package:

```go
// Import required packages
import (
	"context"
	"time"
	"github.com/9ssi7/tguard"
)

// Define the data structure
type TestData struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	// Implement the identity checker
	identityChecker := func(id string, data TestData) bool {
		return id == data.Id
	}

	// Create the configuration
	config := tguard.Config[TestData]{
		Fallback:        func(data TestData) {},
		IdentityChecker: identityChecker,
		DefaultTTL:      time.Minute * 5,
		Interval:        time.Second * 10,
	}

	// Create and start the service
	g := tguard.New(config)
	ctx := context.Background()
	go g.Connect(ctx)

	// Manage data
	data := TestData{
		Id:   "1",
		Name: "test",
	}
	_ = g.Start(ctx, data)
}
```

## Potential Use-Cases

### Real-World Example: Ticketing System

The `tguard` package can be used in a ticketing system to manage reservations with a specific time limit. Here's a Mermaid diagram illustrating the flow:

```mermaid
graph TD
    A[Start] --> B[Request Ticket]
    B --> C[Add to tguard]
    C --> D[Process Payment]
    D --> E[Confirm Reservation]
    E --> F[End]
    C -->|Timeout| G[Cancel Reservation]
    G --> F
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the Apache License. See [LICENSE](LICENSE) for more details.