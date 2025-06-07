package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/slangeres/Vypaar/backend_API/internal/config"
	"github.com/slangeres/Vypaar/backend_API/internal/https/handler"
	"github.com/slangeres/Vypaar/backend_API/internal/storage/sqllite"
)

func main() {

	//!Configuration
	cnf := config.MustDone()

	//!Configure the DB
	db, err := sqllite.ConfigSQL(cnf)

	if err != nil {
		slog.Warn("databse connection issue .....")

	}

	//^Making the router
	router := http.NewServeMux()

	//!hander functions

	router.HandleFunc("GET /api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from vypaar api"))
	})

	router.HandleFunc("POST /api/v1/product", handler.PostProduct(db))

	router.HandleFunc("GET /api/v1/product", handler.GetProduct())

	router.HandleFunc("GET /api/v1/product/{id}", handler.GetProductById())

	router.HandleFunc("PATCH /api/v1/product/{id}", handler.UpdateProduct())

	router.HandleFunc("DELETE /api/v1/product/{id}", handler.DeleteProduct())

	server := &http.Server{
		Addr:    cnf.HttpServer.Addr,
		Handler: router,
	}
	slog.Info("Server is up and running")

	//! gracefully shutdown

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server Failed")
			os.Exit(1)
		}
	}()

	<-done

	slog.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		slog.Info("Failed to shut down")
	} else {
		slog.Info("Server Shutdown ..")
	}

}
