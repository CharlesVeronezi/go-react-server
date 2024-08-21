FROM golang:alpine as build

WORKDIR /app

# Copy Go module files
COPY go.* ./

# Download dependencies
RUN go mod download

RUN go generate ./...

# Copy source files
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY .env .env

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/gosocket ./cmd/goserver/main.go

FROM alpine:3.14.10

EXPOSE 3000

# Copy files from builder stage
COPY --from=build /app/bin/gosocket .

# Increase GC percentage and limit the number of OS threads
ENV GOGC 1000
ENV GOMAXPROCS 3

ENV GS_DATABASE_PORT=5432
ENV GS_DATABASE_USER="postgres"
ENV GS_DATABASE_PASSWORD="pgpassword"
ENV GS_DATABASE_NAME="goserver"
ENV GS_DATABASE_HOST="db"

# Run binary
CMD ["/gosocket"]