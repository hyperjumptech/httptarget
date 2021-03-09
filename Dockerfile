FROM golang:1.16
ENV GOPATH /go
ENV GO111MODULE on
ENV GOOS linux
ENV GOARCH amd64

# Prepare all the dirs.
RUN mkdir -p $GOPATH/src/github.com/hyperjumptech/httptarget
# Copy the build content.
COPY . $GOPATH/src/github.com/hyperjumptech/httptarget
# Checkout the go-resource to auto generate statics into go codes.
WORKDIR $GOPATH/src/github.com/hyperjumptech/httptarget
# Compile the proje ct
RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o httptarget.app cmd/main.go

# Now use the deployment image.
FROM alpine:latest
ENV GOPATH /go
RUN apk --no-cache add ca-certificates
# Copy the built binary to the new image.
WORKDIR /root/
COPY --from=0 $GOPATH/src/github.com/hyperjumptech/httptarget/httptarget.app .
# Expose port.
EXPOSE 51423
# Execute
CMD ["./httptarget.app"]