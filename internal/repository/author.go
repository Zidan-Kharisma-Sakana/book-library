package repository

import (
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"

	"gorm.io/gorm"
)

type AuthorRepository struct {
	db *gorm.DB
}

var _ interfaces.AuthorRepository = (*AuthorRepository)(nil)

func NewAuthorRepository(db *gorm.DB) *AuthorRepository {
	return &AuthorRepository{db: db}
}

func (r *AuthorRepository) Create(author *models.Author) error {
	return r.db.Create(author).Error
}

func (r *AuthorRepository) GetByID(id int) (*models.Author, error) {
	var author models.Author
	err := r.db.First(&author, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil if not found
		}
		return nil, err
	}
	return &author, nil
}

func (r *AuthorRepository) Update(author *models.Author) error {
	return r.db.Save(author).Error
}

func (r *AuthorRepository) Delete(id int) error {
	return r.db.Delete(&models.Author{}, id).Error
}

func (r *AuthorRepository) List(filter models.AuthorFilter) ([]models.Author, int64, error) {
	var authors []models.Author
	var count int64

	query := r.db.Model(&models.Author{})

	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Nationality != "" {
		query = query.Where("nationality ILIKE ?", "%"+filter.Nationality+"%")
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	err = query.Offset(offset).Limit(filter.PageSize).Find(&authors).Error
	if err != nil {
		return nil, 0, err
	}

	return authors, count, nil
}

func (r *AuthorRepository) GetWithBooks(id int) (*models.Author, error) {
	var author models.Author
	err := r.db.Preload("Books").First(&author, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &author, nil
}
