# StreamPulse

A creator sponsorship analytics REST API built in Go. Tracks deals between content creators and sponsors, processes events asynchronously through a goroutine pipeline, and exposes aggregated analytics.

## Stack

| Layer | Technology |
|-------|-----------|
| HTTP framework | [Gin](https://github.com/gin-gonic/gin) |
| Auth | JWT (HS256) via [golang-jwt/jwt](https://github.com/golang-jwt/jwt) |
| Storage | In-memory (DynamoDB-shaped interfaces — swap in prod) |
| Async queue | Buffered Go channel (SQS-shaped interface — swap in prod) |
| Event processing | Goroutine worker pool pipeline |
| Password hashing | bcrypt |

---

## Getting Started

### Prerequisites

- Go 1.22+

```bash
# macOS
brew install go
```

### Run

```bash
git clone <repo-url>
cd StreamPulse

make tidy   # download dependencies
make run    # start server on :8080
```

You should see:

```
[pipeline] deal tracker started with 4 workers
StreamPulse API listening on :8080
```

### Other commands

```bash
make build  # compile binary to ./bin/streampulse
make test   # run all tests with race detector
make vet    # static analysis
make clean  # remove build artifacts
```

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server listen port |
| `JWT_SECRET` | `streampulse-dev-secret-change-in-prod` | HMAC signing key — **change in production** |
| `WORKER_POOL_SIZE` | `4` | Number of pipeline goroutines consuming deal events |
| `QUEUE_BUFFER` | `100` | Buffered channel capacity for the event queue |

---

## Architecture

```
HTTP Request
     │
     ▼
 Gin Router  ──(JWT middleware)──▶ Handlers
     │                                │
     │                                ▼
     │                           Services
     │                           │       │
     │                      Repository  Queue
     │                      (in-memory) (channel)
     │                                    │
     │                                    ▼
     │                           Pipeline Workers
     │                           (goroutine pool)
     │                                    │
     │                                    ▼
     │                        Analytics side-effects
```

### Deal State Machine

Deals follow a strict status lifecycle. Invalid transitions are rejected at the service layer.

```
pending ──▶ active ──▶ completed
   │           │
   └───────────┴──▶ cancelled
```

### Event Pipeline

When a deal status changes, a `DealEvent` is published to the queue. A pool of goroutine workers consumes these events:

| Event | Worker action |
|-------|--------------|
| `deal.created` | Log new deal |
| `deal.activated` | Seed an empty analytics record for the deal |
| `deal.completed` | Log final deal value |
| `deal.cancelled` | Log cancellation |

---

## API Reference

All `/api/v1/*` routes require:
```
Authorization: Bearer <token>
```

### Auth

#### Register
```http
POST /auth/register
Content-Type: application/json

{
  "email": "creator@example.com",
  "password": "securepassword"
}
```
```json
{
  "token": "<jwt>",
  "user": { "id": "...", "email": "creator@example.com", "created_at": "..." }
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "creator@example.com",
  "password": "securepassword"
}
```

---

### Creators

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/creators` | List all creators |
| `POST` | `/api/v1/creators` | Create a creator |
| `GET` | `/api/v1/creators/:id` | Get creator by ID |
| `PUT` | `/api/v1/creators/:id` | Update creator |
| `DELETE` | `/api/v1/creators/:id` | Delete creator |
| `GET` | `/api/v1/creators/:id/analytics` | All analytics across creator's deals |

**Create creator payload:**
```json
{
  "name": "Jane Doe",
  "platform": "youtube",
  "channel_url": "https://youtube.com/@janedoe",
  "subscriber_count": 500000,
  "category": "tech",
  "email": "jane@example.com"
}
```

Supported platforms: `youtube`, `twitch`, `instagram`, `tiktok`

---

### Sponsors

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/sponsors` | List all sponsors |
| `POST` | `/api/v1/sponsors` | Create a sponsor |
| `GET` | `/api/v1/sponsors/:id` | Get sponsor by ID |
| `PUT` | `/api/v1/sponsors/:id` | Update sponsor |
| `DELETE` | `/api/v1/sponsors/:id` | Delete sponsor |

**Create sponsor payload:**
```json
{
  "name": "Acme Corp",
  "industry": "software",
  "website": "https://acme.com",
  "contact_email": "deals@acme.com",
  "budget": 50000.00
}
```

---

### Deals

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/deals` | List all deals |
| `POST` | `/api/v1/deals` | Create a deal |
| `GET` | `/api/v1/deals/:id` | Get deal by ID |
| `PATCH` | `/api/v1/deals/:id/status` | Advance deal status |
| `DELETE` | `/api/v1/deals/:id` | Delete deal |
| `GET` | `/api/v1/deals/:id/analytics` | Get analytics for a deal |
| `POST` | `/api/v1/deals/:id/analytics` | Record analytics snapshot |

**Create deal payload:**
```json
{
  "creator_id": "<uuid>",
  "sponsor_id": "<uuid>",
  "value": 5000.00,
  "start_date": "2026-06-01T00:00:00Z",
  "end_date": "2026-06-30T00:00:00Z",
  "description": "30-second mid-roll sponsorship"
}
```

**Update deal status:**
```json
{ "status": "active" }
```

**Record analytics:**
```json
{
  "views": 120000,
  "clicks": 3400,
  "revenue": 4800.00
}
```
`conversion_rate` is computed automatically: `(clicks / views) * 100`

---

### Analytics

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/v1/analytics/summary` | Aggregated platform-wide stats |

**Response:**
```json
{
  "total_views": 500000,
  "total_clicks": 14200,
  "total_revenue": 22500.00,
  "avg_conversion_rate": 2.84,
  "active_deals": 3,
  "completed_deals": 7,
  "total_deals": 12
}
```

---

## Project Structure

```
StreamPulse/
├── cmd/
│   └── server/
│       └── main.go               # entrypoint — wires all dependencies
├── internal/
│   ├── auth/
│   │   ├── jwt.go                # token generation and verification
│   │   └── middleware.go         # Gin Bearer token middleware
│   ├── config/
│   │   └── config.go             # env-based configuration
│   ├── models/                   # domain structs and request/response types
│   ├── repository/
│   │   ├── interfaces.go         # DynamoDB-shaped storage contracts
│   │   └── memory/               # in-memory implementations (thread-safe)
│   ├── queue/
│   │   ├── interface.go          # SQS-shaped queue contract + event types
│   │   └── memory/               # buffered channel implementation
│   ├── pipeline/
│   │   └── deal_tracker.go       # goroutine worker pool — processes DealEvents
│   ├── service/                  # business logic (deal state machine, analytics)
│   ├── handlers/                 # Gin HTTP handlers
│   └── router/
│       └── router.go             # route registration
├── Makefile
└── go.mod
```

---

## Replacing In-Memory Stubs with AWS

The repository and queue interfaces are designed for this. In production:

1. **DynamoDB** — implement `repository.CreatorRepository`, `DealRepository`, etc. using the AWS SDK v2 DynamoDB client and point `main.go` at the new constructors.
2. **SQS** — implement `queue.Queue` using `sqs.SendMessage` / long-poll `ReceiveMessage` and swap it in `main.go`.

No handler, service, or pipeline code needs to change.
