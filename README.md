# Simple Load Balancer in Go

A lightweight, efficient HTTP load balancer implementation in Go that distributes incoming traffic across multiple backend servers using a round-robin algorithm.

## Features

- **Round-Robin Load Balancing**: Distributes requests evenly across all healthy backend servers
- **Health Checking**: Periodic health checks to detect and avoid unhealthy backends
- **Graceful Shutdown**: Proper cleanup of resources on termination
- **Thread-Safe Operations**: Uses atomic operations and mutexes for concurrent access
- **Health Status API**: Monitor backend server health via `/health` endpoint
- **Automatic Failover**: Automatically routes traffic away from failed backends
- **Error Handling**: Comprehensive error handling with proper logging
- **Configurable Timeouts**: Read, write, and idle timeouts for better resource management

## Architecture

The load balancer consists of several components:

- **Backend**: Represents a backend server with its URL, health status, and reverse proxy
- **ServerPool**: Manages a collection of backend servers and round-robin selection
- **Health Checker**: Periodically checks backend server availability
- **Load Balancer Handler**: Routes incoming requests to healthy backends

## Installation

### Prerequisites

- Go 1.25.3 or higher
- Backend servers running on ports specified in configuration

### Build

```bash
go build -o lb
```

## Usage

### Starting the Load Balancer

```bash
./lb
```

The load balancer will start on port `8080` by default.

### Configuration

Backend servers are configured in `main.go`:

```go
serverList := []string{
    "http://localhost:8081",
    "http://localhost:8082",
}
```

You can modify this list to add or remove backend servers.

### Testing Backend Servers

To test the load balancer, you need to start backend servers. Here's a simple example using Python:

**Backend Server 1 (Port 8081):**
```bash
python3 -m http.server 8081
```

**Backend Server 2 (Port 8082):**
```bash
python3 -m http.server 8082
```

Or use Go:

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Response from server on port 8081")
    })
    log.Fatal(http.ListenAndServe(":8081", nil))
}
```

### Making Requests

Once the load balancer and backend servers are running:

```bash
# Make a request through the load balancer
curl http://localhost:8080/

# Check backend health status
curl http://localhost:8080/health
```

### Health Check Response

The `/health` endpoint returns JSON with backend status:

```json
{
    "total": 2,
    "alive": 2,
    "dead": 0,
    "backends": [
        {
            "url": "http://localhost:8081",
            "alive": true
        },
        {
            "url": "http://localhost:8082",
            "alive": true
        }
    ]
}
```

## Configuration Options

| Constant | Default Value | Description |
|----------|--------------|-------------|
| `Port` | `:8080` | Load balancer listening port |
| `HealthCheckInterval` | `10s` | Interval between health checks |
| `ReadTimeout` | `15s` | Maximum duration for reading requests |
| `WriteTimeout` | `15s` | Maximum duration for writing responses |
| `IdleTimeout` | `60s` | Maximum idle time for connections |

## Code Structure

```
.
├── main.go           # Entry point and server setup
├── models/
│   └── model.go      # Data structures (Backend, ServerPool)
├── routes/
│   └── routes.go     # HTTP handlers
├── helper/
│   └── helper.go     # Utility functions (health checks, peer selection)
└── arch/
    └── document.txt  # Architecture documentation
```

## How It Works

1. **Request Handling**: When a request arrives at the load balancer, it uses the round-robin algorithm to select the next healthy backend server
2. **Health Monitoring**: A background goroutine continuously checks backend health every 10 seconds
3. **Failover**: If a backend becomes unhealthy, it's automatically excluded from the rotation
4. **Error Recovery**: When a backend recovers, health checks detect it and add it back to the pool

## Thread Safety

The implementation uses:
- **Atomic Operations**: For thread-safe counter increments
- **RWMutex**: For read-write locking on backend health status
- **Proper Locking**: Minimizes lock contention while ensuring data consistency

## Improvements Over Previous Version

- ✅ Fixed naming conventions (Go-style camelCase)
- ✅ Added thread-safe operations using atomic counters
- ✅ Improved mutex usage (RWMutex for better read performance)
- ✅ Added graceful shutdown handling
- ✅ Added health status API endpoint
- ✅ Added comprehensive logging
- ✅ Added proper error handling in reverse proxy
- ✅ Added configurable timeouts
- ✅ Improved health check with context cancellation
- ✅ Better code organization and documentation

## Future Enhancements

Potential improvements for future versions:

- [ ] Support for weighted round-robin algorithm
- [ ] Least connections load balancing algorithm
- [ ] Configuration file support (YAML/JSON)
- [ ] Metrics and Prometheus integration
- [ ] Request retry logic with backoff
- [ ] Circuit breaker pattern
- [ ] TLS/HTTPS support
- [ ] Rate limiting
- [ ] Request/response logging middleware
- [ ] Docker support
- [ ] Kubernetes deployment manifests

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the MIT License.

## Author

Mayur Ubarhande

## Acknowledgments

Built with Go's standard library `net/http` and `net/http/httputil` packages.
