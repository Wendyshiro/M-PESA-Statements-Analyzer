# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p uploads && chmod -R 755 uploads

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install pdftotext and qpdf for PDF processing
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    poppler-utils \
    qpdf && \
    update-ca-certificates

# Create non-root user
RUN adduser -D -u 1001 appuser

# Copy binary and migrations
COPY --from=builder /app/main .
COPY --from=builder /app/migrations/ ./migrations/

# Create uploads directory
RUN mkdir -p uploads && chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

CMD ["./main"]