package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	FCM      FCMConfig      `mapstructure:"fcm"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Retry    RetryConfig    `mapstructure:"retry"`
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Env  string `mapstructure:"env"`
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// KafkaConfig holds Kafka connection and consumer configuration
type KafkaConfig struct {
	Brokers            []string      `mapstructure:"brokers"`
	GroupID            string        `mapstructure:"group_id"`
	Topics             TopicsConfig  `mapstructure:"topics"`
	AutoOffsetReset    string        `mapstructure:"auto_offset_reset"`
	EnableAutoCommit   bool          `mapstructure:"enable_auto_commit"`
}

// TopicsConfig holds Kafka topic names
type TopicsConfig struct {
	TaskCreated   string `mapstructure:"task_created"`
	NewMessage    string `mapstructure:"new_message"`
	IncomingCall  string `mapstructure:"incoming_call"`
}

// FCMConfig holds Firebase Cloud Messaging configuration
type FCMConfig struct {
	CredentialsPath string `mapstructure:"credentials_path"`
	ProjectID       string `mapstructure:"project_id"`
}

// LoggerConfig holds logging configuration
type LoggerConfig struct {
	Level string `mapstructure:"level"`
}

// RetryConfig holds retry mechanism configuration
type RetryConfig struct {
	MaxAttempts  int `mapstructure:"max_attempts"`
	DelaySeconds int `mapstructure:"delay_seconds"`
}

// Load reads configuration from .env file and environment variables
func Load() (*Config, error) {
	// Try to load .env file (optional - will use system env vars if not found)
	if err := godotenv.Load(); err != nil {
		// .env file not found, will use system environment variables
		// This is normal on platforms like Render where env vars are set directly
		if !os.IsNotExist(err) {
			fmt.Printf("Warning: Error loading .env file: %v\n", err)
		}
	}

	// Set default values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.env", "development")
	viper.SetDefault("logger.level", "info")
	viper.SetDefault("retry.max_attempts", 3)
	viper.SetDefault("retry.delay_seconds", 5)

	// Enable environment variable reading
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind environment variables to config fields
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.env", "APP_ENV")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("database.sslmode", "DB_SSLMODE")
	viper.BindEnv("kafka.brokers", "KAFKA_BROKERS")
	viper.BindEnv("kafka.group_id", "KAFKA_GROUP_ID")
	viper.BindEnv("kafka.topics.task_created", "KAFKA_TOPIC_TASK_CREATED")
	viper.BindEnv("kafka.topics.new_message", "KAFKA_TOPIC_NEW_MESSAGE")
	viper.BindEnv("kafka.topics.incoming_call", "KAFKA_TOPIC_INCOMING_CALL")
	viper.BindEnv("kafka.auto_offset_reset", "KAFKA_AUTO_OFFSET_RESET")
	viper.BindEnv("kafka.enable_auto_commit", "KAFKA_ENABLE_AUTO_COMMIT")
	viper.BindEnv("fcm.credentials_path", "FCM_CREDENTIALS_PATH")
	viper.BindEnv("fcm.project_id", "FCM_PROJECT_ID")
	viper.BindEnv("logger.level", "LOG_LEVEL")
	viper.BindEnv("retry.max_attempts", "MAX_RETRY_ATTEMPTS")
	viper.BindEnv("retry.delay_seconds", "RETRY_DELAY_SECONDS")

	// Create config struct and populate from environment
	var config Config
	
	// Manually set values from environment variables
	config.Server.Port = viper.GetInt("server.port")
	config.Server.Env = viper.GetString("server.env")
	
	config.Database.Host = viper.GetString("database.host")
	config.Database.Port = viper.GetInt("database.port")
	config.Database.User = viper.GetString("database.user")
	config.Database.Password = viper.GetString("database.password")
	config.Database.DBName = viper.GetString("database.dbname")
	config.Database.SSLMode = viper.GetString("database.sslmode")
	
	// Handle KAFKA_BROKERS which can be comma-separated
	if brokersStr := viper.GetString("kafka.brokers"); brokersStr != "" {
		config.Kafka.Brokers = strings.Split(brokersStr, ",")
		// Trim whitespace from each broker
		for i, broker := range config.Kafka.Brokers {
			config.Kafka.Brokers[i] = strings.TrimSpace(broker)
		}
	}
	config.Kafka.GroupID = viper.GetString("kafka.group_id")
	config.Kafka.Topics.TaskCreated = viper.GetString("kafka.topics.task_created")
	config.Kafka.Topics.NewMessage = viper.GetString("kafka.topics.new_message")
	config.Kafka.Topics.IncomingCall = viper.GetString("kafka.topics.incoming_call")
	config.Kafka.AutoOffsetReset = viper.GetString("kafka.auto_offset_reset")
	config.Kafka.EnableAutoCommit = viper.GetBool("kafka.enable_auto_commit")
	
	config.FCM.CredentialsPath = viper.GetString("fcm.credentials_path")
	config.FCM.ProjectID = viper.GetString("fcm.project_id")
	
	config.Logger.Level = viper.GetString("logger.level")
	
	config.Retry.MaxAttempts = viper.GetInt("retry.max_attempts")
	config.Retry.DelaySeconds = viper.GetInt("retry.delay_seconds")

	return &config, nil
}

// GetDSN returns the PostgreSQL connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

