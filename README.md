# Go-Cinch Auth Service Template

A production-ready authentication and authorization microservice template generator based on [Kratos](https://github.com/go-kratos/kratos) framework with Wire dependency injection.

## Features

- 🔐 **JWT Authentication**: Token-based authentication with signature verification
- 👥 **RBAC Authorization**: Role-Based Access Control (User/Role/Permission)
- 🎯 **Fine-grained Permissions**: Action-based resource permissions
- 👨‍👩‍👧‍👦 **User Groups**: Group-based user management
- 💉 **Dependency Injection**: Wire-based automatic code generation
- 🗄️ **Database Support**: PostgreSQL/MySQL with GORM ORM
- 🚀 **Redis Caching**: Multi-layer cache with hotspot optimization
- 🛡️ **Security Features**: Captcha, user lock, whitelist, password hashing (bcrypt)
- 📊 **Observability**: OpenTelemetry tracing support
- 🚦 **Production Ready**: Health checks, middleware, i18n, task scheduling

## Quick Start

### 1. Install scaffold

```bash
go install github.com/hay-kot/scaffold@v0.12.0
```

### 2. Create Auth Service

#### Using Default Preset (Recommended)

Enterprise-grade auth service with all features:

```bash
scaffold new https://github.com/go-cinch/auth-layout-v2 \
  --output-dir=. \
  --run-hooks=always \
  --no-prompt \
  --preset default \
  Project=myauth
```

**Features:**
- ✅ JWT authentication (HS512 signing)
- ✅ User CRUD with password management (bcrypt)
- ✅ Role management
- ✅ Permission management
- ✅ Action module (fine-grained resource permissions)
- ✅ UserGroup module (group-based user management)
- ✅ Enhanced whitelist (permission + JWT + category matching)
- ✅ Hotspot cache (in-memory hot data with go-cache)
- ✅ Captcha (login verification code with Redis)
- ✅ User lock mechanism (auto-lock after 3 wrong passwords)
- ✅ Redis caching
- ✅ Permission middleware
- ✅ Idempotent middleware
- ✅ OpenTelemetry tracing
- ✅ Health check endpoints
- ✅ i18n support for error messages
- 📦 Database tables: 7 (user, role, action, user_group, user_user_group_relation, whitelist, schema_migrations)

**Generated Structure:**
```
myauth/
├── api/auth-proto/      # Auth service protobuf definitions
├── cmd/myauth/          # Application entry point
├── configs/             # Configuration files (db, redis, log, etc.)
├── internal/
│   ├── biz/            # Business logic (User/Role/Action/Permission/Auth)
│   ├── data/           # Data access layer with GORM
│   ├── server/         # HTTP/gRPC servers + middleware
│   └── service/        # Service implementations
└── Makefile            # Build automation
```

## Building and Running

### 1. Generate Code

```bash
cd myauth
make all  # Install tools, generate proto/wire, lint
```

### 2. Database Setup

**Start PostgreSQL (Docker):**
```bash
docker run -d --name postgres \
  -e POSTGRES_USER=root \
  -e POSTGRES_PASSWORD=password \
  -p 5432:5432 \
  postgres:17
```

**Configure in `configs/db.yaml`:**
```yaml
db:
  driver: postgres
  dsn: "host=localhost user=root password=password dbname=myauth port=5432 sslmode=disable TimeZone=UTC"
  migrate: true  # Auto-run migrations
```

### 3. Redis Setup

```bash
docker run -d --name redis \
  -p 6379:6379 \
  redis:7
```

**Configure in `configs/redis.yaml`:**
```yaml
redis:
  dsn: "redis://:password@localhost:6379/0"
```

### 4. JWT Configuration

**Configure in `configs/server.yaml`:**
```yaml
server:
  jwt:
    key: "your-secret-key-min-32-chars-long"  # HS512 signing key
    expires: "24h"  # Token expiration duration
```

### 5. Build

```bash
make build  # Output: ./bin/myauth
```

### 6. Run

```bash
./bin/myauth -c ./configs
```

**Endpoints:**
- HTTP: http://localhost:8080
- gRPC: localhost:8180
- Health: http://localhost:8080/healthz
- Swagger: http://localhost:8080/docs (if enabled)

## Development Workflow

```bash
# Generate API from proto
make api

# Generate Wire dependency injection
make wire

# Run linter
make lint

# Run tests
make test

# Generate GORM models from database
make gen-model

# Complete build pipeline
make all

# Clean generated files
make clean
```

## Authentication Flow

### 1. User Login

**Request:**
```bash
curl -X POST http://localhost:8080/pub/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "password123",
    "platform": "web"
  }'
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9...",
  "expires": "2025-01-16 13:00:00"
}
```

### 2. Authenticated Request

**Request:**
```bash
curl -X GET http://localhost:8080/user/list \
  -H "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9..."
```

### 3. Permission Check

The permission middleware automatically:
1. Verifies JWT token signature
2. Checks if URI is in whitelist (public endpoints)
3. Matches HTTP method + URI against user's permissions
4. Returns 401/403 if unauthorized

## Authorization Models

### Extended RBAC

```
User ──belongsTo──> Role ──has──> Action (permissions)
User ──belongsTo──> UserGroup ──has──> Action (permissions)
User ──has──> Action (direct permissions)
```

**Example:**
- User: `john@example.com`
- Role: `editor`
- UserGroup: `content-team`
- Actions: `article:create`, `article:edit`, `article:delete:own`

## Security Best Practices

### Password Handling

- ✅ Passwords hashed with bcrypt (cost factor: default)
- ✅ Password comparison cached in Redis (short TTL)
- ✅ Never return password in API responses

### JWT Security

- ✅ HS512 signing algorithm (strong symmetric signing)
- ✅ Token cached in Redis with configurable expiration
- ✅ Token contains minimal claims (code, platform)
- ⚠️ Use environment variables for JWT secret key

### User Lock Mechanism

- ✅ Auto-lock after 3 wrong password attempts
- ✅ Configurable lock duration
- ✅ Automatic unlock on expiration

### Captcha Protection

- ✅ Required after 3 wrong password attempts
- ✅ Stored in Redis with expiration
- ✅ One-time use verification

## Project Naming Guidelines

### ⚠️ Important Naming Rules

**Avoid project names ending with 's'** (e.g., `auths`, `users`, `admins`)

**Why?**
GORM's gentool singularizes table names incorrectly:
- `auths` → generates `Auth` ❌ (expected `Auths`)
- `users` → generates `User` ❌ (expected `Users`)

**Recommended Naming:**
- ✅ `auth`, `user`, `admin`
- ❌ `auths`, `users`, `admins` (will cause compilation errors)

**If you must use plural names:**
Manually correct the generated model struct names in `internal/data/model/*.gen.go` after generation.

## Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────┐
│          API (Proto/HTTP)              │
├─────────────────────────────────────────┤
│         Service Layer                   │  ← gRPC/HTTP handlers
├─────────────────────────────────────────┤
│         Business Logic (Biz)            │  ← UseCase, domain logic
├─────────────────────────────────────────┤
│         Data Access (Data)              │  ← Repository implementation
├─────────────────────────────────────────┤
│      Database (GORM) + Redis            │  ← Persistence
└─────────────────────────────────────────┘
```

### Dependency Injection (Wire)

```
cmd/main.go
  └─> wire.go (injector)
       ├─> server.ProviderSet (HTTP/gRPC servers)
       ├─> service.ProviderSet (Service layer)
       ├─> biz.ProviderSet (Business logic)
       └─> data.ProviderSet (Data access)
```

## Database Schema

### Core Tables

```sql
-- User table
CREATE TABLE user (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  code VARCHAR(50) UNIQUE,
  platform VARCHAR(50),
  action TEXT,  -- comma-separated action codes
  role_id BIGINT REFERENCES role(id),
  locked SMALLINT DEFAULT 0,
  lock_expire BIGINT,
  wrong BIGINT DEFAULT 0,
  last_login TIMESTAMP,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Role table
CREATE TABLE role (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  word VARCHAR(50) UNIQUE,
  action TEXT,  -- comma-separated action codes
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Action table (fine-grained permissions)
CREATE TABLE action (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(100) UNIQUE NOT NULL,
  name VARCHAR(255),
  word VARCHAR(100) UNIQUE,
  resource TEXT,  -- permission rules (e.g., "GET,POST|/api/v1/user/*")
  menu TEXT,
  btn TEXT,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- UserGroup table
CREATE TABLE user_group (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  word VARCHAR(100) UNIQUE,
  action TEXT,  -- comma-separated action codes
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- User-UserGroup relation table
CREATE TABLE user_user_group_relation (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES user(id),
  user_group_id BIGINT REFERENCES user_group(id),
  UNIQUE(user_id, user_group_id)
);

-- Whitelist table
CREATE TABLE whitelist (
  id BIGSERIAL PRIMARY KEY,
  category SMALLINT,  -- 0: permission whitelist, 1: JWT whitelist
  resource VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

## Troubleshooting

### Common Issues

**1. Wire generation fails**
```bash
# Install Wire tool
make init
# Regenerate
make wire
```

**2. Proto compilation fails**
```bash
# Install protoc plugins
make init
# Regenerate
make api
```

**3. Database migration fails**
```bash
# Check database connection in configs/db.yaml
# Manually run migrations
make migrate
```

**4. JWT token invalid**
```bash
# Ensure JWT key is at least 32 characters
# Check token expiration in configs/server.yaml
```

## Contributing

Issues and pull requests are welcome!

## License

Apache 2.0
