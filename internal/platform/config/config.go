package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

const (
	defaultPort        = "8080"
	defaultServiceName = "dys-sub-fsc-orders"
	defaultEnvironment = "local"
	defaultLogLevel    = "INFO"
)

type Config struct {
	Port             string
	ServiceName      string
	Environment      string
	LogLevel         string
	SubscriptionName string
	OIDCValidation   bool
	OIDCAudience     string
	OIDCAllowedEmail string
	MongoURI         string
	MongoDBName      string
	FSCBaseURL       string
	FSCAPIKey        string
	HTTPTimeoutMS    int
}

func Load() Config {
	loadDotEnv(".env")

	return Config{
		Port:             getOrDefault("PORT", defaultPort),
		ServiceName:      getOrDefault("SERVICE_NAME", defaultServiceName),
		Environment:      getOrDefault("ENVIRONMENT", defaultEnvironment),
		LogLevel:         strings.ToUpper(getOrDefault("LOG_LEVEL", defaultLogLevel)),
		SubscriptionName: strings.TrimSpace(os.Getenv("PUBSUB_SUBSCRIPTION_NAME")),
		OIDCValidation:   getBoolOrDefault("OIDC_VALIDATION_ENABLED", false),
		OIDCAudience:     strings.TrimSpace(os.Getenv("OIDC_AUDIENCE")),
		OIDCAllowedEmail: strings.TrimSpace(os.Getenv("OIDC_ALLOWED_EMAIL")),
		MongoURI:         strings.TrimSpace(os.Getenv("MONGODB_URI")),
		MongoDBName:      getOrDefault("MONGODB_DB_NAME", "falabellaDB"),
		FSCBaseURL:       getOrDefault("DYS_API_FSC_BASE_URL", "https://dys-api-fsc-785293986978.us-central1.run.app"),
		FSCAPIKey:        strings.TrimSpace(os.Getenv("DYS_API_FSC_API_KEY")),
		HTTPTimeoutMS:    getIntOrDefault("HTTP_TIMEOUT_MS_FSC", 5000),
	}
}

func getOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getBoolOrDefault(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func getIntOrDefault(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return parsed
}

func loadDotEnv(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}

		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
