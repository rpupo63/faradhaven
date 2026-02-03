package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rpupo63/unified-personal-site-backend/database"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_classes"
	"github.com/rpupo63/unified-personal-site-backend/seed/faradhaven_races"
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
	chiRouter.Get("/", rootHandler())

	// Healthcheck endpoint - accessible from any origin
	chiRouter.Get("/healthcheck", healthcheckHandler(router.startupTime))

	// Initialize all handlers
	handlers := initializeHandlers(database)

	// Ensure at least one user exists (from AUTH_EMAIL/AUTH_PASSWORD if table empty)
	EnsureFirstUserExists(database.UserRepo())
	// Seed Faradhaven races (Race, Trait, TraitOption) with full trait data
	if err := faradhaven_races.SeedFaradhavenRaces(database.DB()); err != nil {
		log.Warn().Err(err).Msg("Seed Faradhaven races skipped or failed")
	}
	// Seed Faradhaven classes (Class, ClassLevel 1-20, ClassComponent) if tables are empty
	if err := faradhaven_classes.SeedFaradhavenClassesIfEmpty(database.DB()); err != nil {
		log.Warn().Err(err).Msg("Seed Faradhaven classes skipped or failed")
	}
	// Backfill AbilityScoreImprovement for existing ClassLevel rows (levels 4, 8, 12, 16, 19)
	if err := faradhaven_classes.BackfillAbilityScoreImprovements(database.DB()); err != nil {
		log.Warn().Err(err).Msg("Backfill ability score improvements skipped or failed")
	}

	// Initialize auth middleware (validates Bearer token against user.Token in DB)
	authMiddleware := newAuthMiddleware(database.UserRepo())

	// Setup all route types
	setupFrontendRoutes(chiRouter, handlers, authMiddleware)

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

// rootHandler returns a handler function for the root endpoint
func rootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Powered-By", "Faradhaven")

		quotes := []string{
			"SYSTEM STATUS: ONLINE.\nWelcome to Faradhaven.",
			"The forge awaits your spells.",
			"Beasts stir in the shadows.",
			"Your spellbook is ready.",
		}

		rand.Seed(time.Now().UnixNano())
		selectedQuote := quotes[rand.Intn(len(quotes))]

		userAgent := r.Header.Get("User-Agent")
		if !strings.Contains(userAgent, "curl") {
			w.Header().Set("Content-Type", "text/html")
			html := fmt.Sprintf(`
        <html>
        <body style="background:#1e1e2e; color:#cdd6f4; font-family: monospace; display:flex; align-items:center; justify-content:center; height:100vh;">
            <div style="border: 1px solid #fab387; padding: 20px; border-radius: 5px;">
                <p style="color:#fab387;">faradhaven:~# ./status</p>
                <p>%s</p>
                <span style="animation: blink 1s infinite;">_</span>
            </div>
            <style>@keyframes blink{50%%{opacity:0;}}</style>
        </body>
        </html>`, selectedQuote)
			w.Write([]byte(html))
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(selectedQuote + "\n"))
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

		response := map[string]interface{}{
			"current_time":   time.Now().Format(time.RFC3339),
			"startup_time":   startupTime.Format(time.RFC3339),
			"uptime_seconds": int(time.Since(startupTime).Seconds()),
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error().Err(err).Msg("Error encoding healthcheck response")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
