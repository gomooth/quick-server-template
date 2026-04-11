package consumer

import (
	"context"
	"fmt"
	"os"
	"testing"

	"server-api/app/consumer/internal/helper"
	"server-api/global/kafka/topic"
	"server-api/internal/testhelper"
)

func TestMain(m *testing.M) {
	// Setup test environment with database support
	if err := testhelper.SetupTestWithDB(); err != nil {
		panic(err)
	}

	code := m.Run()
	testhelper.Cleanup()
	os.Exit(code)
}

func TestKafkaPusherExample(t *testing.T) {
	ctx := context.Background()
	pusher := helper.GetProducer()

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf("example data hello %d", i)
		if err := pusher.Produce(ctx, topic.ExampleData, []byte(msg)); nil != err {
			t.Fatal(err)
			return
		}
	}
}
