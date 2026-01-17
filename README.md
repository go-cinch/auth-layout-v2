# Go-Cinch Auth Service Template

A production-ready authentication and authorization microservice template generator based on [Kratos](https://github.com/go-kratos/kratos) framework with Wire dependency injection.

## Features

- 🔐 **JWT Authentication**: Token-based authentication with signature verification
- 👥 **RBAC Authorization**: Role-Based Access Control (User/Role/Permission)
- 🎯 **Flexible Authorization Models**: Optional fine-grained Action-based permissions
- 🔧 **Pluggable Features**: Multiple presets and customizable auth modules
- 💉 **Dependency Injection**: Wire-based automatic code generation
- 🗄️ **Database Support**: PostgreSQL/MySQL with GORM ORM
- 🚀 **Redis Caching**: Multi-layer cache with hotspot optimization (optional)
- 🛡️ **Security Features**: Captcha, user lock, whitelist, password hashing (bcrypt)
- 📊 **Observability**: OpenTelemetry tracing support
- 🚦 **Production Ready**: Health checks, middleware, i18n, task scheduling

## Quick Start

### 1. Install scaffold

```bash
go install github.com/hay-kot/scaffold@v0.12.0
```

### 2. Create Auth Service

#### Using Presets (Recommended)

Presets provide pre-configured auth service settings:

##### Simple Preset - Core RBAC

Simple RBAC with User/Role/Permission and Redis caching:

```bash
scaffold new https://github.com/go-cinch/auth-layout-v2 \
  --output-dir=. \
  --run-hooks=always \
  --no-prompt \
  --preset simple \
  Project=myauth
```

**Features:**
- ✅ JWT authentication (HS512 signing)
- ✅ User CRUD with password management (bcrypt)
- ✅ Role management
- ✅ Permission management
- ✅ Basic whitelist (URI + JWT whitelist)
- ✅ Redis caching
- ✅ Permission middleware
- ✅ Idempotent middleware
- ✅ OpenTelemetry tracing
- ✅ Health check endpoints
- ❌ No Action (fine-grained permissions)
- ❌ No UserGroup
- ❌ No Hotspot cache optimization
- ❌ No Captcha/User lock
- 📦 Database tables: 4 (user, role, permission, whitelist)

**Generated Structure:**
```
myauth/
├── api/auth-proto/      # Auth service protobuf definitions
├── cmd/myauth/          # Application entry point
├── configs/             # Configuration files (db, redis, log, etc.)
├── internal/
│   ├── biz/            # Business logic (User/Role/Permission/Auth)
│   ├── data/           # Data access layer with GORM
│   ├── server/         # HTTP/gRPC servers + middleware
│   └── service/        # Service implementations
└── Makefile            # Build automation
```

##### Default Preset - Enterprise Auth

Enterprise-grade auth service with all advanced features:

```bash
scaffold new https://github.com/go-cinch/auth-layout-v2 \
  --output-dir=. \
  --run-hooks=always \
  --no-prompt \
  --preset default \
  Project=myauth
```

**Features (Simple +):**
- ✅ All Simple preset features
- ✅ Action module (fine-grained resource permissions)
- ✅ UserGroup module (group-based user management)
- ✅ Enhanced whitelist (permission + JWT + category matching)
- ✅ Hotspot cache (in-memory hot data with go-cache)
- ✅ Captcha (login verification code with Redis)
- ✅ User lock mechanism (auto-lock after 3 wrong passwords)
- ✅ i18n support for error messages
- 📦 Database tables: 7 (+ action, user_group, extended whitelist)

### 3. Presets Comparison

| Feature | Simple | Default |
|---------|----------|------|
| **Authentication** | ✅ JWT (HS512) | ✅ |
| **User Management** | ✅ CRUD + Password | ✅ |
| **RBAC** | ✅ User/Role/Permission | ✅ |
| **Fine-grained Permissions** | ❌ | ✅ Action module |
| **User Groups** | ❌ | ✅ UserGroup |
| **Whitelist** | ✅ Basic | ✅ Enhanced |
| **Caching** | ✅ Redis basic | ✅ + Hotspot optimization |
| **Security Features** | ❌ | ✅ Captcha + User Lock |
| **Permission Middleware** | ✅ | ✅ |
| **Idempotent Middleware** | ✅ | ✅ |
| **OpenTelemetry Tracing** | ✅ | ✅ |
| **i18n Support** | ❌ | ✅ |
| **Database Tables** | 4 | 7 |
| **Use Case** | Simple apps | Enterprise/Multi-tenant |

### 4. Customization

#### Override Preset Options

You can override specific preset options:

```bash
# Enable Captcha on simple preset
scaffold new https://github.com/go-cinch/auth-layout-v2 \
  --preset simple \
  Project=myauth \
  enable_captcha=true

# Use MySQL instead of PostgreSQL
scaffold new https://github.com/go-cinch/auth-layout-v2 \
  --preset simple \
  Project=myauth \
  db_type=mysql

# Change HTTP/gRPC ports
scaffold new https://github.com/go-cinch/auth-layout-v2 \
  --preset default \
  Project=myauth \
  http_port=9090 \
  grpc_port=9190

# Enable Action module without default preset
scaffold new https://github.com/go-cinch/auth-layout-v2 \
  --preset simple \
  Project=myauth \
  enable_action=true
```

#### Available Auth-Specific Options

| Option | Values | Default | Description |
|--------|--------|---------|-------------|
| `enable_action` | `true`/`false` | simple=false, default=true | Fine-grained resource permissions |
| `enable_user_group` | `true`/`false` | simple=false, default=true | Group-based user management |
| `enable_whitelist` | `true`/`false` | both=true | Permission and JWT whitelist |
| `enable_hotspot` | `true`/`false` | simple=false, default=true | Hotspot cache optimization |
| `enable_captcha` | `true`/`false` | simple=false, default=true | Captcha verification code |
| `enable_user_lock` | `true`/`false` | simple=false, default=true | Lock user after wrong passwords |

#### General Options

| Option | Values | Default | Description |
|--------|--------|---------|-------------|
| `service_name` | string | Project name | Service name |
| `module_name` | string | service_name | Go module name |
| `http_port` | string | `8080` | HTTP server port |
| `grpc_port` | string | `8180` | gRPC server port |
| `db_type` | `postgres`/`mysql` | `postgres` | Database type |
| `enable_redis` | `true`/`false` | `true` | Enable Redis (required for auth) |
| `enable_cache` | `true`/`false` | `true` | Enable cache layer |
| `enable_idempotent` | `true`/`false` | `true` | Enable idempotent middleware |
| `enable_trace` | `true`/`false` | `true` | Enable OpenTelemetry tracing |
| `enable_ws` | `true`/`false` | `false` | Enable WebSocket support |
| `enable_task` | `true`/`false` | `false` | Enable task/cron scheduler |
| `enable_i18n` | `true`/`false` | default=true | Enable i18n support |

### 5. Interactive Mode

Answer prompts to configure all options:

```bash
scaffold new https://github.com/go-cinch/auth-layout-v2 \
  --output-dir=. \
  --run-hooks=always \
  Project=myauth
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

**Configure in `configs/client.yaml`:**
```yaml
server:
  jwt:
    key: "your-secret-key-min-32-chars-long"  # HS512 signing key
    expires: 7200  # Token expiration in seconds (2 hours)
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
- Health: http://localhost:8080/health
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
curl -X POST http://localhost:8080/api/v1/auth/login \
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
  "expires": "2025-01-16 13:00:00",
  "wrong": 0
}
```

### 2. Authenticated Request

**Request:**
```bash
curl -X GET http://localhost:8080/api/v1/user/1 \
  -H "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9..."
```

### 3. Permission Check

The permission middleware automatically:
1. Verifies JWT token signature
2. Checks if URI is in whitelist (public endpoints)
3. Matches HTTP method + URI against user's permissions
4. Returns 401/403 if unauthorized

## Authorization Models

### Basic RBAC (Simple Preset)

```
User ──belongsTo──> Role ──has──> Permission
```

**Example:**
- User: `john@example.com`
- Role: `editor`
- Permissions: `POST|/api/v1/article`, `GET|/api/v1/article/*`

### Extended RBAC (Default Preset)

```
User ──belongsTo──> Role ──has──> Permission
User ──belongsTo──> UserGroup
Permission ──uses──> Action (fine-grained resource control)
```

**Example:**
- User: `john@example.com`
- Role: `editor`
- UserGroup: `content-team`
- Actions: `article:create`, `article:edit`, `article:delete:own`

## Security Best Practices

### Password Handling

- ✅ Passwords hashed with bcrypt (cost factor: 10)
- ✅ Password comparison cached in Redis (short TTL)
- ✅ Never return password in API responses

### JWT Security

- ✅ HS512 signing algorithm (strong symmetric signing)
- ✅ Token cached in Redis with 10-minute expiration
- ✅ Token contains minimal claims (code, platform)
- ⚠️ Use environment variables for JWT secret key

### User Lock Mechanism (Default Preset)

- ✅ Auto-lock after 3 wrong password attempts
- ✅ Configurable lock duration
- ✅ Automatic unlock on expiration

### Captcha Protection (Default Preset)

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

### Simple Preset Tables

```sql
-- User table
CREATE TABLE user (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  code VARCHAR(50) UNIQUE,
  platform VARCHAR(50),
  role_id BIGINT REFERENCES role(id),
  locked BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Role table
CREATE TABLE role (
  id BIGSERIAL PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL,
  code VARCHAR(50) UNIQUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Permission table
CREATE TABLE permission (
  id BIGSERIAL PRIMARY KEY,
  role_id BIGINT REFERENCES role(id),
  resource VARCHAR(255) NOT NULL,  -- e.g., "POST|/api/v1/user"
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- Whitelist table
CREATE TABLE whitelist (
  id BIGSERIAL PRIMARY KEY,
  category VARCHAR(50),  -- "permission" or "jwt"
  resource VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);
```

### Default Preset Additional Tables

```sql
-- Action table (fine-grained permissions)
CREATE TABLE action (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(100) UNIQUE NOT NULL,  -- e.g., "article:create"
  name VARCHAR(255),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- UserGroup table
CREATE TABLE user_group (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT REFERENCES user(id),
  group_code VARCHAR(100),
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
# Check token expiration in configs/client.yaml
```

## Contributing

Issues and pull requests are welcome!

## License

Apache 2.0
