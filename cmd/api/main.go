package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/isa0-gh/short-urls/internal/db"
	"github.com/isa0-gh/short-urls/internal/handler"
	"github.com/isa0-gh/short-urls/internal/middleware"
	"github.com/isa0-gh/short-urls/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// @title       Shortlink API
// @version     1.0
// @description API for managing users and shortlinks.
// @host        localhost:8080
// @BasePath    /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter a bearer token in the format: Bearer {token}.

var (
	Port       string
	DB         *pgxpool.Pool
	SessionExp time.Duration
)

func main() {
	godotenv.Load()

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080"
	}
	expStr := os.Getenv("SESSION_EXP")
	if expStr == "" {
		expStr = "1h"
	}
	var err error
	SessionExp, err = time.ParseDuration(expStr)
	if err != nil {
		log.Fatalf("Invalid SESSION_EXP: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)
	fmt.Println(dsn)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DB config: %v", err)
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	if err = DB.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}
	queries := db.New(DB)

	defer DB.Close()

	// services
	authSrv := service.NewAuthService(queries, SessionExp)
	userSrv := service.NewUserService(queries)
	linkSrv := service.NewShortLinkService(queries, nil)

	// handlers
	authHandler := handler.NewAuthHandler(userSrv, authSrv)
	userHandler := handler.NewUserHandler(userSrv)
	linkHandler := handler.NewShortLinkHandler(linkSrv)

	// routing
	apiMux := http.NewServeMux()

	// public routes
	apiMux.HandleFunc("POST /login", authHandler.LoginUser)
	apiMux.HandleFunc("POST /register", authHandler.CreateUser)

	authMiddleware := middleware.AuthMiddleware(authSrv)
	protected := func(h http.HandlerFunc) http.Handler {
		return authMiddleware(h)
	}

	// auth protected routes
	apiMux.Handle("GET /profile", protected(userHandler.GetUserDetails))

	apiMux.Handle("POST /create", protected(linkHandler.CreateShortlink))
	apiMux.Handle("DELETE /link/{shortCode}", protected(linkHandler.DeleteShortLink))
	apiMux.Handle("POST /link/{shortCode}", protected(linkHandler.UpdateShortLink))
	apiMux.Handle("GET /links", protected(linkHandler.ListShortLinks))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /s/{shortCode}", linkHandler.RedirectToDestination)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	http.ListenAndServe(":"+Port, mux)
}
