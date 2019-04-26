FROM golang:1.9-alpine3.7 AS build

WORKDIR /go/src/github.com/nikhil-github/sms-api

RUN apk add --no-cache \
            bash~=4.4 \
            git~=2.15 \
            make~=4.2 \
    rm -rf /var/cache/apk/*

RUN go get -u github.com/golang/dep/cmd/dep

COPY Gopkg.toml Gopkg.lock Makefile ./

COPY . ./
RUN make build-all

FROM alpine AS release

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/nikhil-github/sms-api/sms-api /go/bin/sms-api

CMD ["/go/bin/sms-api"]
