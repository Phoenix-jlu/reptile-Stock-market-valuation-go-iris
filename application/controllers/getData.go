package controllers

import (
	"errors"
	"github.com/kataras/iris/v12"
	"myIris/application/libs/response"
	"myIris/service/datacollect"
	"strconv"
)


func GetCollyData(ctx iris.Context)  {
	var ans = make(map[string]string,0)
	var err error
	dataDriver := datacollect.NewDataCollectDriver()
	err,ans = dataDriver.GetData()

	floatShangZheng,_ := strconv.ParseFloat(ans["ShangZhengSZ"],64)
	floatShenZheng,_ := strconv.ParseFloat(ans["ShenZhengSZ"],64)
	floatGNGDP,_ := strconv.ParseFloat(ans["GNGDP"],64)
	ans["股市估值泡沫"] = strconv.FormatFloat((floatShangZheng + floatShenZheng) / floatGNGDP,'f',3,64)
	if err != nil {
		ctx.JSON(response.NewResponse(response.SystemErr.Code,nil,err.Error()))
	}
	if ans == nil {
		ctx.JSON(response.NewResponse(response.DataEmptyErr.Code,nil,errors.New("the collectdata is nil or empty").Error()))
	}
	ctx.JSON(response.NewResponse(response.NoErr.Code,ans,response.NoErr.Msg))
}
