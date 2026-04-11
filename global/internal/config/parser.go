package config

import (
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gomooth/pkg/storage"
	"github.com/gomooth/utils/fsutil"
	"github.com/gomooth/xerror"
)

// ParseConfig 解析配置
func ParseConfig(content []byte) (*ProjectConfig, error) {
	var cfg ProjectConfig
	if _, err := toml.Decode(string(content), &cfg); nil != err {
		return nil, err
	}

	return &cfg, nil
}

const exampleConfigFilename = "config.example.toml" // APP 配置样例文件

func GetConfigFilename(filename string) (string, error) {
	cnfPath := storage.Disk("config")
	basePath, err := cnfPath.Path()
	if err != nil {
		return "", xerror.Wrap(err, "get config path failed")
	}
	exampleFilename := path.Join(basePath, exampleConfigFilename)
	localFilename := strings.ReplaceAll(exampleFilename, ".example.", ".")
	if len(filename) == 0 {
		if fsutil.Exist(localFilename) {
			return localFilename, nil
		}

		// 如果文件不存在，自动复制
		if !fsutil.Exist(exampleFilename) {
			return "", xerror.New("配置模板文件不存在")
		}

		if _, err := fsutil.Copy(exampleFilename, localFilename); nil != err {
			return "", xerror.Wrap(err, "复制配置模板失败")
		}

		return localFilename, nil
	}

	// 如果是远程连接，则从远程下载
	if strings.HasPrefix(filename, "https://") || strings.HasPrefix(filename, "http://") {
		if err := fsutil.Download(localFilename, filename); nil != err {
			return "", xerror.Wrapf(err, "get config from remote failed, url=%s", filename)
		}
		return localFilename, nil
	}

	return filename, nil
}

func ClearConfigExampleFile() {
	cnfPath := storage.Disk("config")
	basePath, err := cnfPath.Path()
	if err != nil {
		return
	}
	exampleFilename := path.Join(basePath, exampleConfigFilename)
	_ = os.Remove(exampleFilename)
}
