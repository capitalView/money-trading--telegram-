FROM golang:1.22.4-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git gcc gettext musl-dev

COPY ["app/go.mod", "app/go.sum", "./"]

RUN go mod download

COPY app ./
RUN go build -o ./bin/app main.go

FROM alpine as runner

COPY --from=builder  /usr/local/src/bin/app /

COPY .env /

CMD ["/app"]
