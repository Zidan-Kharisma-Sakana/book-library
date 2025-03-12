package service

import (
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
	"github.com/go-playground/validator/v10"
)

type AuthorService struct {
	authorRepo interfaces.AuthorRepository
	validator  *validator.Validate
}

func NewAuthorService(validator *validator.Validate, authorRepo interfaces.AuthorRepository) *AuthorService {
	return &AuthorService{
		validator:  validator,
		authorRepo: authorRepo,
	}
}

func (s *AuthorService) Create(input models.CreateAuthorInput) (*models.Author, error) {
	author := &models.Author{
		Name:        input.Name,
		Biography:   input.Biography,
		BirthDate:   input.BirthDate,
		Nationality: input.Nationality,
	}

	if err := s.validator.Struct(author); err != nil {
		return nil, err
	}

	if err := s.authorRepo.Create(author); err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorService) GetByID(id int) (*models.Author, error) {
	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("author not found")
	}
	return author, nil
}

func (s *AuthorService) GetWithBooks(id int) (*models.Author, error) {
	author, err := s.authorRepo.GetWithBooks(id)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("author not found")
	}
	return author, nil
}

func (s *AuthorService) Update(id int, input models.UpdateAuthorInput) (*models.Author, error) {
	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errors.New("author not found")
	}

	if input.Name != nil {
		author.Name = *input.Name
	}
	if input.Biography != nil {
		author.Biography = *input.Biography
	}
	if input.BirthDate != nil {
		author.BirthDate = *input.BirthDate
	}
	if input.Nationality != nil {
		author.Nationality = *input.Nationality
	}

	if err := s.validator.Struct(author); err != nil {
		return nil, err
	}

	if err := s.authorRepo.Update(author); err != nil {
		return nil, err
	}

	return author, nil
}

func (s *AuthorService) Delete(id int) error {
	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		return err
	}
	if author == nil {
		return errors.New("author not found")
	}

	return s.authorRepo.Delete(id)
}

func (s *AuthorService) List(filter models.AuthorFilter) ([]models.Author, int64, error) {
	return s.authorRepo.List(filter)
}
