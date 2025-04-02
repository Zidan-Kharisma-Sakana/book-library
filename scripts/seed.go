package main

import (
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/db"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/config"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/logger"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

func hashPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hash)
}

func SeedDatabase(db *gorm.DB) {
	log.Println("Seeding database...")
	users := []models.User{
		{
			Username:     "admin1",
			Email:        "admin1@example.com",
			PasswordHash: hashPassword("admin1"),
			FirstName:    "Admin",
			LastName:     "User",
			Role:         "admin",
			Active:       true,
		},
		{
			Username:     "librarian1",
			Email:        "librarian1@example.com",
			PasswordHash: hashPassword("librarian1"),
			FirstName:    "Librarian",
			LastName:     "User",
			Role:         "librarian",
			Active:       true,
		},
		{
			Username:     "user1",
			Email:        "user1@example.com",
			PasswordHash: hashPassword("user1"),
			FirstName:    "Regular",
			LastName:     "User",
			Role:         "user",
			Active:       true,
		},
	}
	for _, user := range users {
		if err := db.FirstOrCreate(&user, models.User{Username: user.Username}).Error; err != nil {
			log.Fatalf("failed to seed user: %v", err)
		}
	}

	authors := []models.Author{
		{
			Name:        "Jane Austen",
			Biography:   "British novelist",
			BirthDate:   time.Date(1775, 12, 16, 0, 0, 0, 0, time.UTC),
			Nationality: "British",
		},
		{
			Name:        "George Orwell",
			Biography:   "English novelist and essayist",
			BirthDate:   time.Date(1903, 6, 25, 0, 0, 0, 0, time.UTC),
			Nationality: "English",
		},
	}
	for _, author := range authors {
		if err := db.FirstOrCreate(&author, models.Author{Name: author.Name}).Error; err != nil {
			log.Fatalf("failed to seed author: %v", err)
		}
	}

	// Create Books
	books := []models.Book{
		{
			Title:       "Pride and Prejudice",
			ISBN:        "978-0141439518",
			Description: "A novel by Jane Austen",
			AuthorID:    1, // Assuming Jane Austen's ID is 1
			Publisher:   "T. Egerton",
			PublishedAt: time.Date(1813, 1, 28, 0, 0, 0, 0, time.UTC),
			Pages:       279,
			Available:   true,
		},
		{
			Title:       "Nineteen Eighty-Four",
			ISBN:        "978-0451524935",
			Description: "A dystopian novel by George Orwell",
			AuthorID:    2, // Assuming George Orwell's ID is 2
			Publisher:   "Secker & Warburg",
			PublishedAt: time.Date(1949, 6, 8, 0, 0, 0, 0, time.UTC),
			Pages:       328,
			Available:   true,
		},
	}
	for _, book := range books {
		if err := db.FirstOrCreate(&book, models.Book{ISBN: book.ISBN}).Error; err != nil {
			log.Fatalf("failed to seed book: %v", err)
		}
	}

	log.Println("Database seeding completed.")
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	logger.Initialize(cfg.LogLevel)

	database, err := db.NewConnection(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer func() {
		err := db.Close(database)
		if err != nil {
			logger.Fatal("Failed to close database", "error", err)
		}
	}()

	SeedDatabase(database)
}
