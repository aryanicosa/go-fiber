# dockerfile is use fro building an image
FROM golang:1.17-alpine AS builder

LABEL maintainer="aryanicosa"

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download -x

# Copy the code into the container.
COPY . .

COPY sql $HOME/sql

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o server . # run is the same as run a command in terminal

FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/server", "/build/.env", "/"]

# Command to run when the container started.
#using ENTRYPOINT
# --> run main command/binary. after start will search additional argument
ENTRYPOINT ["/server"]

#using CMD, run any command when container start using full comment
# --> the same with running 'go run main.go' in terminal
# CMD ["go", "run", "main.go"]
