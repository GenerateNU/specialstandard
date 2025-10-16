package config

type S3 struct {
	Bucket    string `env:"AWS_S3_BUCKET, required"`
	Region    string `env:"AWS_REGION, required"`
	AccessKey string `env:"AWS_ACCESS_KEY_ID, required"`
	SecretKey string `env:"AWS_SECRET_ACCESS_KEY, required"`
}
