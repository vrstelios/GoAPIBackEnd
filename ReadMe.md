# GoAPIBackEnd

A RESTful API service built with Go and Gin framework, featuring JWT authentication, GraphQL support, and Swagger documentation.

![APIs.png](APIs.png)

## Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── congig/
│   └── config.go                # Load envaronment variables
├── database/
│   ├── schema.sql               # Schema from database
│   └── schema-changes.sql       # Change the database
├── docs/                        # Swagger documentation
│   ├── docs.go                
│   ├── swagger.json   
│   └── swagger.yaml         
├── internal/
│   ├── api/                     # HTTP handlers
│   │   ├──  handler_Exercise.go  
│   │   └──  handler_Users.go    
│   ├── apperrors/
│   │   ├──  apperror.go         # Error handlering
│   ├── auth/
│   │   └──  middleware.go       # JWT and role-based middleware
│   ├── database/
│   │   ├── database.go          # Connection with database
│   │   └── crud.go              # Database queries
│   ├── graphql/
│   │   └── schema.go            # GraphQL schema definitions
│   ├── http/
│   │   ├── httpserver.go        # Route definitions
│   │   └── httpserver_test.go   # Route tests
│   ├── model/
│   │   └── model.go             # Data models
│   └── server.go                # Server initialization
├── .env                         # Include all variables
├── air.toml                     # Air live reload configuration
├── Makefile                     # Build and run commands
├── go.mod                       # Go module file
└── README.md                    # This file
```

## Features

- RESTful API with Gin framework
- JWT-based authentication
- Role-based access control (admin/user)
- GraphQL endpoint with GraphiQL interface
- Swagger/OpenAPI documentation
- Task management CRUD operations
- Image download functionality
- Server-side filtering and pagination

## Prerequisites

- Go 1.24 or higher
- Make (optional, for using Makefile commands)

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd GoAPIBackEnd
```

2. Install dependencies:
```bash
go mod download
```

## Usage

### Run the application

```bash
# Using Make
make run

# Or directly
go run cmd/api/main.go

### Generate Swagger documentation

```bash
make swagger
```

Then access Swagger UI at: `http://localhost:8080/swagger/index.html`

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/auth/logout` - Logout

### Tasks (Protected - requires JWT)
- `GET /api/tasks` - Get all tasks
- `GET /api/tasks/:id` - Get task by ID
- `POST /api/tasks` - Create task (admin only)
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task
- `POST /api/tasks/query` - Query tasks with filtering and pagination

### GraphQL
- `POST /api/graphql` - GraphQL endpoint (with GraphiQL interface)

### Other
- `POST /api/download/images` - Download images from URLs

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make TestRoutes
```