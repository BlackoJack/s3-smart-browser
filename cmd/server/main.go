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

    // Generated screenshots referenced by manifest.json
    mux.HandleFunc("/static/images/screenshot-desktop.png", serveGeneratedScreenshot(1280, 720))
    mux.HandleFunc("/static/images/screenshot-mobile.png", serveGeneratedScreenshot(390, 844))

    // Static files
	mux.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("web/static"))))

    // Serve service worker at root scope
    mux.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
        // Ensure the service worker is always fresh and has root scope
        w.Header().Set("Content-Type", "application/javascript")
        w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
        w.Header().Set("Pragma", "no-cache")
        w.Header().Set("Expires", "0")
        w.Header().Set("Service-Worker-Allowed", "/")
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

// serveGeneratedScreenshot dynamically creates a placeholder PNG screenshot of the given size.
func serveGeneratedScreenshot(width, height int) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        img := image.NewRGBA(image.Rect(0, 0, width, height))

        // Background color lighter variant of theme
        background := color.RGBA{R: 0xF3, G: 0xF4, B: 0xF6, A: 0xFF} // #F3F4F6
        draw.Draw(img, img.Bounds(), &image.Uniform{C: background}, image.Point{}, draw.Src)

        // Add a simple header bar using the theme color
        headerHeight := height / 10
        headerRect := image.Rect(0, 0, width, headerHeight)
        theme := color.RGBA{R: 0x66, G: 0x7e, B: 0xea, A: 0xff} // #667eea
        draw.Draw(img, headerRect, &image.Uniform{C: theme}, image.Point{}, draw.Src)

        // Add a simple content area card
        cardMargin := width / 12
        cardTop := headerHeight + (height / 20)
        cardRect := image.Rect(cardMargin, cardTop, width-cardMargin, height-(height/12))
        card := color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
        draw.Draw(img, cardRect, &image.Uniform{C: card}, image.Point{}, draw.Src)

        w.Header().Set("Content-Type", "image/png")
        if err := png.Encode(w, img); err != nil {
            http.Error(w, "failed to generate screenshot", http.StatusInternalServerError)
            return
        }
    }
}
