package main

import (
	"context"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type EmailService struct {
	MailgunPrivateAPIKey string
	MailgunDomain        string
	SendGridAPIKey       string
	SparkPostAPIKey      string
}

func main() {
	emailService := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
		SendGridAPIKey:       os.Getenv("SENDGRID_API_KEY"),
	}

	success := emailService.SendEmail("s176492@student.dtu.dk", "Hello", "Testing some Mailgun awesomeness!")
	if success {
		fmt.Println("Email sent successfully")
	} else {
		fmt.Println("Email sending failed")
	}
}

func (es *EmailService) SendEmail(to string, subject string, body string) bool {
	success := es.sendViaMailgun(to, subject, body)

	if !success {
		success = es.sendViaSendGrid(to, subject, body)
	}

	return success
}

func (es *EmailService) sendViaMailgun(to string, subject string, body string) bool {
	mg := mailgun.NewMailgun(es.MailgunDomain, es.MailgunPrivateAPIKey)
	sender := "Excited User <mailgun@sandbox24c6d43652a24e1c8737bd3b91d2a958.mailgun.org>"
	recipient := to

	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, _, err := mg.Send(ctx, message)

	fmt.Printf("Response: %s\n", resp)
	if err != nil {
		log.Println("Mailgun: ", err)
		return false
	}

	return true
}

func (es *EmailService) sendViaSendGrid(to string, subject string, body string) bool {
	from := mail.NewEmail("Example User", "s176492@student.dtu.dk")
	toEmail := mail.NewEmail("Recipient", to)
	message := mail.NewSingleEmail(from, subject, toEmail, body, body)
	client := sendgrid.NewSendClient(es.SendGridAPIKey)
	response, err := client.Send(message)

	if err != nil {
		log.Println("SendGrid:", err)
		return false
	}

	fmt.Println("Status Code:", response.StatusCode)
	fmt.Println("Body:", response.Body)
	fmt.Println("Headers:", response.Headers)

	return response.StatusCode >= 200 && response.StatusCode < 300 // Check for HTTP success status codes (2xx)
}
