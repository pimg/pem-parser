############################
# STEP 1 build executable binary
############################
FROM golang:1.23.2-alpine@sha256:9dd2625a1ff2859b8d8b01d8f7822c0f528942fe56cfe7a1e7c38d3b8d72d679 AS go-build
# Install build tools.
RUN apk add --update git

# Add non-privileged user
RUN adduser \
    -h "/home/appuser" \
    -g "" \
    -s "/sbin/nologin" \
    -D \
    -H \
    -u 1001 \
    appuser

COPY . /go/src/pem-parser
ENV GO111MODULE=on
WORKDIR /go/src/pem-parser
RUN go mod download

WORKDIR /go/src/pem-parser

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/pem-parser
ENTRYPOINT ["/go/bin/pem-parser"]

############################
# STEP 2 build a small image
############################
FROM scratch

# copy appuser passwd file
COPY --from=go-build /etc/passwd /etc/passwd
COPY --from=go-build /etc/group /etc/group

COPY --from=go-build /go/bin/pem-parser /usr/local/bin/pem-parser

WORKDIR /app

USER appuser

CMD ["/usr/local/bin/pem-parser", "serve"]

