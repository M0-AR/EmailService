package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailServiceProvider interface {
	SendEmail(to, subject, body string) error
}

type EmailService struct {
	providers []EmailServiceProvider
}

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewEmailService(providers ...EmailServiceProvider) *EmailService {
	return &EmailService{
		providers: providers,
	}
}

func (es *EmailService) SendEmail(to, subject, body string) error {
	if !isValidEmail(to) {
		return errors.New("Invalid email address")
	}

	for _, provider := range es.providers {
		err := provider.SendEmail(to, subject, body)
		if err == nil {
			return nil // Email sent successfully
		}
		log.Printf("Failed to send email using provider: %v", err)
	}

	// All providers failed, log the email for later retry
	logEmailToFile(to, subject, body)

	return errors.New("all providers failed to send email")
}

type MailgunService struct {
	PrivateKey string
	Domain     string
}

func (mg *MailgunService) SendEmail(to, subject, body string) error {
	mgClient := mailgun.NewMailgun(mg.Domain, mg.PrivateKey)
	sender := "Excited User <mailgun@sandbox.mailgun.org>"
	message := mgClient.NewMessage(sender, subject, body, to)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mgClient.Send(ctx, message)
	return err
}

type SendGridService struct {
	APIKey string
}

func (sg *SendGridService) SendEmail(to, subject, body string) error {
	from := mail.NewEmail("Example User", "example@example.com")
	toEmail := mail.NewEmail("Recipient", to)
	message := mail.NewSingleEmail(from, subject, toEmail, body, body)
	client := sendgrid.NewSendClient(sg.APIKey)
	response, err := client.Send(message)

	if err != nil || (response.StatusCode < 200 || response.StatusCode >= 300) {
		return errors.New("Failed to send email via SendGrid")
	}
	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(".+@.+\\..+")
	return re.MatchString(email)
}

func logEmailToFile(to, subject, body string) {
	// Opening or creating the log file
	file, err := os.OpenFile("failed_emails.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening or creating log file:", err)
		return
	}
	defer file.Close()

	// Logging the email details
	logEntry := fmt.Sprintf("Time: %s, To: %s, Subject: %s, Body: %s\n", time.Now().Format(time.RFC3339), to, subject, body)
	_, err = file.WriteString(logEntry)
	if err != nil {
		log.Println("Error writing to log file:", err)
	}
}

func main() {
	mailgun := &MailgunService{
		PrivateKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		Domain:     os.Getenv("MAILGUN_DOMAIN"),
	}
	sendGrid := &SendGridService{
		APIKey: os.Getenv("SENDGRID_API_KEY"),
	}
	emailService := NewEmailService(mailgun, sendGrid)

	r := mux.NewRouter()
	r.HandleFunc("/send-email", func(w http.ResponseWriter, r *http.Request) {
		var emailRequest EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&emailRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := emailService.SendEmail(emailRequest.To, emailRequest.Subject, emailRequest.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Email sent successfully"))
	}).Methods("POST")

	http.Handle("/", r)
	fmt.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
