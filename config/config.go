package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	// Database
	DBHost     string
	DBUsername string
	DBPassword string
	DBName     string
	DBPort     string
	SSLMode    string

	// App
	AppPort string

	// Middlewares
	SecretJWT string

	// MAIL
	SmtpHost      string `mapstructure:"SMTP_HOST"`
	SmtpPort      int    `mapstructure:"SMTP_PORT"`
	SenderEmail   string `mapstructure:"SENDER_EMAIL"`
	EmailPassword string `mapstructure:"EMAIL_PASSWORD"`

	// Oy
	BaseUrl  string `mapstructure:"BASEURL"`
	Username string `mapstructure:"USERNAME"`
	ApiKey   string `mapstructure:"API_KEY"`
}

var (
	AppConfig Config
)

func LoadConfig() *Config {
	// Baca file .env
	file, err := os.Open(".env")
	if err != nil {
		panic(fmt.Errorf("failed to open .env file: %w", err))
	}
	defer file.Close()

	// Parse file .env
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue // Abaikan baris kosong atau komentar
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Abaikan baris yang tidak valid
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Set nilai ke Config struct
		switch key {
		case "DB_HOST":
			AppConfig.DBHost = value
		case "DB_USER":
			AppConfig.DBUsername = value
		case "DB_PASSWORD":
			AppConfig.DBPassword = value
		case "DB_DBNAME":
			AppConfig.DBName = value
		case "DB_PORT":
			AppConfig.DBPort = value
		case "SSL_MODE":
			AppConfig.SSLMode = value
		case "APP_PORT":
			AppConfig.AppPort = value
		case "SECRET_JWT":
			AppConfig.SecretJWT = value
		case "BaseURL": // Tambahkan Midtrans Server Key
			AppConfig.BaseUrl = value
		case "Username": // Tambahkan Midtrans Client Key
			AppConfig.Username = value
		case "Apikey": // Tambahkan Midtrans Environment
			AppConfig.ApiKey = value
		}
	}

	// Periksa error saat membaca file
	if err := scanner.Err(); err != nil {
		panic(fmt.Errorf("failed to read .env file: %w", err))
	}

	return &AppConfig
}
