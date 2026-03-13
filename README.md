# User Manual

## Project Structure
The project follows a microservices architecture with a CLI-driven entry point.

- `cmd/`: **CLI Layer**. Contains command definitions and configuration logic.
  - `config/`: Reusable configuration components (Viper-based).
  - `kitchen/`: CLI commands for the Kitchen service.
  - `orders/`: CLI commands for the Orders service.
- `services/`: **Business Logic**.
  - `kitchen/`: Implementation of the Kitchen microservice.
  - `orders/`: Implementation of the Orders microservice.
  - `common/`: Shared code, gRPC generated code, and types.
- `internal/`: **Shared Core**. Logging, server abstractions, database utilities, and networking.
- `protobuf/`: Protocol Buffer definitions.
- `migrations/`: Database migration files (managed by `goose`).

## Configuration
The application uses **Viper** for flexible configuration. Settings are resolved in the following priority:
1. **CLI Flags** (e.g., `--db-host`)
2. **Environment Variables** (e.g., `DATABASE_HOST`)
3. **YAML Config Files** (e.g., `root_config.yaml`)
4. **Defaults** hardcoded in the source code.

### Sample Configuration Files
Sample files are provided for each component:
- `root_config-sample.yaml`: Global CLI settings (debug, silent, colors).
- `services/kitchen/config-sample.yaml`: Kitchen service settings (port, test variables).
- `services/orders/config-sample.yaml`: Orders service settings (database, gRPC/HTTP ports).

To use them, copy the sample file to a file named `config.yaml` in the respective directory.

### Environment Variable Overrides
All configuration keys can be overridden via environment variables. Nested keys use underscores as separators:
- `DATABASE_HOST` overrides `database.host`
- `ORDERS_HTTP_PORT` overrides `orders.http.port`
- `KITCHEN_PORT` overrides `kitchen.port`

## Command Line Interface (CLI)
The project is driven by a Cobra-based CLI.

### Global Flags
To get help on any command:
```bash
go run main.go [command] --help
```

### Main Commands
- **Run All Services**: `go run main.go run` (orchestrates all microservices).
- **Service Specific**: `go run main.go [service] run` (e.g., `orders`, `kitchen`).
- **Migrations**: `go run main.go [service] migrate [ARGS]` (e.g., `go run main.go orders migrate up`).

## Project Management (Makefile)
Common tasks are simplified via `Makefile`:
- `make run-all`: Run the entire application.
- `make migrate-all`: Run migrations for all services.
- `make run-orders` / `make run-kitchen`: Run specific services.
- `make gen-go`: Generate Go code from Protobuf.

## Docker Usage
The project is fully containerized.

### Running with Docker Compose
To start the entire stack (including Postgres):
```bash
docker compose up --build
```

### Exposed Ports
By default, the following ports are mapped to the host:
- `8000`: Kitchen HTTP API
- `3082`: Orders HTTP API
- `5440`: PostgreSQL Database

## Maintenance
- **Goroutine Management**: The project uses specialized utilities in `internal/util` to prevent goroutine leaks and ensure graceful shutdown.
- **Logging**: Uses an internal logging system with icon support and debug levels.
