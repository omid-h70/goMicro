# base go image
FROM golang:1.18 as builder

RUN mkdir /app
RUN mkdir /app/vendor

COPY cmd/api /app
COPY go.mod  /app
COPY go.sum  /app

WORKDIR /app

ENV GONOPROXY=github.com/*
ENV GONOSUMDB=github.com/*

RUN CGO_ENABLED=0 go build -mod=readonly -o brokerApp

RUN chmod +x /app/brokerApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/brokerApp /app

CMD [ "/app/brokerApp" ]