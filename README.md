# Email Service

A generic email service built in Golang, using Mailgun as the primary email provider and SendGrid as fallback services.

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

------------------
### The primary features and characteristics of this solution include:

\textbf{Multiple Email Service Providers}: The solution incorporates the use of multiple email service providers, namely Mailgun and SendGrid. It uses an interface (EmailServiceProvider) to define a common contract for all email service providers, allowing the addition of more providers easily.

\textbf{Fallback Mechanism}: When an attempt to send an email fails with one service provider, it automatically falls back to the next provider in the list. This ensures that the system has higher availability and reliability by not being dependent on a single service provider.

\textbf{Input Validation}: The recipient's email address is validated using a regular expression before attempting to send the email. This helps in avoiding unnecessary API calls to the email service providers.

\textbf{Error Handling}: If sending the email fails with all available providers, the application logs the email details to a file (failed_emails.log) for later review or retry. This ensures that no message is lost without trace in case of consistent failure across all service providers.

\textbf{Context Timeout}: When interacting with email service providers' APIs, the solution uses a context with a timeout to ensure that the service doesn’t hang indefinitely in case the external API is slow or unresponsive.

\textbf{Secret Management}: The solution uses environment variables to store sensitive data such as API keys. This is considered a good practice as it doesn't hard code sensitive data and allows for easier configuration across different environments.

\textbf{HTTP Endpoint}: The application exposes an HTTP endpoint (/send-email) that accepts POST requests with JSON payloads containing the recipient's email address, subject, and body. This endpoint could be integrated into other systems as a way to send emails through this service.

\textbf{Modular and Scalable Structure}: The use of interfaces and separation of concerns among different components of the solution means that it can be easily extended and scaled.

\textbf{File-based Logging}: Failed email sending attempts are logged to a file, which can be useful for troubleshooting and ensuring that no data is lost. However, it's important to note that in a production environment, a more sophisticated logging and monitoring solution might be preferred.

Overall, this solution is a robust starting point for an email sending service, with a focus on reliability through the use of multiple service providers and fallback mechanisms. Additionally, the modular structure ensures that it can be extended and customized according to the requirements. However, for production use, considerations regarding logging, monitoring, and possibly a more persistent queue system for handling failed emails would be advisable.
