package helper

import (
	"strings"

	"github.com/gomooth/utils/valutil"
)

type ICMDArgsParser interface {
	Get(key string, alias ...string) string
	GetInt(key string, alias ...string) int
	GetBool(key string, alias ...string) bool
}

type cmdArgs struct {
	args map[string]string
}

func NewCMDArgs(args ...string) ICMDArgsParser {
	res := make(map[string]string)
	for _, arg := range args {
		vals := strings.SplitN(arg, ":", 2)
		if len(vals) == 2 {
			res[vals[0]] = vals[1]
		}
	}

	return &cmdArgs{
		args: res,
	}
}

func (ca cmdArgs) Get(key string, alias ...string) string {
	keys := append([]string{key}, alias...)
	for _, s := range keys {
		if v, ok := ca.args[s]; ok {
			return v
		}
	}
	return ""
}

func (ca cmdArgs) GetInt(key string, alias ...string) int {
	return valutil.Int(ca.Get(key, alias...))
}

func (ca cmdArgs) GetBool(key string, alias ...string) bool {
	return valutil.Bool(ca.Get(key, alias...))
}
