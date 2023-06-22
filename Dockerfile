# Use the official Golang image to create a build artifact.
FROM golang:1.17 as builder

# Copy local code to the container image.
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . ./

# Build the command inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o email_service

# Use a Docker multi-stage build to create a lean production image.
FROM gcr.io/distroless/base-debian10
COPY --from=builder /app/email_service /email_service

# Run the web service on container startup.
CMD ["/email_service"]
