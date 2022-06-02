FROM golang:1.18.3 AS build-image
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/pawang-be
EXPOSE 1234
CMD ["pawang-be"]