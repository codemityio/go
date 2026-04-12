package mongo

import (
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// NewClient mongodb client factory.
func NewClient(config *Config) (*mongo.Client, error) {
	mongoClient, err := mongo.Connect(
		options.Client().
			ApplyURI(config.DSN).
			SetAuth(options.Credential{
				AuthMechanism:           "",
				AuthMechanismProperties: nil,
				AuthSource:              "",
				Username:                config.Username,
				Password:                config.Password,
				PasswordSet:             false,
				OIDCMachineCallback:     nil,
				OIDCHumanCallback:       nil,
			}),
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrClient, err)
	}

	return mongoClient, nil
}
