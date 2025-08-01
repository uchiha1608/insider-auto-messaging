# 📬 Insider Auto Messaging System

This project implements a Golang-based automated messaging system as described in the Insider backend assessment. It sends unsent messages from a PostgreSQL database to a webhook endpoint every 2 minutes, logs successful sends, and caches the message ID and timestamp in Redis.

## ✨ Features

- ✅ Auto-start background message sender (no cron)
- ✅ Sends 2 unsent messages every 2 minutes
- ✅ PostgreSQL-based message queue
- ✅ Redis caching of sent message IDs (Bonus)
- ✅ RESTful API to start/stop scheduler and list sent messages
- ✅ Auto database table migration
- ✅ Swagger API docs with `swag init` support
- ✅ Unit tests for service, controller, and repository layers

---

## 🚀 Tech Stack

- Go 1.24
- PostgreSQL
- Redis
- Docker & Docker Compose
- Swagger (via [swaggo/swag](https://github.com/swaggo/swag))
- Testify, sqlmock, redismock

---

## 📁 Folder Structure

```
insider-auto-messaging/
├── cmd/                    # Entry point
│   └── main.go
├── config/                 # DB, Redis setup
├── controller/             # API endpoints
├── model/                  # Data models
├── repository/             # Database operations
├── scheduler/              # Custom scheduler (no cron)
├── service/                # Business logic
├── docs/                   # Auto-generated Swagger docs
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

---

## 🧑‍💻 Local Development Setup

### 1. Clone the Repository

```bash
git clone https://github.com/<your-username>/insider-auto-messaging.git
cd insider-auto-messaging
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Install `swag` (once only)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Make sure `$GOPATH/bin` is in your system `PATH`.

### 4. Generate Swagger Docs

```bash
swag init
```

### 5. Run Locally

```bash
go run ./cmd/main.go
```

App will start at: [http://localhost:8080](http://localhost:8080)

---

## 🐳 Run with Docker

### 1. Build and Start

```bash
docker-compose up --build
```

This spins up:
- `app` (Golang service)
- `db` (PostgreSQL)
- `redis` (Redis)

### 2. Swagger UI

Access Swagger docs at:  
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## 🔌 API Endpoints

| Method | Endpoint      | Description                      |
|--------|---------------|----------------------------------|
| GET    | `/start`      | Start auto message sending       |
| GET    | `/stop`       | Stop auto message sending        |
| GET    | `/sent`       | Get list of sent messages        |
| GET    | `/swagger/*`  | Swagger UI (auto-generated)      |

---

## 📦 PostgreSQL Schema

Table is auto-created at app startup.

```sql
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    to TEXT NOT NULL,
    content TEXT NOT NULL CHECK (char_length(content) <= 160),
    is_sent BOOLEAN DEFAULT FALSE,
    message_id TEXT,
    sent_at TIMESTAMP
);
```

---

## 📚 Example Message Record

Insert a message to test:

```sql
INSERT INTO messages (to, content) VALUES ('+905551111111', 'Insider - Project Test Message');
```

App will send this to the webhook in the next cycle.

---

## 🎯 Webhook Configuration

The app sends messages to:

```
https://webhook.site/<your-id>
```

Edit the URL in `message_service.go` if needed. You can create a custom one at: [webhook.site](https://webhook.site)

---

## 🧪 Redis Caching

After a message is sent, the app stores:

- Key: `messageId` from webhook response
- Value: send timestamp

Use `redis-cli` to inspect:

```bash
docker exec -it insider-auto-messaging-redis-1 redis-cli
keys *
get <message-id>
```

---

## ✅ Run Unit Tests

This project includes unit tests for:
- `service/message_service.go`
- `controller/message_controller.go`
- `repository/message_repository.go`

### 📦 Install test dependencies

```bash
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock
go get github.com/go-redis/redismock/v9
```

### ▶ Run all tests

```bash
go test ./... -v
```

### ▶ Run tests for a specific package

```bash
go test ./service -v
go test ./controller -v
go test ./repository -v
```

---

## 🛠 Development Tips

- Change webhook URL for testing
- Modify scheduler interval in `scheduler/sender.go` (default: 2 mins)
- Customize DB connection in `config/config.go`
- Swagger annotations live in `controller/`

---

## 🔒 Auth Headers for Webhook (Static)

The app uses:

```http
Header: x-ins-auth-key: INS.me1x9uMcyYGlhKKQVPoc.bO3j9aZwRTOcA2Ywo
```

As required by the assessment brief.

---

## 🧼 Cleanup

Stop Docker:

```bash
docker-compose down
```

Remove volumes:

```bash
docker-compose down -v
```

---

## 📄 License

This project is submitted for Insider's backend assessment only.
