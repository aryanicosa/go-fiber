# Building GO rest api boilerplate using Fiber, PostgreSQL, redis and docker

TODO list:
1. Controller
   - Create consistence response model
2. Utility
   - Create generic error response function handler
3. Testing
   - create correct testing
   - automate testing
4. Swagger doc
  - generate correct and easy to use swagger api
  - consistence response body and response error
5. Dockerfile
  - create correct dockerfile
6. Makefile
   - create correct makefile
7. Docker-compose
  - create correct docker compose file
  - create docker-compose-test.yml file
8. Deploy and host go-fiber web app
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

Run Testing // TODO
