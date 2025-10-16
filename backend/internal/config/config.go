package config

type Config struct {
	Application Application
	DB          DB
	Supabase    Supabase
	S3Bucket    S3
	TestMode    bool
}
