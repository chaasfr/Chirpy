# Chirpy

Chirpy is a Go-based backend service for a microblogging platform. It provides endpoints for serving static content, managing user accounts and chirps, and administrative tasks. The project uses PostgreSQL as its database, with schema migrations managed by [Goose](https://github.com/pressly/goose).

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Database Migrations](#database-migrations)
- [API Endpoints](#api-endpoints)
  - [/app](#app)
  - [/admin](#admin)
  - [/api](#api)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Features

- **Go-Powered Backend:** High performance API service written in Go.
- **User & Chirp Management:** Create and update users, post and manage chirps.
- **JWT Authentication:** Secure endpoints with JWT and refresh tokens.
- **Admin Tools:** Metrics reporting and a development-only endpoint for resetting the database.
- **Static File Serving:** A simple static file server available at `/app`.

## Prerequisites

Before running Chirpy locally, ensure you have the following installed:

- [Go](https://golang.org/) (latest version recommended)
- [PostgreSQL](https://www.postgresql.org/)
- [Goose](https://github.com/pressly/goose) for database migrations

**For contributions:**  
- [SQLC](https://sqlc.dev/) is recommended to generate type-safe Go code from your SQL queries.

## Installation & Setup

1. **Clone the Repository**

```bash
git clone https://github.com/chaasfr/chirpy.git
cd chirpy
```

2. **Configure Environment Variables**
Copy the .env_template to .env and fill in the necessary values, especially:

- USERNAME – Your PostgreSQL username.
- PASSWORD – Your PostgreSQL password.
- JWT_SECRET – A secret key for signing JWTs.
```bash
cp .env_template .env
```

3. **Run the Application**
After setting up the environment and running the database migrations (see below), start the service:
```bash
go run .
```

## Database Migrations
Before running the application, ensure your database schema is up-to-date:

1. **Navigate to the SQL schema directory:**
```bash
cd sql/schema
```

2. **Run the migrations using Goose:**
CONNECTION_URL is the same as in your .env (without `?sslmode=disabled`)
```bash
goose CONNECTION_URL up
```
## API Endpoints
Chirpy organizes its endpoints into three main groups:

### /app
Server static index file.

### /admin
**GET /admin/metrics**

Returns metrics, including insights such as how many times the fileserver was visited.

**POST /admin/reset**

Wipes the entire database.
Note: This endpoint is intended for use only in the development environment (controlled via the .env configuration).

### /api
**POST /api/users**

Creates a new user.


**PUT /api/users**

Updates a user's email and password.
Authentication: Requires the user to be logged in with a valid JWT.


**POST /api/chirps**

Creates a new chirp.
Authentication: Requires a valid JWT.


**GET /api/chirps**

Retrieves all chirps.
Query Parameters:
- author_id – Filter chirps by a specific user.
- sort – Sort chirps by created_at in either ascending (asc) or descending (desc) order.


**GET /api/chirps/{chirpID}**

Retrieves a single chirp by its ID.


**DELETE /api/chirps/{chirpID}**

Deletes a chirp.

Authentication: Requires a valid JWT and the chirp must belong to the authenticated user.

**POST /api/login**
Authenticates a user and returns a valid JWT along with a refresh token upon success.


**POST /api/refresh**
Provides a new JWT when a correct refresh token is provided.


**POST /api/revoke**
Revokes a refresh token.