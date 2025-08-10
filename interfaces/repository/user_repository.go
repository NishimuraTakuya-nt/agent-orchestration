package repository

import (
	"context"
	
	"agent-orchestration/entities"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id int) (*entities.User, error)
	
	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)
	
	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error
	
	// Delete deletes a user by ID
	Delete(ctx context.Context, id int) error
	
	// List retrieves all users
	List(ctx context.Context) ([]*entities.User, error)
}