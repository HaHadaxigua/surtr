package http

import (
	"context"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	. "github.com/HaHadaxigua/surtr/global"
	"github.com/HaHadaxigua/surtr/internal/http/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Service struct {
	server *http.Server
	conf   *Config
}

func New(configPath ...string) (*Service, error) {
	engin := gin.Default()
	engin.Use(middlewares.CrossMiddleware())
	engin.MaxMultipartMemory = 5 << 30

	//engin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	apiRouter := engin.Group("/api")
	routers(apiRouter)

	var (
		conf *Config
		err  error
		path string
	)
	if configPath != nil {
		path = configPath[0]
	}
	conf, err = NewConfig(path)
	if err != nil {
		return nil, err
	}
	return &Service{
		server: &http.Server{
			Addr:           conf.Addr,
			Handler:        engin,
			MaxHeaderBytes: 1 << 20,
		},
		conf: conf,
	}, nil
}

func (s *Service) Start(ctx context.Context) error {
	defer func() {
		fmt.Println()
		logrus.Infof("GoodBye %s", Surtr)
	}()
	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		case <-ctx.Done():
			s.server.Close()
			return
		case <-quit:
			s.server.Close()
			return
		}
	}()

	logrus.Infof("listen on %s", s.conf.Addr)
	return s.server.ListenAndServe()
}

type data struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Content any    `json:"content"`
}

func Ok(val any) *data {
	return &data{
		Code:    0,
		Message: "OK",
		Content: val,
	}
}

func Err(err error) *data {
	return &data{
		Code:    -1,
		Message: err.Error(),
	}
}
