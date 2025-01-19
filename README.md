# Go REST API

> A production-ready REST API built with Go, featuring comprehensive authentication, session management, and modern security practices.
> Features PostgreSQL integration, OAuth support, and robust logging with zero compromises on security and performance.

A robust Go-based REST API that provides secure authentication, session management, and user profile functionality with a focus on security and scalability.

## Features

- 🔐 Comprehensive authentication system
- 📝 Session-based user management
- 🚀 Rate limiting for security
- 🌐 CORS support
- 📚 Swagger UI documentation
- ✅ Request validation
- 📊 PostgreSQL with GORM
- 📝 Structured logging with Zap

## API Endpoints

- `POST /register`: Create new user accounts
- `POST /login`: Authenticate users and create sessions
- `POST /forgot-password`: Password reset functionality
- `POST /magic-link-login`: Passwordless authentication
- `GET /oauth/google`: Google OAuth integration
- `GET /logout`: Session termination
- `POST /invalidate-sessions`: Bulk session management
- `GET /profile`: User profile access

## Prerequisites

- Go 1.20+
- PostgreSQL 13+
- Docker (optional)
- Make (optional)

## Installation

### Via Docker

1. Clone the repository:
```bash
git clone https://github.com/paulgoerge35/go-rest-api.git
cd go-rest-api
```

2. Copy the example environment file:
```bash
cp .env.example .env
```

3. Start the services:
```bash
docker-compose up -d
```

### Manual Setup

1. Install dependencies:
```bash
go mod download
```

2. Set up the database:
```bash
make migrate
```

3. Start the server:
```bash
go run cmd/api/main.go
```

## Usage

### Authentication Flow

#### Standard Login
```http
POST /login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword"
}
```

#### Google OAuth
```http
GET /oauth/google
```

#### Magic Link Login
```http
POST /magic-link-login
Content-Type: application/json

{
  "email": "user@example.com"
}
```

### Session Management

#### Invalidate All Sessions
```http
POST /invalidate-sessions
Content-Type: application/json
```

## Database Schema

### User Table
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR UNIQUE,
  password_hash VARCHAR,
  name VARCHAR,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### Session Table
```sql
CREATE TABLE sessions (
  id UUID PRIMARY KEY,
  device_info VARCHAR,
  session_token VARCHAR UNIQUE,
  is_active BOOLEAN,
  user_id UUID REFERENCES users(id),
  last_accessed_at TIMESTAMP,
  expires_at TIMESTAMP
  created_at TIMESTAMP,
);
```

## Dependencies

##### Core Dependencies
- `gorm.io/gorm`: PostgreSQL ORM
- `go.uber.org/zap`: Structured logging
- `github.com/go-playground/validator/v10`: Request validation
- `golang.org/x/time/rate`: Rate limiting
- `github.com/swaggo/swag`: API documentation

##### Development Dependencies
- `github.com/golang-migrate/migrate/v4`: Database migrations
- `github.com/stretchr/testify`: Testing utilities
- `golang.org/x/tools`: Development tools

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

Paul George - contact@paulgeorge.dev

Project Link: [https://github.com/paulgeorge35/go-rest-api](https://github.com/paulgeorge35/go-rest-api)