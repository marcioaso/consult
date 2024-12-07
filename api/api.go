package api

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/marcioaso/consult/api/kline"
	"github.com/marcioaso/consult/api/top10"
)

func SetupServer(port string) *echo.Echo {
	e := echo.New()

	configureHandlers(e)

	go func() {
		fmt.Printf("Server running on port %s\n", port)
		if err := e.Start(":" + port); err != nil {
			e.Logger.Info("Shutting down the server...")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	fmt.Println("Server stopped.")
	return e
}

func configureHandlers(e *echo.Echo) {
	e.GET("/status", StatusHandler)
	e.GET("/top10", top10.Top10Handler)
	e.GET("/kline", kline.KlineHandler)
}
