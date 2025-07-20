package auth

import "awesomeProject2/internal/pkg/utils"

type AuthService struct {
	Repo UserRepository
}

func NewAuthService(repo UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) Register(email, password string) error {
	_, err := s.Repo.FindByEmail(email)
	if err == nil {
		return ErrUserExists
	}

	hash, _ := utils.HashPassword(password)
	return s.Repo.Create(&User{
		Email:    email,
		Password: hash,
	})
}

func (s *AuthService) Login(email, password string) (*User, error) {
	user, err := s.Repo.FindByEmail(email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}
