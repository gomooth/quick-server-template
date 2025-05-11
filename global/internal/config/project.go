package config

import "server-api/global/internal/logger"

// 项目配置
type ProjectConfig struct {
	// 服务配置
	Server struct {
		HTTP struct {
			Addr string `toml:"addr"`
			Host string `toml:"host"`

			// 上传等静态资源配置
			Resource struct {
				Host string // 资源域名
				Path string // 资源上传目录
			}

			// Cache 缓存配置
			Cache redisConfig `toml:"cache"`
		} `toml:"http"`

		OpenAPI struct {
			Addr string `toml:"addr"`
			Host string `toml:"host"`

			// Cache 缓存配置
			Cache redisConfig `toml:"cache"`
		}
	} `toml:"server"`

	// 数据
	Data struct {
		// 持久存储数据库
		Persistent struct {
			Enabled     bool `toml:"enabled"`      // 是否启用
			AutoMigrate bool `toml:"auto_migrate"` // 是否自动迁移

			Connects []dBConnectConfig `toml:"connects"` // db 连接
		} `toml:"persistent"`

		// redis 配置
		Redis redisConfig

		// es 配置
		ElasticSearch struct {
			URLs         []string `toml:"urls"`
			SniffEnabled bool     `toml:"sniff_enabled"`
			DebugEnabled bool     `toml:"debug_enabled"`
		} `toml:"elasticsearch"`

		// locker 配置
		Locker struct {
			Enabled bool
			Engine  string
			Redis   redisConfig
		}

		// cache 配置
		Cache struct {
			Enabled bool
			Engine  string
			Redis   redisConfig
		}
	}

	// Consumer 消费者配置
	Consumer consumerConfig `toml:"consumer"`
	// Producer
	Producer producerConfig `toml:"producer"`

	// App 配置
	App struct {
		ID                 string `toml:"id"`
		Name               string `toml:"name"`
		Env                string `toml:"env"`                  // 系统环境: prod/production-生产环境，local-本地环境
		ClearExampleFile   bool   `toml:"clear_example_file"`   // 是否自动删除样例文件
		ClearConfigFile    bool   `toml:"clear_config_file"`    // 启动后是否自动删除配置文件
		WatchConfigEnabled bool   `toml:"watch_config_enabled"` // 监听配置文件变更，自动更新配置
		Secret             string `toml:"secret"`               // 密钥：jwt 认证等

		// 日志配置
		Log logger.LogConfig `toml:"log"`

		// 管理后台配置
		Admin struct {
			Account  string // 管理员帐号
			Password string // 管理员密码
		}
	}
}

type dBConnectConfig struct {
	Name        string // 连接名称
	Dsn         string // 连接
	Driver      string `toml:"type"`          // 数据库类型
	MaxIdle     int    `toml:"max_idle"`      // 最大空闲连接数
	MaxOpen     int    `toml:"max_open"`      // 最大连接数
	LogMode     bool   `toml:"log_mode"`      // 是否打印SQL
	MaxLifeTime int    `toml:"max_life_time"` // 连接存活时间
}

type redisConfig struct {
	Enabled        bool   `toml:"enabled"`
	Addr           string `toml:"addr"`            // 地址
	Password       string `toml:"auth"`            // 密码
	DB             int    `toml:"db"`              // 数据库
	Idle           int    `toml:"idle"`            // 最大连接数
	Active         int    `toml:"active"`          // 一次性活跃
	Wait           bool   `toml:"wait"`            // 是否等待空闲连接
	ConnectTimeout int64  `toml:"connect_timeout"` // 连接超时时间， 毫秒
}

type consumerConfig struct {
	Kafka struct {
		Addrs []string `toml:"addrs"`
	}

	Redis redisConfig `toml:"redis"`

	HTTPSQS struct {
		Addr     string `toml:"addr"`    // 地址
		Password string `toml:"auth"`    // 密码
		Timeout  int64  `toml:"timeout"` // 连接超时时间
	} `toml:"httpsqs"`
}

type producerConfig struct {
	Kafka struct {
		Addrs []string `toml:"addrs"`
	}
}
