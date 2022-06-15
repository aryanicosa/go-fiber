FROM golang:1.15-alpine
LABEL maintainer="aryanicosa"

WORKDIR /go/src/app
COPY . .

RUN go build -o go-fiber

EXPOSE 8080
CMD ["./main"]