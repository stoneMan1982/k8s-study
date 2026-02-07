package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"demo/internal/bootstrap"
	"demo/internal/user"
	"demo/pkg/module"
)

func hello(w http.ResponseWriter, req *http.Request) {
	host, _ := os.Hostname()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "not set"
	}
	io.WriteString(w, fmt.Sprintf("[v3]Hello, Kubernetes, From host: %s, DB_URL: %s\n", host, dbURL))
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	bootstrap.RegisterInternalModules(module.DefaultRegistry)
	if err := module.StartModules(ctx); err != nil {
		panic(err)
	}
	if api, ok := module.GetService[*user.API]("user"); ok {
		log.Println("user module status:", api.Status())
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = module.StopModules(shutdownCtx)
	}()

	http.HandleFunc("/", hello)

	server := &http.Server{Addr: ":3000"}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
	log.Println("server exited gracefully")
}
