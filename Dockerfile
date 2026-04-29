# --- Stage 1: Build ---
FROM golang:1.26 AS builder

WORKDIR /app

#COPY go.mod go.sum ./
#RUN go mod download

COPY . .

# Build the binary
# CGO_ENABLED=0 ensures a static binary for the lightweight alpine/scratch image
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

# --- Stage 2: Final Run Image ---
FROM alpine:latest  

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]