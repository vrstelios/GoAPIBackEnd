# GoAPIBackEnd

A production-ready RESTful API built with Go and the Gin Framework, featuring secure JWT authentication, role-based access control, GraphQL integration, CSV data loading, and full Swagger documentation.

> Project Inspiration: This project is based on the Fitness [Workout Tracker](https://roadmap.sh/projects/fitness-workout-tracker) project idea from roadmap.sh

---

## Features

### Core Functionality
- **User Authentication** – Secure JWT-based login & signup
- **Role-Based Access Control** – Admin & User permissions
- **Workout Management** – Full CRUD operations
- **Coach Management** – Admin-only access
- **Exercise Management** – RESTful endpoints
- **CSV Import Support** – Load workouts from CSV files with go-routing
- **Filtering & Pagination** – Server-side query handling
- **API Route Testing**

### Technical Features
- **Gin Framework** – Fast and lightweight
- **JWT Authentication Middleware**
- **Database Integration**
- **GraphQL Endpoint + GraphiQL UI**
- **Swagger / OpenAPI Documentation**
- **Modular Internal Package Structure**
- **Makefile Support for Common Tasks**

---

## Architecture

This project follows a clean and structured architecture for maintainability and scalability.

```
.
├── cmd/
│   └── api/
│       └── main.go               # Application entry point
├── congig/
│   ├──  config.go                # Load envaronment variables
│   ├──  config-development.yml   # Include all variables for dev
│   └──  config-production.yml    # Include all variables for pro
├── csv/                          # Include all csv files
│   ├──  workouts_set_1.csv  
│   ├──  workouts_set_2.csv 
│   └──  workouts_set_2.csv              
├── database/
│   ├── schema.sql                # Schema from database
│   └── schema-changes.sql        # Change the database
├── docs/                         # Swagger documentation
│   ├── docs.go                
│   ├── swagger.json   
│   └── swagger.yaml       
├── internal/
│   ├── api/                      # HTTP handlers
│   │   ├──  api.go
│   │   ├──  client.go
│   │   ├──  handler_Coach.go
│   │   ├──  handler_Exercise.go
│   │   ├──  handler_Users.go
│   │   ├──  handler_WokroutLogs.go  
│   │   ├──  handler_WokroutLogs.go 
│   │   └──  handler_Wokrouts.go    
│   ├── apperrors/
│   │   ├──  apperror.go          # Error handlering
│   ├── database/
│   │   ├──  database.go          # Connection with database
│   │   └──  crud.go              # Database queries
│   ├── graphql/
│   │   └──  schema.go            # GraphQL schema definitions
│   ├── helpers/
│   │   └──  token.go             # Include tools for password/token
│   ├── http/
│   │   ├──  httpserver.go        # Route definitions
│   │   └──  httpserver_test.go   # Route tests
│   ├── middleware/
│   │   ├──  auth.go              # JWT and role-based middleware
│   │   └──  token_provider.go    # Token provider            
│   ├── type/                     # Data models
│   │   ├── misc/    
│   │   │    └──  params.go         
│   │   └── modeles/
│   │   │    ├──  exercises.go
│   │   │    ├──  users.go
│   │   │    └──  workouts.go             
│   └──  server.go                # Server initialization
├──  Makefile                     # Build and run commands
├── go.mod                       # Go module file
└── README.md                    # This file
```

---

## Database Schema
The project includes:
- Users
- Exercises
- Coaches
- Workouts
- Workout Logs

*(Schema files included in `/database` folder with schema & schema changes SQL.)*

![DiagramDataBase.png](DiagramDataBase.png)

---

## Technology Stack
- **Language:** Go (1.24+)
- **Framework:** Gin Web Framework
- **Authentication:** JWT Tokens
- **GraphQL:** GraphQL endpoint with GraphiQL UI
- **Documentation:** Swagger / OpenAPI
- **Database:** PostgreQL-based relational DB
- **Testing:** Built-in route testing
- **Configuration:** YAML-based config per environment

---

## Getting Started

### Prerequisites
- Go **1.24+**
- PostgreSQL **12+**
- Make (optional but recommended)

---

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

3. Set up environment variables Create a .yml file in the src directory:
```bash
server:
  host: localhost
  port: ":8080"

api:
  baseURL: "http://localhost:8080/api/"

database:
  host: localhost
  port: 5432
  user: your_username
  password: your_password
  name: workout_tracker
  sslmode: disable

jwt:
  secret: your_jwt_secret
```

4. Set the environment variable to load development configuration:

5. Start the server
```bash
go run cmd/api/main.go
```

### Swagger Documentation

The API is fully documented using Swagger/OpenAPI 3.0. Once the server is running, you can access:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- OpenAPI JSON: http://localhost:8080/swagger/doc.json
- OpenAPI YAML: Available in src/docs/swagger.yaml

## API Endpoints

### Authentication
- `POST /api/auth/signup` - Signup a new user
- `POST /api/auth/login` - Login and get JWT token
- `GET /api/auth/logout` - Logout

### Tasks (Protected - requires JWT)
- `GET /api/exercises` - Get all exercises
- `GET /api/exercises/:id` - Get exercises by Id
- `POST /api/exercises` - Create Exercise only admin
- `GET /api/coach/:name` - Get Coach only admin
- `POST /api/coach` - Create Coach only admin
- `POST /api/workouts` - Create Workout
- `GET /api/workouts/:id` - Get Workout
- `PUT /api/workouts/:id` - Update Workout
- `DELETE /api/workouts/:id` - Delete Workout
- `POST /api/workouts/query` - Query Workouts (server side filtering & paging)
- `POST /api/load/workouts` - load data from csv files foe workouts

### GraphQL
- `POST /api/graphql/workouts/logs` - GraphQL endpoint (with GraphQL interface)

## Testing
- `POST /api/auth/signup` - Signup a new user
- `POST /api/auth/login` - Login and get JWT token
- `POST /api/auth/signup` - Signup a new admin User
- `POST /api/auth/login` - Login with admin User and get JWT token

```bash
# Run all tests
make test

# Run tests with coverage
make TestRoutes
```

### Security
- JWT Authentication
- Role-Based Access Control
- Protected Routes
- Secure Token Handling
- Centralized Error Handling

### Contributing
- Fork the repo
- Create a feature branch (git checkout -b feature/amazing-feature)
- Commit your changes (git commit -m 'Add some amazing feature')
- Push to the branch (git push origin feature/amazing-feature)
- Open a Pull Request

### Author
[DoctorVerRossi](https://github.com/vrstelios)
---

If you find this project helpful, please give it a star on GitHub!