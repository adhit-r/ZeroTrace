# ZeroTrace Go API

## Overview
The ZeroTrace Go API serves as the central orchestration layer, managing authentication, scan jobs, data aggregation, and real-time communication. It's built with high performance and scalability in mind to handle 100,000+ data points across multiple companies and agents.

## Architecture

### API Structure
```
/api-go
  /cmd
    /api                 # Main application entry point
      main.go
  /internal
    /handlers            # HTTP request handlers
      /auth              # Authentication handlers
      /scans             # Scan management handlers
      /reports           # Report handlers
      /agents            # Agent management handlers
      /companies         # Company management handlers
      /users             # User management handlers
    /middleware          # HTTP middleware
      /auth              # Authentication middleware
      /logging           # Request logging
      /cors              # CORS handling
      /rate_limit        # Rate limiting
      /validation        # Request validation
    /services            # Business logic services
      /auth              # Authentication service
      /scan              # Scan management service
      /enrichment        # Data enrichment service
      /notification      # Notification service
      /report            # Report generation service
    /repository          # Data access layer
      /database          # Database operations
      /cache             # Cache operations
      /queue             # Message queue operations
    /models              # Data models
      /entities          # Database entities
      /dto               # Data transfer objects
      /requests          # Request models
      /responses         # Response models
    /config              # Configuration
      /env               # Environment variables
      /database          # Database configuration
      /redis             # Redis configuration
    /utils               # Utilities
      /crypto            # Cryptographic functions
      /validation        # Validation utilities
      /logging           # Logging utilities
      /pagination        # Pagination utilities
  /pkg
    /constants           # Application constants
    /types               # Type definitions
    /errors              # Error definitions
  /migrations            # Database migrations
  /scripts               # Build and deployment scripts
  /tests                 # Test files
  go.mod
  go.sum
  Dockerfile
  docker-compose.yml
```

## Core Features

### 1. Authentication & Authorization
- JWT-based authentication
- Role-based access control (RBAC)
- Multi-tenant data isolation
- API key management for agents
- Session management

### 2. Scan Management
- Scan job creation and scheduling
- Real-time scan status tracking
- Agent assignment and load balancing
- Scan result aggregation
- Historical scan data

### 3. Data Management
- Multi-tenant data isolation
- Efficient data storage and retrieval
- Caching strategies
- Data validation and sanitization
- Audit logging

### 4. Real-time Communication
- WebSocket connections for real-time updates
- Event-driven architecture
- Message queuing
- Push notifications

## Implementation Details

### 1. HTTP Server Setup

#### Main Application
```go
func main() {
    // Load configuration
    config := config.Load()
    
    // Initialize database
    db := database.New(config.Database)
    
    // Initialize Redis
    redis := cache.New(config.Redis)
    
    // Initialize services
    authService := auth.NewService(db, redis)
    scanService := scan.NewService(db, redis)
    
    // Setup router
    router := gin.New()
    router.Use(middleware.Logger(), middleware.Recovery())
    
    // Setup routes
    setupRoutes(router, authService, scanService)
    
    // Start server
    server := &http.Server{
        Addr:    config.Server.Address,
        Handler: router,
    }
    
    log.Fatal(server.ListenAndServe())
}
```

#### Route Setup
```go
func setupRoutes(router *gin.Engine, authService *auth.Service, scanService *scan.Service) {
    // Public routes
    public := router.Group("/api/v1")
    {
        public.POST("/auth/login", handlers.Login(authService))
        public.POST("/auth/register", handlers.Register(authService))
    }
    
    // Protected routes
    protected := router.Group("/api/v1")
    protected.Use(middleware.Auth(authService))
    {
        // Scan routes
        scans := protected.Group("/scans")
        {
            scans.GET("/", handlers.GetScans(scanService))
            scans.POST("/", handlers.CreateScan(scanService))
            scans.GET("/:id", handlers.GetScan(scanService))
            scans.PUT("/:id", handlers.UpdateScan(scanService))
            scans.DELETE("/:id", handlers.DeleteScan(scanService))
        }
        
        // Report routes
        reports := protected.Group("/reports")
        {
            reports.GET("/", handlers.GetReports(scanService))
            reports.GET("/:id", handlers.GetReport(scanService))
            reports.POST("/:id/export", handlers.ExportReport(scanService))
        }
    }
}
```

### 2. Authentication Service

#### JWT Authentication
```go
type AuthService struct {
    db    *gorm.DB
    redis *redis.Client
    config *config.AuthConfig
}

func (s *AuthService) Login(credentials *LoginCredentials) (*AuthResponse, error) {
    // Validate credentials
    user, err := s.validateCredentials(credentials)
    if err != nil {
        return nil, err
    }
    
    // Generate JWT token
    token, err := s.generateToken(user)
    if err != nil {
        return nil, err
    }
    
    // Store session
    err = s.storeSession(user.ID, token)
    if err != nil {
        return nil, err
    }
    
    return &AuthResponse{
        Token: token,
        User:  user,
    }, nil
}

func (s *AuthService) generateToken(user *User) (string, error) {
    claims := jwt.MapClaims{
        "user_id":    user.ID,
        "company_id": user.CompanyID,
        "role":       user.Role,
        "exp":        time.Now().Add(time.Hour * 24).Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.config.Secret))
}
```

### 3. Scan Management Service

#### Scan Service
```go
type ScanService struct {
    db          *gorm.DB
    redis       *redis.Client
    enrichment  *enrichment.Service
    notification *notification.Service
}

func (s *ScanService) CreateScan(scanData *CreateScanRequest) (*Scan, error) {
    // Validate scan data
    if err := s.validateScanData(scanData); err != nil {
        return nil, err
    }
    
    // Create scan record
    scan := &Scan{
        CompanyID:   scanData.CompanyID,
        Repository:  scanData.Repository,
        Branch:      scanData.Branch,
        Status:      "pending",
        CreatedAt:   time.Now(),
    }
    
    if err := s.db.Create(scan).Error; err != nil {
        return nil, err
    }
    
    // Assign to available agent
    agent, err := s.assignAgent(scan.CompanyID)
    if err != nil {
        return nil, err
    }
    
    // Send scan job to agent
    err = s.sendScanJob(scan, agent)
    if err != nil {
        return nil, err
    }
    
    // Update scan status
    scan.Status = "assigned"
    scan.AgentID = agent.ID
    s.db.Save(scan)
    
    return scan, nil
}

func (s *ScanService) GetScans(filters *ScanFilters) (*PaginatedResponse, error) {
    var scans []Scan
    query := s.db.Model(&Scan{})
    
    // Apply filters
    if filters.CompanyID != "" {
        query = query.Where("company_id = ?", filters.CompanyID)
    }
    if filters.Status != "" {
        query = query.Where("status = ?", filters.Status)
    }
    if filters.AgentID != "" {
        query = query.Where("agent_id = ?", filters.AgentID)
    }
    
    // Apply pagination
    offset := (filters.Page - 1) * filters.Limit
    query = query.Offset(offset).Limit(filters.Limit)
    
    // Execute query
    if err := query.Find(&scans).Error; err != nil {
        return nil, err
    }
    
    // Get total count
    var total int64
    s.db.Model(&Scan{}).Count(&total)
    
    return &PaginatedResponse{
        Data:  scans,
        Total: total,
        Page:  filters.Page,
        Limit: filters.Limit,
    }, nil
}
```

### 4. Data Models

#### Database Entities
```go
type User struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    Email     string    `json:"email" gorm:"uniqueIndex"`
    Password  string    `json:"-" gorm:"not null"`
    Name      string    `json:"name"`
    Role      UserRole  `json:"role"`
    CompanyID string    `json:"company_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type Company struct {
    ID          string    `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name"`
    Domain      string    `json:"domain"`
    Settings    JSON      `json:"settings"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type Scan struct {
    ID             string                 `json:"id" gorm:"primaryKey"`
    CompanyID      string                 `json:"company_id"`
    AgentID        string                 `json:"agent_id"`
    Repository     string                 `json:"repository"`
    Branch         string                 `json:"branch"`
    Commit         string                 `json:"commit"`
    Status         ScanStatus             `json:"status"`
    Progress       int                    `json:"progress"`
    StartTime      *time.Time             `json:"start_time"`
    EndTime        *time.Time             `json:"end_time"`
    Results        JSON                   `json:"results"`
    Metadata       map[string]interface{} `json:"metadata"`
    CreatedAt      time.Time              `json:"created_at"`
    UpdatedAt      time.Time              `json:"updated_at"`
}

type Agent struct {
    ID        string    `json:"id" gorm:"primaryKey"`
    CompanyID string    `json:"company_id"`
    Name      string    `json:"name"`
    Status    string    `json:"status"`
    Version   string    `json:"version"`
    LastSeen  time.Time `json:"last_seen"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 5. Middleware

#### Authentication Middleware
```go
func Auth(authService *auth.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
            c.Abort()
            return
        }
        
        // Remove "Bearer " prefix
        if strings.HasPrefix(token, "Bearer ") {
            token = token[7:]
        }
        
        // Validate token
        claims, err := authService.ValidateToken(token)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        
        // Set user context
        c.Set("user_id", claims["user_id"])
        c.Set("company_id", claims["company_id"])
        c.Set("role", claims["role"])
        
        c.Next()
    }
}
```

#### Rate Limiting Middleware
```go
func RateLimit(redis *redis.Client) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
        
        // Check current count
        count, err := redis.Get(c, key).Int()
        if err != nil && err != redis.Nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Rate limit error"})
            c.Abort()
            return
        }
        
        if count >= 100 { // 100 requests per minute
            c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        
        // Increment counter
        pipe := redis.Pipeline()
        pipe.Incr(c, key)
        pipe.Expire(c, key, time.Minute)
        pipe.Exec(c)
        
        c.Next()
    }
}
```

### 6. WebSocket Support

#### WebSocket Handler
```go
func WebSocketHandler(hub *Hub) gin.HandlerFunc {
    return func(c *gin.Context) {
        conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
        if err != nil {
            return
        }
        
        // Get user from context
        userID := c.GetString("user_id")
        companyID := c.GetString("company_id")
        
        client := &Client{
            Hub:      hub,
            Conn:     conn,
            Send:     make(chan []byte, 256),
            UserID:   userID,
            CompanyID: companyID,
        }
        
        client.Hub.Register <- client
        
        go client.writePump()
        go client.readPump()
    }
}
```

### 7. Database Operations

#### Repository Pattern
```go
type ScanRepository struct {
    db *gorm.DB
}

func (r *ScanRepository) Create(scan *Scan) error {
    return r.db.Create(scan).Error
}

func (r *ScanRepository) GetByID(id string) (*Scan, error) {
    var scan Scan
    err := r.db.Where("id = ?", id).First(&scan).Error
    if err != nil {
        return nil, err
    }
    return &scan, nil
}

func (r *ScanRepository) GetByCompany(companyID string, filters *ScanFilters) ([]Scan, error) {
    var scans []Scan
    query := r.db.Where("company_id = ?", companyID)
    
    if filters.Status != "" {
        query = query.Where("status = ?", filters.Status)
    }
    
    err := query.Find(&scans).Error
    return scans, err
}

func (r *ScanRepository) Update(scan *Scan) error {
    return r.db.Save(scan).Error
}
```

## Performance Optimizations

### 1. Database Optimization
```go
// Connection pooling
func NewDatabase(config *config.DatabaseConfig) *gorm.DB {
    db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatal(err)
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal(err)
    }
    
    // Configure connection pool
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db
}
```

### 2. Caching Strategy
```go
func (s *ScanService) GetScan(id string) (*Scan, error) {
    // Try cache first
    cacheKey := fmt.Sprintf("scan:%s", id)
    cached, err := s.redis.Get(context.Background(), cacheKey).Result()
    if err == nil {
        var scan Scan
        json.Unmarshal([]byte(cached), &scan)
        return &scan, nil
    }
    
    // Get from database
    scan, err := s.repository.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // Cache result
    scanJSON, _ := json.Marshal(scan)
    s.redis.Set(context.Background(), cacheKey, scanJSON, time.Hour)
    
    return scan, nil
}
```

### 3. Pagination
```go
type PaginationParams struct {
    Page  int `form:"page" binding:"min=1"`
    Limit int `form:"limit" binding:"min=1,max=100"`
}

func (r *ScanRepository) GetPaginated(params *PaginationParams, filters *ScanFilters) (*PaginatedResponse, error) {
    var scans []Scan
    var total int64
    
    query := r.db.Model(&Scan{})
    
    // Apply filters
    if filters.CompanyID != "" {
        query = query.Where("company_id = ?", filters.CompanyID)
    }
    
    // Get total count
    query.Count(&total)
    
    // Get paginated results
    offset := (params.Page - 1) * params.Limit
    err := query.Offset(offset).Limit(params.Limit).Find(&scans).Error
    
    return &PaginatedResponse{
        Data:  scans,
        Total: total,
        Page:  params.Page,
        Limit: params.Limit,
    }, err
}
```

## Security Features

### 1. Input Validation
```go
type CreateScanRequest struct {
    Repository string `json:"repository" binding:"required,url"`
    Branch     string `json:"branch" binding:"required"`
    CompanyID  string `json:"company_id" binding:"required,uuid"`
}

func (s *ScanService) validateScanData(data *CreateScanRequest) error {
    if err := validator.New().Struct(data); err != nil {
        return err
    }
    
    // Additional business logic validation
    if !s.isValidRepository(data.Repository) {
        return errors.New("invalid repository URL")
    }
    
    return nil
}
```

### 2. Data Sanitization
```go
func sanitizeInput(input string) string {
    // Remove potentially dangerous characters
    input = strings.ReplaceAll(input, "<script>", "")
    input = strings.ReplaceAll(input, "</script>", "")
    input = strings.ReplaceAll(input, "javascript:", "")
    
    return input
}
```

## Testing Strategy

### 1. Unit Tests
```go
func TestScanService_CreateScan(t *testing.T) {
    // Setup
    db := setupTestDB()
    redis := setupTestRedis()
    service := scan.NewService(db, redis)
    
    // Test data
    scanData := &CreateScanRequest{
        Repository: "https://github.com/example/repo",
        Branch:     "main",
        CompanyID:  "company-123",
    }
    
    // Execute
    result, err := service.CreateScan(scanData)
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "pending", result.Status)
}
```

### 2. Integration Tests
```go
func TestScanAPI_CreateScan(t *testing.T) {
    // Setup
    router := setupTestRouter()
    
    // Test request
    req := httptest.NewRequest("POST", "/api/v1/scans", strings.NewReader(`{
        "repository": "https://github.com/example/repo",
        "branch": "main"
    }`))
    req.Header.Set("Authorization", "Bearer test-token")
    
    // Execute
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
}
```

## Configuration

### Environment Variables
```bash
# Server Configuration
API_PORT=8080
API_HOST=0.0.0.0
API_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=zerotrace
DB_USER=postgres
DB_PASSWORD=password
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h

# Rate Limiting
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

## Deployment

### Docker Configuration
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o api ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/api .
CMD ["./api"]
```

### Local Development
```bash
# Run the API
go run ./cmd/api

# Run with Docker
podman build -t zerotrace-api .
podman run -p 8080:8080 zerotrace-api
```
