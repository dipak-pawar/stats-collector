FROM golang:1.12.6 AS build-env

WORKDIR /go/src/github.com/dipak-pawar/stats-collector

ENV GO111MODULE=on

COPY . /go/src/github.com/dipak-pawar/stats-collector

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o metrics


FROM alpine
LABEL maintainer "Dipak Pawar <dipakpawar231@gmail.com>" \
      author "Dipak Pawar <dipakpawar231@gmail.com>"

RUN apk --no-cache add ca-certificates

EXPOSE 8080

COPY --from=build-env /go/src/github.com/dipak-pawar/stats-collector/metrics /usr/local/bin/
USER 10001

ENTRYPOINT [ "/usr/local/bin/metrics" ]
