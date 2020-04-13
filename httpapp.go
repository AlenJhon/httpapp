package httpapp

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kratos/kratos/pkg/log"
	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

type AppHttp struct {
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
type AppHttpCfg bm.ServerConfig
type AppHttpRouterGroup bm.RouterGroup //
type HandlerFunc bm.HandlerFunc        //type HandlerFunc func(*Context)

var svc AppHttp

//New new a http server
func New(name string, conf *AppHttpCfg) (s *AppHttp, err error) {
	svc.Name = name
	svc.Engine = bm.DefaultServer(conf)

	initRouter(engine)
	return &svc, nil
}

func (s *AppHttp) Run() (err error) {
	err = s.Engine.Start()
	if err != nil {
		return
	}
	log.Info(s.Name, " start success .")
	s.waitNotify(func() {})
	return nil
}

//register router group
func (s *AppHttp) RouterGroup(relativePath string, handlers ...HandlerFunc) (g *AppHttpRouterGroup) {
	g = s.Engine.Group(relativePath, handlers)
	return
}
func (rg *AppHttpRouterGroup) Get(relativePath string, handlers ...HandlerFunc) {
	rg.Get(relativePath, handlers)
}
func (rg *AppHttpRouterGroup) Post(relativePath string, handlers ...HandlerFunc) {
	rg.Post(relativePath, handlers)
}

func (s *AppHttp) Get(relativePath string, handlers ...HandlerFunc) {
	s.Engine.Get(relativePath, handlers)
}
func (s *AppHttp) Post(relativePath string, handlers ...HandlerFunc) {
	s.Engine.Post(relativePath, handlers)
}

func initRouter(e *bm.Engine) {
	e.Get("/", howToStart)
	//g := e.Group("/demo")
	//{
	//	g.GET("/start", howToStart)
	//}
}

// example for http request handler.
func howToStart(c *bm.Context) {
	c.JSON(" hello webs demo.", nil)
}

func (s *AppHttp) waitNotify(closeFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())

		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			log.Info(s.Name, " exit .")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
