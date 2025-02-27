############################
# STEP 1 build executable binary
############################
FROM golang:1.24.0-alpine@sha256:2d40d4fc278dad38be0777d5e2a88a2c6dee51b0b29c97a764fc6c6a11ca893c AS go-build
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

