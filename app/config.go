package app

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig
	Server    ServerConfig
	Resume    ResumeConfig
	Redis     RedisConfig
	WebSocket WebSocketConfig
	Gemini    GeminiConfig
	Supabase  SupabaseConfig
	Database  DatabaseConfig
}

type AppConfig struct {
	Name        string
	Environment string
}

type ServerConfig struct {
	Port                int
	ReadTimeoutSeconds  int
	WriteTimeoutSeconds int
	IdleTimeoutSeconds  int
}

type ResumeConfig struct {
	MaxUploadSizeMB int
}

type RedisConfig struct {
	URL                  string
	SessionTTLHours      int
	WebSocketTTLMinutes  int
}

type WebSocketConfig struct {
	HeartbeatIntervalSeconds int
	PongTimeoutSeconds       int
}

type GeminiConfig struct {
	APIKey         string
	Model          string
	TimeoutSeconds int
	MaxRetries     int
}

type SupabaseConfig struct {
	URL       string
	AnonKey   string
	JWTSecret string
}

type DatabaseConfig struct {
	URL string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read application config: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Name:        viper.GetString("app.name"),
			Environment: viper.GetString("app.environment"),
		},
		Server: ServerConfig{
			Port:                getEnvInt("PORT", viper.GetInt("server.port")),
			ReadTimeoutSeconds:  viper.GetInt("server.read_timeout_seconds"),
			WriteTimeoutSeconds: viper.GetInt("server.write_timeout_seconds"),
			IdleTimeoutSeconds:  viper.GetInt("server.idle_timeout_seconds"),
		},
		Resume: ResumeConfig{
			MaxUploadSizeMB: viper.GetInt("resume.max_upload_size_mb"),
		},
		Redis: RedisConfig{
			URL:                 getEnvString("REDIS_URL"),
			SessionTTLHours:     viper.GetInt("redis.session_ttl_hours"),
			WebSocketTTLMinutes: viper.GetInt("redis.websocket_ttl_minutes"),
		},
		WebSocket: WebSocketConfig{
			HeartbeatIntervalSeconds: viper.GetInt("websocket.heartbeat_interval_seconds"),
			PongTimeoutSeconds:       viper.GetInt("websocket.pong_timeout_seconds"),
		},
		Gemini: GeminiConfig{
			APIKey:         getEnvString("GEMINI_API_KEY"),
			Model:          getEnvString("GEMINI_MODEL"),
			TimeoutSeconds: viper.GetInt("gemini.timeout_seconds"),
			MaxRetries:     viper.GetInt("gemini.max_retries"),
		},
		Supabase: SupabaseConfig{
			URL:       getEnvString("SUPABASE_URL"),
			AnonKey:   getEnvString("SUPABASE_ANON_KEY"),
			JWTSecret: getEnvString("SUPABASE_JWT_SECRET"),
		},
		Database: DatabaseConfig{
			URL: getEnvString("DATABASE_URL"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}