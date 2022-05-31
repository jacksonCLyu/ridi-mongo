package mongoserve

import (
	"log"

	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetClient returns a mongo client
func GetClient(serviceName string) *mongo.Client {
	defer rescueutil.Recover(func(err any) {
		log.Printf("GetClient error: %v\n", err)
	})
	return GetClientWithOptions(serviceName, defaultClientOpts)
}

// GetClientWithOptions returns a mongo client with options
func GetClientWithOptions(serveName string, options *options.ClientOptions) *mongo.Client {
	defer rescueutil.Recover(func(err any) {
		log.Printf("GetClientWithOptions error: %v\n", err)
	})
	poolObj, _ := _mongoClientPoolMap.LoadOrStore(serveName, newClientPool(options))
	pool := poolObj.(*clientPool)
	return pool.getClientFromPool()
}

// ReleaseClient program release a connection
func ReleaseClient(serveName string, client *mongo.Client) {
	defer rescueutil.Recover(func(err any) {
		log.Printf("ReleaseClient error: %v\n", err)
	})
	poolObj, _ := _mongoClientPoolMap.Load(serveName)
	if poolObj == nil {
		return
	}
	pool := poolObj.(*clientPool)
	pool.returnClientToPool(client)
}
