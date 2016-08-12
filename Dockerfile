FROM golang:1.6.3-alpine

COPY . /go/src/github.com/VoltFramework/volt
WORKDIR /go/src/github.com/VoltFramework/volt

RUN apk update

RUN set -ex \
    && apk add --no-cache --virtual .build-deps git \
    && go get github.com/jteeuwen/go-bindata/... \
    && go get github.com/elazarl/go-bindata-assetfs/... \
    && go-bindata-assetfs -pkg api static/... \
    && mv bindata_assetfs.go api/static.go \
    && go install --ldflags '-extldflags "-static"' \
    && apk del .build-deps

ENTRYPOINT ["/go/bin/volt"]