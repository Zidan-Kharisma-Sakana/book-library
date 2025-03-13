package api

import (
	"context"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/api/handlers"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/api/middleware"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/service"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	addr string,
	authService *service.AuthService,
	bookService *service.BookService,
	authorService *service.AuthorService,
	userService *service.UserService,
) *Server {
	router := mux.NewRouter()

	server := &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      router,
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
	// Create API router
	api := s.router.PathPrefix("/api/v1").Subrouter()

	api.Use(middleware.Logger())
	api.Use(middleware.Recoverer())

	bookHandler := handlers.NewBookHandler(s.bookService)
	authorHandler := handlers.NewAuthorHandler(s.authorService)
	userHandler := handlers.NewUserHandler(s.userService)
	authHandler := handlers.NewAuthHandler(s.userService, s.authService)

	bookHandler.RegisterPublicRoutes(api)
	authorHandler.RegisterPublicRoutes(api)
	authHandler.RegisterPublicRoutes(api)

	authRoutes := api.PathPrefix("/").Subrouter()
	authRoutes.Use(middleware.Auth(s.authService))
	userHandler.RegisterUserRoutes(authRoutes)

	librarianRoutes := authRoutes.PathPrefix("/librarians").Subrouter()
	librarianRoutes.Use(middleware.LibrarianOnly())
	bookHandler.RegisterLibrarianRoutes(librarianRoutes)
	authorHandler.RegisterLibrarianRoutes(librarianRoutes)
	userHandler.RegisterLibrarianRoutes(librarianRoutes)

	adminRoutes := authRoutes.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(middleware.AdminOnly())
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
