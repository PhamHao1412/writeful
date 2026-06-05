# ---- Build stage ----
FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /content-service ./cmd/serverd
# ---- Runtime stage ----
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /content-service /app/content-service
COPY .env /app/.env
EXPOSE 8003
USER nonroot
ENTRYPOINT ["/app/content-service"]
