package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Redis     RedisConfig
	JWT       JWTConfig
	OAuth     OAuthConfig
	External  ExternalConfig
	AWS       AWSConfig
	RateLimit RateLimitConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port   string
	Env    string
	AppURL string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL      string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL  string
	Host string
	Port string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret        string
	Expiry        time.Duration
	RefreshExpiry time.Duration
}

// OAuthConfig holds OAuth provider configuration
type OAuthConfig struct {
	Google GoogleOAuthConfig
}

// GoogleOAuthConfig holds Google OAuth configuration
type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// ExternalConfig holds external service configuration
type ExternalConfig struct {
	Zhipu  ZhipuConfig
	Stripe StripeConfig
	Resend ResendConfig
	OpenAI OpenAIConfig
}

// ZhipuConfig holds ZhipuAI API configuration
type ZhipuConfig struct {
	APIKey string
	Model  string // cogvideox-3, cogvideox-flash, etc.
}

// StripeConfig holds Stripe configuration
type StripeConfig struct {
	SecretKey     string
	WebhookSecret string
}

// ResendConfig holds Resend email configuration
type ResendConfig struct {
	APIKey string
}

// OpenAIConfig holds OpenAI configuration
type OpenAIConfig struct {
	APIKey string
}

// AWSConfig holds AWS configuration
type AWSConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	S3Bucket        string
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Requests int
	Window   time.Duration
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if exists
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port:   getEnv("PORT", "8080"),
			Env:    getEnv("ENV", "development"),
			AppURL: getEnv("APP_URL", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			URL:      getEnv("DATABASE_URL", ""),
			Host:     getEnv("DATABASE_HOST", "localhost"),
			Port:     getEnv("DATABASE_PORT", "5432"),
			User:     getEnv("DATABASE_USER", "postgres"),
			Password: getEnv("DATABASE_PASSWORD", "password"),
			Name:     getEnv("DATABASE_NAME", "genvid"),
		},
		Redis: RedisConfig{
			URL:  getEnv("REDIS_URL", "redis://localhost:6379"),
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: getEnv("REDIS_PORT", "6379"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			Expiry:        getDurationEnv("JWT_EXPIRY", 24*time.Hour),
			RefreshExpiry: getDurationEnv("JWT_REFRESH_EXPIRY", 168*time.Hour),
		},
		OAuth: OAuthConfig{
			Google: GoogleOAuthConfig{
				ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
				ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
				RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
			},
		},
		External: ExternalConfig{
			Zhipu: ZhipuConfig{
				APIKey: getEnv("ZHIPU_API_KEY", ""),
				Model:  getEnv("ZHIPU_MODEL", "cogvideox-3"),
			},
			Stripe: StripeConfig{
				SecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
				WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
			},
			Resend: ResendConfig{
				APIKey: getEnv("RESEND_API_KEY", ""),
			},
			OpenAI: OpenAIConfig{
				APIKey: getEnv("OPENAI_API_KEY", ""),
			},
		},
		AWS: AWSConfig{
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			Region:          getEnv("AWS_REGION", "us-east-1"),
			S3Bucket:        getEnv("S3_BUCKET", "genvid-videos"),
		},
		RateLimit: RateLimitConfig{
			Requests: getIntEnv("RATE_LIMIT_REQUESTS", 100),
			Window:   getDurationEnv("RATE_LIMIT_WINDOW", time.Hour),
		},
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	if c.Database.URL != "" {
		return c.Database.URL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
	)
}

// IsDevelopment returns true if in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
