package user

import (
	"context"
	"fmt"

	domain "github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment"
)

type Service struct {
	userRepository *Repository
}

func NewService(userRepository *Repository) *Service {
	return &Service{userRepository: userRepository}
}

func (s *Service) CreateUser(ctx context.Context, in domain.User) error {
	if err := s.userRepository.Create(ctx, in); err != nil {
		return fmt.Errorf("unable to create user: %w", err)
	}

	return nil
}

func (s *Service) LogIn(ctx context.Context, username, password string) (domain.User, error) {
	return s.userRepository.LogIn(ctx, username, password)
}
