package service

import (
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/config"
	"github.com/go-playground/validator/v10"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserService struct {
	userRepo  interfaces.UserRepository
	config    *config.Config
	validator *validator.Validate
}

func NewUserService(validator *validator.Validate, userRepo interfaces.UserRepository) *UserService {
	cfg, _ := config.Load()
	return &UserService{
		userRepo:  userRepo,
		config:    cfg,
		validator: validator,
	}
}

func (s *UserService) Register(input models.CreateUserInput) (*models.User, error) {
	existingUser, err := s.userRepo.GetByUsername(input.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		Username:  input.Username,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Role:      input.Role,
		Active:    true,
	}

	if err := user.SetPassword(input.Password); err != nil {
		return nil, err
	}

	if user.Role == "" {
		user.Role = "user"
	}

	if err := s.validator.Struct(user); err != nil {
		return nil, err
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(input models.LoginInput) (*models.TokenResponse, error) {
	var user *models.User
	var err error

	if input.Username != "" {
		user, err = s.userRepo.GetByUsername(input.Username)
	} else if input.Email != "" {
		user, err = s.userRepo.GetByEmail(input.Email)
	} else {
		return nil, errors.New("username or email is required")
	}

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.Active {
		return nil, errors.New("user is inactive")
	}

	// Check password
	if !user.CheckPassword(input.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &models.TokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    int(s.config.TokenExpiry.Seconds()),
		TokenType:    "Bearer",
		UserID:       int(user.ID),
		Username:     user.Username,
		Role:         user.Role,
	}, nil
}

func (s *UserService) GetByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *UserService) Update(id int, input models.UpdateUserInput) (*models.User, error) {
	// Get existing user
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if input.Username != nil {
		existingUser, err := s.userRepo.GetByUsername(*input.Username)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && int(existingUser.ID) != id {
			return nil, errors.New("username already exists")
		}
		user.Username = *input.Username
	}
	if input.Email != nil {
		existingUser, err := s.userRepo.GetByEmail(*input.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && int(existingUser.ID) != id {
			return nil, errors.New("email already exists")
		}
		user.Email = *input.Email
	}
	if input.Password != nil {
		if err := user.SetPassword(*input.Password); err != nil {
			return nil, err
		}
	}
	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		user.LastName = *input.LastName
	}
	if input.Role != nil {
		user.Role = *input.Role
	}
	if input.Active != nil {
		user.Active = *input.Active
	}

	if err := s.validator.Struct(user); err != nil {
		return nil, err
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(id int) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *UserService) List(page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.List(page, pageSize)
}

func (s *UserService) ValidateToken(tokenString string) (int, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(s.config.JWTSecret), nil
	})
	if err != nil {
		return 0, "", err
	}

	if !token.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid token")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errors.New("invalid token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", errors.New("invalid token")
	}

	return int(userID), role, nil
}

func (s *UserService) generateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(s.config.TokenExpiry * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) generateRefreshToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(s.config.TokenExpiry * 2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
