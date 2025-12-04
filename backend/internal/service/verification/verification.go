package verification

import (
	"fmt"
	"log/slog"
	"math/rand"
	"specialstandard/internal/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/resend/resend-go/v3"
)

// SendVerificationCode handles sending verification codes via email
func (h *Handler) SendVerificationCode(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	email, err := h.authRepo.GetUserEmail(c.Context(), userID) 
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.SendCodeResponse{
			Success: false,
			Error:   "Failed to get user email",
		})
	}

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.SendCodeResponse{
			Success: false,
			Error:   "User email not found",
		})
	}

	// Generate 6-digit verification code
	code := h.generateVerificationCode()

	now := time.Now()
	expiresAt := now.Add(10 * time.Minute)

	// Store verification code in database
	err = h.verificationRepo.CreateVerificationCode(c.Context(), models.VerificationCode{
		UserID:    userID,
		Code:      code,
		ExpiresAt: expiresAt,
		CreatedAt: now,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.SendCodeResponse{
			Success: false,
			Error:   "Failed to store verification code",
		})
	}

	// Send email using Resend
	emailHTML := h.getEmailTemplate(code)
	
	params := &resend.SendEmailRequest{
		From:    h.fromEmail,
		To:      []string{email},
		Subject: "The Special Standard Verification Code",
		Html:    emailHTML,
	}

	sent, err := h.resendClient.Emails.Send(params)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.SendCodeResponse{
			Success: false,
			Error:   "Failed to send email: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.SendCodeResponse{
		Success:   true,
		MessageID: sent.Id,
	})
}

// Temporarily update your VerifyCode function with debugging:

func (h *Handler) VerifyCode(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	
	var req models.VerifyCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.VerifyCodeResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	// Validate code format
	code := strings.TrimSpace(req.Code)
	if len(code) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(models.VerifyCodeResponse{
			Success: false,
			Error:   "Invalid code format",
		})
	}

	// Verify the code
	valid, err := h.verificationRepo.VerifyCode(c.Context(), userID, code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.VerifyCodeResponse{
			Success: false,
			Error:   "Failed to verify code",
		})
	}

	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(models.VerifyCodeResponse{
			Success:  false,
			Verified: false,
			Error:    "Invalid or expired code",
		})
	}

	err = h.authRepo.MarkEmailVerified(c.Context(), userID)
	if err != nil {
        slog.Warn("Failed to update user metadata", slog.Any("err", err))
    }

	return c.Status(fiber.StatusOK).JSON(models.VerifyCodeResponse{
		Success:  true,
		Verified: true,
		Message:  "Email verified successfully",
	})
}

func (h *Handler) ResendCode(c *fiber.Ctx) error {
	return h.SendVerificationCode(c)
}

func (h *Handler) generateVerificationCode() string {
	return fmt.Sprintf("%06d", rand.Intn(900000)+100000)
}

// getEmailTemplate returns the HTML email template
func (h *Handler) getEmailTemplate(code string) string {
	return fmt.Sprintf(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
	</head>
	<body>
		<div style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
			<div style="background-color: #ffffff; border-radius: 10px; padding: 30px; box-shadow: 0 2px 10px rgba(0,0,0,0.1);">
				<h1 style="color: #333333; font-size: 24px; margin-bottom: 10px;">Verify Your Email</h1>
				<p style="color: #666666; font-size: 16px; line-height: 1.5; margin-bottom: 30px;">
					Please use the verification code below to confirm your email address.
				</p>
				
				<div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 20px; text-align: center; border-radius: 8px; margin-bottom: 30px;">
					<div style="background: white; border-radius: 6px; padding: 15px; display: inline-block;">
						<span style="font-size: 32px; font-weight: bold; letter-spacing: 8px; color: #333333;">%s</span>
					</div>
				</div>
				
				<p style="color: #666666; font-size: 14px; line-height: 1.5;">
					This code will expire in <strong>10 minutes</strong>.
				</p>
				
				<hr style="border: none; border-top: 1px solid #eeeeee; margin: 30px 0;">
				
				<p style="color: #999999; font-size: 13px; line-height: 1.5;">
					If you didn't request this verification code, you can safely ignore this email.
				</p>
			</div>
		</div>
	</body>
	</html>
	`, code)
}