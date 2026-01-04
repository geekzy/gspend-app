package service

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"strings"

	"github.com/geekzy/gspend-app/apps/auth-service/internal/config"
)

// EmailService handles sending emails via SMTP
type EmailService struct {
	config *config.Config
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// IsEnabled returns true if email sending is enabled
func (s *EmailService) IsEnabled() bool {
	return s.config.SMTP.Enabled
}

// GenerateToken generates a secure random token for verification/reset
func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SendVerificationEmail sends an email verification link
func (s *EmailService) SendVerificationEmail(to, fullName, token string) error {
	if !s.IsEnabled() {
		log.Printf("[EMAIL] SMTP disabled - would send verification email to %s", to)
		return nil
	}

	subject := "Verify your gSpend account"
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.config.AppURL, token)

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verify Your Email</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f7f7f7;">
    <div style="background-color: #667eea; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; border-radius: 16px 16px 0 0; text-align: center;">
        <h1 style="color: #ffffff; margin: 0; font-size: 28px;">gSpend</h1>
        <p style="color: #eeeeee; margin: 10px 0 0 0;">Family Financial Management</p>
    </div>
    <div style="background: white; padding: 30px; border-radius: 0 0 16px 16px; box-shadow: 0 4px 12px rgba(0,0,0,0.1);">
        <h2 style="color: #333; margin-top: 0;">Welcome, %s! 👋</h2>
        <p style="color: #555; line-height: 1.6;">Thank you for signing up for gSpend. Please verify your email address to activate your account.</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #667eea; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #ffffff; padding: 14px 32px; border-radius: 8px; text-decoration: none; font-weight: bold; display: inline-block; border: 1px solid #667eea;">Verify Email Address</a>
        </div>
        <p style="color: #888; font-size: 14px;">This link will expire in 24 hours.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="color: #888; font-size: 12px;">If you didn't create an account, you can safely ignore this email.</p>
    </div>
</body>
</html>
`, fullName, verifyURL)

	return s.sendEmail(to, subject, htmlBody)
}

// SendPasswordResetEmail sends a password reset link
func (s *EmailService) SendPasswordResetEmail(to, fullName, token string) error {
	if !s.IsEnabled() {
		log.Printf("[EMAIL] SMTP disabled - would send password reset email to %s", to)
		return nil
	}

	subject := "Reset your gSpend password"
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.config.AppURL, token)

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Reset Your Password</title>
</head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f7f7f7;">
    <div style="background-color: #667eea; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; border-radius: 16px 16px 0 0; text-align: center;">
        <h1 style="color: #ffffff; margin: 0; font-size: 28px;">gSpend</h1>
        <p style="color: #eeeeee; margin: 10px 0 0 0;">Family Financial Management</p>
    </div>
    <div style="background: white; padding: 30px; border-radius: 0 0 16px 16px; box-shadow: 0 4px 12px rgba(0,0,0,0.1);">
        <h2 style="color: #333; margin-top: 0;">Password Reset Request</h2>
        <p style="color: #555; line-height: 1.6;">Hi %s, we received a request to reset your password. Click the button below to create a new password:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="%s" style="background-color: #667eea; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #ffffff; padding: 14px 32px; border-radius: 8px; text-decoration: none; font-weight: bold; display: inline-block; border: 1px solid #667eea;">Reset Password</a>
        </div>
        <p style="color: #888; font-size: 14px;">This link will expire in 1 hour.</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="color: #888; font-size: 12px;">If you didn't request a password reset, you can safely ignore this email. Your password will not be changed.</p>
    </div>
</body>
</html>
`, fullName, resetURL)

	return s.sendEmail(to, subject, htmlBody)
}

// sendEmail sends an email using Gmail SMTP with STARTTLS and LOGIN auth
func (s *EmailService) sendEmail(to, subject, htmlBody string) error {
	from := s.config.SMTP.From
	if from == "" {
		from = s.config.SMTP.User
	}

	// Clean whitespace from config just in case
	smtpUser := strings.TrimSpace(s.config.SMTP.User)
	smtpPass := strings.TrimSpace(s.config.SMTP.Password)
	smtpHost := strings.TrimSpace(s.config.SMTP.Host)

	// Build email headers and body
	var msg bytes.Buffer
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	addr := fmt.Sprintf("%s:%d", smtpHost, s.config.SMTP.Port)

	// Connect to SMTP server
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("[EMAIL] Failed to connect to SMTP server: %v", err)
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Send STARTTLS command
	tlsConfig := &tls.Config{
		ServerName: smtpHost,
	}
	if err = client.StartTLS(tlsConfig); err != nil {
		log.Printf("[EMAIL] Failed to start TLS: %v", err)
		return fmt.Errorf("failed to start TLS: %w", err)
	}

	// Authenticate
	log.Printf("[EMAIL] Auth Debug - User: %q (len=%d), Host: %q", smtpUser, len(smtpUser), smtpHost)

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	if err = client.Auth(auth); err != nil {
		log.Printf("[EMAIL] Failed to authenticate: %v", err)
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err = client.Mail(smtpUser); err != nil {
		log.Printf("[EMAIL] Failed to set sender: %v", err)
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err = client.Rcpt(to); err != nil {
		log.Printf("[EMAIL] Failed to set recipient: %v", err)
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send email body
	w, err := client.Data()
	if err != nil {
		log.Printf("[EMAIL] Failed to get data writer: %v", err)
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = w.Write(msg.Bytes())
	if err != nil {
		log.Printf("[EMAIL] Failed to write email body: %v", err)
		return fmt.Errorf("failed to write email body: %w", err)
	}

	err = w.Close()
	if err != nil {
		log.Printf("[EMAIL] Failed to close data writer: %v", err)
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	client.Quit()

	log.Printf("[EMAIL] Successfully sent email to %s", to)
	return nil
}

// RenderTemplate renders an HTML template with data
func RenderTemplate(tmpl string, data interface{}) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// loginAuth implements the LOGIN authentication mechanism
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		command := string(fromServer)
		switch {
		case strings.Contains(command, "Username"):
			return []byte(a.username), nil
		case strings.Contains(command, "Password"):
			return []byte(a.password), nil
		default:
			// Fallback: sometimes server just sends base64.
			// We can try to guess or just fail.
			// For specific Gmail behaviour, matching "Username" usually works on decoded string.
			return nil, fmt.Errorf("unknown fromServer: %s", command)
		}
	}
	return nil, nil
}
