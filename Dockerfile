# JS build stage
FROM oven/bun:alpine AS js-builder

WORKDIR /js

COPY package.json ./
RUN bun install

COPY web/assets/scripts/gen/ web/assets/scripts/gen/

RUN bun build web/assets/scripts/gen/bundle_entry.js \
    --outfile web/assets/scripts/proto_bundle.js \
    --format esm --minify

# Go build stage
FROM golang:alpine AS builder

WORKDIR /build

RUN apk add --no-cache protobuf
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=js-builder /js/web/assets/scripts/proto_bundle.js web/assets/scripts/proto_bundle.js

RUN protoc \
    --proto_path=api/proto \
    --go_opt=paths=source_relative \
    --go_out=internal/generated \
    documents.proto order.proto prices.proto status.proto

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/app
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o size_import ./cmd/size_import

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /build/server .
COPY --from=builder /build/size_import .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

COPY --from=builder /build/web/ web/
COPY api/ api/
COPY migrations/ migrations/

RUN mkdir -p logs data/contracts data/bills data/offers data/documents data/sizes

EXPOSE 8080

CMD ["./server"]
