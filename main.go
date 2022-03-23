package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocolly/colly/v2"
	"myIris/application"
	"myIris/application/libs"
	"myIris/application/libs/logging"
	"myIris/service/datacollect"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var config = flag.String("config", "", "配置路径")
var version = flag.Bool("version", false, "打印版本号")
var Version = "master"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options] [command]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "  -config <path>\n")
		fmt.Fprintf(os.Stderr, "    设置项目配置文件路径，可选\n")
		fmt.Fprintf(os.Stderr, "  -version <true or false> 打印项目版本号，默认为: false\n")
		fmt.Fprintf(os.Stderr, "    打印版本号\n")
		fmt.Fprintf(os.Stderr, "\n")
	}

	flag.Parse()
	if *version {
		fmt.Println(fmt.Sprintf("版本号：%s\n", Version))
	}
	irisServer := application.NewServer(*config)
	if irisServer == nil {
		panic("http server 初始化失败")
	}
	if libs.IsPortInUse(libs.Config.Port) { if !irisServer.Status {
		panic(fmt.Sprintf("端口 %d 已被使用\n", libs.Config.Port))
	}
		irisServer.Stop() // 停止
		}
	go collyData()
		err := irisServer.Start()
		if err != nil {
			panic(fmt.Sprintf("http server 启动失败: %+v", err))
		}
		logging.InfoLogger.Infof("http server %s:%d start", libs.Config.Host, libs.Config.Port)
}
func collyData(){
	c1 := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"), colly.MaxDepth(100))
	c2 := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"), colly.MaxDepth(100))
	c3 := colly.NewCollector(colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36"), colly.MaxDepth(100))
	ansSlic := make(map[string]string,0) //采集器1，获取文章列表
	c1.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "yfx_c_g_u_id_10000042=_ck22022110333818911875217165557; yfx_mr_f_10000042=%3A%3Amarket_type_free_search%3A%3A%3A%3Abaidu%3A%3A%3A%3A%3A%3A%3A%3Awww.baidu.com%3A%3A%3A%3Apmf_from_free_search; yfx_mr_10000042=%3A%3Amarket_type_free_search%3A%3A%3A%3Abaidu%3A%3A%3A%3A%3A%3A%3A%3Awww.baidu.com%3A%3A%3A%3Apmf_from_free_search; yfx_key_10000042=; VISITED_MENU=%5B%2211913%22%2C%228466%22%2C%229692%22%2C%228765%22%5D; yfx_f_l_v_t_10000042=f_t_1645410818713__r_t_1647919472432__v_t_1647919794579__r_c_4")
		r.Headers.Set("Referer","http://www.sse.com.cn/")
		r.Headers.Set("Accept","*/*")
		r.Headers.Set("Connection","keep-alive")
		r.Headers.Set("Accept-Encoding","gzip, deflate")
		r.Headers.Set("Accept-Language","zh-CN,zh;q=0.9")
	})
	c3.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", "_trs_uv=l121nemx_6_ky4e; JSESSIONID=qdCxWaNQJzeR3489FXK3ujOE9Ar3TPnBz4JPiFhtLF1ujBlTSjlP!-403300632; u=5")
		r.Headers.Set("Referer","https://data.stats.gov.cn/easyquery.htm?cn=C01")
		r.Headers.Set("Accept","application/json, text/javascript, */*; q=0.01")
		r.Headers.Set("Connection","keep-alive")
		r.Headers.Set("Accept-Encoding","gzip, deflate, br")
		r.Headers.Set("Accept-Language","zh-CN,zh;q=0.9")

		r.Headers.Set("X-Requested-With","XMLHttpRequest")
		r.Headers.Set("Sec-Fetch-Site","same-origin")
		r.Headers.Set("Sec-Fetch-Mode","cors")
		r.Headers.Set("Sec-Fetch-Dest","empty")
	})

	c3.WithTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},})
	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c2.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c3.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	//c1.OnHTML("div", func(e *colly.HTMLElement) { //e.ForEach("class",func(i int,eChild *colly.HTMLElement) {  // fmt.Println(e.Attr("class"))  //})
	//	if (e.Attr("class")!= ""){
	//		fmt.Println(e.Attr("class"))
	//	}
	//})
	c1.OnResponse(func (resp *colly.Response) {
		result := make(map[string]interface{}, 0)
		//var item = make([]JsonResponse,0)
		err := json.Unmarshal([]byte(strings.Split(strings.Split(string([]byte(resp.Body)),"(")[1],")")[0]), &result)
		if err != nil {
			fmt.Println("c1.OnResponse err:",err)
		}
		dataValue,ok := result["result"].([]interface{})
		if ok{
			for _,j := range dataValue{
				dataValueUIn,okin := j.(map[string]interface{})
				if okin{
					if dataValueUIn["PRODUCT_NAME"] == "股票"{
						ansSlic["ShangZhengSZ"] = dataValueUIn["TOTAL_VALUE"].(string)
					}
				}
			}
		}
	})

	c2.OnResponse(func (resp *colly.Response) {
		result := make([]map[string]interface{}, 0)
		//var item = make([]JsonResponse,0)
		err := json.Unmarshal([]byte(resp.Body), &result)
		if err != nil {
			fmt.Println("c2.OnResponse err:",err)
		}
		dataValue,ok := result[0]["data"].([]interface{})
		if ok{
			for _,j := range dataValue{
				dataValueUIn,okin := j.(map[string]interface{})
				if okin{
					if dataValueUIn["zbmc"] == "股票总市值（亿元）"{
						ansSlic["ShenZhengSZ"] = strings.Split(dataValueUIn["brsz"].(string),",")[0] + strings.Split(dataValueUIn["brsz"].(string),",")[1]
					}
				}
			}
		}
	})

	c3.OnResponse(func (resp *colly.Response) {
		result := make(map[string]interface{}, 0)
		//var item = make([]JsonResponse,0)
		err := json.Unmarshal([]byte(resp.Body), &result)
		if err != nil {
			fmt.Println("c3.OnResponse err:",err)
		}
		dataValuefirst,_ := result["returndata"].(map[string]interface{})
		dataValue,ok := dataValuefirst["datanodes"].([]interface{})
		if ok{
			for _,j := range dataValue{
				dataValueUIn,okin := j.(map[string]interface{})
				if okin{
					if dataValueUIn["code"] == "zb.A020102_sj.2021"{
						finalans,_ := dataValueUIn["data"].(map[string]interface{})
						ansSlic["GNGDP"] = finalans["strdata"].(string)
					}
				}
			}
		}
	})

	c1.AllowURLRevisit = true
	c2.AllowURLRevisit = true
	c3.AllowURLRevisit = true
	for{
		currentTime := time.Now()
		oldTime := currentTime.AddDate(0, 0, -1)
		year := oldTime.Year()
		month := oldTime.Format("01")
		day := oldTime.Format("02")
		yearS := strconv.Itoa(year)
		ymd := yearS + "-" + month + "-" + day
		urlSZ := "http://www.szse.cn/api/report/ShowReport/data?SHOWTYPE=JSON&CATALOGID=1803_after&TABKEY=tab1&txtQueryDate=" + ymd + "&random=0.15087199301827603"

		c1.Visit("http://query.sse.com.cn/commonQuery.do?jsonCallBack=jsonpCallback72734336&sqlId=COMMON_SSE_SJ_GPSJ_GPSJZM_TJSJ_L&PRODUCT_NAME=%E8%82%A1%E7%A5%A8%2C%E4%B8%BB%E6%9D%BF%2C%E7%A7%91%E5%88%9B%E6%9D%BF&type=inParams&_=1647919794638")

		c2.Visit(urlSZ)

		c3.Visit("https://data.stats.gov.cn/easyquery.htm?m=QueryData&dbcode=hgnd&rowcode=zb&colcode=sj&wds=%5B%5D&dfwds=%5B%5D&k1=1648019828080&h=1")
		dataDriver := datacollect.NewDataCollectDriver()
		dataDriver.SyncCollectData(ansSlic)
		time.Sleep(time.Duration(60)*time.Second)
		ansSlic  = make(map[string]string,0)
	}
}


