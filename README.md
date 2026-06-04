# Distributed Cache

A multi-node in-memory cache system built in Go with LRU eviction, TTL expiration, and consistent hashing for efficient key distribution.

## Features

- **Consistent Hashing**: Efficiently distributes keys across multiple cache nodes
- **TTL Support**: Automatic expiration of cached items
- **LRU Eviction**: Least Recently Used eviction policy when cache is full
- **Docker Support**: Easy deployment with Docker Compose
- **HTTP API**: Simple REST API for cache operations
- **Load Balancing**: Nginx for distributing requests
- **Multiple Nodes**: Scalable multi-node architecture

## Architecture

The system consists of:
- Multiple cache server nodes
- Consistent hash ring for key distribution
- Nginx load balancer
- HTTP handlers for GET/SET/DELETE operations

## Getting Started

### Prerequisites

- Go 1.16 or higher (for local development)
- Docker and Docker Compose (for containerized deployment)

### Running with Docker Compose

1. Clone the repository:
```bash
git clone https://github.com/jknate/distributed-cache.git
cd distributed-cache
```

2. Start the cache cluster:
```bash
docker-compose up
```

This will start multiple cache nodes and an Nginx load balancer.

### Local Development

1. Build the project:
```bash
go build -o cache-server main.go
```

2. Run a cache node:
```bash
./cache-server
```

## API Usage

### Set a value
```bash
curl -X POST http://localhost/set \
  -H "Content-Type: application/json" \
  -d '{"key": "mykey", "value": "myvalue", "ttl": 3600}'
```

### Get a value
```bash
curl http://localhost/get?key=mykey
```

### Delete a value
```bash
curl -X DELETE http://localhost/delete?key=mykey
```

## Configuration

- Modify `docker-compose.yml` to adjust the number of cache nodes
- Update `nginx.conf` for load balancer settings
- Configure TTL and cache size in the cache initialization

## Project Structure

- `main.go` - Entry point and server setup
- `cache/` - Cache implementation with LRU and TTL
- `hashring/` - Consistent hashing implementation
- `handlers/` - HTTP request handlers
- `Dockerfile` - Container image definition
- `docker-compose.yml` - Multi-container orchestration
- `nginx.conf` - Load balancer configuration
