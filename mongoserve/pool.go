package mongoserve

import (
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type clientPool struct {
	pool sync.Pool
}

var _mongoClientPoolMap sync.Map

func (cp *clientPool) getClientFromPool() *mongo.Client {
	return cp.pool.Get().(*mongo.Client)
}

func (cp *clientPool) returnClientToPool(client *mongo.Client) {
	Reset(client)
	cp.pool.Put(client)
}

func newClientPool(opts *options.ClientOptions) *clientPool {
	return &clientPool{
		pool: sync.Pool{
			New: func() interface{} {
				return InitClient(opts)
			},
		},
	}
}
