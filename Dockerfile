FROM golang:1.22-alpine3.20 AS builder

ARG CGO_ENABLED=0

WORKDIR /tmp/server

# Copy source files necessary to download dependencies
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Copy source files required for build
COPY main.go ./
RUN go build -o server .

FROM alpine:3.20 AS runtime

COPY --from=builder /tmp/server/server /usr/local/bin/server

RUN chmod -R 777 /usr/local/bin/server

ENTRYPOINT [ "/usr/local/bin/server" ]
