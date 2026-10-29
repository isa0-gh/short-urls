package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/isa0-gh/short-urls/docs"
	"github.com/isa0-gh/short-urls/internal/db"
	"github.com/isa0-gh/short-urls/internal/handler"
	"github.com/isa0-gh/short-urls/internal/middleware"
	"github.com/isa0-gh/short-urls/internal/service"
	webui "github.com/isa0-gh/short-urls/internal/web"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title       Shortlink API
// @version     1.0
// @description API for managing users and shortlinks.
// @BasePath    /
// @securityDefinitions.apikey SessionCookie
// @in cookie
// @name session
// @description HTTP-only session cookie.

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
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:         redisAddr,
		DialTimeout:  100 * time.Millisecond,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
	})
	defer redisClient.Close()

	// services
	authSrv := service.NewAuthService(queries, redisClient, SessionExp)
	userSrv := service.NewUserService(queries, redisClient)
	linkSrv := service.NewShortLinkService(queries, redisClient)

	// handlers
	authHandler := handler.NewAuthHandler(userSrv, authSrv)
	userHandler := handler.NewUserHandler(userSrv)
	linkHandler := handler.NewShortLinkHandler(linkSrv)

	// routing
	apiMux := http.NewServeMux()

	// public routes
	apiMux.HandleFunc("POST /login", authHandler.LoginUser)
	apiMux.HandleFunc("POST /logout", authHandler.LogoutUser)
	apiMux.HandleFunc("POST /register", authHandler.CreateUser)

	authMiddleware := middleware.AuthMiddleware(authSrv)
	protected := func(h http.HandlerFunc) http.Handler {
		return authMiddleware(h)
	}

	// auth protected routes
	apiMux.Handle("GET /profile", protected(userHandler.GetUserDetails))

	apiMux.Handle("POST /create", protected(linkHandler.CreateShortlink))
	apiMux.Handle("DELETE /link/{linkID}", protected(linkHandler.DeleteShortLink))
	apiMux.Handle("POST /link/{linkID}", protected(linkHandler.UpdateShortLink))
	apiMux.Handle("GET /links", protected(linkHandler.ListShortLinks))

	mux := http.NewServeMux()
	webui.RegisterRoutes(mux)
	mux.HandleFunc("GET /s/{shortCode}", linkHandler.RedirectToDestination)
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", apiMux))
	mux.Handle("GET /swagger/", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))

	logOutput := os.Stdout
	if logPath := os.Getenv("LOG_FILE"); logPath != "" {
		logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}
		defer logFile.Close()
		logOutput = logFile
	}

	logger := slog.New(slog.NewJSONHandler(logOutput, nil))
	loggedMux := middleware.LoggingMiddleware(logger)(mux)
	http.ListenAndServe(":"+Port, loggedMux)
}
