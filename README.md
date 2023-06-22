# Email Service

A generic email service built in Golang, using Mailgun as the primary email provider with SparkPost and Amazon SES as fallback services.

## Approach
The service is designed to abstract the email sending process. It first tries to send an email via Mailgun. If this fails, it automatically switches to SparkPost, and finally to Amazon SES if necessary.

## Testing
Run the tests using the following command: go test

## Docker
Build the Docker image: docker build -t email-service .

