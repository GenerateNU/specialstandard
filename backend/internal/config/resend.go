package config

type Resend struct {
	APIKey    string `env:"RESEND_API_KEY,required"`
	FromEmail string `env:"RESEND_FROM_EMAIL,default=Kevin Matula <kevinmatula@plantkeepr.co>"`
}