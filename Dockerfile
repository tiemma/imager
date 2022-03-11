FROM golang:1.16-alpine AS builder

LABEL maintainer="me@bakman.build"
LABEL version="v1"

ENV GIT_SSL_NO_VERIFY 1
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app
ADD . .

RUN go build -trimpath -a -o imager -ldflags="-w -s" cmd/*.go

FROM alpine:latest AS production
COPY --from=builder /app/imager imager

RUN chmod +x imager
CMD ["./imager"]