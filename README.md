# Building GO USER and BOOK CRUD rest api using Fiber, PostgreSQL and redis

TODO list:
1. Testing
  - create correct testing
2. Swagger doc
  - create correct and easy to use swagger api
3. Dockerfile
  - create correct dockerfile
4. Makefile
   - create correct makefile
5. Docker-compose
  - create correct docker compose file
  - create docker-compose-test.yml file
6. Deploy and host go-fiber web app
   - CI/CD
   - Heroku

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
