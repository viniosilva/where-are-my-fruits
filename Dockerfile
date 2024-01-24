# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server

# Deploy the application binary into a lean image
FROM scratch

WORKDIR /app

COPY --from=build /server .

EXPOSE 3001

USER nonroot:nonroot

ENTRYPOINT ["/server"]
