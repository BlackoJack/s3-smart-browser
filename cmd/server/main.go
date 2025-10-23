package main

import (
    "image"
    "image/color"
    "image/draw"
    "image/png"
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

    // Generated icons (ensure PWA required icons exist)
    mux.HandleFunc("/static/images/icon-192x192.png", serveGeneratedIcon(192))
    mux.HandleFunc("/static/images/icon-512x512.png", serveGeneratedIcon(512))
    mux.HandleFunc("/static/images/favicon-32x32.png", serveGeneratedIcon(32))
    mux.HandleFunc("/static/images/favicon-16x16.png", serveGeneratedIcon(16))
    mux.HandleFunc("/static/images/apple-touch-icon.png", serveGeneratedIcon(180))

    // Static files
	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))

    // Serve service worker at root scope
    mux.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/javascript")
        http.ServeFile(w, r, "web/static/js/sw.js")
    })

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

// serveGeneratedIcon dynamically creates a solid-color PNG of the given size.
func serveGeneratedIcon(size int) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        img := image.NewRGBA(image.Rect(0, 0, size, size))
        // Theme color #667eea
        theme := color.RGBA{R: 0x66, G: 0x7e, B: 0xea, A: 0xff}
        draw.Draw(img, img.Bounds(), &image.Uniform{C: theme}, image.Point{}, draw.Src)

        w.Header().Set("Content-Type", "image/png")
        if err := png.Encode(w, img); err != nil {
            http.Error(w, "failed to generate icon", http.StatusInternalServerError)
            return
        }
    }
}
