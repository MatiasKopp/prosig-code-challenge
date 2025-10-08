# BUILD
FROM golang:1.25.2-alpine AS builder

RUN apk add --no-cache build-base sqlite-dev
WORKDIR /app

COPY ./src .
RUN go mod download
COPY posts.db /app

RUN CGO_ENABLED=1 GOOS=linux go build -o server ./cmd/api/main.go

# RUNTIME
FROM alpine:latest

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    addgroup -S guestgroup && \
    adduser -S guestuser -G guestgroup

RUN apk add --no-cache sqlite

WORKDIR /app

# SET ENV VARIABLES
ENV APP_PORT=":8080"
ENV DB_LOCATION="/app/posts.db"

COPY --from=builder /app/posts.db .
COPY --from=builder /app/server .

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./server"]
