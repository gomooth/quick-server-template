package global

import "fmt"

func GetServerRunningFlagKey(mode, name string) string {
	return fmt.Sprintf("serverRunningState:%s:%s", mode, name)
}
