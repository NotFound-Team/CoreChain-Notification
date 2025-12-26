# CoreChain-Notification

Notification service for CoreChain HRM system using Firebase Cloud Messaging (FCM) and Kafka.

## 📋 Overview

This is a Golang-based notification service that:
- Consumes Kafka messages from the NestJS backend
- Sends push notifications via Firebase Cloud Messaging (FCM)
- Stores notification history in PostgreSQL
- Supports multiple notification types (tasks, messages, calls)

**Current Flow**: Task created on server (NestJS) → Kafka → **Notification service (Golang)** → FCM → Mobile app

## 🏗️ Architecture

Built with **Clean Architecture** principles:
- `cmd/server/` - Application entry point
- `internal/domain/` - Business entities and interfaces
- `internal/application/` - Business logic (services & DTOs)
- `internal/infrastructure/` - External integrations (Kafka, FCM, PostgreSQL)
- `internal/delivery/` - Message handlers
- `pkg/` - Public/shared code

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15
- Kafka (or use Docker Compose setup)

### Setup

1. **Clone and setup**
   ```bash
   cd CoreChain-Notification
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Start with Docker Compose** (includes PostgreSQL, Kafka, and service)
   ```bash
   make docker-up
   ```

4. **Or run locally**
   ```bash
   # Make sure PostgreSQL and Kafka are running
   make run
   ```

## 📝 Configuration

Configure via environment variables (see `.env.example`):

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=notification_user
DB_PASSWORD=notification_pass
DB_NAME=notification_db

# Kafka
KAFKA_BROKERS=localhost:9092
KAFKA_GROUP_ID=notification-service-group
KAFKA_TOPIC_TASK_CREATED=task.created

# FCM
FCM_CREDENTIALS_PATH=./google-services.json
FCM_PROJECT_ID=corechain-e1321
```

## 🔧 Available Commands

```bash
make help           # Show all available commands
make deps           # Download dependencies
make build          # Build the application
make run            # Run locally
make test           # Run tests
make docker-up      # Start with Docker Compose
make docker-down    # Stop Docker services
make docker-logs    # View logs
```

## 📦 Project Structure

```
CoreChain-Notification/
├── cmd/server/              # Main application
├── internal/
│   ├── config/             # Configuration management
│   ├── domain/             # Domain models & interfaces
│   ├── application/        # Services & DTOs
│   ├── infrastructure/     # Kafka, FCM, PostgreSQL
│   ├── delivery/           # Message handlers
│   └── utils/              # Logger, errors, validators
├── deployments/
│   └── docker/             # Dockerfile & docker-compose
└── configs/                # Config files
```

## 🔔 Supported Notifications

### ✅ Implemented
- **Task Created** - Notifies when a new task is assigned

### 🚧 Planned
- **Task Updated** - Notifies when a task is modified
- **New Message** - Notifies on incoming messages
- **Incoming Call** - Notifies on incoming calls

## 🛠️ Development

### Adding New Notification Type

1. Add constant in `pkg/constants/notification_types.go`
2. Create DTO in `internal/application/dto/`
3. Add template in `internal/infrastructure/fcm/templates.go`
4. Create handler in `internal/delivery/kafka/`
5. Register handler in `cmd/server/main.go`

### Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
```

## 🐳 Docker

### Build Image
```bash
make docker-build
```

### Docker Compose Services
- **postgres** - PostgreSQL database (port 5432)
- **kafka** - Apache Kafka (port 9092)
- **zookeeper** - Kafka dependency (port 2181)
- **notification-service** - This service (port 8080)

## 📊 Database Schema

### Notifications Table
Stores all notification records with delivery status tracking.

### User FCM Tokens Table
Tracks device tokens for each user.

See `deployments/docker/migrations/001_init.sql` for full schema.

## 🤝 Integration with NestJS

The NestJS backend should publish messages to Kafka in this format:

```json
{
  "event_type": "task.created",
  "timestamp": "2025-12-26T07:29:50Z",
  "data": {
    "_id": "task-id",
    "title": "Task title",
    "assignedTo": "user-id",
    ...
  },
  "metadata": {
    "assignedToUser": {
      "_id": "user-id",
      "fcmToken": "device-token",
      "name": "User Name"
    }
  }
}
```

## 📄 License

MIT License
