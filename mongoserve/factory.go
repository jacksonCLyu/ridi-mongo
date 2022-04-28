package mongoserve

import (
	"github.com/jacksonCLyu/ridi-log/log"
	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetClient returns a mongo client
func GetClient(serviceName string) *mongo.Client {
	defer rescueutil.Recover(func(err any) {
		log.Errorf("GetClient error: %v", err)
	})
	return GetClientWithOptions(serviceName, GetOptions(serviceName))
}

// GetClientWithOptions returns a mongo client with options
func GetClientWithOptions(serveName string, options *options.ClientOptions) *mongo.Client {
	defer rescueutil.Recover(func(err any) {
		log.Errorf("GetClientWithOptions error: %v", err)
	})
	client, _ := mongoClientMap.LoadOrStore(serveName, InitPool(options))
	return client.(*mongo.Client)
}

// ReleaseClient program release a connection
func ReleaseClient(hostStr string, client *mongo.Client) {
}
