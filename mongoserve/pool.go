package mongoserve

import (
	"sync"

	"github.com/jacksonCLyu/ridi-utils/utils/errcheck"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type clientPool struct {
	pool sync.Pool
}

var _mongoClientPoolMap sync.Map

func (cp *clientPool) getClientFromPool() *mongo.Client {
	client := cp.pool.Get().(*mongo.Client)
	if !CheckConnOK(client) {
		errcheck.CheckAndPanic(Conn(client))
	}
	return client
}

func (cp *clientPool) returnClientToPool(client *mongo.Client) {
	errcheck.CheckAndPanic(Reset(client))
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
