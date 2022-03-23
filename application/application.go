package application

import (
	stdContext "context"
	"fmt"
	"myIris/application/libs"
	"myIris/application/libs/logging"
	"myIris/application/middleware"
	"myIris/service/cache"
	"path/filepath"
	"time"
	"myIris/application/controllers"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/pprof"
	"github.com/kataras/iris/v12/middleware/rate"
)

// HttpServer
type HttpServer struct {
	ConfigPath string
	App        *iris.Application
	Models     []interface{}
	Status     bool
}

func NewServer(config string) *HttpServer {
	app := iris.New()
	iris.RegisterOnInterrupt(func() {
		fmt.Println("iris 已经启动！")
	})
	httpServer := &HttpServer{
		ConfigPath: config,
		App:        app,
		Status:     false,
	}

	httpServer._Init()
	// httpServer 初始化后才可以加载到配置文件，感谢 @ren-ming  https://github.com/ren-ming 的提醒
	app.Logger().SetLevel(libs.Config.LogLevel)
	return httpServer
}

// Start
func (s *HttpServer) Start() error {
	if err := s.App.Run(
		iris.Addr(fmt.Sprintf("%s:%d", libs.Config.Host, libs.Config.Port)),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(time.RFC3339),
	); err != nil {
		return err
	}
	s.Status = true
	return nil
}

// Start close the server at 3-6 seconds
func (s *HttpServer) Stop() {
	go func() {
		time.Sleep(3 * time.Second)
		ctx, cancel := stdContext.WithTimeout(stdContext.TODO(), 3*time.Second)
		defer cancel()
		s.App.Shutdown(ctx)
		s.Status = false
	}()
}

func (s *HttpServer) _Init() error {
	err := libs.InitConfig(s.ConfigPath)
	if err != nil {
		logging.ErrorLogger.Errorf("系统配置初始化失败:", err)
		return err
	}
	if libs.Config.Cache.Driver == "redis" {
		cache.InitRedisCluster(libs.GetRedisUris(), libs.Config.Redis.Password)
	}
	s.RouteInit()
	return nil
}

// RouteInit
func (s *HttpServer) RouteInit() {
	s.App.UseRouter(middleware.CrsAuth())
	app := s.App.Party("/").AllowMethods(iris.MethodOptions)
	{

		// 开启 pprof 调试
		if libs.Config.Pprof {
			app.Get("/", func(ctx iris.Context) {
				ctx.HTML("<h1> Please click <a href='/debug/pprof'>here</a>")
			})

			p := pprof.New()
			app.Any("/debug/pprof", p)
			app.Any("/debug/pprof/{action:path}", p)
		}

		app.HandleDir("/uploads", iris.Dir(filepath.Join(libs.CWD(), "uploads")))
		v1 := app.Party("api/v1")
		{
			// 是否开启接口请求频率限制
			if !libs.Config.Limit.Disable {
				limitV1 := rate.Limit(libs.Config.Limit.Limit, libs.Config.Limit.Burst, rate.PurgeEvery(time.Minute, 5*time.Minute))
				v1.Use(limitV1)
			}
			v1.PartyFunc("/admin", func(admin iris.Party) { //casbin for gorm                                                   // <- IMPORTANT, register the middleware.

				admin.Get("/getData", controllers.GetCollyData).Name = "获取数据"

			})
		}
	}
}
