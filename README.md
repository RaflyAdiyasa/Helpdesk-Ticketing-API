# Helpdesk Ticketing API

REST API helpdesk ticketing system built with Go, Fiber, GORM, MySQL, JWT authentication, and a Clean Architecture-style project layout.

This API lets users register, login, create support tickets, view their own tickets, and lets admins view all tickets and update ticket status.

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Run with Docker Compose](#run-with-docker-compose)
- [Configuration](#configuration)
- [Run Locally with Go](#run-locally-with-go)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
- [Usage Flow](#usage-flow)
- [Implementation Notes](#implementation-notes)

## Project Overview

The project is organized around helpdesk tickets:

- A `USER` can register, login, create tickets, and list their own tickets.
- An `ADMIN` can list all tickets and update a ticket status.
- Authentication uses JWT Bearer tokens.
- Data is stored in MySQL through GORM.
- Passwords are hashed with bcrypt.
- Ticket image uploads are stored in MinIO, and the saved ticket stores the MinIO object key in the `image` field.

The server starts from [cmd/api/main.go](cmd/api/main.go). It loads configuration, connects to MySQL, runs migrations, initializes repositories and use cases, registers HTTP routes, and starts Fiber on the configured port.

## Architecture

The application separates delivery, use case, domain, and infrastructure concerns. Handlers translate HTTP requests into use case calls, use cases hold the business rules, and repositories hide the storage details behind interfaces.

## Project Structure

The code follows a Clean Architecture-inspired separation:

```text
cmd/api
  Application entrypoint

internal/config
  Environment loading and app configuration

internal/domain/entity
  Core business models: User and Ticket

internal/domain/repository
  Repository interfaces used by the use cases

internal/usecase
  Business logic for auth and ticket operations

internal/infrastructure
  MySQL, GORM repositories, migrations, and MinIO client setup

internal/delivery/http
  Fiber handlers and middleware

pkg
  Shared helper packages for JWT, UUID generation, and bcrypt
```

### Request Flow

```text
HTTP request
-> Fiber route
-> Middleware, if route is protected
-> Handler
-> Use case
-> Repository interface
-> MySQL or MinIO implementation
-> JSON response
```

### Important Code Paths

- App bootstrap: [cmd/api/main.go](cmd/api/main.go)
- Config loader: [internal/config/config.go](internal/config/config.go)
- Auth routes: [internal/delivery/http/handler/auth_handler.go](internal/delivery/http/handler/auth_handler.go)
- Ticket routes: [internal/delivery/http/handler/ticket_handler.go](internal/delivery/http/handler/ticket_handler.go)
- JWT middleware: [internal/delivery/http/middleware/auth_middleware.go](internal/delivery/http/middleware/auth_middleware.go)
- Auth business logic: [internal/usecase/auth_usecase.go](internal/usecase/auth_usecase.go)
- Ticket business logic: [internal/usecase/ticket_usecase.go](internal/usecase/ticket_usecase.go)
- User model: [internal/domain/entity/user.go](internal/domain/entity/user.go)
- Ticket model: [internal/domain/entity/ticket.go](internal/domain/entity/ticket.go)
- MySQL connection: [internal/infrastructure/database/mysql.go](internal/infrastructure/database/mysql.go)
- Database migration: [internal/infrastructure/database/migration.go](internal/infrastructure/database/migration.go)

## Requirements

- Docker and Docker Compose, recommended for easiest setup
- Go 1.25.5, based on `go.mod`, if running without Docker
- MySQL, if running without Docker
- MinIO for ticket image uploads

## Run with Docker Compose

Docker Compose is the easiest way to run the project because it starts the API, MySQL, and MinIO together.

Start with the default configuration:

```bash
docker compose up --build
```

Or run it in the background:

```bash
docker compose up --build -d
```

The services will be available at:

| Service | URL / Port | Notes |
| --- | --- | --- |
| API | `http://localhost:8080` | Main REST API |
| MySQL | `localhost:3306` | Database data is stored in the `mysql_data` Docker volume |
| MinIO API | `http://localhost:9000` | Object storage API |
| MinIO Console | `http://localhost:9001` | Login with `MINIO_ACCESSKEY` and `MINIO_SECRETKEY` |

Useful Docker commands:

```bash
docker compose logs -f api
docker compose ps
docker compose down
```

To remove containers and delete database/object-storage data:

```bash
docker compose down -v
```

If you previously started Docker before the migration fix and see a `CREATE TABLE users` error containing `fk_tickets_owner`, reset the old development volumes and rebuild the API image:

```bash
docker compose down -v
docker compose build --no-cache api
docker compose up
```

Use the same reset if MySQL logs `Access denied for user 'helpdesk'`. MySQL only applies `MYSQL_USER`, `MYSQL_PASSWORD`, and `MYSQL_DATABASE` when the data volume is first created.

The MySQL database is created automatically from `DB_NAME`, and GORM creates the `users` and `tickets` tables when the API starts.

## Configuration

The app reads `.env` automatically with `github.com/joho/godotenv`. If `.env` does not exist, defaults from [internal/config/config.go](internal/config/config.go) are used.

For Docker Compose, a `.env` file is optional because [docker-compose.yml](docker-compose.yml) already contains safe development defaults. To customize the setup, copy the example file:

```bash
cp .env.example .env
```

Then edit `.env` in the project root:

```env
PORT=8080
JWT_SECRET=change-this-secret

DB_HOST=localhost
DB_PORT=3306
DB_NAME=helpdesk
DB_USER=helpdesk
DB_PASSWORD=helpdesk_password
DB_ROOT_PASSWORD=root_password
MYSQL_HOST_PORT=3306

MINIO_ENDPOINT=localhost:9000
MINIO_ACCESSKEY=minioadmin
MINIO_SECRETKEY=minioadmin123
MINIO_BUCKETNAME=uploads
MINIO_API_PORT=9000
MINIO_CONSOLE_PORT=9001
```

Important Docker Compose behavior:

- The `api` container uses `DB_HOST=mysql`, `DB_PORT=3306`, and `MINIO_ENDPOINT=minio:9000` because containers communicate through Docker service names.
- The root `.env` values `DB_HOST=localhost` and `MINIO_ENDPOINT=localhost:9000` are useful when running the Go app directly on your machine.
- Keep `.env` in the project root. If you run `go run ./cmd/api` from the root directory, that is the `.env` file loaded by the app.

Environment variables:

| Variable | Used by | Default | Description |
| --- | --- | --- | --- |
| `PORT` | API | `8080` | HTTP port used by the API |
| `JWT_SECRET` | API | `change-this-secret` in Docker, `indianman` in code | Secret key for signing JWT tokens |
| `DB_HOST` | API local run | `localhost` | MySQL host when running without Docker |
| `DB_PORT` | API local run | `3306` | MySQL port when running without Docker |
| `DB_NAME` | API, MySQL | `helpdesk` in Docker, `be` in code | MySQL database name |
| `DB_USER` | API, MySQL | `helpdesk` in Docker, `huanlocal` in code | MySQL user |
| `DB_PASSWORD` | API, MySQL | `helpdesk_password` in Docker, `pass123` in code | MySQL password |
| `DB_ROOT_PASSWORD` | MySQL | `root_password` | MySQL root password used by the container |
| `MYSQL_HOST_PORT` | Docker Compose | `3306` | Host port mapped to the MySQL container |
| `MINIO_ENDPOINT` | API local run | `localhost:9000` | MinIO endpoint when running without Docker |
| `MINIO_ACCESSKEY` | API, MinIO | `minioadmin` in Docker | MinIO username/access key |
| `MINIO_SECRETKEY` | API, MinIO | `minioadmin123` in Docker | MinIO password/secret key |
| `MINIO_BUCKETNAME` | API, MinIO init | `uploads` | Bucket created by the `minio-init` service |
| `MINIO_API_PORT` | Docker Compose | `9000` | Host port mapped to MinIO API |
| `MINIO_CONSOLE_PORT` | Docker Compose | `9001` | Host port mapped to MinIO Console |

If running without Docker, create the database before starting the app:

```sql
CREATE DATABASE helpdesk;
```

If you use the code defaults instead of `.env`, create:

```sql
CREATE DATABASE be;
```

On startup, GORM auto-migrates these tables:

- `users`
- `tickets`

## Run Locally with Go

For local Go development without Docker, start MySQL first and make sure the root `.env` points to your local services:

```env
DB_HOST=localhost
DB_PORT=3306
DB_NAME=helpdesk
DB_USER=helpdesk
DB_PASSWORD=helpdesk_password

MINIO_ENDPOINT=localhost:9000
MINIO_ACCESSKEY=minioadmin
MINIO_SECRETKEY=minioadmin123
```

Install dependencies:

```bash
go mod download
```

Run the API:

```bash
go run ./cmd/api
```

Default base URL:

```text
http://localhost:8080
```

Health check route is not currently implemented, so use one of the documented endpoints to test the server.

## Authentication

Protected ticket routes require this header:

```http
Authorization: Bearer <token>
```

The JWT contains:

- `user_id`
- `role`
- `exp`
- `iat`

Token expiry is currently set to 12 hours in code.

## API Endpoints

| Method | Endpoint | Auth | Role | Description |
| --- | --- | --- | --- | --- |
| `POST` | `/api/v1/register` | No | Public | Register a new user |
| `POST` | `/api/v1/login` | No | Public | Login and receive JWT token |
| `POST` | `/api/v1/tickets/` | Yes | `USER` or `ADMIN` | Create ticket |
| `GET` | `/api/v1/tickets/my-tickets` | Yes | `USER` or `ADMIN` | Get tickets owned by current user |
| `GET` | `/api/v1/tickets/admin/all` | Yes | `ADMIN` | Get all tickets |
| `PUT` | `/api/v1/tickets/admin/:id/status` | Yes | `ADMIN` | Update ticket status |

### 1. Register

Creates a new user. The `role` must be `USER` or `ADMIN`. The `is_remote` field is currently parsed from a string, so send `"true"` or `"false"`.

```bash
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jane",
    "email": "jane@example.com",
    "password": "pass123",
    "department": "IT department",
    "role": "USER",
    "is_remote": "false"
  }'
```

Success response:

```http
HTTP/1.1 201 Created
Content-Type: application/json
```

```json
{
  "message": "User registered successfully",
  "user": {
    "user_id": "USER-48050fe4-aeb5-403c-b121-91bf116c43be",
    "username": "jane",
    "email": "jane@example.com",
    "role": "USER",
    "profile_pict": "",
    "department": "IT department",
    "remote": false,
    "created_at": "2026-01-11T14:21:21.11+07:00",
    "updated_at": "2026-01-11T14:21:21.11+07:00"
  }
}
```

Common error responses:

```json
{
  "error": "User sudah ada"
}
```

```json
{
  "error": "Role must be either 'ADMIN' or 'USER'"
}
```

### 2. Login

Logs in with `username` and `password`, then returns a JWT token.

```bash
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "jane",
    "password": "pass123"
  }'
```

Success response:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Common error response:

```json
{
  "error": "invalid credentials: username tidak ditemukan"
}
```

### 3. Create Ticket

Creates a ticket for the authenticated user. The initial status is always `OPEN`.

To upload an image to MinIO, send `multipart/form-data` with `title`, `description`, and an `image` file field. The API also accepts `file` as a fallback file field name.

```bash
curl -X POST http://localhost:8080/api/v1/tickets/ \
  -H "Authorization: Bearer <token>" \
  -F "title=Printer Rusak" \
  -F "description=Printer di ruang HR tidak bisa connect ke PC" \
  -F "image=@./printer-error.jpg"
```

If you do not want to upload a file, JSON still works and `image` is stored as the plain string you send.

Success response:

```http
HTTP/1.1 201 Created
Content-Type: application/json
```

```json
{
  "ticket_id": "tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc",
  "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
  "title": "Printer Rusak",
  "description": "Printer di ruang HR tidak bisa connect ke PC",
  "image": "tickets/tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc/image-4b638c13-36c4-4ff1-bc60-2935873d6f9a.jpg",
  "status": "OPEN",
  "created_at": "2026-01-11T14:29:24.299+07:00",
  "updated_at": "2026-01-11T14:29:24.299+07:00",
  "deleted_at": null,
  "owner": null
}
```

Common error responses:

```json
{
  "error": "Authorization header required"
}
```

```json
{
  "error": "Invalid token"
}
```

```json
{
  "error": "title is required"
}
```

```json
{
  "error": "descripton is empty"
}
```

### 4. Get My Tickets

Returns tickets owned by the authenticated user.

```bash
curl -X GET http://localhost:8080/api/v1/tickets/my-tickets \
  -H "Authorization: Bearer <token>"
```

Success response:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
[
  {
    "ticket_id": "tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc",
    "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
    "title": "Printer Rusak",
    "description": "Printer di ruang HR tidak bisa connect ke PC",
    "image": "tickets/tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc/image-4b638c13-36c4-4ff1-bc60-2935873d6f9a.jpg",
    "status": "OPEN",
    "created_at": "2026-01-11T14:29:24.299+07:00",
    "updated_at": "2026-01-11T14:29:24.299+07:00",
    "deleted_at": null,
    "owner": null
  }
]
```

### 5. Get All Tickets, Admin Only

Returns all tickets and preloads each ticket owner.

```bash
curl -X GET http://localhost:8080/api/v1/tickets/admin/all \
  -H "Authorization: Bearer <admin-token>"
```

Success response:

```http
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "Jumlah": 1,
  "tickets": [
    {
      "ticket_id": "tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc",
      "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
      "title": "Printer Rusak",
      "description": "Printer di ruang HR tidak bisa connect ke PC",
      "image": "tickets/tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc/image-4b638c13-36c4-4ff1-bc60-2935873d6f9a.jpg",
      "status": "OPEN",
      "created_at": "2026-01-11T14:29:24.299+07:00",
      "updated_at": "2026-01-11T14:29:24.299+07:00",
      "deleted_at": null,
      "owner": {
        "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
        "username": "jane",
        "email": "jane@example.com",
        "role": "USER",
        "profile_pict": "",
        "department": "IT department",
        "remote": false,
        "created_at": "2026-01-11T14:21:21.11+07:00",
        "updated_at": "2026-01-11T14:21:21.11+07:00"
      }
    }
  ]
}
```

Common error response when the token is valid but not admin:

```json
{
  "error": "Insufficient permissions"
}
```

### 6. Update Ticket Status, Admin Only

Updates a ticket status. Allowed values are:

- `OPEN`
- `IN_PROGRESS`
- `DONE`

```bash
curl -X PUT http://localhost:8080/api/v1/tickets/admin/tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <admin-token>" \
  -d '{
    "status": "DONE"
  }'
```

Success response:

```http
HTTP/1.1 202 Accepted
Content-Type: application/json
```

```json
{
  "message": "update success",
  "ticket": {
    "ticket_id": "tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc",
    "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
    "title": "Printer Rusak",
    "description": "Printer di ruang HR tidak bisa connect ke PC",
    "image": "tickets/tick-1152cdd4-9b79-4b82-bd53-e9dc215ae9dc/image-4b638c13-36c4-4ff1-bc60-2935873d6f9a.jpg",
    "status": "DONE",
    "created_at": "2026-01-11T14:29:24.299+07:00",
    "updated_at": "2026-01-11T14:47:14.184+07:00",
    "deleted_at": null,
    "owner": {
      "user_id": "USER-2415f042-fd6e-4e61-b467-b50e0fd9e122",
      "username": "jane",
      "email": "jane@example.com",
      "role": "USER",
      "profile_pict": "",
      "department": "IT department",
      "remote": false,
      "created_at": "2026-01-11T14:21:21.11+07:00",
      "updated_at": "2026-01-11T14:21:21.11+07:00"
    }
  }
}
```

Common error responses:

```json
{
  "error": "Status must be  'DONE' or 'IN_PROGRESS'"
}
```

```json
{
  "error": "Ticket tidak ditemukan"
}
```

## Usage Flow

1. Register a normal user with `role` set to `USER`.
2. Register an admin user with `role` set to `ADMIN`.
3. Login as the normal user and save the returned token.
4. Create a ticket with the user token.
5. List the user's tickets with `/api/v1/tickets/my-tickets`.
6. Login as the admin user and save the returned token.
7. View all tickets with `/api/v1/tickets/admin/all`.
8. Update a ticket status with `/api/v1/tickets/admin/:id/status`.

## Implementation Notes

- `Register` checks duplicate users by email.
- `Register` accepts `role` from the request, so an admin account can currently be created through the public register endpoint.
- `Login` looks up users by `username`.
- User IDs are generated with the role prefix, for example `USER-<uuid>`.
- Ticket IDs are generated with the `tick-<uuid>` format.
- Ticket creation validates only `title` and `description` in the use case.
- The `validate` struct tags are present on request structs, but no validation middleware/library is currently called.
- `GET /api/v1/tickets/admin/all` returns an object with `Jumlah` and `tickets`.
- Multipart ticket creation uploads the `image` file to MinIO and stores the object key in the ticket `image` field.
- JSON ticket creation still accepts `image` as a plain string for clients that do not upload a file.
