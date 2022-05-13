package mongoserve

import (
	"context"
	"errors"
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

// Init init module
func Init(opts ...Option) error {
	initOptions := &initOptions{
		configer: config.L(),
	}
	for _, opt := range opts {
		opt.apply(initOptions)
	}
	if initOptions.configer == nil {
		return errors.New("config is nil")
	}
	config.SetDefaultConfig(initOptions.configer)
	return nil
}

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
	if !config.ContainsKey(HostKey) {
		panic("mongo host not found")
	}
	hostStr := assignutil.Assign(config.GetString(HostKey))
	hosts := strings.Split(hostStr, ",")
	if len(hosts) == 0 {
		panic("mongo.hosts is empty")
	}
	authMechanism := assignutil.Assign(config.GetString(AuthMechanismKey))
	username := assignutil.Assign(config.GetString(UsernameKey))
	password := assignutil.Assign(config.GetString(PasswordKey))
	authSource := assignutil.Assign(config.GetString(AuthSourceKey))
	var minPoolSize uint64
	if config.ContainsKey(MinPoolSizeKey) {
		minPoolSize = assignutil.Assign(config.GetUint64(MinPoolSizeKey))
	} else {
		minPoolSize = 1
	}
	var maxPoolSize uint64
	if config.ContainsKey(MaxPoolSizeKey) {
		maxPoolSize = assignutil.Assign(config.GetUint64(MaxPoolSizeKey))
	} else {
		maxPoolSize = 10
	}
	var serverSelectionTimeout uint64
	if config.ContainsKey(ServerSelectionTimeoutKey) {
		serverSelectionTimeout = assignutil.Assign(config.GetUint64(ServerSelectionTimeoutKey))
	} else {
		serverSelectionTimeout = 3000
	}
	var connectTimeout uint64
	if config.ContainsKey(ConnectTimeoutKey) {
		connectTimeout = assignutil.Assign(config.GetUint64(ConnectTimeoutKey))
	} else {
		connectTimeout = 3000
	}
	var socketTimeout uint64
	if config.ContainsKey(SocketTimeoutKey) {
		socketTimeout = assignutil.Assign(config.GetUint64(SocketTimeoutKey))
	} else {
		socketTimeout = 120000
	}
	var maxConnIdleTime uint64
	if config.ContainsKey(MaxConnIdleTimeKey) {
		maxConnIdleTime = assignutil.Assign(config.GetUint64(MaxConnIdleTimeKey))
	} else {
		maxConnIdleTime = 180000
	}
	serverSelectionTo := time.Duration(serverSelectionTimeout) * time.Millisecond
	connectTo := time.Duration(connectTimeout) * time.Millisecond
	socketTo := time.Duration(socketTimeout) * time.Millisecond
	maxConnIdleTo := time.Duration(maxConnIdleTime) * time.Millisecond
	var readPref *readpref.ReadPref
	if !config.ContainsKey(ReadPreferenceKey) {
		readPref = readpref.SecondaryPreferred()
	} else {
		preferPrimary := assignutil.Assign(config.GetBool(ReadPreferenceKey))
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
