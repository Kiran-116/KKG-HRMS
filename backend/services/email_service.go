package services

import (
	"fmt"

	"hrms/config"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
	SendMagicLinkEmail(to, name, magicLink string) error
}

type emailService struct {
	enabled bool
	from    string
	dialer  *gomail.Dialer
}

func NewEmailService() EmailService {
	cfg := config.AppConfig.SMTP
	if !cfg.Enabled {
		return &emailService{enabled: false}
	}

	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)

	return &emailService{
		enabled: true,
		from:    cfg.From,
		dialer:  dialer,
	}
}

func (s *emailService) SendEmail(to, subject, body string) error {
	if !s.enabled {
		// Log that email would be sent in production
		fmt.Printf("Email would be sent to %s: %s\n", to, subject)
		return nil
	}

	message := gomail.NewMessage()
	message.SetHeader("From", s.from)
	message.SetHeader("To", to)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", body)

	return s.dialer.DialAndSend(message)
}

func (s *emailService) SendMagicLinkEmail(to, name, magicLink string) error {
	subject := "Welcome to HRMS - Set Your Password"
	
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="IE=edge">
	<title>Welcome to HRMS</title>
	<style>
		* {
			margin: 0;
			padding: 0;
			box-sizing: border-box;
		}
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
			line-height: 1.6;
			color: #1f2937;
			background-color: #f3f4f6;
			padding: 20px;
		}
		.email-container {
			max-width: 600px;
			margin: 0 auto;
			background-color: #ffffff;
			border-radius: 8px;
			overflow: hidden;
			box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		}
		.email-header {
			background: linear-gradient(135deg, #4F46E5 0%%, #7C3AED 100%%);
			color: #ffffff;
			padding: 40px 30px;
			text-align: center;
		}
		.email-header h1 {
			font-size: 28px;
			font-weight: 700;
			margin-bottom: 10px;
			letter-spacing: -0.5px;
		}
		.email-header p {
			font-size: 16px;
			opacity: 0.95;
			font-weight: 400;
		}
		.email-body {
			padding: 40px 30px;
			background-color: #ffffff;
		}
		.greeting {
			font-size: 18px;
			font-weight: 600;
			color: #1f2937;
			margin-bottom: 20px;
		}
		.content-text {
			font-size: 16px;
			color: #4b5563;
			margin-bottom: 16px;
			line-height: 1.7;
		}
		.button-container {
			text-align: center;
			margin: 35px 0;
		}
		.cta-button {
			display: inline-block;
			padding: 16px 40px;
			background: linear-gradient(135deg, #4F46E5 0%%, #7C3AED 100%%);
			color: #ffffff;
			text-decoration: none;
			border-radius: 6px;
			font-size: 16px;
			font-weight: 600;
			letter-spacing: 0.3px;
			box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
			transition: transform 0.2s, box-shadow 0.2s;
		}
		.cta-button:hover {
			transform: translateY(-2px);
			box-shadow: 0 6px 16px rgba(79, 70, 229, 0.4);
		}
		.info-box {
			background-color: #f0f9ff;
			border-left: 4px solid #3b82f6;
			padding: 16px 20px;
			margin: 30px 0;
			border-radius: 4px;
		}
		.info-box-title {
			font-size: 14px;
			font-weight: 600;
			color: #1e40af;
			margin-bottom: 8px;
		}
		.info-box-text {
			font-size: 14px;
			color: #1e3a8a;
			line-height: 1.6;
		}
		.link-fallback {
			background-color: #f9fafb;
			border: 1px solid #e5e7eb;
			border-radius: 6px;
			padding: 16px;
			margin: 25px 0;
		}
		.link-fallback-title {
			font-size: 13px;
			font-weight: 600;
			color: #6b7280;
			margin-bottom: 8px;
			text-transform: uppercase;
			letter-spacing: 0.5px;
		}
		.link-fallback-url {
			font-size: 13px;
			color: #4F46E5;
			word-break: break-all;
			font-family: 'Courier New', monospace;
		}
		.security-note {
			background-color: #fef3c7;
			border-left: 4px solid #f59e0b;
			padding: 16px 20px;
			margin: 30px 0;
			border-radius: 4px;
		}
		.security-note-title {
			font-size: 14px;
			font-weight: 600;
			color: #92400e;
			margin-bottom: 8px;
		}
		.security-note-text {
			font-size: 14px;
			color: #78350f;
			line-height: 1.6;
		}
		.email-footer {
			background-color: #f9fafb;
			padding: 30px;
			text-align: center;
			border-top: 1px solid #e5e7eb;
		}
		.footer-text {
			font-size: 13px;
			color: #6b7280;
			line-height: 1.6;
			margin-bottom: 8px;
		}
		.footer-text:last-child {
			margin-bottom: 0;
		}
		.divider {
			height: 1px;
			background-color: #e5e7eb;
			margin: 30px 0;
		}
		@media only screen and (max-width: 600px) {
			.email-header {
				padding: 30px 20px;
			}
			.email-header h1 {
				font-size: 24px;
			}
			.email-body {
				padding: 30px 20px;
			}
			.cta-button {
				padding: 14px 32px;
				font-size: 15px;
			}
		}
	</style>
</head>
<body>
	<div class="email-container">
		<div class="email-header">
			<h1>Welcome to HRMS</h1>
			<p>Human Resources Management System</p>
		</div>
		
		<div class="email-body">
			<div class="greeting">Hello %s,</div>
			
			<p class="content-text">
				Your HRMS account has been successfully created! We're excited to have you on board.
			</p>
			
			<p class="content-text">
				To get started, you'll need to set up your password. Click the button below to access your account and create a secure password.
			</p>
			
			<div class="button-container">
				<a href="%s" class="cta-button">Set Your Password</a>
			</div>
			
			<div class="info-box">
				<div class="info-box-title">📋 What to Expect</div>
				<div class="info-box-text">
					After clicking the button above, you'll be redirected to a secure page where you can set your password. 
					Make sure to choose a strong password with at least 10 characters, including numbers and symbols.
				</div>
			</div>
			
			<div class="security-note">
				<div class="security-note-title">⏰ Important: Link Expiration</div>
				<div class="security-note-text">
					This link will expire in <strong>24 hours</strong> for security reasons. 
					If you need a new link, please contact your administrator.
				</div>
			</div>
			
			<div class="link-fallback">
				<div class="link-fallback-title">Button not working?</div>
				<div class="link-fallback-url">%s</div>
			</div>
			
			<div class="divider"></div>
			
			<p class="content-text" style="font-size: 14px; color: #6b7280;">
				If you didn't request this account or have any questions, please contact your administrator immediately.
			</p>
		</div>
		
		<div class="email-footer">
			<p class="footer-text">
				This is an automated message from HRMS. Please do not reply to this email.
			</p>
			<p class="footer-text">
				For support, please contact your system administrator.
			</p>
		</div>
	</div>
</body>
</html>
	`, name, magicLink, magicLink)

	return s.SendEmail(to, subject, body)
}
