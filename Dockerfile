FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server cmd/server/main.go

FROM alpine:latest

# Install any necessary system dependencies
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the build stage
COPY --from=build /app/server .

# Set the binary as the entrypoint
ENTRYPOINT ["./server"]