package service

import (
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/errs"
	"github.com/go-playground/validator/v10"
)

type BookService struct {
	bookRepo   interfaces.BookRepository
	authorRepo interfaces.AuthorRepository
	validator  *validator.Validate
}

func NewBookService(validator *validator.Validate, bookRepo interfaces.BookRepository, authorRepo interfaces.AuthorRepository) *BookService {
	return &BookService{
		validator:  validator,
		bookRepo:   bookRepo,
		authorRepo: authorRepo,
	}
}

func (s *BookService) Create(input models.CreateBookInput) (*models.Book, error) {
	if err := s.validator.Struct(input); err != nil {
		return nil, errs.NewBadRequestError().SetError(err)
	}

	book := &models.Book{
		Title:       input.Title,
		ISBN:        input.ISBN,
		Description: input.Description,
		AuthorID:    input.AuthorID,
		Publisher:   input.Publisher,
		PublishedAt: input.PublishedAt,
		Pages:       input.Pages,
		Available:   input.Available,
	}

	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) GetByID(id int) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookService) Update(id int, input models.UpdateBookInput) (*models.Book, error) {
	if err := s.validator.Struct(input); err != nil {
		return nil, errs.NewBadRequestError().SetError(err)
	}

	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errs.NewNotFoundError()
	}

	book.AuthorID = *input.AuthorID
	book.ISBN = *input.ISBN
	book.Title = *input.Title
	book.Description = *input.Description
	book.Publisher = *input.Publisher
	book.PublishedAt = *input.PublishedAt
	book.Pages = *input.Pages
	book.Available = *input.Available

	if err := s.bookRepo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookService) Delete(id int) error {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return err
	}
	if book == nil {
		return errs.NewNotFoundError()
	}

	return s.bookRepo.Delete(id)
}

func (s *BookService) List(filter models.BookFilter) ([]models.Book, int64, error) {
	return s.bookRepo.List(filter)
}
