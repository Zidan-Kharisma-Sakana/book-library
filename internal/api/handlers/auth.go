package handlers

import (
	"net/http"

	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/service"
	"github.com/gorilla/mux"
)

type AuthHandler struct {
	userService *service.UserService
	authService *service.AuthService
}

func NewAuthHandler(userService *service.UserService, authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *AuthHandler) RegisterPublicRoutes(r *mux.Router) {
	r.HandleFunc("/auth/register", routeWrapper(h.Register)).Methods("POST")
	r.HandleFunc("/auth/login", routeWrapper(h.Login)).Methods("POST")
	r.HandleFunc("/auth/refresh", routeWrapper(h.RefreshToken)).Methods("POST")
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	input, err := decodeBody[models.CreateUserInput](r)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.Register(input)
	if err != nil {
		return nil, err
	}

	w.WriteHeader(http.StatusCreated)
	return user, nil
}

func (h *AuthHandler) Login(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	input, err := decodeBody[models.LoginInput](r)
	if err != nil {
		return nil, err
	}

	token, err := h.userService.Login(input)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (h *AuthHandler) RefreshToken(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	input, err := decodeBody[struct {
		RefreshToken string `json:"refresh_token"`
	}](r)
	if err != nil {
		return nil, err
	}

	userID, _, err := h.authService.ValidateToken(input.RefreshToken)
	if err != nil {
		return nil, err
	}

	token, err := h.userService.RefreshToken(userID)
	if err != nil {
		return nil, err
	}
	return token, nil
}
