package auth

import (
	"fmt"

	"github.com/supabase-community/gotrue-go/types"
	supabase "github.com/supabase-community/supabase-go"
)

type AuthService struct {
	db *supabase.Client
}

func NewAuthService(db *supabase.Client) *AuthService {
	return &AuthService{
		db: db,
	}
}

func (s *AuthService) SignIn(email string, password string) (types.Session, error) {
	session, err := s.db.SignInWithEmailPassword(
		email,
		password,
	)

	if err != nil {
		return types.Session{}, err
	}

	return session, nil
}

func (s *AuthService) SignUp(email string, password string) (*types.SignupResponse, error) {
	resp, err := s.db.Auth.Signup(types.SignupRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *AuthService) GetUserByToken(token string) (*types.UserResponse, error) {
	if s.db == nil || s.db.Auth == nil {
		return nil, fmt.Errorf("auth client not initialized")
	}
	authedClient := s.db.Auth.WithToken(token)
	resp, err := authedClient.GetUser()
	if err != nil {
		return nil, err
	}
	return resp, nil
}