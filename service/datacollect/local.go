package datacollect

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"time"
)

type tokens []string
type skeys []string

var localCache *cache.Cache

type LocalDataDriver struct {
	Cache *cache.Cache
}

func NewLocalAuth() *LocalDataDriver {
	if localCache == nil {
		localCache = cache.New(4*time.Hour, 24*time.Minute)
	}
	return &LocalDataDriver{
		Cache: localCache,
	}
}

func (lc *LocalDataDriver)SyncCollectData(data map[string]string) error{
	if data == nil || len(data) == 0{
		return  errors.New("the collectdata is nil or empty")
	}
	var dataKey = CollectDataFlag + "collectData"
	lc.Cache.Set(dataKey,data,0)
	return nil
}

func (lc *LocalDataDriver)GetData() (error,map[string]string) {
	ans := make(map[string]string,0)
	var datas interface{}
	var uFound bool
	var dataKey = CollectDataFlag + "collectData"
	if datas,uFound = lc.Cache.Get(dataKey); !uFound{
		return errors.New("the collectdata is nil or empty"),nil
	}
	ans = datas.(map[string]string)
	return nil,ans
}

