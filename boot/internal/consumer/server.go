package consumer

import (
	"context"

	"github.com/gomooth/pkg/framework/app"
)

func New(ctx context.Context) app.IApp {
	// Note: 按需初始化
	return newSingle()
	//return newKafka("consumer")
}
