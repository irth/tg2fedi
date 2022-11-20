FROM golang:1.19-alpine AS builder
RUN apk add --no-cache make
WORKDIR /src
COPY ./go.mod ./go.sum /src/
RUN go mod download
COPY . /src/
RUN make build-linux

FROM alpine:latest
COPY --from=builder /src/bin/tg2fedi /tg2fedi
WORKDIR /

CMD ["/feeder"]
