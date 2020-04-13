package httpapp

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

type HttpApp struct {
	Name   string
	Engine *bm.Engine //http server
	//cfg     bm.ServerConfig
}

//http server
//// ServerConfig is the bm server config model
//type ServerConfig struct {
//	Network      string         `dsn:"network"`
//	Addr         string         `dsn:"address"`
//	Timeout      xtime.Duration `dsn:"query.timeout"`
//	ReadTimeout  xtime.Duration `dsn:"query.readTimeout"`
//	WriteTimeout xtime.Duration `dsn:"query.writeTimeout"`
//}
type Cfg = bm.ServerConfig
type RouterGroup = bm.RouterGroup //
type HandlerFunc = bm.HandlerFunc //type HandlerFunc func(*Context)

var svc HttpApp

//New new a http server
func New(name string, conf *Cfg) (s *HttpApp, err error) {
	if name == "" {
		svc.Name = "demo"
	} else {
		svc.Name = name
	}
	if conf != nil {
		svc.Engine = bm.DefaultServer(conf)
	} else {
		var conf Cfg
		conf.Network = "tcp"
		conf.Addr = "0.0.0.0:8000"
		conf.Timeout = 1
		svc.Engine = bm.DefaultServer(&conf)
	}

	initRouter(svc.Engine)
	return &svc, nil
}

func (svc *HttpApp) Run() (err error) {
	err = svc.Engine.Start()
	if err != nil {
		return
	}
	log.Info("%s start success .", svc.Name)
	svc.waitNotify(func() {})
	return nil
}

func (svc *HttpApp) RouterGroupGet(relativeRouterGroupPath string, relativePath string, handlers ...HandlerFunc) {
	rg := svc.Engine.Group(relativeRouterGroupPath)
	rg.GET(relativePath, handlers...)
}
func (svc *HttpApp) RouterGroupPost(relativeRouterGroupPath string, relativePath string, handlers ...HandlerFunc) {
	rg := svc.Engine.Group(relativeRouterGroupPath)
	rg.POST(relativePath, handlers...)
}

func (svc *HttpApp) Get(relativePath string, handlers ...HandlerFunc) {
	svc.Engine.GET(relativePath, handlers...)
}
func (svc *HttpApp) Post(relativePath string, handlers ...HandlerFunc) {
	svc.Engine.POST(relativePath, handlers...)
}

func initRouter(e *bm.Engine) {
	e.GET("/", howToStart)
	//g := e.Group("/demo")
	//{
	//	g.GET("/start", howToStart)
	//}
}

// example for http request handler.
func howToStart(c *bm.Context) {
	c.JSON(" hello webs demo.", nil)
}

func (svc *HttpApp) waitNotify(closeFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())

		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info("%s exit .", svc.Name)
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
