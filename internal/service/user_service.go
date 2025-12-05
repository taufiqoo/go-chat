package service

import (
	"errors"

	"github.com/taufiqoo/go-chat/internal/config"
	"github.com/taufiqoo/go-chat/internal/domain"
	"github.com/taufiqoo/go-chat/internal/repository"
	"github.com/taufiqoo/go-chat/internal/utils"
)

type UserService interface {
	Register(req *domain.UserRegisterRequest) (*domain.UserResponse, error)
	Login(req *domain.UserLoginRequest) (*domain.UserResponse, error)
	GetUserByID(id uint) (*domain.User, error)
	GetAllUsers() ([]domain.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewUserService(userRepo repository.UserRepository, cfg *config.Config) UserService {
	return &userService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (u *userService) Register(req *domain.UserRegisterRequest) (*domain.UserResponse, error) {
	existingUser, _ := u.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	existingUser, _ = u.userRepo.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already taken")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Fullname: req.Fullname,
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	token, err := utils.GenerateToken(user.ID, u.cfg.JWTSecret, u.cfg.JWTExpiration)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

func (u *userService) Login(req *domain.UserLoginRequest) (*domain.UserResponse, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, u.cfg.JWTSecret, u.cfg.JWTExpiration)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:       user.ID,
		Fullname: user.Fullname,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

func (u *userService) GetUserByID(id uint) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userService) GetAllUsers() ([]domain.User, error) {
	return u.userRepo.GetAll()
}
