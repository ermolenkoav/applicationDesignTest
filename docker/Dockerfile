ARG GO_VER=1.22.6
ARG ALPINE_VER=3.20
ARG VERSION=0.0.1

FROM golang:${GO_VER}-alpine${ALPINE_VER} AS builder
LABEL stage=gobuilder
RUN apk update --no-cache && apk add --no-cache tzdata
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -ldflags="-s -w" -o app ./cmd/app/main.go

# Run the tests in the container
FROM builder
RUN go test -v ./...

FROM alpine:${ALPINE_VER}
COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York
ENV TZ=America/New_York
RUN apk update --no-cache && apk add --no-cache ca-certificates
WORKDIR /app
EXPOSE 8080
#COPY --from=builder /src/api api/.
COPY --from=builder /src/app .
CMD ["./app"]