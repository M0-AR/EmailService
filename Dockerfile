# https://hub.docker.com/_/golang
FROM golang:1.17 as builder

# Copy local code to the container image.
WORKDIR /go/src/app
COPY . .

# Build the binary.
RUN go get -d -v ./...
RUN go build -v -o app

# Use a minimal alpine image for the runtime container
FROM alpine:3.14

# Copy the binary to the production image from the builder stage.
COPY --from=builder /go/src/app/app /app

# Run the web service on container startup.
CMD ["/app"]
