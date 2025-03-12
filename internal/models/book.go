package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string    `json:"title" gorm:"size:255;not null" validate:"required"`
	ISBN        string    `json:"isbn" gorm:"size:20;unique;not null" validate:"required,isbn"`
	Description string    `json:"description" gorm:"type:text"`
	AuthorID    int       `json:"author_id" gorm:"not null" validate:"required"`
	Author      Author    `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	Publisher   string    `json:"publisher" gorm:"size:255"`
	PublishedAt time.Time `json:"published_at"`
	Pages       int       `json:"pages" validate:"gte=0"`
	Available   bool      `json:"available" gorm:"default:true"`
}

func (Book) TableName() string {
	return "books"
}

type CreateBookInput struct {
	Title       string    `json:"title" validate:"required"`
	ISBN        string    `json:"isbn" validate:"required,isbn"`
	Description string    `json:"description"`
	AuthorID    int       `json:"author_id" validate:"required"`
	Publisher   string    `json:"publisher"`
	PublishedAt time.Time `json:"published_at"`
	Pages       int       `json:"pages" validate:"gte=0"`
	Available   bool      `json:"available"`
}

type UpdateBookInput struct {
	Title       *string    `json:"title"`
	ISBN        *string    `json:"isbn" validate:"omitempty,isbn"`
	Description *string    `json:"description"`
	AuthorID    *int       `json:"author_id"`
	Publisher   *string    `json:"publisher"`
	PublishedAt *time.Time `json:"published_at"`
	Pages       *int       `json:"pages" validate:"omitempty,gte=0"`
	Available   *bool      `json:"available"`
}

type BookFilter struct {
	Title     string `form:"title"`
	AuthorID  int    `form:"author_id"`
	Publisher string `form:"publisher"`
	Available *bool  `form:"available"`
	Page      int    `form:"page" default:"1"`
	PageSize  int    `form:"page_size" default:"10"`
}
