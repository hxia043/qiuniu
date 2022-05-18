package service

import (
	"fmt"
	"github/hxia043/qiuniu/internal/app/config"
	"github/hxia043/qiuniu/internal/app/controller"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

type Service struct {
	CloseCh chan struct{}
	engine  *gin.Engine
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "OPTIONS OK")
		}

		c.Next()
	}
}

func (s *Service) Run() error {
	c := controller.New()

	s.engine.Use(Cors())
	c.RegisterRouter(s.engine)

	go func() {
		err := s.engine.Run(fmt.Sprintf("%s:%s", config.Config.ServiceIp, config.Config.ServicePort))
		if err != nil {
			close(s.CloseCh)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	select {
	case <-s.CloseCh:
	case <-signals:
	}

	return nil
}

func New() *Service {
	return &Service{
		CloseCh: make(chan struct{}),
		engine:  gin.Default(),
	}
}
