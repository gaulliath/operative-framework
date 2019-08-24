############################
# Builder image
############################
ARG GOLANG_BUILDER_VERSION=1.13rc1-alpine
FROM golang:${GOLANG_BUILDER_VERSION} AS builder

# Here's a oneliner for your Dockerfile that fails if the Alpine image is vulnerable.
# RUN apk add --no-network --no-cache --repositories-file /dev/null "apk-tools>2.10.1"

# install pre-requisites
RUN apk update && \
	apk add --no-cache --no-progress build-base git tzdata ca-certificates sqlite-dev && \
	update-ca-certificates && \
	go get github.com/Masterminds/glide && \
	go get github.com/golang/dep/cmd/dep

# copy sources
COPY . /go/src/github.com/graniet/operative-framework
WORKDIR /go/src/github.com/graniet/operative-framework

# fetch dependencies
# RUN yes no | glide create && glide install --strip-vendor && go build -o /opf .
# RUN go get -d -v ./...
RUN glide install --strip-vendor && \
		GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o /opf .

############################
# Runtime image
############################
FROM alpine:3.10
EXPOSE 8888
RUN apk update && \
	apk add --no-cache --no-progress ca-certificates && \
	rm -rf /var/cache/apk/*

COPY --from=builder /opf /opf
ENTRYPOINT [ "/opf" ]

