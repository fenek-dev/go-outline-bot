package server

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CancelDeposit(ctx context.Context, transactionExternalID string) (err error)
	ConfirmDeposit(ctx context.Context, transactionExternalID string) (err error)
}

type HTTPServer struct {
	port       string
	service    Service
	router     *gin.Engine
	log        *slog.Logger
	stopSignal chan os.Signal
}

type Option func(*HTTPServer)

func New(port string, service Service, stopSignal chan os.Signal, opts ...Option) *HTTPServer {
	server := &HTTPServer{
		port:       port,
		service:    service,
		router:     gin.Default(),
		stopSignal: stopSignal,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *HTTPServer) {
		s.log = logger
	}
}

func (s *HTTPServer) Run() error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-s.stopSignal
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")

	return nil
}

func (s *HTTPServer) Handle(method, path string, handler gin.HandlerFunc) {
	s.router.Handle(method, path, handler)
}

func (s *HTTPServer) AddMiddleware(middleware gin.HandlerFunc) {
	s.router.Use(middleware)
}
