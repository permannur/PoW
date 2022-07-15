# Step 1: Modules caching
FROM golang:1.18-alpine3.16 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.18-alpine3.16 as builder
COPY --from=modules /go/pkg /go/pkg
COPY ./server /app/server
COPY go.mod go.sum /app/
WORKDIR /app/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/tcp/

# Step 3: Final
FROM scratch
COPY --from=builder /app/server/config/config.yml /server/config/config.yml
COPY --from=builder /bin/app /app
CMD ["/app"]
