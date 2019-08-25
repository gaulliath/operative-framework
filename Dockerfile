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
FROM alpine:3.10 AS runtime

# Install tini to /usr/local/sbin
ADD https://github.com/krallin/tini/releases/download/v0.18.0/tini-muslc-amd64 /usr/local/sbin/tini

# Install runtime dependencies & create runtime user
RUN \
	apk update && \
	apk add --no-cache --no-progress ca-certificates && \
	rm -rf /var/cache/apk/* && \
		\
		chmod +x /usr/local/sbin/tini && \
		mkdir -p /opt && \
 			\
	 		adduser -D opf -h /opt/operative-framework -s /bin/sh && \
 			su opf -c 'cd /opt/operative-framework; mkdir -p bin config data services'

# Switch to user context
USER opf
WORKDIR /opt/operative-framework

COPY --from=builder /opf /opt/operative-framework/bin/opf
ENV PATH $PATH:/opt/operative-framework/bin

# Container configuration
EXPOSE 8888
VOLUME ["/opt/operative-framework/data"]
# ENTRYPOINT ["tini", "-g", "--"]
ENTRYPOINT [ "/opt/operative-framework/bin/opf" ]
# CMD [ "/opt/operative-framework/bin/opf" ]