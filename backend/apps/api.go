package apps

import (
	"context"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-core/config"
	"github.com/crawlab-team/crawlab-core/controllers"
	"github.com/crawlab-team/crawlab-core/middlewares"
	"github.com/crawlab-team/crawlab-core/models"
	"github.com/crawlab-team/crawlab-core/routes"
	"github.com/crawlab-team/crawlab-db/mongo"
	"github.com/crawlab-team/crawlab-db/redis"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"time"
)

type Api struct {
	BaseApp
	app *gin.Engine
	srv *http.Server
}

func (app *Api) init() {
	// initialize config
	_ = app.initModule("config", config.InitConfig)

	// initialize mongo
	_ = app.initModule("mongo", mongo.InitMongo)

	// initialize redis
	_ = app.initModule("redis", redis.InitRedis)

	// initialize model services
	_ = app.initModule("mode-services", models.InitModelServices)

	// initialize controllers
	_ = app.initModule("controllers", controllers.InitControllers)

	// initialize middlewares
	_ = app.initModuleWithApp("middlewares", middlewares.InitMiddlewares)

	// initialize routes
	_ = app.initModuleWithApp("routes", routes.InitRoutes)
}

func (app *Api) run() {
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")
	address := net.JoinHostPort(host, port)
	app.srv = &http.Server{
		Handler: app.app,
		Addr:    address,
	}
	if err := app.srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Error("run server error:" + err.Error())
		} else {
			log.Info("server graceful down")
		}
	}
}

func (app *Api) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := app.srv.Shutdown(ctx); err != nil {
		log.Error("run server error:" + err.Error())
	}
}

func (app *Api) initModuleWithApp(name string, fn func(app *gin.Engine) error) (err error) {
	return app.initModule(name, func() error {
		return fn(app.app)
	})
}

func NewApi() *Api {
	app := gin.New()
	return &Api{
		app: app,
	}
}