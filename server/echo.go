package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type (
	IServer interface {
		Start(pctx context.Context)
	}

	server struct {
		app *echo.Echo
		db  *gorm.DB
	}
)

func NewEchoServer(db *gorm.DB) IServer {
	return &server{
		app: echo.New(),
		db:  db,
	}
}

func (s *server) gracefulShutdown(pctx context.Context, close <-chan os.Signal) {
	resClose := <-close

	if resClose != nil {
		log.Println("Shutting down server...")

		ctx, cancel := context.WithTimeout(pctx, time.Second*10)
		defer cancel()

		if err := s.app.Shutdown(ctx); err != nil {
			log.Fatalf("Failed to shutdown")
			panic(err)
		}

	}

	log.Println("Shutting down Server...")
}

func (s *server) Start(pctx context.Context) {

	s.app.Use(middleware.Logger())

	close := make(chan os.Signal, 1)
	signal.Notify(close, syscall.SIGINT, syscall.SIGTERM)

	s.addDataService()

	go s.gracefulShutdown(pctx, close)

	s.app.Start(":8081")

}
