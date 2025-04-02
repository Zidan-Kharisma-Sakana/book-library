package handlers

import (
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/service"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/errs"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) RegisterUserRoutes(r *mux.Router) {
	r.HandleFunc("/profile", routeWrapper(h.GetProfile)).Methods("GET")
	r.HandleFunc("/profile", routeWrapper(h.UpdateProfile)).Methods("PUT")
}

func (h *UserHandler) RegisterLibrarianRoutes(r *mux.Router) {
	r.HandleFunc("/users", routeWrapper(h.ListUsers)).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", routeWrapper(h.GetUser)).Methods("GET")
}

func (h *UserHandler) RegisterAdminRoutes(r *mux.Router) {
	r.HandleFunc("/users/{id:[0-9]+}", routeWrapper(h.UpdateUser)).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", routeWrapper(h.DeleteUser)).Methods("DELETE")
}

func (h *UserHandler) GetUser(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.GetByID(int(id))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *UserHandler) UpdateUser(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	input, err := decodeBody[models.UpdateUserInput](r)
	if err != nil {
		return nil, err
	}

	role, ok := r.Context().Value("role").(string)
	if !ok {
		return nil, errs.NewUnauthorized()
	}

	user, err := h.userService.Update(int(id), role, input)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	if err := h.userService.Delete(int(id)); err != nil {
		return nil, err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func (h *UserHandler) ListUsers(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	query := r.URL.Query()

	page := 1
	pageSize := 10

	if pageStr := query.Get("page"); pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		parsedPageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	users, count, err := h.userService.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(count) / pageSize
	if int(count)%pageSize > 0 {
		totalPages++
	}

	return struct {
		Users      []models.User `json:"users"`
		TotalCount int64         `json:"total_count"`
		Page       int           `json:"page"`
		PageSize   int           `json:"page_size"`
		TotalPages int           `json:"total_pages"`
	}{
		Users:      users,
		TotalCount: count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (h *UserHandler) GetProfile(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		return nil, errs.NewUnauthorized()
	}

	user, err := h.userService.GetByID(int(userID))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *UserHandler) UpdateProfile(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	userID, ok := r.Context().Value("userID").(uint)
	if !ok {
		return nil, errs.NewUnauthorized()
	}
	role, ok := r.Context().Value("role").(string)
	if !ok {
		return nil, errs.NewUnauthorized()
	}
	input, err := decodeBody[models.UpdateUserInput](r)
	if err != nil {
		return nil, err
	}

	input.Role = nil

	user, err := h.userService.Update(int(userID), role, input)
	if err != nil {
		return nil, err
	}

	return user, nil
}
