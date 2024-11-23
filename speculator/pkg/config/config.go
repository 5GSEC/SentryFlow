package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/5gsec/sentryflow/speculator/pkg/util"
)

type Exchange struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Durable    bool   `json:"durable,omitempty"`
	AutoDelete bool   `json:"autoDelete,omitempty"`
}

type RabbitMQ struct {
	Host      string    `json:"host"`
	Port      string    `json:"port"`
	User      string    `json:"user,omitempty"`
	Password  string    `json:"password,omitempty"`
	Exchange  *Exchange `json:"exchange"`
	QueueName string    `json:"queueName"`
}

type Database struct {
	LogLevel string `json:"logLevel,omitempty"`
	Uri      string `json:"uri,omitempty"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name,omitempty"`
}

type Configuration struct {
	RabbitMQ *RabbitMQ `json:"rabbitmq"`
	Database *Database `json:"database"`
}

func (c *Configuration) validate() error {
	if c.RabbitMQ == nil {
		return fmt.Errorf("configuration does not contain a valid RabbitMQ configuration")
	}
	if c.RabbitMQ.Host == "" {
		return fmt.Errorf("configuration does not contain a valid RabbitMQ host")
	}
	if c.RabbitMQ.Port == "" || len(c.RabbitMQ.Port) > 5 {
		return fmt.Errorf("configuration does not contain a valid RabbitMQ port")
	}
	if c.RabbitMQ.Exchange == nil {
		return fmt.Errorf("configuration does not contain a valid RabbitMQ exchange")
	}
	if c.RabbitMQ.Exchange.Name == "" {
		return fmt.Errorf("configuration does not contain a valid RabbitMQ exchange name")
	}
	if c.RabbitMQ.QueueName == "" {
		return fmt.Errorf("configuration does not contain a valid RabbitMQ queue name")
	}

	if c.Database == nil {
		return fmt.Errorf("configuration does not contain a valid database configuration")
	}
	if c.Database.Uri == "" {
		return fmt.Errorf("configuration does not contain a valid database URI")
	}
	if c.Database.User == "" {
		return fmt.Errorf("configuration does not contain a valid database user")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("configuration does not contain a valid database password")
	}

	return nil
}

const DefaultConfigFilePath = "config/default.yaml"

func New(configFilePath string, logger *zap.SugaredLogger) (*Configuration, error) {
	if configFilePath == "" {
		configFilePath = DefaultConfigFilePath
		logger.Warnf("using default configfile path: %s", configFilePath)
	}

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Configuration{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	if config.Database.LogLevel == "" {
		config.Database.LogLevel = util.LevelInfo
		logger.Warnf("using default `INFO` database log level: %s", config.Database.LogLevel)
	}

	dbUser := config.Database.User
	dbPassword := config.Database.Password

	config.Database.User = ""
	config.Database.Password = ""

	bytes, err := json.Marshal(config)
	if err != nil {
		logger.Errorf("failed to marshal config file: %v", err)
	}
	logger.Debugf("configuration: %s", string(bytes))

	config.Database.User = dbUser
	config.Database.Password = dbPassword

	return config, nil
}
