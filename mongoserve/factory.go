package mongoserve

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetClient returns a mongo client
func GetClient(hostStr string) *mongo.Client {
	return GetClientWithOptions(hostStr, DefaultOptions())
}

// GetClientWithOptions returns a mongo client with options
func GetClientWithOptions(serveName string, options *options.ClientOptions) *mongo.Client {
	client, _ := mongoClientMap.LoadOrStore(serveName, InitPool(options))
	return client.(*mongo.Client)
}

// ReleaseClient program release a connection
func ReleaseClient(hostStr string, client *mongo.Client) {
}
