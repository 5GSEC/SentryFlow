package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/5gsec/sentryflow/speculator/pkg/config"
	"github.com/5gsec/sentryflow/speculator/pkg/util"
)

type Database interface {
	//GetApiSpec() (libopenapi.Document, error)
	//PutApiSpec(document libopenapi.Document) error
	//DeleteProvidedApiSpec() error
	//DeleteApprovedAPISpec() error
}

type Handler struct {
	Database   *mongo.Database
	Disconnect func() error
}

func New(ctx context.Context, dbConfig *config.Database) (*Handler, error) {
	logger := util.GetLogger()

	client, err := mongo.Connect(ctx, clientOptions(dbConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Infof("connecting to %s database", dbConfig.Name)
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	logger.Info("connected to database")

	return &Handler{
		Database: client.Database(dbConfig.Name),
		Disconnect: func() error {
			return client.Disconnect(ctx)
		},
	}, nil
}

func clientOptions(dbConfig *config.Database) *options.ClientOptions {
	dbLoggerOptions := &options.LoggerOptions{
		ComponentLevels: map[options.LogComponent]options.LogLevel{},
	}
	switch dbConfig.LogLevel {
	case util.LevelInfo:
		dbLoggerOptions.SetComponentLevel(options.LogComponentCommand, options.LogLevelInfo)
	case util.LevelDebug:
		dbLoggerOptions.SetComponentLevel(options.LogComponentCommand, options.LogLevelDebug)
	}

	return options.Client().
		ApplyURI(dbConfig.Uri).
		SetAuth(options.Credential{
			Username: dbConfig.User,
			Password: dbConfig.Password,
		}).
		SetLoggerOptions(dbLoggerOptions)
}
