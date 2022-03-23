package libs

import (
	"fmt"
	"github.com/jinzhu/configor"
	"path/filepath"
	"strings"
)

var Config = struct {
	LogLevel string `env:"Loglevel" default:"info"`
	Host     string `env:"Host" default:"127.0.0.1" `
	Port     int64  `env:"Port" default:"80"`
	MaxSize  int64  `env:"MaxSize" default:""`
	Pprof    bool   `env:"Pprof" default:"false"`
	Cache    struct {
		Driver string `env:"CacheDriver" default:"local"`
	}
	Redis struct {
		Host     string `env:"RedisHost" default:"localhost"`
		Port     string `env:"RedisPort" default:"6379"`
		Password string `env:"RedisPassword" default:""`
	}
	Limit struct {
		Disable bool    `env:"LimitDisable" default:"true"`
		Limit   float64 `env:"LimitLimit" default:"1"`
		Burst   int     `env:"LimitBurst" default:"5"`
	}
	Qiniu struct {
		Enable    bool   `env:"QiniuEnable" default:"false"`
		Host      string `env:"QiniuHost" default:""`
		Accesskey string `env:"QiniuAccesskey" default:""`
		Secretkey string `env:"QiniuSecretkey" default:""`
		Bucket    string `env:"QiniuBucket" default:""`
	}
}{}

func InitConfig(config string) error {
	path := filepath.Join(CWD(), "application.yml")
	if config != "" {
		path = config
	}

	if err := configor.Load(&Config, path); err != nil {
		return err
	}

	if Config.MaxSize == 0 {
		Config.MaxSize = 5 << 20
	}

	return nil
}

// redis 集群
func GetRedisUris() []string {
	addrs := make([]string, 0, 0)
	hosts := strings.Split(Config.Redis.Host, ";")
	ports := strings.Split(Config.Redis.Port, ";")
	for _, h := range hosts {
		for _, p := range ports {
			addrs = append(addrs, fmt.Sprintf("%s:%s", h, p))
		}
	}
	return addrs
}
