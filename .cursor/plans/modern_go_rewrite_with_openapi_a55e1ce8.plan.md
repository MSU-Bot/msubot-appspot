---
name: Modern Go Rewrite with OpenAPI
overview: Complete rewrite of the codebase using modern Go patterns, OpenAPI specification, structured logging, dependency injection, and comprehensive error handling to replace the legacy implementation.
todos:
  - id: upgrade-go
    content: Upgrade Go version to 1.22+ in go.mod and update CircleCI config
    status: pending
  - id: openapi-spec
    content: Create comprehensive OpenAPI 3.0 specification for all 7 endpoints
    status: pending
  - id: project-structure
    content: Set up modern project structure (cmd/, internal/, pkg/, api/)
    status: pending
  - id: config-management
    content: Implement configuration management with envconfig
    status: pending
  - id: logging-migration
    content: Migrate from logrus to standard library slog
    status: pending
  - id: error-handling
    content: Create error handling package and remove all panic() calls
    status: pending
  - id: dependency-injection
    content: Set up dependency injection container and service interfaces
    status: pending
  - id: firebase-service
    content: Refactor Firebase operations into service layer with interfaces
    status: pending
  - id: scraper-service
    content: Refactor scraper logic into service with proper HTTP client management
    status: pending
  - id: messaging-service
    content: Refactor messaging (Plivo/Email) into unified service
    status: pending
  - id: chi-router
    content: Implement chi router with middleware stack (logging, recovery, CORS, request ID)
    status: pending
  - id: handler-refactor
    content: Rewrite all handlers to use dependency injection and modern patterns
    status: pending
  - id: http-client-factory
    content: Create HTTP client factory with timeouts and retry logic
    status: pending
  - id: update-dependencies
    content: Update all dependencies to latest versions and remove +incompatible
    status: pending
  - id: unit-tests
    content: Write comprehensive unit tests for services and handlers
    status: pending
  - id: integration-tests
    content: Add integration tests for critical flows
    status: pending
  - id: linting-config
    content: Update linting configuration (revive.toml, add golangci-lint)
    status: pending
  - id: ci-cd-update
    content: Update CircleCI with modern Go version and additional quality checks
    status: pending
  - id: documentation
    content: Update README with new architecture and API documentation
    status: pending
isProject: false
---

# Modern Go Rewrite with OpenAPI Specification

## Overview

This plan outlines a complete rewrite of the msubot-appspot codebase using modern Go practices, including OpenAPI specification, structured logging, dependency injection, and improved error handling.

## Current State Analysis

### Critical Issues Identified

- **Go Version**: Using Go 1.20 (should upgrade to 1.22+)
- **Panics**: 4 instances of `panic()` calls that should be proper error returns
- **Resource Leaks**: Missing `defer response.Body.Close()` in several places
- **HTTP Client**: Using `http.DefaultClient` without timeouts (potential hangs)
- **Error Handling**: Inconsistent patterns, some functions return `nil` instead of errors
- **Logging**: Mix of `log.Printf` and structured logging (should standardize on `slog`)
- **No Router**: Using raw `http.HandleFunc` instead of modern router (chi/gorilla)
- **No API Documentation**: No OpenAPI/Swagger spec
- **Dependencies**: Outdated versions, some with `+incompatible` flag
- **Testing**: Minimal test coverage
- **CI/CD**: Using Go 1.16 in CircleCI (outdated)

## Architecture Decisions

### 1. OpenAPI Specification

- Create comprehensive OpenAPI 3.0 specification documenting all endpoints
- Use `oapi-codegen` or `go-swagger` for code generation
- Generate client SDKs for consumers

### 2. Modern HTTP Stack

- **Router**: Use `chi` router for middleware support and better routing
- **Middleware**: Request logging, error recovery, CORS, request ID tracing
- **Context**: Proper context propagation throughout handlers

### 3. Dependency Injection

- Replace global state with dependency injection
- Use constructor functions for services
- Enable testability and modularity

### 4. Structured Logging

- Migrate from `logrus` to standard library `slog` (Go 1.21+)
- Consistent log levels and structured fields
- Request-scoped logging with trace IDs

### 5. Error Handling

- Remove all `panic()` calls
- Use `errors.Join` and `errors.Is/As` (Go 1.20+)
- Custom error types for better error handling
- Proper HTTP error responses

### 6. HTTP Client Management

- Create HTTP client factory with proper timeouts
- Context-aware requests
- Retry logic for external API calls

## Implementation Plan

### Phase 1: Foundation & Specification

1. **Upgrade Go Version**
  - Update `go.mod` to Go 1.22 or 1.23
  - Update CircleCI config to use latest Go version
2. **Create OpenAPI Specification**
  - Document all 7 endpoints:
    - `GET /sections` - Scrape course sections
    - `POST /welcomeuser` - Welcome user handler
    - `GET /checktrackedsections` - Check tracked sections (cron)
    - `POST /prunesections` - Prune expired sections (cron)
    - `POST /receivemessage` - Receive SMS from Plivo
    - `GET /healthcheck` - Health check endpoint
    - `POST /updatedepartments` - Update department data
  - Define request/response schemas
  - Add authentication/authorization if needed
  - Include error response schemas
3. **Project Structure**

```
   /api
     openapi.yaml
   /cmd
     /server
       main.go
   /internal
     /config
     /handlers
     /services
     /models
     /middleware
     /client
   /pkg
     /errors
   

```

### Phase 2: Core Infrastructure

1. **Configuration Management**
  - Use `envconfig` or similar for environment variables
  - Validate required config on startup
  - Type-safe configuration struct
2. **Dependency Injection Setup**
  - Create service interfaces
  - Implement dependency container
  - Wire dependencies in `main.go`
3. **Structured Logging**
  - Replace `logrus` with `slog`
  - Create logger factory
  - Add request ID middleware
4. **Error Handling Package**
  - Define custom error types
  - HTTP error response helpers
  - Error wrapping utilities

### Phase 3: Service Layer Refactoring

1. **Firebase Service**
  - Extract to `internal/services/firestore.go`
  - Interface-based design
  - Proper connection pooling
  - Context-aware operations
2. **Scraper Service**
  - Extract to `internal/services/scraper.go`
  - HTTP client with timeouts
  - Retry logic for external calls
  - Proper resource cleanup
3. **Messaging Service**
  - Extract to `internal/services/messenger.go`
  - Plivo client abstraction
  - Email service abstraction
  - Unified notification interface
4. **Section Tracking Service**
  - Extract to `internal/services/tracker.go`
  - Business logic separation
  - Cleaner worker pool implementation

### Phase 4: HTTP Layer Modernization

1. **Router Setup**
  - Implement chi router
  - Add middleware stack:
    - Request logging
    - Error recovery
    - CORS (configurable, not hardcoded `*`)
    - Request ID
    - Timeout middleware
2. **Handler Refactoring**
  - Convert handlers to use dependency injection
  - Consistent error handling
  - Proper HTTP status codes
  - Request/response validation
3. **OpenAPI Integration**
  - Generate server code from OpenAPI spec
  - Or implement handlers matching spec
  - Add Swagger UI endpoint

### Phase 5: Testing & Quality

1. **Unit Tests**
  - Service layer tests with mocks
  - Handler tests with httptest
  - Test utilities package
2. **Integration Tests**
  - Test critical flows
  - Use testcontainers for Firebase if possible
3. **Linting & Formatting**
  - Update `revive.toml` for modern rules
  - Add `golangci-lint` configuration
  - Ensure `go fmt` compliance

### Phase 6: Dependencies & CI/CD

1. **Dependency Updates**
  - Update all dependencies to latest stable versions
  - Remove `+incompatible` dependencies
  - Use `go mod tidy` and verify
2. **CI/CD Updates**
  - Update CircleCI to latest Go version
  - Add linting step
  - Add security scanning
  - Improve test reporting
3. **Documentation**
  - Update README with new architecture
  - Add API documentation links
  - Migration guide from old to new

## Key Files to Create/Modify

### New Files

- `api/openapi.yaml` - OpenAPI 3.0 specification
- `internal/config/config.go` - Configuration management
- `internal/middleware/*.go` - Middleware functions
- `internal/services/*.go` - Service layer implementations
- `pkg/errors/*.go` - Error handling utilities
- `.golangci.yml` - Linter configuration

### Files to Rewrite

- `main.go` → `cmd/server/main.go` - Complete rewrite with DI
- `server/*/*.go` → `internal/handlers/*.go` - Modern handlers
- `server/serverutils/*.go` → `internal/services/*.go` - Service layer
- `server/models/*.go` → `internal/models/*.go` - Models (may need updates)

### Files to Update

- `go.mod` - Upgrade Go version, update dependencies
- `.circleci/config.yml` - Update Go version, add steps
- `revive.toml` - Modern linting rules
- `README.md` - Architecture documentation

## Migration Strategy

1. **Parallel Development**: Build new implementation alongside old
2. **Feature Parity**: Ensure all existing functionality is preserved
3. **Testing**: Comprehensive testing before switchover
4. **Gradual Rollout**: Deploy to staging first, then production

## Benefits

- **Maintainability**: Clear separation of concerns, testable code
- **Documentation**: OpenAPI spec provides clear API contract
- **Reliability**: Proper error handling, no panics, resource management
- **Observability**: Structured logging with request tracing
- **Security**: Proper CORS, input validation, no hardcoded secrets
- **Performance**: HTTP client timeouts, connection pooling
- **Developer Experience**: Modern Go patterns, better tooling support

