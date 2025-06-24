package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
	"github.com/slangeres/Vypaar/backend_API/internal/config"
	"github.com/slangeres/Vypaar/backend_API/internal/https/handler"
	"github.com/slangeres/Vypaar/backend_API/internal/https/middleware"
	"github.com/slangeres/Vypaar/backend_API/internal/storage/sqllite"
	"github.com/slangeres/Vypaar/backend_API/internal/token"
)

func main() {

	//!Configuration
	cnf := config.MustDone()

	//!Configure the DB
	db, err := sqllite.ConfigSQL(cnf)

	if err != nil {
		slog.Warn("databse connection issue .....")

	}

	//JWT Maker
	jwtMaker := token.NewJwtMaker(cnf.JwtSecrateKey)

	userDb, err := sqllite.InitUserDb(cnf)

	if err != nil {
		slog.Warn("user db connection issue")
	}

	//^Making the router
	router := http.NewServeMux()

	//!hander functions

	router.HandleFunc("GET /api/v1/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from vypaar api"))
	})

	//! Authenticate handler

	router.HandleFunc("POST /api/v1/signup", handler.SignupUser(userDb))

	router.HandleFunc("POST /api/v1/login", handler.LoginUser(userDb, jwtMaker))

	router.Handle("POST /api/v1/products", middleware.AuthMiddleware(jwtMaker)(handler.PostProduct(db)))

	router.Handle("GET /api/v1/products", middleware.AuthMiddleware(jwtMaker)(handler.GetProduct(db)))

	router.Handle("GET /api/v1/products/{id}", middleware.AuthMiddleware(jwtMaker)(handler.GetProductById(db)))

	router.Handle("PATCH /api/v1/products/{id}", middleware.AuthMiddleware(jwtMaker)(handler.UpdateProduct(db)))

	router.Handle("DELETE /api/v1/products/{id}", middleware.AuthMiddleware(jwtMaker)(handler.DeleteProduct(db)))

	// ! Adding cors to the api
	corsReq := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //!for the testing purpose i put cors origin * ...... Later change
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-type", "Authorization"},
		AllowCredentials: true,
	})
	//! Wrapping cors header in the handler ..............
	requestHandler := corsReq.Handler(router)

	server := &http.Server{
		Addr:    cnf.HttpServer.Addr,
		Handler: requestHandler,
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
