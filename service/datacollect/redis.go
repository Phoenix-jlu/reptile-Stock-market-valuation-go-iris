package datacollect

import (
	"myIris/service/cache"
)

type RedisDriver struct {
	Conn *cache.RedisCluster
}

func NewRedisDriver() *RedisDriver {
	return &RedisDriver{
		Conn: cache.GetRedisClusterClient(),
	}
}

func (ra *RedisDriver)SyncCollectData(data map[string]string) error{
	return nil
}

func (ra *RedisDriver)GetData() (error,map[string]string) {

	return nil,nil
}