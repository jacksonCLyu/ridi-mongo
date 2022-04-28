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

// GetOptions return options
func GetOptions(serviceName string) *options.ClientOptions {
	if !config.ContainsKey(serviceName) {
		panic("service config not found")
	}
	subConfig := assignutil.Assign(config.GetSection(serviceName))
	if !subConfig.ContainsKey("mongo.hostStr") {
		panic("mongo host not found")
	}
	hostStr := assignutil.Assign(subConfig.GetString("mongo.hostStr"))
	hosts := strings.Split(hostStr, ",")
	if len(hosts) == 0 {
		panic("mongo.hosts is empty")
	}
	authMechanism := assignutil.Assign(subConfig.GetString("mongo.auth.authMechanism"))
	username := assignutil.Assign(subConfig.GetString("mongo.auth.username"))
	password := assignutil.Assign(subConfig.GetString("mongo.auth.password"))
	authSource := assignutil.Assign(subConfig.GetString("mongo.auth.authSource"))
	var minPoolSize uint64
	if subConfig.ContainsKey("mongo.minPoolSizePerHost") {
		minPoolSize = assignutil.Assign(subConfig.GetUint64("mongo.minPoolSizePerHost"))
	} else {
		minPoolSize = 1
	}
	var maxPoolSize uint64
	if subConfig.ContainsKey("mongo.maxPoolSizePerHost") {
		maxPoolSize = assignutil.Assign(subConfig.GetUint64("mongo.maxPoolSizePerHost"))
	} else {
		maxPoolSize = 10
	}
	var serverSelectionTimeout uint64
	if subConfig.ContainsKey("mongo.serverSelectionTimeout") {
		serverSelectionTimeout = assignutil.Assign(subConfig.GetUint64("mongo.serverSelectionTimeout"))
	} else {
		serverSelectionTimeout = 3000
	}
	var connectTimeout uint64
	if subConfig.ContainsKey("mongo.connectTimeout") {
		connectTimeout = assignutil.Assign(subConfig.GetUint64("mongo.connectTimeout"))
	} else {
		connectTimeout = 3000
	}
	var socketTimeout uint64
	if subConfig.ContainsKey("mongo.socketTimeout") {
		socketTimeout = assignutil.Assign(subConfig.GetUint64("mongo.socketTimeout"))
	} else {
		socketTimeout = 120000
	}
	var maxConnIdleTime uint64
	if subConfig.ContainsKey("mongo.maxConnIdleTime") {
		maxConnIdleTime = assignutil.Assign(subConfig.GetUint64("mongo.maxConnIdleTime"))
	} else {
		maxConnIdleTime = 180000
	}
	serverSelectionTo := time.Duration(serverSelectionTimeout) * time.Millisecond
	connectTo := time.Duration(connectTimeout) * time.Millisecond
	socketTo := time.Duration(socketTimeout) * time.Millisecond
	maxConnIdleTo := time.Duration(maxConnIdleTime) * time.Millisecond
	var readPref *readpref.ReadPref
	if !subConfig.ContainsKey("mongo.preferPrimary") {
		readPref = readpref.SecondaryPreferred()
	} else {
		preferPrimary := assignutil.Assign(subConfig.GetBool("mongo.preferPrimary"))
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
