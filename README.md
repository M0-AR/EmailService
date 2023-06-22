# Email Service

A generic email service built in Golang, using Mailgun as the primary email provider with MailGun and SendGrid as fallback services.

## Approach
The service is designed to abstract the email sending process. It first tries to send an email via SendGrid. 

## Testing
Run the tests using the following command: go test

## Docker
Build the Docker image: docker build -t email-service .

