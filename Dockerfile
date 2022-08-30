# Builder
FROM golang:alpine3.16 AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o pawang-api .

# Deploy
FROM alpine:latest

WORKDIR /app

COPY .env .
COPY --from=builder ["/build/pawang-api", "/build/.env", "/build", "/"]

EXPOSE 1234

ENTRYPOINT [ "/pawang-api" ]