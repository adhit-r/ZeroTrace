# ZeroTrace API

The core backend API for the ZeroTrace vulnerability scanning platform.

## Features

- **Authentication**: JWT-based authentication with role-based access control
- **Scan Management**: Create, read, update, and delete vulnerability scans
- **Company Management**: Multi-tenant support with company isolation
- **Dashboard Data**: Real-time vulnerability statistics and trends
- **RESTful API**: Clean, documented REST endpoints
- **High Performance**: Optimized for handling 100,000+ data points

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Redis 7+

### Installation

1. **Clone and navigate to the API directory:**
   ```bash
   cd api-go
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Run the API:**
   ```bash
   go run cmd/api/main.go
   ```

The API will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### Scans
- `GET /api/v1/scans` - List scans (with pagination)
- `POST /api/v1/scans` - Create new scan
- `GET /api/v1/scans/:id` - Get scan details
- `PUT /api/v1/scans/:id` - Update scan
- `DELETE /api/v1/scans/:id` - Delete scan

### Companies
- `GET /api/v1/companies/:id` - Get company details
- `PUT /api/v1/companies/:id` - Update company

### Dashboard
- `GET /api/v1/dashboard/overview` - Dashboard overview
- `GET /api/v1/dashboard/trends` - Vulnerability trends

### Health Check
- `GET /health` - API health status

## Testing

### Test Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@zerotrace.com",
    "password": "password"
  }'
```

### Create Scan (with authentication)
```bash
curl -X POST http://localhost:8080/api/v1/scans \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "repository": "https://github.com/example/repo",
    "branch": "main",
    "scan_type": "full"
  }'
```

## Development

### Project Structure
```
api-go/
├── cmd/api/           # Application entry point
├── internal/          # Private application code
│   ├── config/        # Configuration management
│   ├── handlers/      # HTTP request handlers
│   ├── middleware/    # HTTP middleware
│   ├── models/        # Data models
│   ├── services/      # Business logic
│   └── repository/    # Data access layer (TODO)
├── pkg/               # Public packages
├── migrations/        # Database migrations (TODO)
└── tests/             # Test files
```

### Adding New Endpoints

1. **Create handler** in `internal/handlers/`
2. **Add service logic** in `internal/services/`
3. **Register route** in `cmd/api/main.go`
4. **Add tests** in `tests/`

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `API_PORT` | Server port | `8080` |
| `API_HOST` | Server host | `0.0.0.0` |
| `API_MODE` | Debug mode | `debug` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `zerotrace` |
| `JWT_SECRET` | JWT signing key | `dev-secret-key` |

## TODO

- [ ] Database integration with PostgreSQL
- [ ] Redis caching implementation
- [ ] Database migrations
- [ ] Comprehensive test suite
- [ ] API documentation with Swagger
- [ ] Rate limiting middleware
- [ ] Logging and monitoring
- [ ] Docker containerization
- [ ] CI/CD pipeline

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License
