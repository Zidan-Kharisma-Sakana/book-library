package models

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Name        string    `json:"name" gorm:"size:255;not null" validate:"required"`
	Biography   string    `json:"biography" gorm:"type:text"`
	BirthDate   time.Time `json:"birth_date"`
	Nationality string    `json:"nationality" gorm:"size:100"`
	Books       []Book    `json:"books,omitempty" gorm:"foreignKey:AuthorID"`
}

func (Author) TableName() string {
	return "authors"
}

type CreateAuthorInput struct {
	Name        string    `json:"name" validate:"required"`
	Biography   string    `json:"biography"`
	BirthDate   time.Time `json:"birth_date"`
	Nationality string    `json:"nationality"`
}

type UpdateAuthorInput struct {
	Name        *string    `json:"name"`
	Biography   *string    `json:"biography"`
	BirthDate   *time.Time `json:"birth_date"`
	Nationality *string    `json:"nationality"`
}

type AuthorFilter struct {
	Name        string `form:"name"`
	Nationality string `form:"nationality"`
	Page        int    `form:"page" default:"1"`
	PageSize    int    `form:"page_size" default:"10"`
}
