############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
COPY config.yml /go/bin/config.yml
RUN go mod download
RUN go mod verify
ENV CGO_ENABLED=0
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/hello
############################
# STEP 2 build a small image
############################
FROM scratch
# Import from builder.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/hello /go/bin/hello
COPY --from=builder /go/bin/config.yml /go/bin/config.yml
# Use an unprivileged user.
# Run the hello binary.
WORKDIR "/go/bin/"
ENTRYPOINT ["/go/bin/hello"]
