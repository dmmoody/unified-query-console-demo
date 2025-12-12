# ACH Concourse

A microservices-based ACH (Automated Clearing House) processing system POC built with Go, PostgreSQL, and Docker.

## Architecture

This system consists of 5 independent microservices:

- **ODFI** (Originating Depository Financial Institution) - Manages outgoing ACH entries
- **RDFI** (Receiving Depository Financial Institution) - Manages incoming ACH entries
- **Ledger** - Simple posting and balance tracking system
- **EIP** (Exception/Investigation Platform) - Case management for ACH exceptions
- **Console** - **Unified API Gateway** that provides ALL CRUD operations for ALL services through a single endpoint

Each service has its own PostgreSQL database, demonstrating true microservice independence.

### 🎯 **API Gateway Pattern**

The Console service acts as a complete API Gateway:
- **All operations** for all services available at `http://localhost:8080`
- Single entry point for clients (Postman, UI, etc.)
- Service discovery handled by gateway
- Consistent error handling and responses
- Future-proof for authentication, rate limiting, etc.

See [docs/GATEWAY.md](docs/GATEWAY.md) for complete API Gateway documentation.

## Tech Stack

- **Language**: Go 1.22+
- **Database**: PostgreSQL (latest)
- **HTTP Router**: Chi v5
- **Database Driver**: pgx v5
- **Containerization**: Docker + Docker Compose

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Go 1.22+ (for local development)

### Running the System

```bash
# Build all services
docker-compose build

# Start all services
docker-compose up

# Start in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### Service Ports

| Service | Host Port | Description |
|---------|-----------|-------------|
| Console | 8080 | Unified API gateway |
| ODFI | 8081 | Origination service |
| RDFI | 8082 | Receiving service |
| Ledger | 8083 | Ledger/posting service |
| EIP | 8084 | Exception/Investigation Platform |

### Seeding Test Data

After starting the services, seed the databases with realistic test data:

```bash
# Fast Go-based seeding (recommended - seeds 620 records in seconds)
make seed

# Alternative bash script version
./seed.sh
```

This will create:
- **150 ODFI entries** (origination) with various statuses
- **150 RDFI entries** (receiving) including some returns
- **200 Ledger postings** (debits and credits)
- **120 EIP cases** (exception tracking) with different statuses

Total: **620 realistic ACH records** for demo purposes!

### Database Ports

| Database | Host Port |
|----------|-----------|
| ODFI DB | 5433 |
| RDFI DB | 5434 |
| Ledger DB | 5435 |
| EIP DB | 5436 |

## API Documentation

### 🚀 **Recommended: Use the Unified Gateway**

**All operations for all services are available through the Console Gateway at port 8080:**

```bash
# ODFI operations via gateway
curl http://localhost:8080/api/v1/odfi/entries
curl -X POST http://localhost:8080/api/v1/odfi/entries -H "Content-Type: application/json" -d '{...}'

# RDFI operations via gateway
curl http://localhost:8080/api/v1/rdfi/entries
curl -X POST http://localhost:8080/api/v1/rdfi/entries -H "Content-Type: application/json" -d '{...}'

# Ledger operations via gateway
curl http://localhost:8080/api/v1/ledger/postings
curl http://localhost:8080/api/v1/ledger/balances

# EIP operations via gateway
curl http://localhost:8080/api/v1/eip/cases
```

**📖 See [docs/GATEWAY.md](docs/GATEWAY.md) for complete API Gateway documentation with all endpoints and examples.**

### Direct Service Access (for debugging)

Services can also be accessed directly at their individual ports:

### Console Service (Unified Gateway) - Port 8080

The console service provides a unified interface to query and operate on ACH entries across services.

#### Get All ACH Items

```bash
GET http://localhost:8080/api/v1/ach-items
```

**Query Parameters:**
- `side` (optional): Filter by "ODFI" or "RDFI"
- `status` (optional): Filter by status
- `trace_number` (optional): Filter by trace number

**Example:**
```bash
curl "http://localhost:8080/api/v1/ach-items?side=ODFI&status=PENDING"
```

**Response:**
```json
[
  {
    "side": "ODFI",
    "source": "odfi",
    "entry_id": "uuid",
    "trace_number": "123456789",
    "amount_cents": 10000,
    "status": "PENDING",
    "extra": {
      "company_name": "ACME Corp",
      "sec_code": "PPD"
    }
  }
]
```

#### Get Single ACH Item

```bash
GET http://localhost:8080/api/v1/ach-items/{side}/{id}
```

**Example:**
```bash
curl "http://localhost:8080/api/v1/ach-items/ODFI/uuid-here"
```

#### Return an RDFI Entry

```bash
POST http://localhost:8080/api/v1/ach-items/RDFI/{id}/return
```

**Request Body:**
```json
{
  "reason": "R01"
}
```

**Example:**
```bash
curl -X POST "http://localhost:8080/api/v1/ach-items/RDFI/uuid-here/return" \
  -H "Content-Type: application/json" \
  -d '{"reason": "R01"}'
```

---

### ODFI Service - Port 8081

Manages originating (outgoing) ACH entries.

#### Create Entry

```bash
POST http://localhost:8081/api/v1/entries
```

**Request Body:**
```json
{
  "trace_number": "123456789",
  "company_name": "ACME Corp",
  "sec_code": "PPD",
  "amount_cents": 10000
}
```

**Example:**
```bash
curl -X POST "http://localhost:8081/api/v1/entries" \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "123456789",
    "company_name": "ACME Corp",
    "sec_code": "PPD",
    "amount_cents": 10000
  }'
```

#### List Entries

```bash
GET http://localhost:8081/api/v1/entries?status=PENDING&trace_number=123456789
```

#### Get Single Entry

```bash
GET http://localhost:8081/api/v1/entries/{id}
```

#### Update Entry Status

```bash
PATCH http://localhost:8081/api/v1/entries/{id}/status
```

**Request Body:**
```json
{
  "status": "SENT"
}
```

Valid statuses: `PENDING`, `SENT`, `CANCELLED`

#### Health Check

```bash
GET http://localhost:8081/healthz
```

---

### RDFI Service - Port 8082

Manages receiving (incoming) ACH entries.

#### Create Entry

```bash
POST http://localhost:8082/api/v1/entries
```

**Request Body:**
```json
{
  "trace_number": "987654321",
  "receiver_name": "John Doe",
  "amount_cents": 5000
}
```

**Example:**
```bash
curl -X POST "http://localhost:8082/api/v1/entries" \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "987654321",
    "receiver_name": "John Doe",
    "amount_cents": 5000
  }'
```

#### List Entries

```bash
GET http://localhost:8082/api/v1/entries?status=RECEIVED&trace_number=987654321
```

#### Get Single Entry

```bash
GET http://localhost:8082/api/v1/entries/{id}
```

#### Return Entry

```bash
POST http://localhost:8082/api/v1/entries/{id}/return
```

**Request Body:**
```json
{
  "reason": "R01"
}
```

Valid statuses: `RECEIVED`, `POSTED`, `RETURNED`

#### Health Check

```bash
GET http://localhost:8082/healthz
```

---

### Ledger Service - Port 8083

Simple ledger for tracking debits and credits.

#### Create Posting

```bash
POST http://localhost:8083/api/v1/postings
```

**Request Body:**
```json
{
  "ach_side": "ODFI",
  "trace_number": "123456789",
  "amount_cents": 10000,
  "direction": "DEBIT",
  "description": "Payment to vendor"
}
```

**Example:**
```bash
curl -X POST "http://localhost:8083/api/v1/postings" \
  -H "Content-Type: application/json" \
  -d '{
    "ach_side": "ODFI",
    "trace_number": "123456789",
    "amount_cents": 10000,
    "direction": "DEBIT",
    "description": "Payment to vendor"
  }'
```

Valid values:
- `ach_side`: `ODFI`, `RDFI`
- `direction`: `DEBIT`, `CREDIT`

#### List Postings

```bash
GET http://localhost:8083/api/v1/postings?ach_side=ODFI&trace_number=123456789
```

#### Get Balances

```bash
GET http://localhost:8083/api/v1/balances
```

**Response:**
```json
{
  "total_debits": 100000,
  "total_credits": 50000,
  "net_balance": -50000
}
```

#### Health Check

```bash
GET http://localhost:8083/healthz
```

---

### EIP Service - Port 8084

Exception and investigation case management.

#### Create Case

```bash
POST http://localhost:8084/api/v1/cases
```

**Request Body:**
```json
{
  "side": "RDFI",
  "trace_number": "987654321",
  "type": "RETURN_REVIEW",
  "notes": "Customer called to dispute charge"
}
```

**Example:**
```bash
curl -X POST "http://localhost:8084/api/v1/cases" \
  -H "Content-Type: application/json" \
  -d '{
    "side": "RDFI",
    "trace_number": "987654321",
    "type": "RETURN_REVIEW",
    "notes": "Customer called to dispute charge"
  }'
```

Valid values:
- `side`: `ODFI`, `RDFI`
- `type`: `RETURN_REVIEW`, `NOC_REVIEW`, `CUSTOMER_DISPUTE`

#### List Cases

```bash
GET http://localhost:8084/api/v1/cases?status=OPEN&side=RDFI&trace_number=987654321
```

#### Get Single Case

```bash
GET http://localhost:8084/api/v1/cases/{id}
```

#### Update Case Status

```bash
PATCH http://localhost:8084/api/v1/cases/{id}/status
```

**Request Body:**
```json
{
  "status": "RESOLVED"
}
```

Valid statuses: `OPEN`, `IN_PROGRESS`, `RESOLVED`

#### Health Check

```bash
GET http://localhost:8084/healthz
```

---

## Development

### Local Development (without Docker)

1. Install dependencies:
```bash
go mod download
```

2. Start PostgreSQL databases (adjust connection strings as needed)

3. Run a service:
```bash
# Set environment variables
export PORT=8081
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=odfi_user
export DB_PASSWORD=odfi_pass
export DB_NAME=odfi_db
export DB_SSLMODE=disable

# Run ODFI service
go run cmd/odfi/main.go
```

### Project Structure

```
ach-concourse/
├── cmd/                    # Service entry points
│   ├── odfi/
│   ├── rdfi/
│   ├── ledger/
│   ├── eip/
│   └── console/
├── internal/               # Internal packages
│   ├── common/            # Shared utilities
│   │   ├── db/           # Database connection helper
│   │   └── http/         # HTTP response helpers
│   ├── odfi/             # ODFI service logic
│   ├── rdfi/             # RDFI service logic
│   ├── ledger/           # Ledger service logic
│   ├── eip/              # EIP service logic
│   └── console/          # Console service logic
├── Dockerfile.odfi
├── Dockerfile.rdfi
├── Dockerfile.ledger
├── Dockerfile.eip
├── Dockerfile.console
├── docker-compose.yml
├── go.mod
└── README.md
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/odfi/...
```

## Example Workflow

Here's a complete workflow demonstrating the system:

```bash
# 1. Create an ODFI (outgoing) entry
curl -X POST "http://localhost:8081/api/v1/entries" \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "1234567890",
    "company_name": "ACME Corp",
    "sec_code": "PPD",
    "amount_cents": 50000
  }'
# Save the returned ID

# 2. Create an RDFI (incoming) entry
curl -X POST "http://localhost:8082/api/v1/entries" \
  -H "Content-Type: application/json" \
  -d '{
    "trace_number": "0987654321",
    "receiver_name": "Jane Smith",
    "amount_cents": 25000
  }'
# Save the returned ID

# 3. Query all entries via Console
curl "http://localhost:8080/api/v1/ach-items"

# 4. Create a ledger posting for the ODFI entry
curl -X POST "http://localhost:8083/api/v1/postings" \
  -H "Content-Type: application/json" \
  -d '{
    "ach_side": "ODFI",
    "trace_number": "1234567890",
    "amount_cents": 50000,
    "direction": "DEBIT",
    "description": "Outgoing payment"
  }'

# 5. Check balances
curl "http://localhost:8083/api/v1/balances"

# 6. Create an exception case for RDFI entry
curl -X POST "http://localhost:8084/api/v1/cases" \
  -H "Content-Type: application/json" \
  -d '{
    "side": "RDFI",
    "trace_number": "0987654321",
    "type": "CUSTOMER_DISPUTE",
    "notes": "Customer disputes transaction"
  }'

# 7. Return the RDFI entry via Console
curl -X POST "http://localhost:8080/api/v1/ach-items/RDFI/{rdfi-entry-id}/return" \
  -H "Content-Type: application/json" \
  -d '{"reason": "R10"}'
```

## Troubleshooting

### Services won't start
- Check if ports are already in use: `lsof -i :8080`
- Verify Docker is running: `docker ps`
- Check logs: `docker-compose logs [service-name]`

### Database connection errors
- Ensure databases are healthy: `docker-compose ps`
- Check database logs: `docker-compose logs odfi-db`
- Verify environment variables in `docker-compose.yml`

### Console can't reach other services
- Ensure all services are on the same Docker network
- Check service names match environment variables
- Use service hostnames (e.g., `http://odfi:8080`) not `localhost` in Docker

## License

MIT

## Documentation

📖 **Complete documentation available in the [`docs/`](docs/) directory:**

- [API Gateway Reference](docs/GATEWAY.md) - All 26 endpoints
- [Demo Guide](docs/DEMO.md) - Step-by-step scenarios  
- [Curl Commands](docs/CURL_COMMANDS.md) - Ready for Postman
- [Implementation Details](docs/IMPLEMENTATION.md) - Technical deep dive
- [More...](docs/README.md) - See full docs index

## Contributing

This is a POC/demo project. Feel free to fork and extend!

