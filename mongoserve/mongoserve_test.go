package mongoserve

import (
	"context"
	"testing"

	"github.com/jacksonCLyu/ridi-utils/utils/assignutil"
	"github.com/jacksonCLyu/ridi-utils/utils/errcheck"
	"github.com/jacksonCLyu/ridi-utils/utils/rescueutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetClient(t *testing.T) {
	defer rescueutil.Recover(func(err any) {
		t.Errorf("TestGetClient error: %v\n", err)
	})
	if err := Init(WithClientOpts(&options.ClientOptions{})); err != nil {
		t.Errorf("Init error: %v", err)
	}
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
