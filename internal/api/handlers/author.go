package handlers

import (
	"net/http"
	"strconv"

	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/service"
	"github.com/gorilla/mux"
)

type AuthorHandler struct {
	authorService *service.AuthorService
}

func NewAuthorHandler(authorService *service.AuthorService) *AuthorHandler {
	return &AuthorHandler{
		authorService: authorService,
	}
}

func (h *AuthorHandler) RegisterPublicRoutes(r *mux.Router) {
	r.HandleFunc("/authors", routeWrapper(h.ListAuthors)).Methods("GET")
	r.HandleFunc("/authors/{id:[0-9]+}", routeWrapper(h.GetAuthor)).Methods("GET")
	r.HandleFunc("/authors/{id:[0-9]+}/with-books", routeWrapper(h.GetAuthorWithBooks)).Methods("GET")
}

func (h *AuthorHandler) RegisterLibrarianRoutes(r *mux.Router) {
	r.HandleFunc("/authors", routeWrapper(h.CreateAuthor)).Methods("POST")
	r.HandleFunc("/authors/{id:[0-9]+}", routeWrapper(h.UpdateAuthor)).Methods("PUT")
	r.HandleFunc("/authors/{id:[0-9]+}", routeWrapper(h.DeleteAuthor)).Methods("DELETE")
}

func (h *AuthorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	input, err := decodeBody[models.CreateAuthorInput](r)
	if err != nil {
		return nil, err
	}

	author, err := h.authorService.Create(input)
	if err != nil {
		return nil, err
	}

	w.WriteHeader(http.StatusCreated)
	return author, nil
}

func (h *AuthorHandler) GetAuthor(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	author, err := h.authorService.GetByID(int(id))
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (h *AuthorHandler) GetAuthorWithBooks(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	author, err := h.authorService.GetWithBooks(int(id))
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (h *AuthorHandler) UpdateAuthor(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	input, err := decodeBody[models.UpdateAuthorInput](r)
	if err != nil {
		return nil, err
	}

	author, err := h.authorService.Update(int(id), input)
	if err != nil {
		return nil, err
	}

	return author, nil
}

func (h *AuthorHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	if err := h.authorService.Delete(int(id)); err != nil {
		return nil, err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func (h *AuthorHandler) ListAuthors(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	query := r.URL.Query()

	filter := models.AuthorFilter{
		Name:        query.Get("name"),
		Nationality: query.Get("nationality"),
		Page:        1,
		PageSize:    10,
	}

	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err == nil && page > 0 {
			filter.Page = page
		}
	}

	if pageSizeStr := query.Get("page_size"); pageSizeStr != "" {
		pageSize, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSize > 0 {
			filter.PageSize = pageSize
		}
	}

	authors, count, err := h.authorService.List(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(count) / filter.PageSize
	if int(count)%filter.PageSize > 0 {
		totalPages++
	}

	return struct {
		Authors    []models.Author `json:"authors"`
		TotalCount int64           `json:"total_count"`
		Page       int             `json:"page"`
		PageSize   int             `json:"page_size"`
		TotalPages int             `json:"total_pages"`
	}{
		Authors:    authors,
		TotalCount: count,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}
