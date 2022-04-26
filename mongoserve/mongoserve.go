package mongoserve

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jacksonCLyu/ridi-config/pkg/config"
	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/jacksonCLyu/ridi-utils/utils/errcheck"
	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var mongoClientMap sync.Map

// InitPool init mongo driver client
func InitPool(opts *options.ClientOptions) *mongo.Client {
	defer rescueutil.Recover(func(err any) {
		fmt.Printf("InitPool error: %v", err)
	})
	client := assignutil.Assign(mongo.NewClient(opts))
	// connect
	ctx, cancel := context.WithTimeout(context.Background(), *opts.ConnectTimeout)
	defer cancel()
	errcheck.CheckAndPanic(client.Connect(ctx))
	// test connection
	pingCtx, pingCancel := context.WithTimeout(context.Background(), *opts.ServerSelectionTimeout)
	defer pingCancel()
	errcheck.CheckAndPanic(client.Ping(pingCtx, nil))
	return client
}

// DestoryPool destroy mongo client pool
func DestoryPool(client *mongo.Client) {
	defer rescueutil.Recover(func(err any) {
		fmt.Printf("DestoryPool error: %v", err)
	})
	errcheck.CheckAndPanic(client.Disconnect(context.Background()))
}

// DefaultOptions return default options
func DefaultOptions() *options.ClientOptions {
	hostStr := assignutil.Assign(config.GetString("mongo.hostStr"))
	hosts := strings.Split(hostStr, ",")
	if len(hosts) == 0 {
		panic("mongo.hosts is empty")
	}
	authMechanism := assignutil.Assign(config.GetString("mongo.auth.authMechanism"))
	username := assignutil.Assign(config.GetString("mongo.auth.username"))
	password := assignutil.Assign(config.GetString("mongo.auth.password"))
	authSource := assignutil.Assign(config.GetString("mongo.auth.authSource"))
	minPoolSize := assignutil.Assign(config.GetUint64("mongo.minPoolSizePerHost"))
	if minPoolSize == 0 {
		minPoolSize = 1
	}
	maxPoolSize := assignutil.Assign(config.GetUint64("mongo.maxPoolSizePerHost"))
	if maxPoolSize == 0 {
		maxPoolSize = 10
	}
	serverSelectionTimeout := assignutil.Assign(config.GetUint64("mongo.serverSelectionTimeout"))
	if serverSelectionTimeout == 0 {
		serverSelectionTimeout = 3000
	}
	connectTimeout := assignutil.Assign(config.GetUint64("mongo.connectTimeout"))
	if connectTimeout == 0 {
		connectTimeout = 3000
	}
	socketTimeout := assignutil.Assign(config.GetUint64("mongo.socketTimeout"))
	if socketTimeout == 0 {
		socketTimeout = 120000
	}
	maxConnIdleTime := assignutil.Assign(config.GetUint64("mongo.maxConnIdleTime"))
	if maxConnIdleTime == 0 {
		maxConnIdleTime = 180000
	}
	serverSelectionTo := time.Duration(serverSelectionTimeout) * time.Millisecond
	connectTo := time.Duration(connectTimeout) * time.Millisecond
	socketTo := time.Duration(socketTimeout) * time.Millisecond
	maxConnIdleTo := time.Duration(maxConnIdleTime) * time.Millisecond
	var readPref *readpref.ReadPref
	if !config.ContainsKey("mongo.preferPrimary") {
		readPref = readpref.SecondaryPreferred()
	} else {
		preferPrimary := assignutil.Assign(config.GetBool("mongo.preferPrimary"))
		if preferPrimary {
			readPref = readpref.PrimaryPreferred()
		} else {
			readPref = readpref.SecondaryPreferred()
		}
	}
	return &options.ClientOptions{
		Hosts:                  hosts,
		Auth:                   &options.Credential{AuthMechanism: authMechanism, Username: username, AuthSource: authSource, Password: password},
		MinPoolSize:            &minPoolSize,
		MaxPoolSize:            &maxPoolSize,
		ServerSelectionTimeout: &serverSelectionTo,
		ConnectTimeout:         &connectTo,
		SocketTimeout:          &socketTo,
		MaxConnIdleTime:        &maxConnIdleTo,
		ReadPreference:         readPref,
	}
}
