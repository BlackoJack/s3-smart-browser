package main

import (
	"log"
	"net/http"

	"s3-smart-browser/internal/config"
	"s3-smart-browser/internal/handlers"
	"s3-smart-browser/internal/s3"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация S3 клиента
	s3Client, err := s3.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	// Инициализация обработчиков
	h := handlers.NewHandlers(s3Client)

	// Настройка маршрутов
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/list", h.ListDirectory)
	mux.HandleFunc("/api/download", h.DownloadFile)
	mux.HandleFunc("/api/version", h.GetVersion)

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))

	// UI
	mux.HandleFunc("/", h.ServeUI)

	// Запуск сервера
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	log.Fatal(server.ListenAndServe())
}
