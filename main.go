package main

import (
	"context"
	"errors"
	"fmt"
	"gin-test/database"
	"gin-test/handler"
	"gin-test/middleware"
	"gin-test/scanner"
	"gin-test/service/job"
	"gin-test/util/flag"
	"gin-test/util/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	UUID = "uuid"
)

func main() {
	database.Init()
	apps := Init()
	err := Start(apps)
	if err != nil {
		log.Log(fmt.Sprintf("start failed: %s\n", err))
	}

	f, _ := os.Create(flag.LogPath)
	gin.DefaultWriter = io.MultiWriter(f)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	//路由上加入session中间件
	r.Use(sessions.Sessions(UUID, store))

	r.POST("/signup", handler.SignUp())
	r.POST("/login", handler.Login())
	r.GET("/hello", handler.GetHello())

	middleware.InitMiddlewares(r)
	r.POST("/hello", handler.Hello())
	r.POST("/logout", handler.Logout())

	r.POST("/add_safe", handler.AddSafe())
	r.POST("/say_safe", handler.SaySafe())
	r.GET("/get_safe", handler.GetSafe())

	srv := &http.Server{
		//0.0.0.0:8080
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		log.Logger("server started on: [%s]\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Logger("listen: %s\n", err)
		}
	}()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	sig := <-c
	log.Logger("[main],app stopping,receive:%v\n", sig)

	// stop main apps
	Stop(apps)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer ctx.Done()
	if err := srv.Shutdown(ctx); err != nil {
		log.Logger("Server forced to shutdown:%v\n", err)
	}
	log.Logger("[main],app stopped,receive:%v\n", sig)
}

func Init() (apps []App) {
	newScanner := scanner.NewScanner()
	newScanner.AddService("0 0 6 * * ?", job.NewSafeJobService())

	apps = append(apps, newScanner)
	return apps
}

func Start(apps []App) error {
	for _, a := range apps {
		err := a.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

func Stop(apps []App) {

}

type App interface {
	Start() error
	Stop() error
}
