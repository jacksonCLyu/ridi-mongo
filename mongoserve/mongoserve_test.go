package mongoserve

import (
	"context"
	"testing"

	"github.com/jacksonCLyu/ridi-config/pkg/config"
	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/jacksonCLyu/ridi-utils/utils/errcheck"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetClient(t *testing.T) {
	config.Init()
	cfg := assignutil.Assign(config.NewConfig(config.WithFilePath("./testdata/testConfig.yaml")))
	errcheck.CheckAndPanic(config.Init(config.WithConfigurable(cfg)))
	client := GetClient("test")
	collection := client.Database("flight").Collection("a_temp")
	cursor := assignutil.Assign(collection.Find(context.Background(), bson.M{}))
	for cursor.Next(context.Background()) {
		var result bson.M
		errcheck.CheckAndPanic(cursor.Decode(&result))
		t.Log(result)
	}
	oClient := GetClient("test")
	if client != oClient {
		t.Error("client is not equal")
	}
}
