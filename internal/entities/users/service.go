package users

import (
	"errors"

	"family_budget/internal/auth"
	"family_budget/internal/dto"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(req dto.RegisterRequest) (*User, error)
	Login(req dto.LoginRequest) (accessToken string, refreshToken string, err error)
	RefreshToken(refreshToken string) (newAccessToken string, newRefreshToken string, err error)
}

type service struct {
	storage UserStorage
	jwtService auth.JWTService 
}

func NewUserService(storage UserStorage, jwtService auth.JWTService) UserService {
	return &service{storage: storage, jwtService: jwtService}
}

func (s *service) Register(req dto.RegisterRequest) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser := &User{
		Name:     req.Name,
		Surname:  req.Surname,
		Login:    req.Login,
		Password: string(hashedPassword),
		Email:    req.Email,
		Phone:    req.Phone,
	}

	err = s.storage.Create(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *service) Login(req dto.LoginRequest) (string, string, error) {
	user, err := s.storage.GetByLogin(req.Login)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	return s.jwtService.GenerateTokens(user.ID, user.FamilyID, user.RoleID)
}

func (s *service) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := s.jwtService.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	user, err := s.storage.GetByID(claims.UserID)
	if err != nil {
		return "", "", errors.New("user not found")
	}

	return s.jwtService.GenerateTokens(user.ID, user.FamilyID, user.RoleID)
}
