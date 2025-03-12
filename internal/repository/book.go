package repository

import (
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

var _ interfaces.BookRepository = (*BookRepository)(nil)

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) Create(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) GetByID(id int) (*models.Book, error) {
	var book models.Book
	err := r.db.Preload("Author").First(&book, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) GetByISBN(isbn string) (*models.Book, error) {
	var book models.Book
	err := r.db.Where("isbn = ?", isbn).Preload("Author").First(&book).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) Update(book *models.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) Delete(id int) error {
	return r.db.Delete(&models.Book{}, id).Error
}

func (r *BookRepository) List(filter models.BookFilter) ([]models.Book, int64, error) {
	var books []models.Book
	var count int64

	query := r.db.Model(&models.Book{})

	if filter.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Title+"%")
	}
	if filter.AuthorID != 0 {
		query = query.Where("author_id = ?", filter.AuthorID)
	}
	if filter.Publisher != "" {
		query = query.Where("publisher ILIKE ?", "%"+filter.Publisher+"%")
	}
	if filter.Available != nil {
		query = query.Where("available = ?", *filter.Available)
	}

	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.PageSize
	err = query.Offset(offset).Limit(filter.PageSize).Preload("Author").Find(&books).Error
	if err != nil {
		return nil, 0, err
	}

	return books, count, nil
}

func (r *BookRepository) GetBooksByAuthorID(authorID int) ([]models.Book, error) {
	var books []models.Book
	err := r.db.Where("author_id = ?", authorID).Find(&books).Error
	return books, err
}
