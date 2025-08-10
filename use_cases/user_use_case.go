package use_cases

import (
	"context"
	"time"
	
	"agent-orchestration/entities"
	"agent-orchestration/interfaces/repository"
)

// UserUseCase handles user business logic
type UserUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (uc *UserUseCase) CreateUser(ctx context.Context, name, email string) (*entities.User, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, entities.ErrUserAlreadyExists
	}
	
	// Create new user
	user := &entities.User{
		Name:    name,
		Email:   email,
		Created: time.Now(),
		Updated: time.Now(),
	}
	
	// Validate user
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	// Save user
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(ctx context.Context, id int) (*entities.User, error) {
	if id <= 0 {
		return nil, entities.ErrInvalidID
	}
	
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// UpdateUser updates an existing user
func (uc *UserUseCase) UpdateUser(ctx context.Context, id int, name, email string) (*entities.User, error) {
	if id <= 0 {
		return nil, entities.ErrInvalidID
	}
	
	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	// Update user data
	if name != "" {
		if err := user.UpdateName(name); err != nil {
			return nil, err
		}
	}
	
	if email != "" {
		if err := user.UpdateEmail(email); err != nil {
			return nil, err
		}
	}
	
	// Save updated user
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	
	return user, nil
}

// DeleteUser deletes a user by ID
func (uc *UserUseCase) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return entities.ErrInvalidID
	}
	
	// Check if user exists
	_, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	// Delete user
	return uc.userRepo.Delete(ctx, id)
}

// ListUsers retrieves all users
func (uc *UserUseCase) ListUsers(ctx context.Context) ([]*entities.User, error) {
	return uc.userRepo.List(ctx)
}