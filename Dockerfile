# Build stage
FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/app

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY web/ web/
COPY api/ api/
COPY migrations/ migrations/

RUN mkdir -p logs data/contracts data/bills data/offers data/documents data/sizes

EXPOSE 8080

CMD ["./server"]
