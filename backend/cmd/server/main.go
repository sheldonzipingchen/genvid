package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/genvid/backend/internal/config"
	"github.com/genvid/backend/internal/handler"
	"github.com/genvid/backend/internal/middleware"
	"github.com/genvid/backend/internal/repository"
	"github.com/genvid/backend/internal/service"
	"github.com/genvid/backend/internal/zhipu"
	"github.com/genvid/backend/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := sqlx.Connect("postgres", cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	jwtService := auth.NewJWTService(cfg.JWT)

	zhipuClient := zhipu.NewClient(cfg.External.Zhipu.APIKey)

	profileRepo := repository.NewProfileRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	authService := service.NewAuthService(profileRepo, jwtService, cfg)
	projectService := service.NewProjectService(projectRepo, profileRepo, authService, zhipuClient, cfg)

	authHandler := handler.NewAuthHandler(authService)
	projectHandler := handler.NewProjectHandler(projectService)
	avatarHandler := handler.NewAvatarHandler()
	paymentHandler := handler.NewPaymentHandler(cfg)
	uploadHandler := handler.NewUploadHandler("./uploads", cfg.Server.AppURL)

	r := chi.NewRouter()

	allowedOrigins := []string{cfg.Server.AppURL}
	if cfg.IsDevelopment() {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000", "http://127.0.0.1:3000")
	}

	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.CORSMiddleware(allowedOrigins))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	r.Route("/api", func(r chi.Router) {
		r.Post("/auth/register", authHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.RefreshToken)
		r.Post("/payments/webhook", paymentHandler.HandleWebhook)

		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtService))

			r.Get("/user/profile", authHandler.GetProfile)
			r.Patch("/user/profile", authHandler.UpdateProfile)

			r.Get("/projects", projectHandler.List)
			r.Post("/projects", projectHandler.Create)
			r.Get("/projects/{id}", projectHandler.GetByID)
			r.Delete("/projects/{id}", projectHandler.Delete)
			r.Post("/projects/{id}/generate", projectHandler.GenerateVideo)

			r.Get("/avatars", avatarHandler.List)
			r.Get("/avatars/{id}", avatarHandler.GetByID)

			r.Post("/upload", uploadHandler.Upload)
			r.Delete("/upload", uploadHandler.Delete)

			r.Post("/payments/checkout", paymentHandler.CreateCheckoutSession)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}
