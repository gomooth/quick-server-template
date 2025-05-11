package consumer

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"server-api/app/consumer/internal/helper"
	"server-api/global"
	"server-api/global/kafka/topic"
	"testing"

	"github.com/gomooth/utils/fsutil"
)

func init() {
	if !fsutil.PathExist("storage") {
		cmd := exec.Command("ln", "-s", "../../storage/", "storage")
		if err := cmd.Run(); nil != err {
			log.Fatal(err)
		}
	}

	// 加载配置
	if err := global.ParseConfig("../../config/config.toml"); nil != err {
		log.Fatal(err)
	}

	// 初始化日志
	if err := global.InitLogger("test"); err != nil {
		log.Fatal(err)
	}

	// 初始化db
	if err := global.InitDataBase(); err != nil {
		log.Fatal(err)
	}
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
