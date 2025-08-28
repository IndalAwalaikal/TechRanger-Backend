// utils/email.go
package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendResetPasswordEmail(toEmail, token string) error {
	SMTPHost := os.Getenv("SMTP_HOST")
	SMTPPort := os.Getenv("SMTP_PORT")
	EmailSender := os.Getenv("SMTP_USER")
	EmailPassword := os.Getenv("SMTP_PASS")
	FrontendURL := os.Getenv("FRONTEND_URL")

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", FrontendURL, token)

	subject := "Reset Password Akun Anda"
	body := fmt.Sprintf(`
		<h1>Permintaan Reset Password</h1>
		<p>Kami menerima permintaan untuk mereset password akun Anda.</p>
		<p>Silakan klik link di bawah untuk mengatur password baru:</p>
		<p><a href="%s" style="color: #007BFF;">Reset Password</a></p>
		<p>Link ini akan kedaluwarsa dalam 1 jam.</p>
		<br>
		<p>Jika Anda tidak meminta reset password, abaikan email ini.</p>
		<br>
		<p>Salam,<br>Tim COCONUT</p>
	`, resetLink)

	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
		"%s",
		EmailSender, toEmail, subject, body)

	auth := smtp.PlainAuth("", EmailSender, EmailPassword, SMTPHost)

	err := smtp.SendMail(SMTPHost+":"+SMTPPort, auth, EmailSender, []string{toEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("gagal kirim email reset password: %v", err)
	}

	return nil
}

func SendVerificationEmail(toEmail, token string) error {
    SMTPHost := os.Getenv("SMTP_HOST")
	SMTPPort := os.Getenv("SMTP_PORT")
	EmailSender := os.Getenv("SMTP_USER")
	EmailPassword := os.Getenv("SMTP_PASS")
	FrontendURL := os.Getenv("FRONTEND_URL")

    verificationLink := fmt.Sprintf("%s/verify?token=%s", FrontendURL, token)

    subject := "Verifikasi Email Anda"
    body := fmt.Sprintf(`
        <h2>Verifikasi Akun</h2>
        <p>Klik link berikut untuk verifikasi akun Anda:</p>
        <a href="%s" style="color: #007BFF;">Verifikasi Akun</a>
        <br><br>
        <p>Jika Anda tidak mendaftar, abaikan email ini.</p>
    `, verificationLink)

    msg := fmt.Sprintf("From: %s\r\n"+
        "To: %s\r\n"+
        "Subject: %s\r\n"+
        "MIME-Version: 1.0\r\n"+
        "Content-Type: text/html; charset=UTF-8\r\n\r\n"+
        "%s",
        EmailSender, toEmail, subject, body)

    auth := smtp.PlainAuth("", EmailSender, EmailPassword, SMTPHost)

    return smtp.SendMail(SMTPHost+":"+SMTPPort, auth, EmailSender, []string{toEmail}, []byte(msg))
}
