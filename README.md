# Solid Coffee Backend

A RESTful API backend service for the Solid Coffee application, built with Go and following clean architecture principles.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Database](#database)
- [Development](#development)
- [API Documentation](#api-documentation)
- [Deployment](#deployment)
- [Project Structure](#project-structure)

## Overview

This backend service provides a comprehensive API for managing a coffee shop application, including user authentication, product management, order processing, and payment handling. The application is built using clean architecture principles to ensure maintainability, testability, and scalability.

## Architecture

The project follows a layered architecture pattern:

- **Controller Layer**: Handles HTTP requests and responses
- **Service Layer**: Contains business logic
- **Repository Layer**: Manages data access
- **Model Layer**: Defines data structures

## Tech Stack

- **Runtime**: Go 1.25+
- **Web Framework**: Gin-Gonic
- **Database**: PostgreSQL 18+
- **Cache**: Redis 8.4+
- **Authentication**: JWT (JSON Web Tokens)
- **Migration**: golang-migrate
- **Documentation**: Swagger/OpenAPI

## Prerequisites

Before running this application, ensure you have the following installed:

- Go 1.25 or higher
- PostgreSQL 18 or higher
- Redis 8.4 or higher
- golang-migrate CLI
- Make (optional, but recommended)

## Getting Started

### Installation

1.Clone the repository:

```bash
git clone https://github.com/NugrahaPancaWibisana/solid-coffee-be.git
cd solid-coffee-be
```

2.Install Go dependencies:

```bash
go mod download
```

3.Install golang-migrate CLI
[How To Download](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Quick Start

1.Copy the environment template:

```bash
cp .env.example .env
```

2.Configure your environment variables in `.env`

3.Run database migrations:

```bash
make migrate-up
```

5.Start the development server:

```bash
make run
```

The API will be available at `http://localhost:8080`

## Configuration

Create a `.env` file in the root directory with the following variables:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USERNAME=username
DB_PASSWORD=password
DB_NAME=database

RDB_HOST=localhost
RDB_PORT=6380
RDB_USERNAME=username
RDB_PASSWORD=password
RDB_NAME=0

RDB_KEY=key # example: name | it will be like this in the project (name:)

JWT_SECRET=secret
JWT_ISSUER=username
```

## Database

### Migrations

Run all pending migrations:

```bash
make migrate-up
```

Rollback the last migration:

```bash
make migrate-down
```

Create a new migration:

```bash
make migrate-create NAME=migration-name
```

Run Seeder

```bash
make seeder
```

### Database Schema

The application uses the following main tables:

- `users` - User accounts and authentication
- `products` - Product catalog
- `categories` - Product categories
- `orders` - Order management
- `dt_order` - Order details
- `payments` - Payment records
- `reviews` - Product reviews
- `menus` - Menu items
- `product_images` - Product image references
- `product_categories` - Product-category relationships
- `product_size` - Product size options
- `product_type` - Product type classifications

## Development

### Running Locally

**Standard mode:**

```bash
go run cmd/main.go
```

**With hot reload using Air:**

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
air
```

**Using Make:**

```bash
make run
# or with hot-reload
make dev
```

### Building

Build the production binary:

```bash
go build -o bin/app cmd/main.go
./bin/app
```

### Docker

Build the Docker image:

```bash
docker build -t solid-coffee-be .
```

Run the container:

```bash
docker run -p 8080:8080 --env-file .env solid-coffee-be
```

### Debugging with VS Code

The project includes VS Code debug configurations. Press F5 to start debugging, or use the Run and Debug panel.

Debug configuration (`.vscode/launch.json`):

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Go App",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/main.go",
      "envFile": "${workspaceFolder}/.env",
      "output": "${workspaceFolder}/tmp/__debug_bin"
    }
  ]
}
```

## API Documentation

### Swagger Documentation

Once the application is running, access the Swagger UI at:

```sh
http://localhost:8080/swagger/index.html
```

### Generating Swagger Documentation

Install swag:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Generate documentation:

```bash
swag init -g cmd/main.go -o docs
```

### API Endpoints

_**Authentication**_

- `POST /auth` - User login
- `POST /auth/new` - Register new user
- `DELETE /auth` - User logout (user role required)
- `POST /auth/forgot-password` - Request password reset
- `POST /auth/forgot-password/update` - Update password after reset

_**Users**_

- `GET /user` - Get current user profile (user/admin role required)
- `PATCH /user` - Update user profile (user/admin role required)
- `PATCH /user/password` - Update user password (user/admin role required)
- `POST /admin/user` - Create new user (admin role required)
- `GET /admin/user` - Get all users (admin role required)
- `PATCH /admin/user/:id` - Update user profile (admin role required)
- `DELETE /admin/user/:id` - Delete user (admin role required)

_**Products**_

- `GET /products` - List all products
- `GET /products/:id` - Get product by ID
- `GET /products/product-sizes` - List all product sizes
- `GET /products/product-types` - List all product types
- `POST /admin/products` - Create new product (admin role required)
- `GET /admin/products/:id` - Get product by ID (admin role required)
- `PATCH /admin/products/:id` - Update product (admin role required)
- `DELETE /admin/products/:id` - Delete product (admin role required)
- `DELETE /admin/products/image/:id` - Delete product image (admin role required)

_**Orders**_

- `POST /orders` - Create new order (user role required)
- `GET /orders/history` - List user order history (user role required)
- `GET /orders/history/:id` - Get order details (user/admin role required)
- `POST /orders/review` - Add a review to an order (user role required)
- `GET /admin/orders` - List all orders (admin role required)
- `PATCH /admin/orders` - Update order status (admin role required)

_**Menu**_

- `GET /admin/menu` - List menu items (admin role required)
- `GET /admin/menu/:id` - Get menu item details (admin role required)
- `POST /admin/menu` - Create new menu item (admin role required)
- `PATCH /admin/menu/:id` - Update menu item (admin role required)
- `DELETE /admin/menu/:id` - Delete menu item (admin role required)

## Deployment

### Production Build

Build the binary:

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app cmd/main.go
```

2.Set environment to production:

```env
ENV=production
```

3.Run migrations on production database

4.Deploy the binary with appropriate environment variables

### Docker Deployment

1.Build production image:

```bash
docker build -t solid-coffee-be:latest .
```

2.Run with docker-compose or your orchestration tool of choice

### Environment-Specific Configurations

- Development: `.env` or `.env.local`
- Staging: `.env.staging`
- Production: `.env.production`

## Project Structure

```bash
solid-coffee-be/
├── cmd/                    # Application entry points
│   └── main.go             # Main application
├── db/
│   └── migration/          # Database migration files
├── docs/                   # API documentation (Swagger)
├── internal/               # Private application code
│   ├── apperror/           # Application-specific errors
│   ├── cache/              # Caching layer (Redis)
│   ├── config/             # Configuration management
│   ├── controller/         # HTTP request handlers
│   ├── dto/                # Data Transfer Objects
│   ├── middleware/         # HTTP middlewares
│   ├── model/              # Domain models
│   ├── repository/         # Data access layer
│   ├── response/           # Response utilities
│   ├── router/             # Route definitions
│   └── service/            # Business logic layer
├── pkg/                    # Public libraries
│   ├── hash/               # Password hashing utilities
│   └── jwt/                # JWT token management
├── public/                 # Static files and uploads
│   ├── products/           # Product images
│   └── profile/            # User profile pictures
├── tmp/                    # Temporary files (gitignored)
├── .env.example            # Environment template
├── .gitignore              # Git ignore rules
├── Dockerfile              # Docker configuration
├── Makefile                # Build automation
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── README.md               # This file
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Coding Standards

- Follow Go standard formatting (`gofmt`, `goimports`)
- Write unit tests for new features
- Update documentation as needed
- Follow the existing project structure

## Authors

- Nugraha Panca Wibisana - Backend Developer & Github Master
- Ari Ramadhan - Backend Developer

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Go community for excellent tooling and libraries
- Contributors and maintainers of dependencies
