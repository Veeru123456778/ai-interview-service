package app

import (
	"fmt"
	"os"
	"strconv"
)

func getEnvString(key string) string {
	return os.Getenv(key)
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return number
}

func (c *Config) validate() error {
	required := map[string]string{
		"DATABASE_URL":        c.Database.URL,
		"REDIS_URL":           c.Redis.URL,
		"SUPABASE_URL":        c.Supabase.URL,
		"SUPABASE_ANON_KEY":   c.Supabase.AnonKey,
		"SUPABASE_JWT_SECRET": c.Supabase.JWTSecret,
		"GEMINI_API_KEY":      c.Gemini.APIKey,
		"GEMINI_MODEL":        c.Gemini.Model,
	}

	for key, value := range required {
		if value == "" {
			return fmt.Errorf("missing required environment variable: %s", key)
		}
	}

	return nil
}