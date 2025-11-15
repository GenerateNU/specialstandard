package verification

import (
	"fmt"
	"math/rand"
	"specialstandard/internal/models"
	"specialstandard/internal/storage"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/resend/resend-go/v3"
)

type Handler struct {
	verificationRepo storage.VerificationRepository
	db               *pgxpool.Pool
	resendClient     *resend.Client
	fromEmail        string
}

type SendCodeResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Error     string `json:"error,omitempty"`
}

type VerifyCodeRequest struct {
	Code string `json:"code"`
}

type VerifyCodeResponse struct {
	Success  bool   `json:"success"`
	Verified bool   `json:"verified,omitempty"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Createing a new verification handler
func NewHandler(verificationRepo storage.VerificationRepository, db *pgxpool.Pool, resendApiKey, fromEmail string) *Handler {
	resendClient := resend.NewClient(resendApiKey)
	
	return &Handler{
		verificationRepo: verificationRepo,
		db:               db,
		resendClient:     resendClient,
		fromEmail:        fromEmail,
	}
}

// extractUserIDFromToken extracts the user ID from a JWT token
func (h *Handler) extractUserIDFromToken(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("user ID not found in token")
	}

	return userID, nil
}

// SendVerificationCode handles sending verification codes via email
func (h *Handler) SendVerificationCode(c *fiber.Ctx) error {
	// Extract JWT from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(SendCodeResponse{
			Success: false,
			Error:   "No authentication token provided",
		})
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)
	
	// Parse the JWT to get the user ID
	userID, err := h.extractUserIDFromToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(SendCodeResponse{
			Success: false,
			Error:   "Invalid authentication token",
		})
	}

	// Get user email from auth.users table
	var email string
	query := `SELECT email FROM auth.users WHERE id = $1`
	err = h.db.QueryRow(c.Context(), query, userID).Scan(&email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(SendCodeResponse{
			Success: false,
			Error:   "Failed to get user email",
		})
	}

	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(SendCodeResponse{
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
		return c.Status(fiber.StatusInternalServerError).JSON(SendCodeResponse{
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
		return c.Status(fiber.StatusInternalServerError).JSON(SendCodeResponse{
			Success: false,
			Error:   "Failed to send email: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(SendCodeResponse{
		Success:   true,
		MessageID: sent.Id,
	})
}

// Temporarily update your VerifyCode function with debugging:

func (h *Handler) VerifyCode(c *fiber.Ctx) error {
	// Extract JWT from Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(VerifyCodeResponse{
			Success: false,
			Error:   "No authentication token provided",
		})
	}

	token := strings.Replace(authHeader, "Bearer ", "", 1)
	userID, err := h.extractUserIDFromToken(token)
	if err != nil {
		fmt.Printf("Error extracting user ID from token: %v\n", err)
		return c.Status(fiber.StatusUnauthorized).JSON(VerifyCodeResponse{
			Success: false,
			Error:   "Invalid authentication token",
		})
	}
	
	fmt.Printf("Extracted user ID: %s\n", userID)

	var req VerifyCodeRequest
	if err := c.BodyParser(&req); err != nil {
		fmt.Printf("Error parsing body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(VerifyCodeResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}
	
	fmt.Printf("Received code: %s\n", req.Code)

	// Validate code format
	code := strings.TrimSpace(req.Code)
	if len(code) != 6 {
		return c.Status(fiber.StatusBadRequest).JSON(VerifyCodeResponse{
			Success: false,
			Error:   "Invalid code format",
		})
	}

	// Verify the code
	valid, err := h.verificationRepo.VerifyCode(c.Context(), userID, code)
	if err != nil {
		fmt.Printf("Error verifying code: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(VerifyCodeResponse{
			Success: false,
			Error:   "Failed to verify code",
		})
	}
	
	fmt.Printf("Code valid: %v\n", valid)

	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(VerifyCodeResponse{
			Success:  false,
			Verified: false,
			Error:    "Invalid or expired code",
		})
	}

	// Update user metadata to mark as verified
	updateQuery := `
		UPDATE auth.users 
		SET raw_user_meta_data = jsonb_set(
			COALESCE(raw_user_meta_data, '{}'::jsonb),
			'{email_verified}',
			'true'
		),
		updated_at = NOW()
		WHERE id = $1
	`
	
	_, err = h.db.Exec(c.Context(), updateQuery, userID)
	if err != nil {
		fmt.Printf("Warning: Failed to update user metadata: %v\n", err)
	}

	return c.Status(fiber.StatusOK).JSON(VerifyCodeResponse{
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