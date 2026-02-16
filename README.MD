# 🧾 M-PESA Statement Analyzer

> A secure, production-grade financial analysis platform built with Go. Transform your M-PESA statements into actionable insights with AI-powered categorization and interactive dashboards.

 • [Report Bug](https://github.com/Wendyshiro/M-PESA-Statements-Analyzer/issues) • [Request Feature](https://github.com/Wendyshiro/M-PESA-Statements-Analyzer/issues)

---

## 📖 Table of Contents

- [About The Project](#-about-the-project)
- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Getting Started](#-getting-started)
- [Project Structure](#-project-structure)
- [API Documentation](#-api-documentation)
- [Development Roadmap](#-development-roadmap)
- [Contributing](#-contributing)
- [License](#-license)
- [Contact](#-contact)
- [Acknowledgments](#-acknowledgments)

---

## 🎯 About The Project

M-PESA Statement Analyzer is a comprehensive financial management tool that helps users understand their spending patterns by parsing M-PESA PDF statements, automatically categorizing transactions, and presenting beautiful visualizations.

### The Problem
- Manual analysis of M-PESA statements is time-consuming
- Difficult to track spending across categories
- No built-in tools for financial trend analysis
- Hard to identify spending patterns and anomalies

### The Solution
This application provides:
- **Automated PDF parsing** - Extract all transactions from your statements instantly
- **AI categorization** - Intelligently classify transactions (Food, Transport, Bills, etc.)
- **Visual analytics** - Interactive charts and graphs to understand your finances
- **Secure storage** - Your financial data is encrypted and protected
- **Multi-user support** - Each user has their own private workspace

---

## ✨ Features

### 🔐 Security & Authentication
- JWT-based authentication system
- Bcrypt password hashing
- Secure file validation and sanitization
- CORS and security headers protection
- Input validation on all endpoints
- Encrypted data storage (coming in Phase 4)

### 📊 Financial Analysis
- PDF statement parsing and text extraction
- AI-powered transaction categorization using OpenAI
- Summary reports (total inflows, outflows, balance)
- Category-wise spending breakdown
- Time-series trend analysis (coming in Phase 3)
- Custom date range filtering (coming in Phase 3)

### 🏗️ Architecture & Performance
- Clean architecture with repository pattern
- PostgreSQL with optimized indexes
- Database migration system
- Redis caching (coming in Phase 4)
- Async job processing (coming in Phase 3)
- Docker containerization
- Horizontal scalability ready

### 📱 User Experience
- RESTful API design
- Interactive web dashboard
- Real-time job status tracking (coming in Phase 3)
- Export to CSV/Excel (coming in Phase 5)
- Mobile-responsive interface

---

## 🛠️ Tech Stack

### Backend
- **Go 1.21+** - Main programming language
- **PostgreSQL 15** - Primary database
- **Redis 7** - Caching layer (Phase 4)
- **JWT** - Authentication tokens
- **Bcrypt** - Password encryption

### Frontend
- **HTML5/CSS3** - Dashboard UI
- **Chart.js** - Data visualization
- **Vanilla JavaScript** - Client-side logic

### DevOps & Tools
- **Docker & Docker Compose** - Containerization
- **golang-migrate** - Database migrations
- **Air** - Hot reload in development
- **GitHub Actions** - CI/CD (Phase 5)

### External APIs
- **OpenAI GPT** - Transaction categorization
- **pdftotext** - PDF parsing

---

## 🚀 Getting Started

### Prerequisites

Make sure you have the following installed:
- [Go 1.21+](https://go.dev/dl/)
- [Docker & Docker Compose](https://www.docker.com/get-started)
- [Git](https://git-scm.com/downloads)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/Wendyshiro/M-PESA-Statements-Analyzer.git
   cd M-PESA-Statements-Analyzer
   ```

2. **Install Go dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` with your configuration:
   ```env
   # Server
   PORT=8080
   ENV=development
   
   # Database
   DATABASE_URL=postgresql://mpesa_user:mpesa_pass@localhost:5432/mpesa_analyzer?sslmode=disable
   
   # Security (generate secure keys!)
   JWT_SECRET=your-secret-key-minimum-32-characters
   ENCRYPTION_KEY=exactly-32-characters-required!!
   
   # OpenAI
   OPENAI_API_KEY=sk-your-openai-api-key
   ```

4. **Start infrastructure with Docker**
   ```bash
   docker-compose up -d
   ```
   
   This starts:
   - PostgreSQL on port 5432
   - Redis on port 6379
   - pgAdmin on port 5050 (optional)

5. **Run database migrations**
   ```bash
   # Install migrate tool
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   
   # Apply migrations
   make migrate-up
   # OR
   migrate -path migrations -database "${DATABASE_URL}" up
   ```

6. **Start the application**
   ```bash
   # Using Make
   make run
   
   # OR directly with Go
   go run cmd/api/main.go
   ```

7. **Access the application**
   - API: http://localhost:8080
   - Dashboard: http://localhost:8080/dashboard
   - Health Check: http://localhost:8080/health

### Quick Test

```bash
# Register a new user
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# You'll get a token in the response - copy it!

# Upload a statement (replace YOUR_TOKEN)
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@path/to/statement.pdf"
```

---

## 📁 Project Structure

```
M-PESA-Statements-Analyzer/
│
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
│
├── handlers/                    # HTTP request handlers
│   ├── auth.go                  # Registration & login
│   ├── upload.go                # File upload handling
│   └── health.go                # Health checks
│
├── models/                      # Data models
│   ├── user.go                  # User entity
│   ├── job.go                   # Processing job entity
│   └── transaction.go           # Transaction entity
│
├── repository/                  # Database layer (Data Access)
│   ├── user.go                  # User DB operations
│   ├── job.go                   # Job DB operations
│   └── transaction.go           # Transaction DB operations
│
├── middleware/                  # HTTP middleware
│   ├── auth.go                  # JWT authentication
│   ├── security.go              # Security headers
│   └── validation.go            # Input validation
│
├── auth/                        # Authentication logic
│   └── jwt.go                   # Token management
│
├── database/                    # Database connection
│   └── database.go              # Pool & health checks
│
├── config/                      # Configuration
│   └── config.go                # Env variable loader
│
├── utils/                       # Utility functions
│   ├── pdf_parser.go            # PDF extraction
│   └── categorizer.go           # AI categorization
│
├── migrations/                  # Database migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_jobs_table.up.sql
│   └── ...
│
├── dashboard/                   # Frontend
│   └── index.html
│
├── .env.example                 # Environment template
├── .gitignore
├── docker-compose.yml           # Docker services
├── Dockerfile                   # Application container
├── Makefile                     # Development commands
├── go.mod
├── go.sum
└── README.md
```

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                     Client Layer                        │
│              (Browser / Mobile / Postman)               │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/HTTPS
                     ↓
┌─────────────────────────────────────────────────────────┐
│                  API Gateway Layer                      │
│              (Middleware: Auth, CORS, etc)              │
└────────────────────┬────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────┐
│                   Handlers Layer                        │
│           (auth.go, upload.go, health.go)               │
└────────────────────┬────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────┐
│                  Service Layer                          │
│        (Business Logic: PDF Parsing, AI, etc)           │
└────────────────────┬────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────┐
│                 Repository Layer                        │
│              (Database Operations)                      │
└────────────────────┬────────────────────────────────────┘
                     │
                     ↓
┌─────────────────────────────────────────────────────────┐
│                  Database Layer                         │
│            PostgreSQL + Redis Cache                     │
└─────────────────────────────────────────────────────────┘
```

---

## 📚 API Documentation

### Base URL
```
http://localhost:8080
```

### Authentication

Most endpoints require authentication. Include the JWT token in the Authorization header:
```
Authorization: Bearer <your-token>
```

### Endpoints

#### 🟢 Public Endpoints

<details>
<summary><b>POST /register</b> - Register a new user</summary>

**Request:**
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2026-02-16T10:30:00Z",
    "updated_at": "2026-02-16T10:30:00Z"
  }
}
```

**Validation Rules:**
- Email must be valid format
- Password must be at least 8 characters
- Email must be unique

</details>

<details>
<summary><b>POST /login</b> - Login existing user</summary>

**Request:**
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "created_at": "2026-02-16T10:30:00Z",
    "updated_at": "2026-02-16T10:30:00Z"
  }
}
```

**Error Responses:**
- 401: Invalid credentials
- 400: Missing email or password

</details>

<details>
<summary><b>GET /health</b> - Health check</summary>

**Request:**
```bash
curl http://localhost:8080/health
```

**Response (200 OK):**
```json
{
  "status": "healthy",
  "services": {
    "database": "healthy"
  }
}
```

</details>

#### 🔒 Protected Endpoints

<details>
<summary><b>POST /upload</b> - Upload M-PESA statement</summary>

**Request:**
```bash
curl -X POST http://localhost:8080/upload \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -F "file=@statement.pdf"
```

**Response (202 Accepted):**
```json
{
  "job_id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "message": "File uploaded successfully and queued for processing",
  "filename": "mpesa_statement.pdf"
}
```

**File Requirements:**
- Must be a PDF file
- Maximum size: 10MB
- Must be a valid M-PESA statement

**Error Responses:**
- 401: Missing or invalid token
- 400: Invalid file type or size
- 413: File too large

</details>

### Error Response Format

All errors follow this format:
```json
{
  "error": "Human-readable error message",
  "code": "ERROR_CODE"
}
```

**Common Error Codes:**
- `UNAUTHORIZED` - Missing or invalid authentication
- `INVALID_INPUT` - Validation failed
- `USER_EXISTS` - Email already registered
- `INVALID_FILE` - File validation failed
- `INTERNAL_ERROR` - Server error

---

## 🗺️ Development Roadmap

### ✅ Phase 1: Basic Security (Completed)
**Status:** ✅ Complete (Week 1-3)

- [x] File validation and sanitization
- [x] JWT-based authentication
- [x] Security headers (CORS, CSP, XSS protection)
- [x] Password hashing with bcrypt
- [x] Input validation middleware
- [x] Error handling

### 🔄 Phase 2: Database Integration (In Progress)
**Status:** 🔄 In Progress (Week 4)

- [x] PostgreSQL setup with Docker
- [x] Database migration system
- [x] User repository implementation
- [ ] Job repository implementation
- [ ] Transaction repository implementation
- [ ] Database indexes and optimization
- [ ] Connection pooling

### 📋 Phase 3: Async Processing (Upcoming)
**Status:** 📅 Planned (Week 5-6)

- [ ] Redis integration
- [ ] Message queue setup (RabbitMQ or Redis Queue)
- [ ] Background worker for PDF processing
- [ ] Job status tracking and updates
- [ ] Retry logic for failed jobs
- [ ] Job result storage
- [ ] Webhook notifications

### ⚡ Phase 4: Caching & Performance (Upcoming)
**Status:** 📅 Planned (Week 7-8)

- [ ] Redis caching layer
- [ ] Category result caching
- [ ] Database query optimization
- [ ] API response caching
- [ ] Rate limiting implementation
- [ ] Query result pagination
- [ ] Database connection pooling tuning

### 🧪 Phase 5: Testing & CI/CD (Upcoming)
**Status:** 📅 Planned (Week 9-10)

- [ ] Unit tests for all layers
- [ ] Integration tests
- [ ] End-to-end tests
- [ ] GitHub Actions CI/CD pipeline
- [ ] Code coverage reporting (target: 80%+)
- [ ] Security scanning (gosec)
- [ ] Dependency vulnerability scanning

### 🚀 Phase 6: Production Ready (Upcoming)
**Status:** 📅 Planned (Week 11-12)

- [ ] Prometheus metrics
- [ ] Structured logging
- [ ] Distributed tracing (OpenTelemetry)
- [ ] Kubernetes deployment manifests
- [ ] Horizontal Pod Autoscaling
- [ ] Database backup automation
- [ ] Disaster recovery procedures
- [ ] Documentation site

### 🔮 Future Enhancements

- [ ] Multi-language support (i18n)
- [ ] PDF statement templates for other providers
- [ ] Budget planning and alerts
- [ ] Spending predictions with ML
- [ ] Mobile app (Flutter)
- [ ] Batch processing for multiple statements
- [ ] Export to CSV, Excel, PDF
- [ ] Email reports
- [ ] Social features (anonymous spending comparisons)

---

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific package
go test ./handlers/...

# Verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Manual Testing with Postman

1. Download [Postman](https://www.postman.com/downloads/)
2. Import the collection from `postman_collection.json` (if available)
3. Set up environment variables:
   - `baseUrl`: `http://localhost:8080`
   - `token`: (auto-filled after login)

**Testing Workflow:**
1. Register → Get token
2. Login → Verify token
3. Upload → Check file validation
4. Health → Verify services

See [Postman Testing Guide](docs/postman-testing-guide.md) for detailed instructions.

---

## 🗄️ Database

### Schema Overview

The application uses PostgreSQL with the following main tables:

**users** - User accounts
- `id` (UUID) - Primary key
- `email` (VARCHAR) - Unique email
- `password_hash` (VARCHAR) - Bcrypt hash
- Timestamps

**jobs** - Processing jobs
- `id` (UUID) - Primary key
- `user_id` (UUID) - Foreign key to users
- `file_path` (VARCHAR) - File location
- `status` (ENUM) - queued, processing, completed, failed
- Timestamps

**transactions** - Parsed M-PESA transactions
- `id` (UUID) - Primary key
- `job_id` (UUID) - Foreign key to jobs
- `receipt_no` (VARCHAR) - M-PESA receipt
- `amount` (DECIMAL) - Transaction amount
- `category` (VARCHAR) - AI-assigned category
- Timestamps

### Migrations

```bash
# Create new migration
migrate create -ext sql -dir migrations -seq migration_name

# Apply migrations
make migrate-up

# Rollback
make migrate-down

# Check version
make migrate-version
```

---

## 🐳 Docker

### Development

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild
docker-compose up -d --build
```

### Production

```bash
# Build production image
docker build -t mpesa-analyzer:v1.0.0 .

# Run container
docker run -d \
  -p 8080:8080 \
  --env-file .env.production \
  mpesa-analyzer:v1.0.0

# With docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

---

## 🤝 Contributing

Contributions are what make the open-source community amazing! Any contributions you make are **greatly appreciated**.

### How to Contribute

1. **Fork the Project**
   ```bash
   # Click the "Fork" button on GitHub
   ```

2. **Clone your Fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/M-PESA-Statements-Analyzer.git
   cd M-PESA-Statements-Analyzer
   ```

3. **Create a Feature Branch**
   ```bash
   git checkout -b feature/AmazingFeature
   ```

4. **Make your Changes**
   - Write clean, documented code
   - Follow Go best practices
   - Add tests for new features
   - Update documentation

5. **Commit your Changes**
   ```bash
   git add .
   git commit -m "Add some AmazingFeature"
   ```

6. **Push to your Fork**
   ```bash
   git push origin feature/AmazingFeature
   ```

7. **Open a Pull Request**
   - Go to the original repository
   - Click "New Pull Request"
   - Select your feature branch
   - Describe your changes

### Contribution Guidelines

- **Code Style**: Follow Go conventions (run `gofmt`)
- **Testing**: Add tests for new features
- **Documentation**: Update README and code comments
- **Commits**: Use clear, descriptive commit messages
- **Issues**: Check existing issues before creating new ones

### Good First Issues

Look for issues tagged with `good first issue` or `help wanted`:
- Add input validation for edge cases
- Improve error messages
- Write documentation
- Add unit tests
- Fix typos

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

```
MIT License

Copyright (c) 2026 Wendyshiro


## 📞 Contact

**Wendyshiro**
- GitHub: [@Wendyshiro](https://github.com/Wendyshiro)
- Project Link: [https://github.com/Wendyshiro/M-PESA-Statements-Analyzer](https://github.com/Wendyshiro/M-PESA-Statements-Analyzer)
- Live Demo: Coming soon!!!

---

## 🙏 Acknowledgments

This project was built with the help of amazing open-source tools and resources:

### Technologies & Libraries
- [Go](https://go.dev/) - Programming language
- [PostgreSQL](https://www.postgresql.org/) - Database
- [Redis](https://redis.io/) - Caching
- [JWT-Go](https://github.com/golang-jwt/jwt) - Authentication
- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [Chart.js](https://www.chartjs.org/) - Data visualization
- [OpenAI](https://openai.com/) - AI categorization

### Learning Resources
- [Go by Example](https://gobyexample.com/)
- [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests/)
- [PostgreSQL Tutorial](https://www.postgresqltutorial.com/)
- [Docker Documentation](https://docs.docker.com/)

### Inspiration
- M-PESA for revolutionizing mobile money in Kenya and Africa
- The Go community for excellent documentation and support
- Open-source contributors worldwide

---
---

<div align="center">

### ⭐ Star this repository if you found it helpful!

Made with ❤️ by [Wendyshiro](https://github.com/Wendyshiro)

</div>