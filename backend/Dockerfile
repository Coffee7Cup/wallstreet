FROM go:1.25.5-alpine 

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build statically
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 2. Production Stage: Minimal Alpine runtime
FROM alpine:3.19

WORKDIR /app

# Copy compiled binary from builder stage
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
