package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"net/http"
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

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	emailService := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
		SendGridAPIKey:       os.Getenv("SENDGRID_API_KEY"),
	}

	r := mux.NewRouter()
	r.HandleFunc("/send-email", func(w http.ResponseWriter, r *http.Request) {
		var emailRequest EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&emailRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		success := emailService.SendEmail(emailRequest.To, emailRequest.Subject, emailRequest.Body)
		if success {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Email sent successfully"))
		} else {
			http.Error(w, "Email sending failed", http.StatusInternalServerError)
		}
	}).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
