FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY . .
COPY /config /tempEnv
COPY /migrations /tempMigration
RUN go build -o . main.go

FROM alpine
WORKDIR /build
COPY --from=builder /build/main /build/main
COPY --from=builder /tempEnv /build/config
COPY --from=builder /tempMigration /build/migrations
CMD ["./main"]