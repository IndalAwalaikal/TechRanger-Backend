// utils/email.go
package utils

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

// helper: create SMTP auth and send email
func sendEmail(toEmail, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	emailSender := os.Getenv("EMAIL_FROM") // Misal: no-reply@coconut.ac.id
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPort == "" || emailSender == "" || smtpUser == "" || smtpPass == "" {
		return fmt.Errorf("konfigurasi SMTP tidak lengkap di environment")
	}

	message := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n\r\n"+
		"%s",
		emailSender, toEmail, subject, body)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	// Buat koneksi dengan timeout
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("gagal terhubung ke SMTP server: %v", err)
	}
	defer client.Quit()

	if err := client.StartTLS(nil); err != nil {
		return fmt.Errorf("gagal memulai TLS: %v", err)
	}

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("gagal autentikasi SMTP: %v", err)
	}

	if err := client.Mail(emailSender); err != nil {
		return fmt.Errorf("gagal set pengirim: %v", err)
	}

	// Pastikan email penerima valid
	if !strings.Contains(toEmail, "@") {
		return fmt.Errorf("email penerima tidak valid: %s", toEmail)
	}
	if err := client.Rcpt(toEmail); err != nil {
		return fmt.Errorf("gagal set penerima: %v", err)
	}

	// Kirim body
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("gagal buka data writer: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("gagal tulis email: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("gagal tutup data writer: %v", err)
	}

	return nil
}

// Fungsi: Kirim Email Verifikasi
func SendVerificationEmail(toEmail, token string) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		return fmt.Errorf("FRONTEND_URL tidak diatur")
	}

	verificationLink := fmt.Sprintf("%s/verify?token=%s", frontendURL, token)

	subject := "Verifikasi Email Anda - COCONUT"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<style>
				body { font-family: Arial, sans-serif; color: #333; }
				.container { max-width: 600px; margin: auto; padding: 20px; }
				.header { background: #007BFF; color: white; padding: 15px; text-align: center; border-radius: 8px 8px 0 0; }
				.content { padding: 20px; background: #f9f9f9; border: 1px solid #ddd; border-top: none; }
				.button { display: inline-block; padding: 12px 24px; margin: 15px 0; background: #007BFF; color: white; text-decoration: none; border-radius: 5px; }
				.footer { text-align: center; margin-top: 20px; color: #777; font-size: 0.9em; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>COCONUT</h1>
					<p>Computer Club Oriented Network, Utility & Technology</p>
				</div>
				<div class="content">
					<h2>Verifikasi Akun Anda</h2>
					<p>Klik tombol di bawah untuk memverifikasi alamat email Anda:</p>
					<a href="%s" class="button">Verifikasi Email</a>
					<p>Link ini akan kedaluwarsa dalam <strong>1 jam</strong>.</p>
					<p>Jika Anda tidak mendaftar di COCONUT, abaikan email ini.</p>
				</div>
				<div class="footer">
					<p>© 2025 COCONUT. Semua hak dilindungi undang-undang.</p>
				</div>
			</div>
		</body>
		</html>
	`, verificationLink)

	return sendEmail(toEmail, subject, body)
}

// Fungsi: Kirim Email Reset Password
func SendResetPasswordEmail(toEmail, token string) error {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		return fmt.Errorf("FRONTEND_URL tidak diatur")
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	subject := "Reset Password Akun Anda - COCONUT"
	body := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<style>
				body { font-family: Arial, sans-serif; color: #333; }
				.container { max-width: 600px; margin: auto; padding: 20px; }
				.header { background: #28a745; color: white; padding: 15px; text-align: center; border-radius: 8px 8px 0 0; }
				.content { padding: 20px; background: #f9f9f9; border: 1px solid #ddd; border-top: none; }
				.button { display: inline-block; padding: 12px 24px; margin: 15px 0; background: #28a745; color: white; text-decoration: none; border-radius: 5px; }
				.footer { text-align: center; margin-top: 20px; color: #777; font-size: 0.9em; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="header">
					<h1>COCONUT</h1>
					<p>Computer Club Oriented Network, Utility & Technology</p>
				</div>
				<div class="content">
					<h2>Reset Password Anda</h2>
					<p>Kami menerima permintaan reset password untuk akun Anda.</p>
					<p>Klik tombol di bawah untuk mengatur ulang password:</p>
					<a href="%s" class="button">Reset Password</a>
					<p>Link ini akan kedaluwarsa dalam <strong>1 jam</strong>.</p>
					<p>Jika Anda tidak meminta reset, abaikan email ini.</p>
				</div>
				<div class="footer">
					<p>© 2025 COCONUT. Semua hak dilindungi undang-undang.</p>
				</div>
			</div>
		</body>
		</html>
	`, resetLink)

	return sendEmail(toEmail, subject, body)
}