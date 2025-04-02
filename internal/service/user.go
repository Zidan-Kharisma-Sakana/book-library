package service

import (
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/models"
	"github.com/Zidan-Kharisma-Sakana/book-library/internal/repository/interfaces"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/config"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/errs"
	"github.com/go-playground/validator/v10"
)

type UserService struct {
	userRepo    interfaces.UserRepository
	authService AuthService
	config      *config.Config
	validator   *validator.Validate
}

func NewUserService(cfg *config.Config, validator *validator.Validate, userRepo interfaces.UserRepository, authService AuthService) *UserService {
	return &UserService{
		userRepo:    userRepo,
		config:      cfg,
		validator:   validator,
		authService: authService,
	}
}

func (s *UserService) Register(input models.CreateUserInput) (*models.User, error) {
	if err := s.validator.Struct(input); err != nil {
		return nil, errs.NewBadRequestError().SetError(err)
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

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(input models.LoginInput) (*models.TokenResponse, error) {
	if err := s.validator.Struct(input); err != nil {
		return nil, errs.NewBadRequestError().SetError(err)
	}

	var user *models.User
	var err error

	if input.Username != "" {
		user, err = s.userRepo.GetByUsername(input.Username)
	} else if input.Email != "" {
		user, err = s.userRepo.GetByEmail(input.Email)
	} else {
		return nil, errs.NewBadRequestError().SetMessage("username or email is required")
	}

	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.NewBadRequestError().SetMessage("invalid credentials")
	}

	if !user.Active {
		return nil, errs.NewBadRequestError().SetMessage("user is inactive")
	}

	if !user.CheckPassword(input.Password) {
		return nil, errs.NewBadRequestError().SetMessage("invalid credentials")
	}

	token, err := s.authService.generateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.authService.generateRefreshToken(user)
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

func (s *UserService) RefreshToken(userId int) (*models.TokenResponse, error) {
	user, err := s.userRepo.GetByID(userId)
	if err != nil {
		return nil, err
	}
	token, err := s.authService.generateToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.authService.generateRefreshToken(user)
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

func (s *UserService) Update(id int, role string, input models.UpdateUserInput) (*models.User, error) {
	// Get existing user
	if err := s.validator.Struct(input); err != nil {
		return nil, errs.NewBadRequestError().SetError(err)
	}
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errs.NewNotFoundError()
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
	if input.Role != nil && role == "admin" {
		user.Role = *input.Role
	}
	if input.Active != nil && role == "admin" {
		user.Active = *input.Active
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
