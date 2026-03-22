ARG VERSION=0.0.1
ARG PORT=8080

# Build the binary
FROM golang:alpine AS builder
LABEL stage=gobuilder
RUN apk update --no-cache && apk add --no-cache tzdata
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /src
COPY . .
RUN go mod download

# Run tests — build fails here if tests fail
RUN go test -v ./...

RUN go build -ldflags="-s -w" -o app ./cmd/app/main.go

# Produce the final image
FROM alpine:latest
COPY --from=builder /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
ENV TZ=Europe/Moscow
RUN apk update --no-cache && apk add --no-cache ca-certificates
WORKDIR /app
EXPOSE ${PORT}
COPY --from=builder /src/app .
CMD ["./app"]