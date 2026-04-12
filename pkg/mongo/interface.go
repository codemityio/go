package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// Clienter interface.
type Clienter interface {
	Connect(opts ...*options.ClientOptions) (*mongo.Client, error)
	Disconnect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	StartSession(opts ...options.Lister[options.SessionOptions]) (*mongo.Session, error)
	Database(name string, opts ...options.Lister[options.DatabaseOptions]) *mongo.Database
	ListDatabases(
		ctx context.Context,
		filter any,
		opts ...options.Lister[options.ListDatabasesOptions],
	) (mongo.ListDatabasesResult, error)
}
