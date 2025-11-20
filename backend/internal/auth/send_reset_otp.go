package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"specialstandard/internal/config"
	"specialstandard/internal/errs"
	"sync"
	"time"
)

type OTPStore struct {
	sync.Mutex
	otps map[string]OTPData
}

type OTPData struct {
	OTP       string
	Email     string
	ExpiresAt time.Time
}

var store = &OTPStore{
	otps: make(map[string]OTPData),
}

func generateOTP() string {
	const digits = "0123456789"
	b := make([]byte, 6)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(10))
		b[i] = digits[num.Int64()]
	}
	return string(b)
}

func SendResetOTP(cfg *config.Supabase, email string) error {
	otp := generateOTP()

	store.Lock()
	store.otps[email] = OTPData{
		OTP:       otp,
		Email:     email,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	store.Unlock()

	// Send OTP email
	subject := "Password Reset OTP"
	body := fmt.Sprintf(`
		<h2>Password Reset</h2>
		<p>Your OTP is: <strong>%s</strong></p>
		<p>This OTP expires in 15 minutes.</p>
	`, otp)

	// Use your email service to send
	return sendEmail(email, subject, body)
}

func VerifyOTPAndResetPassword(cfg *config.Supabase, email, otp, newPassword string) error {
	store.Lock()
	otpData, exists := store.otps[email]
	store.Unlock()

	if !exists {
		return errs.BadRequest("No OTP found for this email")
	}

	if time.Now().After(otpData.ExpiresAt) {
		store.Lock()
		delete(store.otps, email)
		store.Unlock()
		return errs.BadRequest("OTP has expired")
	}

	if otpData.OTP != otp {
		return errs.BadRequest("Invalid OTP")
	}

	// OTP is valid, update password using admin API
	err := updatePasswordViaAdmin(cfg, email, newPassword)
	if err != nil {
		return err
	}

	// Clear OTP
	store.Lock()
	delete(store.otps, email)
	store.Unlock()

	return nil
}

func updatePasswordViaAdmin(cfg *config.Supabase, email, newPassword string) error {
	// Use Supabase admin API to update password
	// Similar to the earlier admin API approach
	// You'll need to implement this based on your Supabase setup
	return nil
}
