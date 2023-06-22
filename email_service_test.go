package main

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/sendgrid/sendgrid-go"
	"os"
	"testing"
)

func TestSendViaMailgun(t *testing.T) {
	es := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
	}
	// Attempt to send an email using the sendViaMailgun method
	success := es.sendViaMailgun("s176492@student.dtu.dk", "Test", "This is a test email")

	// Check if the email sending was not successful
	if !success {
		t.Error("Expected success, but got failure")
	}
}

func TestSendViaMailgunInvalidAPIKey(t *testing.T) {
	es := &EmailService{
		MailgunPrivateAPIKey: "invalid_api_key",
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
	}
	success := es.sendViaMailgun("s176492@student.dtu.dk", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to invalid API Key, but got success")
	}
}

func TestSendViaMailgunInvalidDomain(t *testing.T) {
	es := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        "invalid_domain",
	}
	success := es.sendViaMailgun("s176492@student.dtu.dk", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to invalid domain, but got success")
	}
}

func TestSendViaMailgunInvalidEmailAddress(t *testing.T) {
	es := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
	}
	success := es.sendViaMailgun("invalid_email", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to invalid email address, but got success")
	}
}

func TestSendViaMailgunNoEnvVars(t *testing.T) {
	os.Unsetenv("MAILGUN_PRIVATE_API_KEY")
	os.Unsetenv("MAILGUN_DOMAIN")

	es := &EmailService{
		MailgunPrivateAPIKey: os.Getenv("MAILGUN_PRIVATE_API_KEY"),
		MailgunDomain:        os.Getenv("MAILGUN_DOMAIN"),
	}
	success := es.sendViaMailgun("s176492@student.dtu.dk", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to missing environment variables, but got success")
	}
}

func TestSendViaSendGrid(t *testing.T) {
	// Mocking SendGrid's API
	httpmock.ActivateNonDefault(sendgrid.DefaultClient.HTTPClient)
	defer httpmock.DeactivateAndReset()

	emailService := &EmailService{
		SendGridAPIKey: "FAKE_API_KEY",
	}

	// Test 1: Successful Email Sending
	httpmock.RegisterResponder("POST", "https://api.sendgrid.com/v3/mail/send",
		httpmock.NewStringResponder(202, ""))

	success := emailService.sendViaSendGrid("s176492@student.dtu.dk", "Test", "This is a test email")
	if !success {
		t.Error("Expected success, but got failure")
	}

	// Test 2: Invalid API Key
	httpmock.RegisterResponder("POST", "https://api.sendgrid.com/v3/mail/send",
		httpmock.NewStringResponder(401, ""))

	success = emailService.sendViaSendGrid("test@example.com", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to invalid API key, but got success")
	}

	// Test 3: Simulate Network Error
	httpmock.RegisterResponder("POST", "https://api.sendgrid.com/v3/mail/send",
		httpmock.NewErrorResponder(fmt.Errorf("simulated network error")))

	success = emailService.sendViaSendGrid("test@example.com", "Test", "This is a test email")
	if success {
		t.Error("Expected failure due to network error, but got success")
	}
}
