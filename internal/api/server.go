package api

import (
	"context"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/api/handlers"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/api/middleware"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/service"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/config"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"golang.org/x/time/rate"
)

type Server struct {
	httpServer    *http.Server
	router        *mux.Router
	authService   *service.AuthService
	bookService   *service.BookService
	authorService *service.AuthorService
	userService   *service.UserService
}

// New creates a new httpServer
func NewServer(
	cfg *config.Config,
	authService *service.AuthService,
	bookService *service.BookService,
	authorService *service.AuthorService,
	userService *service.UserService,
) *Server {
	router := mux.NewRouter()
	// Apply CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	var handler http.Handler = router
	handler = c.Handler(handler)
	handler = middleware.Timeout(10 * time.Second)(handler)

	if cfg.RateLimit.Enabled {
		limiter := rate.NewLimiter(cfg.RateLimit.Limit, cfg.RateLimit.Burst)
		handler = middleware.RateLimiter(limiter)(handler)
	}

	server := &Server{
		httpServer: &http.Server{
			Addr:         cfg.ServerAddress,
			Handler:      handler,
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
		},
		router:        router,
		authService:   authService,
		bookService:   bookService,
		authorService: authorService,
		userService:   userService,
	}

	server.setupRoutes()

	return server
}

func (s *Server) setupRoutes() {
	// main router
	api := s.router.PathPrefix("/api/v1").Subrouter()

	api.Use(middleware.Logger())
	api.Use(middleware.Recoverer())
	api.Use(middleware.Timeout(10 * time.Second))

	bookHandler := handlers.NewBookHandler(s.bookService)
	authorHandler := handlers.NewAuthorHandler(s.authorService)
	userHandler := handlers.NewUserHandler(s.userService)
	authHandler := handlers.NewAuthHandler(s.userService, s.authService)

	bookHandler.RegisterPublicRoutes(api)
	authorHandler.RegisterPublicRoutes(api)
	authHandler.RegisterPublicRoutes(api)

	// This will be private route that need authentication
	authRoutes := api.PathPrefix("/").Subrouter()
	authRoutes.Use(middleware.Authentication(s.authService))
	userHandler.RegisterUserRoutes(authRoutes)

	librarianRoutes := authRoutes.PathPrefix("").Subrouter()
	librarianRoutes.Use(middleware.FilterRole("librarian", "admin"))
	bookHandler.RegisterLibrarianRoutes(librarianRoutes)
	authorHandler.RegisterLibrarianRoutes(librarianRoutes)
	userHandler.RegisterLibrarianRoutes(librarianRoutes)

	adminRoutes := authRoutes.PathPrefix("").Subrouter()
	adminRoutes.Use(middleware.FilterRole("admin"))
	userHandler.RegisterAdminRoutes(adminRoutes)

	s.router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
