# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Create necessary directories and set permissions
RUN mkdir -p uploads dashboard \
    && chmod -R 755 uploads dashboard \
    && touch output.txt \
    && chmod 666 output.txt

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install required system dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    poppler-utils && \
    update-ca-certificates

# Create a non-root user
RUN adduser -D -u 1001 appuser \
    && chown -R appuser:appuser /app

# Copy the binary and necessary files from builder
COPY --from=builder /app/main .
COPY --from=builder /app/output.txt .
COPY --from=builder /app/uploads/ ./uploads/
COPY --from=builder /app/dashboard/ ./dashboard/
COPY --from=builder /app/handlers/ ./handlers/
COPY --from=builder /app/models/ ./models/
COPY --from=builder /app/utils/ ./utils/

# Set ownership of copied files
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
