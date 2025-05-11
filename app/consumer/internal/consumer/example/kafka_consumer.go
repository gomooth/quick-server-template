package example

import (
	"server-api/global"
)

func KafkaConsumer(topicName string, msg []byte) error {
	global.Log.Debugf("example kafka consumer handle, only print. topic=%s, msg=%s", topicName, msg)
	return nil
}
