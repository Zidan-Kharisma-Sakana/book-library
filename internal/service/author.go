package service

import (
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/errs"
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
	if err := s.validator.Struct(input); err != nil {
		return nil, err
	}

	author := &models.Author{
		Name:        input.Name,
		Biography:   input.Biography,
		BirthDate:   input.BirthDate,
		Nationality: input.Nationality,
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
	return author, nil
}

func (s *AuthorService) GetWithBooks(id int) (*models.Author, error) {
	author, err := s.authorRepo.GetWithBooks(id)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (s *AuthorService) Update(id int, input models.UpdateAuthorInput) (*models.Author, error) {
	if err := s.validator.Struct(input); err != nil {
		return nil, err
	}

	author, err := s.authorRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if author == nil {
		return nil, errs.NewNotFoundError()
	}

	author.Name = *input.Name
	author.Biography = *input.Biography
	author.BirthDate = *input.BirthDate
	author.Nationality = *input.Nationality

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
		return errs.NewNotFoundError()
	}

	return s.authorRepo.Delete(id)
}

func (s *AuthorService) List(filter models.AuthorFilter) ([]models.Author, int64, error) {
	return s.authorRepo.List(filter)
}
