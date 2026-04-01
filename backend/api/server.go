package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rs/zerolog/log"
)

type Server struct {
	*http.Server
	startupTime time.Time
}

func NewServer(database database.Database) (Server, error) {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	address := "0.0.0.0:" + port // Bind to 0.0.0.0 for external access

	// Capture startup time
	startupTime := time.Now()

	router := newRouter(database, withStartupTime(startupTime))

	// Hardcoded timeout values
	readTimeout := 180 * time.Second
	writeTimeout := 180 * time.Second
	idleTimeout := 180 * time.Second

	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return Server{server, startupTime}, nil
}

type router struct {
	startupTime time.Time
}

func withStartupTime(startupTime time.Time) func(*router) {
	return func(r *router) {
		r.startupTime = startupTime
	}
}

func newRouter(database database.Database, opts ...func(*router)) *chi.Mux {
	var router router
	for _, opt := range opts {
		opt(&router)
	}

	chiRouter := chi.NewRouter()
	chiRouter.Use(LogInternalServerErrors)

	// Apply CORS middleware (must be before routes)
	acceptedOrigins := strings.Split(os.Getenv("ACCEPTED_ORIGINS"), ",")
	chiRouter.Use(CORSCheckMiddleware(acceptedOrigins))
	chiRouter.Use(corsMiddleware(acceptedOrigins))

	// Root endpoint - Message of the Day
	chiRouter.Get("/", rootHandler(router.startupTime))

	// Healthcheck endpoint - accessible from any origin
	chiRouter.Get("/healthcheck", healthcheckHandler(router.startupTime))

	// Initialize WebSocket Hub
	hub := NewHub()
	go hub.Run()

	// Initialize all handlers
	handlers := initializeHandlers(database, hub)

	// Ensure at least one user exists (from AUTH_EMAIL/AUTH_PASSWORD if table empty)
	EnsureFirstUserExists(database.UserRepo())

	// Initialize auth middleware (validates Bearer token against user.Token in DB)
	authMiddleware := newAuthMiddleware(database.UserRepo())

	// Setup all route types
	setupFrontendRoutes(chiRouter, handlers, authMiddleware, hub)

	return chiRouter
}

func (s Server) Start(errChannel chan<- error) {
	log.Info().Msgf("Server started on: %s", s.Addr)
	errChannel <- s.ListenAndServe()
}

func (s Server) ShutdownGracefully(timeout time.Duration) {
	log.Info().Msg("Gracefully shutting down...")

	gracefullCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.Shutdown(gracefullCtx); err != nil {
		log.Error().Msgf("Error shutting down the server: %v", err)
	} else {
		log.Info().Msg("HttpServer gracefully shut down")
	}
}

// formatUptime returns a human-readable uptime string
func formatUptime(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// rootHandler returns a handler function for the root endpoint
func rootHandler(startupTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "Faradhaven")

		quotes := []string{
			"The forge awaits your spells.",
			"Beasts stir in the shadows.",
			"Your spellbook is ready.",
			"Adventure calls from distant lands.",
		}

		rand.Seed(time.Now().UnixNano())
		selectedQuote := quotes[rand.Intn(len(quotes))]

		uptime := time.Since(startupTime)
		uptimeStr := formatUptime(uptime)

		userAgent := r.Header.Get("User-Agent")
		if !strings.Contains(userAgent, "curl") {
			w.Header().Set("Content-Type", "text/html")
			html := fmt.Sprintf(`
        <html>
        <body style="background:#1e1e2e; color:#cdd6f4; font-family: monospace; display:flex; align-items:center; justify-content:center; height:100vh;">
            <div style="border: 1px solid #fab387; padding: 20px; border-radius: 5px; min-width: 400px;">
                <p style="color:#fab387;">faradhaven:~# ./status</p>
                <p style="color:#a6e3a1;">● SYSTEM STATUS: ONLINE</p>
                <p style="margin: 4px 0; color:#89b4fa;">Started: %s</p>
                <p style="margin: 4px 0; color:#89b4fa;">Uptime:  %s</p>
                <p style="margin: 4px 0; color:#89b4fa;">Runtime: Go %s</p>
                <hr style="border-color:#45475a; margin: 12px 0;">
                <p style="color:#f9e2af; font-style:italic;">"%s"</p>
                <span style="animation: blink 1s infinite;">_</span>
            </div>
            <style>@keyframes blink{50%%{opacity:0;}}</style>
        </body>
        </html>`,
				startupTime.Format("2006-01-02 15:04:05 MST"),
				uptimeStr,
				runtime.Version()[2:], // Strip "go" prefix
				selectedQuote)
			w.Write([]byte(html))
			return
		}

		// CLI/curl output
		cliOutput := fmt.Sprintf(`SYSTEM STATUS: ONLINE
Started: %s
Uptime:  %s
Runtime: Go %s

"%s"
`,
			startupTime.Format("2006-01-02 15:04:05 MST"),
			uptimeStr,
			runtime.Version()[2:],
			selectedQuote)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(cliOutput))
	}
}

// healthcheckHandler returns a handler function for the healthcheck endpoint
func healthcheckHandler(startupTime time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Get memory stats
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		uptime := time.Since(startupTime)

		response := map[string]interface{}{
			"status":       "healthy",
			"current_time": time.Now().Format(time.RFC3339),
			"startup_time": startupTime.Format(time.RFC3339),
			"uptime":       formatUptime(uptime),
			"uptime_seconds": int(uptime.Seconds()),
			"runtime": map[string]interface{}{
				"go_version":   runtime.Version(),
				"os":           runtime.GOOS,
				"arch":         runtime.GOARCH,
				"num_cpu":      runtime.NumCPU(),
				"num_goroutine": runtime.NumGoroutine(),
			},
			"memory": map[string]interface{}{
				"alloc_mb":       float64(memStats.Alloc) / 1024 / 1024,
				"total_alloc_mb": float64(memStats.TotalAlloc) / 1024 / 1024,
				"sys_mb":         float64(memStats.Sys) / 1024 / 1024,
				"num_gc":         memStats.NumGC,
			},
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error().Err(err).Msg("Error encoding healthcheck response")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
