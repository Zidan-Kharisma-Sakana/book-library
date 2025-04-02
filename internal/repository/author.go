package repository

import (
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/errs"
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
	if err := r.db.Create(author).Error; err != nil {
		return errs.FromDatabase(err)
	}
	return nil
}

func (r *AuthorRepository) GetByID(id int) (*models.Author, error) {
	var author models.Author
	err := r.db.First(&author, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.FromDatabase(err)
	}
	return &author, nil
}

func (r *AuthorRepository) Update(author *models.Author) error {
	if err := r.db.Save(author).Error; err != nil {
		return errs.FromDatabase(err)
	}
	return nil
}

func (r *AuthorRepository) Delete(id int) error {
	if err := r.db.Delete(&models.Author{}, id).Error; err != nil {
		return errs.FromDatabase(err)
	}
	return nil
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
		return nil, 0, errs.FromDatabase(err)
	}

	offset := (filter.Page - 1) * filter.PageSize
	err = query.Offset(offset).Limit(filter.PageSize).Find(&authors).Error
	if err != nil {
		return nil, 0, errs.FromDatabase(err)
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
		return nil, errs.FromDatabase(err)
	}
	return &author, nil
}
