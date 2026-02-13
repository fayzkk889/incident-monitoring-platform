package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"Incident_Monitoring_Project/internal/store"
)

func main() {
	_ = godotenv.Load()

	dbURL := getenv("DATABASE_URL", "postgres://incident:incidentpassword@localhost:5432/incidentdb?sslmode=disable")
	mlServiceURL := getenv("ML_SERVICE_URL", "http://localhost:8000")

	ctx := context.Background()
	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbpool.Close()

	if err := store.RunMigrations(ctx, dbpool); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	repo := store.NewRepository(dbpool)

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	handler := NewHandler(repo, mlServiceURL)

	e.POST("/api/logs", handler.IngestLogs)
	e.GET("/api/health", handler.Health)
	e.GET("/api/incidents", handler.ListIncidents)
	e.GET("/api/summary/:incident_id", handler.GetIncidentSummary)

	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}

	srv := &http.Server{
		Addr:         addr,
		Handler:      e,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("Go API listening on %s (ML service: %s)", addr, mlServiceURL)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

