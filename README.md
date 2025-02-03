# Chirpy
Part of my Go learnings from boot dev.

### Install
- Install go 1.22+
- Install Postgres v15+
- Install Goose `go install github.com/pressly/goose/v3/cmd/goose@latest`
- Install SQLC `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- Create a .env file at the root of the project containing `DB_URL="postgres postgres://username:password@localhost:5432/chirpy?sslmode=disable"`


### How to
In a terminal, start the server `go build -o out && ./out` or `go run .` then go to your [localhost:8080](http://localhost:8080/)