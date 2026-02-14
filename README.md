# Simple Load Balancer in Go

A lightweight, efficient HTTP load balancer implementation in Go that distributes incoming traffic across multiple backend servers using a round-robin algorithm.

## Features

- **Round-Robin Load Balancing**: Distributes requests evenly across all healthy backend servers
- **Health Checking**: Periodic health checks to detect and avoid unhealthy backends
- **Graceful Shutdown**: Proper cleanup of resources on termination
- **Thread-Safe Operations**: Uses atomic operations and mutexes for concurrent access
- **Health Status API**: Monitor backend server health via `/health` endpoint
- **Metrics API**: Track request statistics via `/metrics` endpoint
- **Request Logging**: Comprehensive logging of all requests with duration tracking
- **Automatic Failover**: Automatically routes traffic away from failed backends
- **Error Handling**: Comprehensive error handling with proper logging
- **Configurable Timeouts**: Read, write, and idle timeouts for better resource management
- **JSON Configuration**: Easy configuration via JSON file
- **Docker Support**: Ready-to-use Docker and docker-compose setup
- **Makefile**: Convenient commands for building, testing, and running

## Architecture

The load balancer consists of several components:

- **Backend**: Represents a backend server with its URL, health status, and reverse proxy
- **ServerPool**: Manages a collection of backend servers and round-robin selection
- **Health Checker**: Periodically checks backend server availability
- **Load Balancer Handler**: Routes incoming requests to healthy backends
- **Middleware**: Logging and metrics collection for monitoring
- **Configuration**: JSON-based configuration management

## Installation

### Prerequisites

- Go 1.25.3 or higher
- Backend servers running on ports specified in configuration
- Make (optional, for using Makefile commands)
- Docker (optional, for containerized deployment)

### Quick Start with Make

```bash
# Build the project
make build

# Run tests
make test

# Run the load balancer
make run

# See all available commands
make help
```

### Manual Build

```bash
go build -o lb
```

## Usage

### Starting the Load Balancer

```bash
# Using default config (config.json)
./lb

# Using custom config file
./lb -config myconfig.json
```

The load balancer will start on port `8080` by default.

### Configuration

Create a `config.json` file (see `config.example.json` for reference):

```json
{
  "port": ":8080",
  "health_check_interval": 10000000000,
  "read_timeout": 15000000000,
  "write_timeout": 15000000000,
  "idle_timeout": 60000000000,
  "backends": [
    "http://localhost:8081",
    "http://localhost:8082"
  ]
}
```

Note: Time values are in nanoseconds (10000000000 = 10 seconds)

```go
serverList := []string{
    "http://localhost:8081",
    "http://localhost:8082",
}
```

You can modify this list to add or remove backend servers.

### Testing Backend Servers

#### Option 1: Using the Provided Example

Start the example backend servers in separate terminals:

```bash
# Terminal 1
go run examples/backend.go -port 8081

# Terminal 2
go run examples/backend.go -port 8082
```

Or use Make:

```bash
# Terminal 1
make backend1

# Terminal 2  
make backend2
```

#### Option 2: Using Python

```bash
# Terminal 1
python3 -m http.server 8081

# Terminal 2
python3 -m http.server 8082
```

#### Option 3: Using Docker Compose

```bash
docker-compose up
```

This starts the load balancer and two backend servers automatically.

### Making Requests

Once the load balancer and backend servers are running:

```bash
# Make a request through the load balancer
curl http://localhost:8080/

# Check backend health status
curl http://localhost:8080/health

# Check request metrics
curl http://localhost:8080/metrics
```

### Automated Testing

Use the provided test script:

```bash
./test_lb.sh
```

This script will:
- Verify the load balancer is running
- Test all endpoints
- Send multiple requests to verify load balancing
- Show final metrics

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

### Metrics Response

The `/metrics` endpoint returns request statistics:

```json
{
    "total_requests": 150,
    "success_requests": 148,
    "failed_requests": 2
}
```

## Configuration Options

Configure the load balancer using a JSON file. All durations are in nanoseconds.

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `port` | string | `:8080` | Load balancer listening port |
| `health_check_interval` | int64 | `10000000000` | Interval between health checks (10s) |
| `read_timeout` | int64 | `15000000000` | Maximum duration for reading requests (15s) |
| `write_timeout` | int64 | `15000000000` | Maximum duration for writing responses (15s) |
| `idle_timeout` | int64 | `60000000000` | Maximum idle time for connections (60s) |
| `backends` | []string | See example | List of backend server URLs |

### Time Conversion Reference

- 1 second = 1,000,000,000 nanoseconds
- 10 seconds = 10,000,000,000 nanoseconds
- 1 minute = 60,000,000,000 nanoseconds

## API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/` | ANY | Proxy requests to backend servers |
| `/health` | GET | Get health status of all backends |
| `/metrics` | GET | Get request statistics |

## Code Structure

```
.
├── main.go              # Entry point and server setup
├── config/
│   └── config.go        # Configuration management
├── models/
│   ├── model.go         # Data structures (Backend, ServerPool)
│   └── model_test.go    # Model tests
├── routes/
│   └── routes.go        # HTTP handlers
├── helper/
│   ├── helper.go        # Utility functions (health checks, peer selection)
│   └── helper_test.go   # Helper tests
├── middleware/
│   └── middleware.go    # HTTP middleware (logging, metrics)
├── examples/
│   └── backend.go       # Example backend server
├── arch/
│   └── document.txt     # Architecture documentation
├── Dockerfile           # Docker build configuration
├── docker-compose.yml   # Multi-container Docker setup
├── Makefile            # Build and test automation
└── test_lb.sh          # Automated testing script
```

## How It Works

1. **Request Handling**: When a request arrives at the load balancer, it uses the round-robin algorithm to select the next healthy backend server
2. **Health Monitoring**: A background goroutine continuously checks backend health at configurable intervals
3. **Failover**: If a backend becomes unhealthy, it's automatically excluded from the rotation
4. **Error Recovery**: When a backend recovers, health checks detect it and add it back to the pool
5. **Metrics Collection**: Middleware tracks all requests and their outcomes
6. **Request Logging**: Each request is logged with method, path, status code, and duration

## Thread Safety

The implementation uses:
- **Atomic Operations**: For thread-safe counter increments (request counting, round-robin index)
- **RWMutex**: For read-write locking on backend health status (allows concurrent reads)
- **Proper Locking**: Minimizes lock contention while ensuring data consistency

## Testing

### Run All Tests

```bash
make test
# or
go test ./...
```

### Run Tests with Coverage

```bash
make test-coverage
```

This generates an HTML coverage report.

### Run Benchmarks

```bash
make bench
# or
go test -bench=. -benchmem ./...
```

## Docker Usage

### Build Docker Image

```bash
docker build -t simple-lb .
```

### Run with Docker

```bash
docker run -p 8080:8080 simple-lb
```

### Run with Docker Compose

```bash
# Start all services
docker-compose up

# Run in background
docker-compose up -d

# Stop all services
docker-compose down
```

## Development

### Prerequisites

```bash
# Install dependencies
make install

# Format code
make fmt

# Run linters
make lint
```

### Project Commands

```bash
# Show all available commands
make help

# Build the project
make build

# Run tests
make test

# Clean build artifacts
make clean

# Run everything (clean, lint, test, build)
make all
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development guidelines.

## Improvements Over Previous Version

- ✅ Fixed naming conventions (Go-style camelCase)
- ✅ Added thread-safe operations using atomic counters
- ✅ Improved mutex usage (RWMutex for better read performance)
- ✅ Added graceful shutdown handling with context cancellation
- ✅ Added health status API endpoint
- ✅ Added metrics API endpoint for monitoring
- ✅ Added comprehensive request logging with duration tracking
- ✅ Added proper error handling in reverse proxy with auto-marking down
- ✅ Added configurable timeouts for better resource management
- ✅ Improved health check with context cancellation
- ✅ Added JSON configuration file support
- ✅ Added comprehensive unit tests
- ✅ Added Docker and docker-compose support
- ✅ Added Makefile for build automation
- ✅ Added example backend server
- ✅ Added automated testing script
- ✅ Better code organization with middleware pattern
- ✅ Comprehensive documentation (README, CONTRIBUTING)

## Future Enhancements

Potential improvements for future versions:

- [ ] Support for weighted round-robin algorithm
- [ ] Least connections load balancing algorithm
- [ ] Request retry logic with exponential backoff
- [ ] Circuit breaker pattern for fault tolerance
- [ ] TLS/HTTPS support for secure connections
- [ ] Rate limiting per client/IP
- [ ] Request/response body logging (optional)
- [ ] Prometheus metrics export
- [ ] WebSocket support
- [ ] Admin dashboard UI
- [ ] Hot reload of configuration
- [ ] Session persistence (sticky sessions)
- [ ] Request routing based on path patterns
- [ ] Multiple load balancing strategies per route

## Troubleshooting

### Load Balancer Won't Start

- Check if port 8080 is already in use: `lsof -i :8080` (Linux/Mac) or `netstat -ano | findstr :8080` (Windows)
- Verify configuration file is valid JSON
- Check if backends are reachable

### All Backends Showing as Down

- Verify backend servers are actually running
- Check backend URLs in configuration
- Ensure firewall isn't blocking connections
- Check health check interval isn't too aggressive

### High Memory Usage

- Adjust `IdleTimeout` to close idle connections sooner
- Check for connection leaks in backends
- Monitor metrics endpoint for unusual patterns

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Quick Contribution Steps

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is open source and available under the MIT License.

## Author

Mayur Ubarhande

## Acknowledgments

- Built with Go's standard library `net/http` and `net/http/httputil` packages
- Inspired by production load balancers like HAProxy and NGINX
- Community contributions and feedback
