# Email Service

A generic email service built in Golang, using Mailgun as the primary email provider with MailGun and SendGrid as fallback services.

## Approach
The service is designed to abstract the email sending process. It first tries to send an email via SendGrid. 

## Testing
Run the tests using the following command: go test

You can now use Postman or any other API client to send a POST request to http://localhost:8080/send-email with a JSON payload containing the to, subject, and body fields.
{
    "to": "s176492@student.dtu.dk",
    "subject": "Hello",
    "body": "Testing email sending!"
}

## Docker
Build the Docker image: docker build -t email-service .

