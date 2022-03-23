<h1 align="center">reptile-Stock-market-valuation-go-iris</h1>

> 简单项目仅供学习，欢迎指点！
本项目为了学习golang的colly和iris框架，通过爬取上海证券交易所、深圳证券交易所和国家统计局数据官网，获得股市估值相关数据，并提供接口供前端使用。

本项目涉及框架：
iris
golang的Web 框架，号称宇宙最快的go语言Web 框架
colly
colly是用go实现的网络爬虫框架
go-cache
go-cache是一款类似于memached 的key/value 缓存软件。它比较适用于单机执行的应用程序。
---

- 安装项目依赖

>加载依赖管理包 (解决国内下载依赖太慢问题)
>使用国内七牛云的 go module 镜像。
>
>参考 https://github.com/goproxy/goproxy.cn。
>
>阿里： https://mirrors.aliyun.com/goproxy/
>
>官方： https://goproxy.io/
>
>中国：https://goproxy.cn
>
>其他：https://gocenter.io

##### golang 1.13+ 可以直接执行：
```shell script
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```
- 复制配置文件
```
cp application.example.yml application.yml
```

>  修改配置文件 `application.yml` 

- 运行项目
>如果想使用 `go run main.go --config ` 命令运行,注意不用 --config 指定配置路径，将无法加载配置文件
```
# --config 指定配置文件绝对路径
 go run main.go --config /Users/Phoenix/go/src/github.com/reptile-Stock-market-valuation-go-iris/application.yml
```



#### 感谢 
项目目录结构参考了https://github.com/snowlyg/iris-admin项目

