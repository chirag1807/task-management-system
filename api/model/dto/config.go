package dto

type Config struct {
	Port     uint     `mapstructure:"PORT"`
	Database Database `mapstructure:",squash"`
	Redis    Redis    `mapstructure:",squash"`
	RabbitMQ RabbitMQ `mapstructure:",squash"`
	SMTP     SMTP     `mapstructure:",squash"`
}

type Database struct {
	Username         string `mapstructure:"DATABASE_USERNAME"`
	Password         string `mapstructure:"DATABASE_PASSWORD"`
	Port             string `mapstructure:"DATABASE_PORT"`
	Name             string `mapstructure:"DATABASE_NAME"`
	TestDatabaseName string `mapstructure:"TEST_DATABASE_NAME"`
	SSLMode          string `mapstructure:"DATABASE_SSLMODE"`
}

type Redis struct {
	Port     string `mapstructure:"REDIS_PORT"`
	Password string `mapstructure:"REDIS_PASSWORD"`
	DB       int    `mapstructure:"REDIS_DB"`
}

type RabbitMQ struct {
	Username string `mapstructure:"RABBITMQ_USERNAME"`
	Password string `mapstructure:"RABBITMQ_PASSWORD"`
	Port     string `mapstructure:"RABBITMQ_PORT"`
}

type JWTSecret struct {
	SecretKey string `json:"secretkey"`
}

type SMTP struct {
	EmailFrom     string `mapstructure:"SMTP_EMAIL_FROM"`
	EmailPassword string `mapstructure:"SMTP_EMAIL_PASSWORD"`
	Host          string `mapstructure:"SMTP_HOST"`
	Port          string `mapstructure:"SMTP_PORT"`
}
