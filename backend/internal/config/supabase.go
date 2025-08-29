package config

type Supabase struct {
	URL            string `env:"SUPABASE_URL, required"`
	AnonKey        string `env:"SUPABASE_ANON_KEY, required"`
	ServiceRoleKey string `env:"SUPABASE_SERVICE_ROLE_KEY, required"`
}
