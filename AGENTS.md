You are an expert Principal level Go backend engineer.
Build a small ACH-themed microservice POC with four services and a unified console API, all in Go, using Postgres and Docker.

High-level goals
	•	Repository root: ach-concourse
	•	Language: Go (use a recent stable version, e.g., Go 1.22+)
	•	Datastore: Postgres (use a stable image like postgres:latest in Docker)
	•	Services (all HTTP, JSON APIs):
	•	odfi – origination (outgoing ACH)
	•	rdfi – receiving (incoming ACH)
	•	ledger – very simple ledger/posting demo
	•	eip – Exception / Investigation Platform (cases on top of entries)
	•	console – unified query & orchestration API for Postman / UI to talk to
	•	Each service runs as its own container.
	•	Use docker-compose to run all services + their Postgres instances locally.
	•	All APIs must be callable easily from Postman (JSON request/response, clear URLs).

⸻

1. Repository & module layout

Create this directory structure:

ach-concourse/
  go.mod
  go.sum
  cmd/
    odfi/
      main.go
    rdfi/
      main.go
    ledger/
      main.go
    eip/
      main.go
    console/
      main.go
  internal/
    common/
      db/
        db.go          # shared DB connection helper
      http/
        response.go    # helpers for JSON responses and error handling
    odfi/
      handlers.go
      models.go
      repository.go
      service.go
    rdfi/
      handlers.go
      models.go
      repository.go
      service.go
    ledger/
      handlers.go
      models.go
      repository.go
      service.go
    eip/
      handlers.go
      models.go
      repository.go
      service.go
    console/
      handlers.go
      models.go
      service.go       # calls other services over HTTP
  Dockerfile.odfi
  Dockerfile.rdfi
  Dockerfile.ledger
  Dockerfile.eip
  Dockerfile.console
  docker-compose.yml
  README.md

Go module
	•	Single module at the root: module ach-concourse
	•	Use Go modules to manage dependencies.

⸻

2. Common libraries

internal/common/db/db.go

Implement a small DB helper:
	•	Function NewPostgresConnectionFromEnv() that:
	•	Reads env vars: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE
	•	Connects to Postgres (use github.com/jackc/pgx/v5/stdlib or github.com/jackc/pgx/v5 through database/sql)
	•	Returns *sql.DB with a sensible connection pool config.

Expose a helper to run CREATE TABLE IF NOT EXISTS statements on startup for each service.

internal/common/http/response.go

Helpers:
	•	JSON(w http.ResponseWriter, status int, v any)
	•	Error(w http.ResponseWriter, status int, message string)
	•	Common error struct { "error": "..." }

Use net/http + a router like github.com/go-chi/chi/v5.

⸻

3. Domain models & APIs (POC-level)

This is a spike, so keep models minimal but ACH-flavored.

3.1 ODFI service (origination)

Base URL: http://localhost:8081

Table: odfi_entries

Columns:
	•	id (UUID, primary key)
	•	trace_number (TEXT, not null)
	•	company_name (TEXT)
	•	sec_code (TEXT)
	•	amount_cents (BIGINT)
	•	status (TEXT) — e.g. PENDING, SENT, CANCELLED
	•	created_at (TIMESTAMP WITH TIME ZONE, default now)
	•	updated_at (TIMESTAMP WITH TIME ZONE, default now)

Go struct: ODFIEntry

Endpoints (JSON):
	•	POST /api/v1/entries
	•	Request: { "trace_number": "...", "company_name": "...", "sec_code": "PPD", "amount_cents": 12345 }
	•	Response: 201 + created entry JSON
	•	GET /api/v1/entries
	•	Optional query params: status, trace_number
	•	Response: list of entries
	•	GET /api/v1/entries/{id}
	•	Response: single entry or 404
	•	PATCH /api/v1/entries/{id}/status
	•	Request: { "status": "SENT" }
	•	Response: updated entry

Add GET /healthz returning { "status": "ok" }.

3.2 RDFI service (receiving)

Base URL: http://localhost:8082

Table: rdfi_entries

Columns:
	•	id (UUID, primary key)
	•	trace_number (TEXT, not null)
	•	receiver_name (TEXT)
	•	amount_cents (BIGINT)
	•	status (TEXT) — e.g. RECEIVED, POSTED, RETURNED
	•	return_reason (TEXT, nullable)
	•	created_at
	•	updated_at

Endpoints:
	•	POST /api/v1/entries
	•	Request: { "trace_number": "...", "receiver_name": "...", "amount_cents": 12345 }
	•	GET /api/v1/entries
	•	Filter by status, trace_number
	•	GET /api/v1/entries/{id}
	•	POST /api/v1/entries/{id}/return
	•	Request: { "reason": "R01" }
	•	Sets status="RETURNED", return_reason="R01"

/healthz as above.

3.3 Ledger service

Base URL: http://localhost:8083

Very simple “ledger entries” simulation.

Table: ledger_entries

Columns:
	•	id (UUID)
	•	ach_side (TEXT) — "ODFI" or "RDFI"
	•	trace_number (TEXT)
	•	amount_cents (BIGINT)
	•	direction (TEXT) — "DEBIT" or "CREDIT"
	•	description (TEXT)
	•	created_at

Endpoints:
	•	POST /api/v1/postings
	•	Request: { "ach_side": "ODFI", "trace_number": "...", "amount_cents": 12345, "direction": "DEBIT", "description": "Test posting" }
	•	GET /api/v1/postings
	•	Filter by ach_side, trace_number
	•	Optional: GET /api/v1/balances (return simple total of credits minus debits)

/healthz as above.

3.4 EIP service (exceptions / cases)

Base URL: http://localhost:8084

This models exceptions on top of entries.

Table: eip_cases

Columns:
	•	id (UUID)
	•	side (TEXT) — "ODFI" or "RDFI"
	•	trace_number (TEXT)
	•	status (TEXT) — "OPEN", "IN_PROGRESS", "RESOLVED"
	•	type (TEXT) — "RETURN_REVIEW", "NOC_REVIEW", "CUSTOMER_DISPUTE"
	•	notes (TEXT)
	•	created_at
	•	updated_at

Endpoints:
	•	POST /api/v1/cases
	•	Request: { "side": "RDFI", "trace_number": "...", "type": "RETURN_REVIEW", "notes": "Customer called in" }
	•	GET /api/v1/cases
	•	Filter by status, side, trace_number
	•	GET /api/v1/cases/{id}
	•	PATCH /api/v1/cases/{id}/status
	•	Request: { "status": "RESOLVED" }

/healthz as above.

⸻

4. Console service (unified ACH view)

Base URL: http://localhost:8080

This is the unified query interface POC.
It does not talk directly to any DB. It calls the other services via HTTP and merges their results in memory.

Use environment variables to discover other services:
	•	ODFI_BASE_URL (e.g., http://odfi:8080 inside Docker network)
	•	RDFI_BASE_URL
	•	LEDGER_BASE_URL
	•	EIP_BASE_URL

Console model

Define a simple unified view struct:

type UnifiedAchItem struct {
    Side         string  `json:"side"`          // "ODFI" or "RDFI"
    Source       string  `json:"source"`        // "odfi", "rdfi"
    EntryID      string  `json:"entry_id"`
    TraceNumber  string  `json:"trace_number"`
    AmountCents  int64   `json:"amount_cents"`
    Status       string  `json:"status"`
    Extra        any     `json:"extra,omitempty"` // Optional map or struct for service-specific fields
}

Console endpoints
	1.	GET /api/v1/ach-items
	•	Query params: side, status, trace_number (all optional).
	•	Behaviour:
	•	Calls ODFI /api/v1/entries and RDFI /api/v1/entries (using appropriate filters).
	•	Normalizes results into []UnifiedAchItem.
	•	Returns merged list as JSON.
	2.	GET /api/v1/ach-items/{side}/{id}
	•	side is "ODFI" or "RDFI".
	•	id is the service’s entry id.
	•	Calls the appropriate service’s GET /api/v1/entries/{id}.
	•	Returns unified representation.
	3.	POST /api/v1/ach-items/{side}/{id}/return
	•	For this POC, only support side="RDFI".
	•	Request: { "reason": "R01" }
	•	Console calls RDFI service’s POST /api/v1/entries/{id}/return.
	•	Returns the RDFI entry response (or a unified form).

/healthz as above.

The console service demonstrates that the UI only needs one API to query and operate on both RDFI and ODFI data, without a shared DB.

⸻

5. HTTP server & routing

For all services:
	•	Use github.com/go-chi/chi/v5 as router.
	•	In each cmd/<service>/main.go:
	•	Parse env vars for port and DB connection (where applicable).
	•	Initialize DB via common/db.NewPostgresConnectionFromEnv() (except console).
	•	Run CREATE TABLE IF NOT EXISTS ... to ensure schema.
	•	Wire up handlers in internal/<service> with chi routes.
	•	Start HTTP server with graceful shutdown.

Ports (host-mapped for Postman):
	•	console: 8080
	•	odfi: 8081
	•	rdfi: 8082
	•	ledger: 8083
	•	eip: 8084

Inside Docker, services can listen on 0.0.0.0:8080; docker-compose will map to host.

⸻

6. Dockerfiles

Create one Dockerfile per service (Dockerfile.odfi, etc.), following this pattern (multi-stage):

# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder
WORKDIR /app

# Install git if needed
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build specific service binary; replace "odfi" with the service name
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/odfi ./cmd/odfi

FROM alpine:3.20
WORKDIR /app

ENV PORT=8080

COPY --from=builder /bin/odfi /app/odfi

EXPOSE 8080

CMD ["/app/odfi"]

Duplicate/adjust for rdfi, ledger, eip, console (changing binary name and build path accordingly).

⸻

7. Docker Compose

Create docker-compose.yml at the repo root with:
	•	A Docker network (default).
	•	Postgres containers for each service:
	•	odfi-db on 5433 (host)
	•	rdfi-db on 5434
	•	ledger-db on 5435
	•	eip-db on 5436
	•	Service containers:
	•	odfi, rdfi, ledger, eip, console

Example (simplified, you fill in fully):

version: "3.9"

services:
  odfi-db:
    image: postgres:latest
    environment:
      POSTGRES_USER: odfi_user
      POSTGRES_PASSWORD: odfi_pass
      POSTGRES_DB: odfi_db
    ports:
      - "5433:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U odfi_user -d odfi_db"]
      interval: 5s
      timeout: 5s
      retries: 5

  odfi:
    build:
      context: .
      dockerfile: Dockerfile.odfi
    environment:
      PORT: 8080
      DB_HOST: odfi-db
      DB_PORT: 5432
      DB_USER: odfi_user
      DB_PASSWORD: odfi_pass
      DB_NAME: odfi_db
      DB_SSLMODE: disable
    depends_on:
      odfi-db:
        condition: service_healthy
    ports:
      - "8081:8080"

  rdfi-db:
    image: postgres:latest
    environment:
      POSTGRES_USER: rdfi_user
      POSTGRES_PASSWORD: rdfi_pass
      POSTGRES_DB: rdfi_db
    ports:
      - "5434:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U rdfi_user -d rdfi_db"]
      interval: 5s
      timeout: 5s
      retries: 5

  rdfi:
    build:
      context: .
      dockerfile: Dockerfile.rdfi
    environment:
      PORT: 8080
      DB_HOST: rdfi-db
      DB_PORT: 5432
      DB_USER: rdfi_user
      DB_PASSWORD: rdfi_pass
      DB_NAME: rdfi_db
      DB_SSLMODE: disable
    depends_on:
      rdfi-db:
        condition: service_healthy
    ports:
      - "8082:8080"

  ledger-db:
    image: postgres:latest
    environment:
      POSTGRES_USER: ledger_user
      POSTGRES_PASSWORD: ledger_pass
      POSTGRES_DB: ledger_db
    ports:
      - "5435:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ledger_user -d ledger_db"]
      interval: 5s
      timeout: 5s
      retries: 5

  ledger:
    build:
      context: .
      dockerfile: Dockerfile.ledger
    environment:
      PORT: 8080
      DB_HOST: ledger-db
      DB_PORT: 5432
      DB_USER: ledger_user
      DB_PASSWORD: ledger_pass
      DB_NAME: ledger_db
      DB_SSLMODE: disable
    depends_on:
      ledger-db:
        condition: service_healthy
    ports:
      - "8083:8080"

  eip-db:
    image: postgres:latest
    environment:
      POSTGRES_USER: eip_user
      POSTGRES_PASSWORD: eip_pass
      POSTGRES_DB: eip_db
    ports:
      - "5436:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U eip_user -d eip_db"]
      interval: 5s
      timeout: 5s
      retries: 5

  eip:
    build:
      context: .
      dockerfile: Dockerfile.eip
    environment:
      PORT: 8080
      DB_HOST: eip-db
      DB_PORT: 5432
      DB_USER: eip_user
      DB_PASSWORD: eip_pass
      DB_NAME: eip_db
      DB_SSLMODE: disable
    depends_on:
      eip-db:
        condition: service_healthy
    ports:
      - "8084:8080"

  console:
    build:
      context: .
      dockerfile: Dockerfile.console
    environment:
      PORT: 8080
      ODFI_BASE_URL: http://odfi:8080
      RDFI_BASE_URL: http://rdfi:8080
      LEDGER_BASE_URL: http://ledger:8080
      EIP_BASE_URL: http://eip:8080
    depends_on:
      - odfi
      - rdfi
      - ledger
      - eip
    ports:
      - "8080:8080"


⸻

8. README

Create a README.md with:
	•	How to build and run:

docker-compose build
docker-compose up


	•	List of Postman-friendly endpoints:
	•	ODFI: http://localhost:8081/api/v1/entries
	•	RDFI: http://localhost:8082/api/v1/entries
	•	Ledger: http://localhost:8083/api/v1/postings
	•	EIP: http://localhost:8084/api/v1/cases
	•	Console unified queries:
	•	GET http://localhost:8080/api/v1/ach-items
	•	GET http://localhost:8080/api/v1/ach-items/{side}/{id}
	•	POST http://localhost:8080/api/v1/ach-items/RDFI/{id}/return
	•	Example JSON bodies for each POST/PATCH.
