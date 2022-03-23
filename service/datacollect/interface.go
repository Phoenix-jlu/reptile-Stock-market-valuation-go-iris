package datacollect

import (
	"errors"
	"time"
)

const (
	CollectDataFlag   =   "CDF"
)

var (
	ErrTokenInvalid               = errors.New("token is invalid")
	ZxwSessionUserMaxTokenDefault = 10
)

const (
	NoneScope uint64 = iota
	AdminScope
)

const (
	NoAuth int = iota
	AuthPwd
	AuthCode
	AuthThirdParty
)

const (
	LoginTypeWeb int = iota
	LoginTypeApp
	LoginTypeWx
	LoginTypeAlipay
	LoginApplet
)

var (
	RedisSessionTimeoutWeb    = 30 * time.Minute
	RedisSessionTimeoutApp    = 24 * time.Hour
	RedisSessionTimeoutWx     = 5 * 52 * 168 * time.Hour
	RedisSessionTimeoutApplet = 7 * 24 * time.Hour
)

type Session struct {
	UserId       string `json:"user_id" redis:"user_id"`
	LoginType    int    `json:"login_type" redis:"login_type"`
	AuthType     int    `json:"auth_type" redis:"auth_type"`
	CreationDate int64  `json:"creation_data" redis:"creation_data"`
	ExpiresIn    int    `json:"expires_in" redis:"expires_in"`
	Scope        uint64 `json:"scope" redis:"scope"`
}

type DataDriver interface {
	SyncCollectData(data map[string]string) error
	GetData() (error,map[string]string)
}
