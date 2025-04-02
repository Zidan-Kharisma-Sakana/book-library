package handlers

import (
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BookHandler struct {
	bookService *service.BookService
}

func NewBookHandler(bookService *service.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) RegisterPublicRoutes(r *mux.Router) {
	r.HandleFunc("/books", routeWrapper(h.ListBooks)).Methods("GET")
	r.HandleFunc("/books/{id:[0-9]+}", routeWrapper(h.GetBook)).Methods("GET")
}

func (h *BookHandler) RegisterLibrarianRoutes(r *mux.Router) {
	r.HandleFunc("/books", routeWrapper(h.CreateBook)).Methods("POST")
	r.HandleFunc("/books/{id:[0-9]+}", routeWrapper(h.UpdateBook)).Methods("PUT")
	r.HandleFunc("/books/{id:[0-9]+}", routeWrapper(h.DeleteBook)).Methods("DELETE")
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	input, err := decodeBody[models.CreateBookInput](r)
	if err != nil {
		return nil, err
	}

	book, err := h.bookService.Create(input)
	if err != nil {
		return nil, err
	}

	w.WriteHeader(http.StatusCreated)
	return book, nil
}

func (h *BookHandler) GetBook(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	book, err := h.bookService.GetByID(int(id))
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (h *BookHandler) UpdateBook(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	input, err := decodeBody[models.UpdateBookInput](r)
	if err != nil {
		return nil, err
	}

	book, err := h.bookService.Update(int(id), input)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return nil, err
	}

	if err := h.bookService.Delete(int(id)); err != nil {
		return nil, err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil, nil
}

func (h *BookHandler) ListBooks(_ http.ResponseWriter, r *http.Request) (interface{}, error) {
	query := r.URL.Query()

	filter := models.BookFilter{
		Title:     query.Get("title"),
		Publisher: query.Get("publisher"),
		Page:      1,
		PageSize:  10,
	}

	if authorIDStr := query.Get("author_id"); authorIDStr != "" {
		authorID, err := strconv.ParseInt(authorIDStr, 10, 32)
		if err == nil {
			filter.AuthorID = int(authorID)
		}
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

	if availableStr := query.Get("available"); availableStr != "" {
		available := availableStr == "true"
		filter.Available = &available
	}

	books, count, err := h.bookService.List(filter)
	if err != nil {
		return nil, err
	}

	totalPages := int(count) / filter.PageSize
	if int(count)%filter.PageSize > 0 {
		totalPages++
	}

	return struct {
		Books      []models.Book `json:"books"`
		TotalCount int64         `json:"total_count"`
		Page       int           `json:"page"`
		PageSize   int           `json:"page_size"`
		TotalPages int           `json:"total_pages"`
	}{
		Books:      books,
		TotalCount: count,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}
