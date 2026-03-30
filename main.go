package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/greatbody/local-service-registry/internal/checker"
	"github.com/greatbody/local-service-registry/internal/handler"
	"github.com/greatbody/local-service-registry/internal/store"
)

func main() {
	addr := flag.String("addr", ":8500", "HTTP listen address")
	dbPath := flag.String("db", "registry.db", "SQLite database file path")
	interval := flag.Duration("interval", 5*time.Minute, "Health check interval")
	flag.Parse()

	// --- store ---
	st, err := store.New(*dbPath)
	if err != nil {
		log.Fatalf("failed to open store: %v", err)
	}
	defer st.Close()

	// --- health checker ---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chk := checker.New(st, *interval)
	go chk.Run(ctx)

	// --- HTTP server ---
	h := handler.New(st, chk)
	srv := &http.Server{
		Addr:    *addr,
		Handler: h,
	}

	// graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("shutting down...")
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	log.Printf("local-service-registry listening on %s", *addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
	log.Println("server stopped")
}
