package datacollect

import "myIris/application/libs"

var dataDriver DataDriver

func NewDataCollectDriver()  DataDriver {
	if dataDriver != nil {
		return dataDriver
	}
	switch libs.Config.Cache.Driver {
	case "reids":
		return NewRedisDriver()
	case "local":
		return NewLocalAuth()
	default:
		return NewLocalAuth()
	}
}