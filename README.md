# Building GO rest api using Fiber, PostgreSQL, redis and docker

# Run Service Locally
1. git clone [this repo](https://github.com/aryanicosa/go-fiber)
2. Install Make, Docker, Docker Compose and the following useful Go tools to your system:
   - golang-migrate/migrate for apply migrations
   - github.com/swaggo/swag for auto-generating Swagger API docs
   - github.com/securego/gosec for checking Go security issues
   - github.com/go-critic/go-critic for checking Go the best practice issues
   - github.com/golangci/golangci-lint for checking Go linter issues
3. swag init
4. docker-compose up --build

**Result**
![img.png](img.png)


# Development Flow, run with testing:
- Create some changes
- in terminal run: docker-compose up

- Deployment and host go-fiber web app
  - CI/CD - github action
  - Docker image saved to Docker Hub
  - aws -> http://ec2-52-91-120-237.compute-1.amazonaws.com/swagger/
  ![img_1.png](img_1.png)
