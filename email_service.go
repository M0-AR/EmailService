package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

type EmailService struct {
	MailgunPrivateAPIKey string
	MailgunDomain        string
	SparkPostAPIKey      string
	AWSAccessKey         string
	AWSSecretKey         string
	AWSRegion            string
}

func main() {
	emailService := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
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

	//if !success {
	//	success = es.sendViaSparkPost(to, subject, body)
	//}


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