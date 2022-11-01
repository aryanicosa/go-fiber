# Building GO USER and BOOK CRUD rest api using Fiber, PostgreSQL and redis

TODO list:
- refactor db connection
- create correct docker compose file
- create correct makefile
- create correct dockerfile
- create correct testing
- make sure every api work as expected

Run Service Locally
1. git clone git@github.com:aryanicosa/go-fiber.git
2. Install Docker and the following useful Go tools to your system:
   - golang-migrate/migrate for apply migrations
   - github.com/swaggo/swag for auto-generating Swagger API docs
   - github.com/securego/gosec for checking Go security issues
   - github.com/go-critic/go-critic for checking Go the best practice issues
   - github.com/golangci/golangci-lint for checking Go linter issues
3. swag init
4. make run-dependencies
5. go run main.go
