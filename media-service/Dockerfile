# ---- Build stage ----
FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /media-service ./cmd/serverd
# ---- Runtime stage ----
FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app
COPY --from=build /media-service /app/media-service
COPY .env /app/.env
EXPOSE 8005
USER nonroot
ENTRYPOINT ["/app/media-service"]
