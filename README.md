# Building GO rest api boilerplate using Fiber, PostgreSQL, redis and docker

TODO list:
- Database
   - Simplify database connection and init to be once for all
- Controller
   - refactor code to use token expire and credential checker utility (p2)
   - add basic auth to user and books api
- Utility
   - Create token expire and credential checker (p2)
- Testing
   - create clean and correct testing (p1) : only open testing db once
   - automate testing with Github action (p1)
- Swagger doc
   - generate correct and easy to use swagger api (p1)
   - consistence response body and response error (done)
- Dockerfile
   - create correct dockerfile (p1)
- Makefile
   - create correct makefile (p1)
- Docker-compose
   - create correct docker compose file
   - create docker-compose-test.yml file (done)
   - create docker-compose-dependencies.yml file for local development (done)
- Response
  - create constants for type of error response
- Deploy and host go-fiber web app
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
