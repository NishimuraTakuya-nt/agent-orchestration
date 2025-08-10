package database

import (
	"context"
	"sync"

	"agent-orchestration/entities"
	"agent-orchestration/interfaces/repository"
)

// InMemoryUserRepository is an in-memory implementation for testing
type InMemoryUserRepository struct {
	users  map[int]*entities.User
	emails map[string]*entities.User
	nextID int
	mutex  sync.RWMutex
}

// NewInMemoryUserRepository creates a new in-memory user repository
func NewInMemoryUserRepository() repository.UserRepository {
	return &InMemoryUserRepository{
		users:  make(map[int]*entities.User),
		emails: make(map[string]*entities.User),
		nextID: 1,
	}
}

// Create creates a new user
func (r *InMemoryUserRepository) Create(ctx context.Context, user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	// Check if email already exists
	if _, exists := r.emails[user.Email]; exists {
		return entities.ErrUserAlreadyExists
	}
	
	// Assign new ID
	user.ID = r.nextID
	r.nextID++
	
	// Store user
	r.users[user.ID] = user
	r.emails[user.Email] = user
	
	return nil
}

// GetByID retrieves a user by ID
func (r *InMemoryUserRepository) GetByID(ctx context.Context, id int) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.users[id]
	if !exists {
		return nil, entities.ErrUserNotFound
	}
	
	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

// GetByEmail retrieves a user by email
func (r *InMemoryUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	user, exists := r.emails[email]
	if !exists {
		return nil, entities.ErrUserNotFound
	}
	
	// Return a copy to prevent external modifications
	userCopy := *user
	return &userCopy, nil
}

// Update updates an existing user
func (r *InMemoryUserRepository) Update(ctx context.Context, user *entities.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	existing, exists := r.users[user.ID]
	if !exists {
		return entities.ErrUserNotFound
	}
	
	// Remove old email mapping
	delete(r.emails, existing.Email)
	
	// Check if new email already exists (for different user)
	if emailUser, emailExists := r.emails[user.Email]; emailExists && emailUser.ID != user.ID {
		// Restore old email mapping on error
		r.emails[existing.Email] = existing
		return entities.ErrUserAlreadyExists
	}
	
	// Update user
	r.users[user.ID] = user
	r.emails[user.Email] = user
	
	return nil
}

// Delete deletes a user by ID
func (r *InMemoryUserRepository) Delete(ctx context.Context, id int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	
	user, exists := r.users[id]
	if !exists {
		return entities.ErrUserNotFound
	}
	
	// Remove user from both maps
	delete(r.users, id)
	delete(r.emails, user.Email)
	
	return nil
}

// List retrieves all users
func (r *InMemoryUserRepository) List(ctx context.Context) ([]*entities.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	users := make([]*entities.User, 0, len(r.users))
	for _, user := range r.users {
		// Return copies to prevent external modifications
		userCopy := *user
		users = append(users, &userCopy)
	}
	
	return users, nil
}