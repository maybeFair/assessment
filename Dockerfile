FROM golang:1.19.1-alpine AS builder

RUN apk add build-base

WORKDIR /app

ENV CGO_ENABLED=0
ENV TZ=Asia/Bangkok

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build ./app/server.go

FROM alpine:3.16

COPY --from=builder /app/server /server
COPY .git .git

EXPOSE 2565

CMD ["/server"]
