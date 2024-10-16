package config

import (
	"os"
	"strconv"
)

type Config struct {
	Version          string
	ServerBindPort   int
	ServerBindHost   string
	ApiRootPath      string
	KeycloakUrl      string
	KeycloakUser     string
	KeycloakPassword string
	ApiDescription   string
	ApiDomain        string
}

func (c *Config) GetServerBindAddress() string {
	return c.ServerBindHost + ":" + strconv.Itoa(c.ServerBindPort)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultValue
}

var AppSettings = &Config{}

func Setup() {
	AppSettings = &Config{
		Version:          getEnv("VERSION", "0.0.1"),
		ServerBindPort:   getEnvInt("SERVER_BIND_PORT", 8080),
		ServerBindHost:   getEnv("SERVER_BIND_HOST", "0.0.0.0"),
		ApiRootPath:      getEnv("API_ROOT_PATH", ""),
		KeycloakUrl:      getEnv("KC_URL", "https://auth.openepi.io"),
		KeycloakUser:     getEnv("KC_USER", "admin"),
		KeycloakPassword: getEnv("KC_PASSWORD", "admin"),
	}
}
