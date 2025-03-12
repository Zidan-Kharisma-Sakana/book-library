package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" gorm:"size:50;unique;not null" validate:"required"`
	Email        string `json:"email" gorm:"size:100;unique;not null" validate:"required,email"`
	PasswordHash string `json:"-" gorm:"size:255;not null"`
	FirstName    string `json:"first_name" gorm:"size:100"`
	LastName     string `json:"last_name" gorm:"size:100"`
	Role         string `json:"role" gorm:"size:20;default:'user'"`
	Active       bool   `json:"active" gorm:"default:true"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

type CreateUserInput struct {
	Username  string `json:"username" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role" validate:"omitempty,oneof=user admin librarian"`
}

type UpdateUserInput struct {
	Username  *string `json:"username"`
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=8"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      *string `json:"role" validate:"omitempty,oneof=user admin librarian"`
	Active    *bool   `json:"active"`
}

type LoginInput struct {
	Username string `json:"username" validate:"required_without=Email"`
	Email    string `json:"email" validate:"required_without=Username,omitempty,email"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
	Role         string `json:"role"`
}
