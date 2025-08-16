# Golang Gin Boilerplate

A production-ready, feature-rich Go web application boilerplate built with [Gin](https://github.com/gin-gonic/gin) framework. This boilerplate provides a solid foundation for building scalable REST APIs with modern Go practices and best practices.

## 🚀 Features

- **🛠️ Gin Web Framework** - Fast HTTP web framework for Go
- **🗄️ PostgreSQL Database** - With GORM ORM and automatic migrations
- **🔐 JWT Authentication** - Built-in JWT-based authentication system
- **📚 Swagger Documentation** - Auto-generated API documentation
- **📝 Structured Logging** - With Zap logger and OpenTelemetry support
- **🧪 Testing Support** - Comprehensive testing utilities and examples
- **🐳 Docker Support** - Docker Compose for easy development setup
- **🔄 Database Migrations** - Automatic database schema management
- **🌍 CORS Support** - Cross-origin resource sharing enabled
- **📊 Health Checks** - Built-in health check endpoints
- **🔧 Environment Configuration** - Flexible environment-based configuration
- **⚡ Performance Optimized** - Gzip compression, profiling support
- **🛡️ Security Features** - Input validation, secure headers

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

- **Go 1.24.2+** - [Download Go](https://golang.org/dl/)
- **PostgreSQL 10.3+** - [Download PostgreSQL](https://www.postgresql.org/download/)
- **Docker & Docker Compose** (optional) - [Download Docker](https://www.docker.com/products/docker-desktop)
- **Make** - For using the provided Makefile commands
- **golang-migrate** - For database migrations

### Installing golang-migrate

#### Option 1: Using Go Install (Recommended)
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### Option 2: Using Homebrew (macOS)
```bash
brew install golang-migrate
```

#### Option 3: Using Binary Download
```bash
# For Linux/macOS
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.3/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# For Windows
# Download from https://github.com/golang-migrate/migrate/releases
```

#### Option 4: Using Docker
```bash
# Run migrate commands using Docker
docker run -v $(pwd)/assets/migrations:/migrations --network host migrate/migrate -path=/migrations -database "postgres://postgres:postgres@localhost:5432/app_db?sslmode=disable" up
```

## 🛠️ Installation & Setup

### 1. Clone the Repository

```bash
git clone https://github.com/amahdian/golang-gin-boilerplate.git
cd golang-gin-boilerplate
```

### 2. Environment Configuration

Create environment files for your configuration:

```bash
# Create base environment file
cp .env.example .env

# For development
cp .env.example .env.dev

# For testing
cp .env.example .env.test
```

Configure your environment variables in the `.env` file:

```env
# Server Configuration
GIN_MODE=debug
LOG_LEVEL=debug
LOG_FORMAT=text
HTTP_PORT=8090
SWAGGER_HOST_ADDR=localhost:8090
ASSETS_DIR=./assets
JWT_SECRET=your-super-secret-jwt-key

# Database Configuration
DB_DSN=postgres://postgres:postgres@localhost:5432/app_db?sslmode=disable
DB_LOG_LEVEL=error

# Profile (optional)
PROFILE=dev
```

### 3. Database Setup

#### Option A: Using Docker (Recommended for Development)

```bash
# Start PostgreSQL with Docker Compose
docker-compose up -d postgres

# Create database and run migrations
make create-db
make migrate-up
```

#### Option B: Using Local PostgreSQL

```bash
# Create database manually
createdb app_db

# Run migrations (requires golang-migrate to be installed)
make migrate-up
```

> **Note**: The migration commands (`make migrate-up`, `make migrate-down`, etc.) require `golang-migrate` to be installed. See the [Prerequisites](#-prerequisites) section for installation instructions.

### 4. Install Dependencies

```bash
# Install Go dependencies
make vendor
```

### 5. Generate Documentation

```bash
# Generate Swagger documentation
make docs
```

## 🚀 Running the Application

### Development Mode

```bash
# Run in development mode with hot reload
make dev
```

### Production Mode

```bash
# Build the application
make build

# Run the built binary
./build/app-bin
```

### Using Go Run

```bash
# Run directly with Go
make run
```

## 📚 Available Make Commands

The project includes a comprehensive Makefile with useful commands:

```bash
# Database Management
make create-db          # Create database if it doesn't exist
make drop-db            # Drop database if it exists
make migrate-up         # Apply all migrations
make migrate-down       # Rollback all migrations
make migrate-one-up     # Apply one migration
make migrate-one-down   # Rollback one migration
make new-migration name='migration_name'  # Create new migration

# Development
make vendor             # Tidy dependencies and update vendor
make docs               # Generate Swagger documentation
make build              # Build binary
make run                # Run main process
make dev                # Run with development setup
make test-all           # Run all tests with coverage
```

## 🏗️ Project Structure

```
golang-gin-boilerplate/
├── assets/                 # Static assets and migrations
│   └── migrations/         # Database migration files
├── docs/                   # Auto-generated Swagger documentation
├── domain/                 # Domain models and contracts
│   ├── contracts/          # Interface definitions
│   └── model/              # Data models
├── global/                 # Global configurations and utilities
│   ├── env/                # Environment configuration
│   ├── errs/               # Error definitions
│   └── test/               # Test utilities
├── pkg/                    # Reusable packages
│   ├── fileutil/           # File utilities
│   ├── logger/             # Logging utilities
│   └── msg/                # Message utilities
├── server/                 # HTTP server components
│   ├── binding/            # Request/response bindings
│   ├── middleware/         # HTTP middleware
│   ├── router/             # Route definitions
│   └── utils/              # Server utilities
├── storage/                # Data storage layer
│   └── pg/                 # PostgreSQL implementation
├── svc/                    # Business logic services
│   └── auth/               # Authentication service
├── testutil/               # Testing utilities
├── version/                # Version information
├── docker-compose.yaml     # Docker services configuration
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── main.go                 # Application entry point
└── Makefile                # Build and development commands
```

## 🔐 Authentication

The boilerplate includes JWT-based authentication:

- JWT tokens for API authentication
- Secure token validation
- User management system
- Role-based access control (ready to implement)

## 📖 API Documentation

Once the application is running, you can access:

- **Swagger UI**: `http://localhost:8090/swagger/index.html`
- **API Documentation**: `http://localhost:8090/swagger/doc.json`

## 🧪 Testing

```bash
# Run all tests with coverage
make test-all

# Run specific test files
go test ./pkg/...
go test ./svc/...
```

## 🐳 Docker Deployment

### Development with Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Production Docker Build

```bash
# Build production image
docker build -t golang-gin-boilerplate .

# Run container
docker run -p 8090:8090 golang-gin-boilerplate
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `GIN_MODE` | Gin framework mode | `debug` | No |
| `LOG_LEVEL` | Logging level | `debug` | No |
| `LOG_FORMAT` | Log format (text/json) | `text` | No |
| `HTTP_PORT` | Server port | `8090` | No |
| `SWAGGER_HOST_ADDR` | Swagger host address | - | No |
| `ASSETS_DIR` | Assets directory path | - | Yes |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `DB_DSN` | Database connection string | - | Yes |
| `DB_LOG_LEVEL` | Database log level | `error` | No |

### Profiles

The application supports multiple environment profiles:

- **Default**: Uses `.env` file
- **Development**: Uses `.env.dev` file (set `PROFILE=dev`)
- **Testing**: Uses `.env.test` file (set `PROFILE=test`)

## 📊 Monitoring & Observability

- **Structured Logging**: JSON and text formats with Zap logger
- **OpenTelemetry**: Distributed tracing and metrics
- **Health Checks**: Built-in health check endpoints
- **Profiling**: Pprof integration for performance analysis

## 🔒 Security Features

- JWT token validation
- CORS configuration
- Input validation with Gin validator
- Secure headers middleware
- SQL injection prevention with GORM

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues or have questions:

1. Check the [Issues](https://github.com/amahdian/golang-gin-boilerplate/issues) page
2. Create a new issue with detailed information
3. Review the documentation and examples

## 🙏 Acknowledgments

- [Gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [GORM](https://gorm.io/) - ORM library
- [Zap](https://github.com/uber-go/zap) - Logging library
- [Swagger](https://swagger.io/) - API documentation
- [OpenTelemetry](https://opentelemetry.io/) - Observability framework

---

**Happy Coding! 🚀**
