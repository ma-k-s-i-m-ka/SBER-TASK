package server

import (
	"Sber/app/internal/cache"
	"Sber/app/internal/task"
	"Sber/app/pkg/config"
	"Sber/app/pkg/logger"
	_ "Sber/docs"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/browser"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"time"
)

type Server struct {
	srv     *http.Server
	log     *logger.Logger
	cfg     *config.Config
	handler *httprouter.Router
	cache   *cache.Cache
}

func NewServer(cfg *config.Config, handler *httprouter.Router, log *logger.Logger, cache *cache.Cache) *Server {
	return &Server{
		srv: &http.Server{
			Handler:      handler,
			WriteTimeout: time.Duration(cfg.HTTP.WriteTimeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.HTTP.ReadTimeout) * time.Second,
			Addr:         fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port),
		},
		log:     log,
		cfg:     cfg,
		handler: handler,
		cache:   cache,
	}
}

func (s *Server) Run(dbConn *pgx.Conn) error {

	reqTimeout := s.cfg.PostgreSQL.RequestTimeout

	taskStorage := task.NewStorage(dbConn, reqTimeout, s.cache)
	taskService := task.NewService(taskStorage, *s.log)
	taskHandler := task.NewHandler(*s.log, taskService, s.cache)
	taskHandler.Register(s.handler)
	s.log.Info("Initialized task routes")

	s.handler.Handler(http.MethodGet, "/docs/*any", httpSwagger.WrapHandler)
	s.log.Info("Initialized task documentation")

	err := browser.OpenURL("http://localhost:" + s.cfg.HTTP.Port + "/docs/")
	if err != nil {
		s.log.Error("Failed to open documentation in the browser:", err)
	}
	fs := http.FileServer(http.Dir("/public"))
	//fs := http.FileServer(http.Dir("public"))
	s.handler.Handler(http.MethodGet, "/", fs)
	s.handler.Handler(http.MethodGet, "/index.html", fs)
	err = browser.OpenURL("http://localhost:" + s.cfg.HTTP.Port + "/")
	if err != nil {
		s.log.Error("Failed to open website in the browser:", err)
	}
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
