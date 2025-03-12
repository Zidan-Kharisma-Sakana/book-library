package interfaces

import "github.com/Zidan-Kharisma-Sakana/book-library/internal/models"

type BookRepository interface {
	Create(book *models.Book) error
	GetByID(id int) (*models.Book, error)
	GetByISBN(isbn string) (*models.Book, error)
	Update(book *models.Book) error
	Delete(id int) error
	List(filter models.BookFilter) ([]models.Book, int64, error)
	GetBooksByAuthorID(authorID int) ([]models.Book, error)
}

type AuthorRepository interface {
	Create(author *models.Author) error
	GetByID(id int) (*models.Author, error)
	Update(author *models.Author) error
	Delete(id int) error
	List(filter models.AuthorFilter) ([]models.Author, int64, error)
	GetWithBooks(id int) (*models.Author, error)
}

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id int) error
	List(page, pageSize int) ([]models.User, int64, error)
}
