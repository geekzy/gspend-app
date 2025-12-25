# gRPC & Protocol Buffers - Implementation Guide

## Overview

This document provides the Protocol Buffer definitions and gRPC implementation details for inter-service communication in the Family Finance Application.

## Communication Pattern

- **Frontend to Backend**: HTTP/REST through Nginx
- **Backend to Backend**: gRPC for high-performance inter-service calls

## Protocol Buffer Definition

### `proto/auth/v1/auth.proto`

```protobuf
syntax = "proto3";

package auth.v1;

option go_package = "github.com/yourorg/family-finance-app/pkg/proto/auth/v1;authv1";

import "google/protobuf/timestamp.proto";

// AuthService handles authentication and user management via gRPC
service AuthService {
  // ValidateToken validates a JWT token and returns user information
  rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
  
  // GetUserProfile retrieves user profile by user ID
  rpc GetUserProfile(GetUserProfileRequest) returns (GetUserProfileResponse);
  
  // CheckUserExists verifies if a user exists
  rpc CheckUserExists(CheckUserExistsRequest) returns (CheckUserExistsResponse);
}

// ValidateToken Messages
message ValidateTokenRequest {
  string token = 1;
}

message ValidateTokenResponse {
  bool valid = 1;
  string user_id = 2;
  string email = 3;
  string error_message = 4;
}

// GetUserProfile Messages
message GetUserProfileRequest {
  string user_id = 1;
}

message GetUserProfileResponse {
  User user = 1;
}

message User {
  string id = 1;
  string email = 2;
  string full_name = 3;
  int32 family_size = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}

// CheckUserExists Messages
message CheckUserExistsRequest {
  string user_id = 1;
}

message CheckUserExistsResponse {
  bool exists = 1;
}
```

## gRPC Service Implementation

### Auth Service - gRPC Server

**File**: `apps/auth-service/internal/grpc/auth_grpc_service.go`

```go
package grpc

import (
    "context"
    authv1 "github.com/yourorg/family-finance-app/pkg/proto/auth/v1"
    "github.com/yourorg/family-finance-app/apps/auth-service/internal/service"
)

type AuthGRPCService struct {
    authv1.UnimplementedAuthServiceServer
    authService *service.AuthService
}

func NewAuthGRPCService(authService *service.AuthService) *AuthGRPCService {
    return &AuthGRPCService{
        authService: authService,
    }
}

func (s *AuthGRPCService) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
    // Validate JWT token
    claims, err := s.authService.ValidateJWT(req.Token)
    if err != nil {
        return &authv1.ValidateTokenResponse{
            Valid: false,
            ErrorMessage: err.Error(),
        }, nil
    }
    
    return &authv1.ValidateTokenResponse{
        Valid:  true,
        UserId: claims.UserID,
        Email:  claims.Email,
    }, nil
}

func (s *AuthGRPCService) GetUserProfile(ctx context.Context, req *authv1.GetUserProfileRequest) (*authv1.GetUserProfileResponse, error) {
    user, err := s.authService.GetUserByID(req.UserId)
    if err != nil {
        return nil, err
    }
    
    return &authv1.GetUserProfileResponse{
        User: &authv1.User{
            Id:         user.ID,
            Email:      user.Email,
            FullName:   user.FullName,
            FamilySize: int32(user.FamilySize),
        },
    }, nil
}

func (s *AuthGRPCService) CheckUserExists(ctx context.Context, req *authv1.CheckUserExistsRequest) (*authv1.CheckUserExistsResponse, error) {
    exists, err := s.authService.CheckUserExists(req.UserId)
    if err != nil {
        return nil, err
    }
    
    return &authv1.CheckUserExistsResponse{
        Exists: exists,
    }, nil
}
```

**File**: `apps/auth-service/internal/grpc/server.go`

```go
package grpc

import (
    "fmt"
    "net"
    
    "google.golang.org/grpc"
    authv1 "github.com/yourorg/family-finance-app/pkg/proto/auth/v1"
)

type GRPCServer struct {
    server *grpc.Server
    port   int
}

func NewGRPCServer(authGRPCService *AuthGRPCService, port int) *GRPCServer {
    server := grpc.NewServer()
    authv1.RegisterAuthServiceServer(server, authGRPCService)
    
    return &GRPCServer{
        server: server,
        port:   port,
    }
}

func (s *GRPCServer) Start() error {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
    if err != nil {
        return err
    }
    
    return s.server.Serve(lis)
}

func (s *GRPCServer) Stop() {
    s.server.GracefulStop()
}
```

### Financial Service - gRPC Client

**File**: `apps/financial-service/internal/grpc/client/auth_client.go`

```go
package client

import (
    "context"
    "time"
    
    authv1 "github.com/yourorg/family-finance-app/pkg/proto/auth/v1"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

type AuthGRPCClient struct {
    client authv1.AuthServiceClient
    conn   *grpc.ClientConn
}

func NewAuthGRPCClient(authServiceAddr string) (*AuthGRPCClient, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    conn, err := grpc.DialContext(ctx, authServiceAddr, 
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),
    )
    if err != nil {
        return nil, err
    }
    
    client := authv1.NewAuthServiceClient(conn)
    
    return &AuthGRPCClient{
        client: client,
        conn:   conn,
    }, nil
}

func (c *AuthGRPCClient) ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error) {
    return c.client.ValidateToken(ctx, &authv1.ValidateTokenRequest{
        Token: token,
    })
}

func (c *AuthGRPCClient) GetUserProfile(ctx context.Context, userID string) (*authv1.GetUserProfileResponse, error) {
    return c.client.GetUserProfile(ctx, &authv1.GetUserProfileRequest{
        UserId: userID,
    })
}

func (c *AuthGRPCClient) CheckUserExists(ctx context.Context, userID string) (*authv1.CheckUserExistsResponse, error) {
    return c.client.CheckUserExists(ctx, &authv1.CheckUserExistsRequest{
        UserId: userID,
    })
}

func (c *AuthGRPCClient) Close() error {
    return c.conn.Close()
}
```

## Code Generation

### `proto/Makefile`

```makefile
.PHONY: generate
generate:
\t@echo "Generating Protocol Buffer code..."
\tprotoc --go_out=../apps/auth-service/pkg/proto \\
\t       --go_opt=paths=source_relative \\
\t       --go-grpc_out=../apps/auth-service/pkg/proto \\
\t       --go-grpc_opt=paths=source_relative \\
\t       auth/v1/*.proto
\t
\tprotoc --go_out=../apps/financial-service/pkg/proto \\
\t       --go_opt=paths=source_relative \\
\t       --go-grpc_out=../apps/financial-service/pkg/proto \\
\t       --go-grpc_opt=paths=source_relative \\
\t       auth/v1/*.proto
\t@echo "Done!"

.PHONY: clean
clean:
\t@echo "Cleaning generated files..."
\tfind ../apps -name "*.pb.go" -type f -delete
\t@echo "Cleaned!"

.PHONY: install-tools
install-tools:
\t@echo "Installing protoc plugins..."
\tgo install google.golang.org/protobuf/cmd/protoc-gen-go@latest
\tgo install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
\t@echo "Done!"
```

### Usage

```bash
## Install required tools
cd proto
make install-tools

# Generate Go code from proto files
make generate

# Clean generated files
make clean
```

## Server Initialization

### Auth Service Main

```go
// apps/auth-service/cmd/server/main.go
func main() {
    // ... existing initialization ...
    
    // Start gRPC server in goroutine
    authGRPCService := grpc.NewAuthGRPCService(authService)
    grpcServer := grpc.NewGRPCServer(authGRPCService, 9091)
    
    go func() {
        if err := grpcServer.Start(); err != nil {
            log.Fatalf("Failed to start gRPC server: %v", err)
        }
    }()
    
    // Start HTTP server
    // ... existing HTTP server code ...
}
```

### Financial Service Main

```go
// apps/financial-service/cmd/server/main.go
func main() {
    // ... existing initialization ...
    
    // Initialize gRPC client for Auth Service
    authClient, err := client.NewAuthGRPCClient("auth-service:9091")
    if err != nil {
        log.Fatalf("Failed to connect to Auth Service: %v", err)
    }
    defer authClient.Close()
    
    // Inject gRPC client into middleware or services
    authMiddleware := middleware.NewAuthMiddleware(authClient)
    
    // Start HTTP server with middleware
    // ... existing HTTP server code ...
}
```

## Docker Compose Updates

Update `docker-compose.yml` to expose gRPC ports:

```yaml
  auth-service:
    ...
    ports:
      - "${AUTH_SERVICE_PORT:-8081}:8081"      # HTTP
      - "${AUTH_SERVICE_GRPC_PORT:-9091}:9091"  # gRPC
    ...

  financial-service:
    ...
    ports:
      - "${FINANCIAL_SERVICE_PORT:-8082}:8082"      #HTTP
      - "${FINANCIAL_SERVICE_GRPC_PORT:-9092}:9092"  # gRPC
    environment:
      ...
      - AUTH_SERVICE_GRPC_ADDR=auth-service:9091
    ...
```

## Benefits of gRPC for Inter-Service Communication

1. **Performance**: Binary protocol, HTTP/2 multiplexing
2. **Type Safety**: Strong typing with Protocol Buffers
3. **Code Generation**: Automated client/server code
4. **Streaming**: Built-in support for streaming RPCs
5. **Language Agnostic**: Can add services in other languages later
6. **Versioning**: Easy API versioning with package versions

## Testing gRPC Services

### Example Test

```go
func TestAuthGRPCService_ValidateToken(t *testing.T) {
    // Setup
    authService := setupMockAuthService()
    grpcService := NewAuthGRPCService(authService)
    
    // Test
    resp, err := grpcService.ValidateToken(context.Background(), &authv1.ValidateTokenRequest{
        Token: "valid-token",
    })
    
    // Assert
    assert.NoError(t, err)
    assert.True(t, resp.Valid)
    assert.NotEmpty(t, resp.UserId)
}
```
