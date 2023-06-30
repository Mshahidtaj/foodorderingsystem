# Stage 1: Build stage
FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Stage 2: Final stage
FROM alpine:latest

RUN addgroup -S app && adduser -S -G app app

RUN apk --no-cache add ca-certificates

WORKDIR /home/app

COPY --from=builder /app/app .

RUN chown -R app:app /home/app

USER app

EXPOSE 8080

CMD ["./app"]
