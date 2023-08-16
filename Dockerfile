FROM golang:1.20-alpine AS builder

LABEL author="Raihan hamdani" \
      title="todolist_api" \
      website="https://github.com/reyhanhmdani/Todo_Gin_Gorm"


RUN apk update && apk add --no-cache git

WORKDIR /app


# Copy the Go modules files to the container
COPY go.mod go.sum ./

RUN go mod download && go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o API .

FROM alpine:latest

# Copy MySQL configuration file
#COPY my.cnf /etc/mysql/my.cnf

WORKDIR /app

COPY --from=builder /app/API .
COPY --from=builder /app/database/migrations/ ./database/migrations/

CMD ["./API"]
