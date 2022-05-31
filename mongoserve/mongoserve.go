package mongoserve

import (
	"context"
	"fmt"

	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/jacksonCLyu/ridi-utils/utils/errcheck"
	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Init init module
func Init(opts ...Option) error {
	initOptions := &initOptions{}
	for _, opt := range opts {
		opt.apply(initOptions)
	}
	if initOptions.clientOptions != nil {
		defaultClientOpts = initOptions.clientOptions
	}
	return nil
}

// InitClient init mongo driver client
func InitClient(opts *options.ClientOptions) *mongo.Client {
	defer rescueutil.Recover(func(err any) {
		fmt.Printf("InitPool error: %v", err)
	})
	client := assignutil.Assign(mongo.NewClient(opts))
	// connect
	errcheck.CheckAndPanic(Conn(client))
	return client
}

// Conn connect to mongo server
func Conn(client *mongo.Client) error {
	// connect
	ctx, cancel := context.WithTimeout(context.Background(), *defaultClientOpts.ConnectTimeout)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		return err
	}
	pingCtx, pingCancel := context.WithTimeout(context.Background(), *defaultClientOpts.ServerSelectionTimeout)
	defer pingCancel()
	return client.Ping(pingCtx, nil)
}

func CheckConnOK(client *mongo.Client) bool {
	pingCtx, pingCancel := context.WithTimeout(context.Background(), *defaultClientOpts.ServerSelectionTimeout)
	defer pingCancel()
	return client.Ping(pingCtx, nil) == nil
}

// Disconn disconnect mongo client pool
func Disconn(client *mongo.Client) error {
	disconnCtx, disconnCancel := context.WithCancel(context.Background())
	defer disconnCancel()
	return client.Disconnect(disconnCtx)
}

// Reset reset mongo client connection state
func Reset(client *mongo.Client) error {
	_ = Disconn(client)
	return Conn(client)
}
