# Chat Application with Golang + WebSocket

Real-time chat application built with Golang, Gin Framework, WebSocket, and MySQL.

## Features

- 🔐 User authentication (Register/Login) with JWT
- 💬 Real-time messaging with WebSocket
- 📝 Chat history
- ✅ Read receipts
- 👥 User management
- 🐳 Docker support
- ☁️ Ready for GCP deployment

## Tech Stack

- **Language**: Golang 1.21+
- **Framework**: Gin
- **Database**: MySQL 8.0
- **WebSocket**: gorilla/websocket
- **ORM**: GORM
- **Authentication**: JWT
- **Containerization**: Docker

## Project Structure

```
chat-app/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── config/          # Configuration
│   ├── domain/          # Business entities
│   ├── repository/      # Data access layer
│   ├── usecase/         # Business logic
│   ├── delivery/        # Presentation layer (HTTP, WebSocket)
│   └── utils/           # Helper functions
├── pkg/                 # Public libraries
├── migrations/          # Database migrations
├── docker/              # Docker configuration
└── scripts/             # Utility scripts
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0
- Docker (optional)

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd chat-app
```

2. Install dependencies
```bash
make install-deps
```

3. Setup environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run with Docker (recommended)
```bash
make docker-up
```

Or run locally:
```bash
# Make sure MySQL is running
make build
make run
```

### API Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user

#### User
- `GET /api/v1/profile` - Get user profile (protected)
- `GET /api/v1/users` - Get all users (protected)

#### Messages
- `POST /api/v1/messages` - Send message (protected)
- `GET /api/v1/messages/:userId` - Get chat history (protected)
- `PATCH /api/v1/messages/:messageId/read` - Mark as read (protected)
- `GET /api/v1/messages/unread/count` - Get unread count (protected)

#### WebSocket
- `GET /api/v1/ws?user_id=<id>` - WebSocket connection

### WebSocket Message Format

Send message:
```json
{
  "type": "message",
  "receiver_id": 2,
  "content": "Hello!"
}
```

## Development

### Running Tests
```bash
make test
```

### Building
```bash
make build
```

### Database Migration
```bash
make migrate
```

## Deployment to GCP

### Using Cloud Run

1. Build and push Docker image
```bash
# Replace YOUR_PROJECT_ID with your GCP project ID
make docker-build-gcp
make docker-push-gcp
```

2. Deploy to Cloud Run
```bash
gcloud run deploy chat-app \
  --image gcr.io/YOUR_PROJECT_ID/chat-app:latest \
  --platform managed \
  --region asia-southeast1 \
  --allow-unauthenticated \
  --set-env-vars="DB_HOST=YOUR_DB_HOST,DB_USER=YOUR_DB_USER,DB_PASSWORD=YOUR_DB_PASSWORD"
```

### Using GKE (Google Kubernetes Engine)

1. Create cluster
```bash
gcloud container clusters create chat-cluster \
  --num-nodes=3 \
  --zone=asia-southeast1-a
```

2. Deploy application
```bash
kubectl apply -f k8s/
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Server port | 8080 |
| DB_HOST | MySQL host | localhost |
| DB_PORT | MySQL port | 3306 |
| DB_USER | MySQL user | root |
| DB_PASSWORD | MySQL password | - |
| DB_NAME | Database name | chat_app |
| JWT_SECRET | JWT secret key | - |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

```

---

## Cara Menjalankan Aplikasi:

### 1. Setup dengan Docker (Termudah)
```bash
# Copy environment file
cp .env.example .env

# Start semua services (MySQL + App)
make docker-up

# Lihat logs
make docker-logs

# Stop services
make docker-down
```

### 2. Setup Manual
```bash
# Install dependencies
go mod download

# Setup database MySQL terlebih dahulu
# Kemudian jalankan migrasi
make migrate

# Build aplikasi
make build

# Run aplikasi
make run
```

### 3. Testing API dengan cURL

**Register:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"password123"}'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'
```

**Send Message:**
```bash
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"receiver_id":2,"content":"Hello!"}'
```

### 4. WebSocket Connection (JavaScript)
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?user_id=1');

ws.onopen = () => {
  console.log('Connected to WebSocket');
  
  // Send message
  ws.send(JSON.stringify({
    type: 'message',
    receiver_id: 2,
    content: 'Hello via WebSocket!'
  }));
};

ws.onmessage = (event) => {
  console.log('Received:', JSON.parse(event.data));
};
```

## Dependencies yang Perlu Diinstall:

```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/gorilla/websocket
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/joho/godotenv
go get -u golang.org/x/crypto/bcrypt
```


## Penjelasan Struktur:

### 1. **cmd/server/**
Entry point aplikasi. Berisi main.go yang akan menginisialisasi semua dependencies dan menjalankan server.

### 2. **internal/**
Package private yang tidak bisa diimport oleh project lain.

- **config/**: Konfigurasi aplikasi (database, server, JWT, dll)
- **domain/**: Entity/model bisnis (User, Message)
- **repository/**: Layer untuk akses data (interface dan implementasi MySQL)
- **usecase/**: Business logic aplikasi
- **delivery/**: Layer presentasi (HTTP handlers, WebSocket, middleware, routing)
- **utils/**: Helper functions (JWT, password hashing, response format)

### 3. **pkg/**
Package yang bisa digunakan oleh project lain. Berisi database connection setup.

### 4. **migrations/**
SQL migration files untuk setup database schema.

### 5. **docker/**
File-file terkait Docker (Dockerfile, docker-compose).

### 6. **scripts/**
Bash scripts untuk automation (migrasi, build, deploy).

## Arsitektur:
Menggunakan **Clean Architecture** dengan pembagian layer:
1. **Domain Layer** - Entity bisnis
2. **Repository Layer** - Data access
3. **Service Layer** - Business logic
4. **Delivery Layer** - Presentation (HTTP/WebSocket)

## Benefits:
- ✅ Separation of concerns yang jelas
- ✅ Mudah untuk testing (mockable interfaces)
- ✅ Scalable dan maintainable
- ✅ Independent dari framework
- ✅ Siap untuk deployment dengan Docker
- ✅ Mengikuti Golang best practices

## Dependencies Utama:
```bash
go get -u github.com/gin-gonic/gin
go get -u github.com/gorilla/websocket
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/joho/godotenv
go get -u golang.org/x/crypto/bcrypt
```

## Next Steps:
1. Implementasi setiap layer sesuai kebutuhan
2. Setup Docker dan docker-compose
3. Buat migration files
4. Implementasi WebSocket hub untuk broadcast messages
5. Setup deployment ke GCP (Cloud Run atau GKE)