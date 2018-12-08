FROM golang:alpine

ENV CGO_ENABLED 0
ENV GOOS linux

WORKDIR /go/src/migrate

COPY main.go Gopkg.lock Gopkg.toml /go/src/migrate/

RUN apk --update add curl git && \
	mkdir -p /go/bin && \
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
	dep ensure

ENTRYPOINT ["/usr/local/go/bin/go"]
