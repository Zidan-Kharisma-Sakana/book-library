package repository

import (
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/errs"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var _ interfaces.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return errs.FromDatabase(err)
	}
	return nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.FromDatabase(err)
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.FromDatabase(err)
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errs.FromDatabase(err)
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return errs.FromDatabase(err)
	}
	return nil
}

func (r *UserRepository) Delete(id int) error {
	if err := r.db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		return errs.FromDatabase(err)
	}
	return nil
}

func (r *UserRepository) List(page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var count int64

	err := r.db.Model(&models.User{}).Count(&count).Error
	if err != nil {
		return nil, 0, errs.FromDatabase(err)
	}

	offset := (page - 1) * pageSize
	err = r.db.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, errs.FromDatabase(err)
	}

	return users, count, nil
}
