package service

import (
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
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
	author, err := s.authorRepo.GetByID(input.AuthorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("author not found")
	}

	existingBook, err := s.bookRepo.GetByISBN(input.ISBN)
	if err != nil {
		return nil, err
	}
	if existingBook != nil {
		return nil, errors.New("book with this ISBN already exists")
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

	if err := s.validator.Struct(book); err != nil {
		return nil, err
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
	if book == nil {
		return nil, errors.New("book not found")
	}
	return book, nil
}

func (s *BookService) Update(id int, input models.UpdateBookInput) (*models.Book, error) {
	book, err := s.bookRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	if input.Title != nil {
		book.Title = *input.Title
	}
	if input.ISBN != nil {
		existingBook, err := s.bookRepo.GetByISBN(*input.ISBN)
		if err != nil {
			return nil, err
		}
		if existingBook != nil && int(existingBook.ID) != id {
			return nil, errors.New("book with this ISBN already exists")
		}
		book.ISBN = *input.ISBN
	}
	if input.Description != nil {
		book.Description = *input.Description
	}
	if input.AuthorID != nil {
		// Check if author exists
		author, err := s.authorRepo.GetByID(int(*input.AuthorID))
		if err != nil {
			return nil, err
		}
		if author == nil {
			return nil, errors.New("author not found")
		}
		book.AuthorID = *input.AuthorID
	}
	if input.Publisher != nil {
		book.Publisher = *input.Publisher
	}
	if input.PublishedAt != nil {
		book.PublishedAt = *input.PublishedAt
	}
	if input.Pages != nil {
		book.Pages = *input.Pages
	}
	if input.Available != nil {
		book.Available = *input.Available
	}

	if err := s.validator.Struct(book); err != nil {
		return nil, err
	}

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
		return errors.New("book not found")
	}

	return s.bookRepo.Delete(id)
}

func (s *BookService) List(filter models.BookFilter) ([]models.Book, int64, error) {
	return s.bookRepo.List(filter)
}

func (s *BookService) GetBooksByAuthorID(authorID int) ([]models.Book, error) {
	author, err := s.authorRepo.GetByID(authorID)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("author not found")
	}

	return s.bookRepo.GetBooksByAuthorID(authorID)
}
