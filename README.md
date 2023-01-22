# Building GO rest api boilerplate using Fiber, PostgreSQL, redis and docker

TODO list:
1. Controller
   - refactor code to use token expire and credential checker utility (p2)
   - add basic auth to user and books api
2. Utility
   - Create token expire and credential checker (p2)
3. Testing
   - create clean and correct testing (p1) : only open testing db once
   - automate testing with Github action (p1)
4. Swagger doc
   - generate correct and easy to use swagger api (p1)
   - consistence response body and response error (done)
5. Dockerfile
   - create correct dockerfile (p1)
6. Makefile
   - create correct makefile (p1)
7. Docker-compose
   - create correct docker compose file
   - create docker-compose-test.yml file (done)
   - create docker-compose-dependencies.yml file for local development (done)
8. Response
  - create constants for type of error response
9. Deploy and host go-fiber web app
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
