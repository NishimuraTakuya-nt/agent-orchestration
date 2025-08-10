package testutils

import (
	"fmt"
	"time"

	"agent-orchestration/entities"
)

// TestUserBuilder provides a builder pattern for creating test users
type TestUserBuilder struct {
	user *entities.User
}

// NewTestUserBuilder creates a new test user builder with default values
func NewTestUserBuilder() *TestUserBuilder {
	return &TestUserBuilder{
		user: &entities.User{
			ID:      1,
			Name:    "Test User",
			Email:   "test@example.com",
			Created: time.Now(),
			Updated: time.Now(),
		},
	}
}

// WithID sets the user ID
func (b *TestUserBuilder) WithID(id int) *TestUserBuilder {
	b.user.ID = id
	return b
}

// WithName sets the user name
func (b *TestUserBuilder) WithName(name string) *TestUserBuilder {
	b.user.Name = name
	return b
}

// WithEmail sets the user email
func (b *TestUserBuilder) WithEmail(email string) *TestUserBuilder {
	b.user.Email = email
	return b
}

// WithCreated sets the created timestamp
func (b *TestUserBuilder) WithCreated(created time.Time) *TestUserBuilder {
	b.user.Created = created
	return b
}

// WithUpdated sets the updated timestamp
func (b *TestUserBuilder) WithUpdated(updated time.Time) *TestUserBuilder {
	b.user.Updated = updated
	return b
}

// Build returns the built user
func (b *TestUserBuilder) Build() *entities.User {
	// Return a copy to prevent modifications to the builder
	user := *b.user
	return &user
}

// BuildPointer returns a pointer to the built user
func (b *TestUserBuilder) BuildPointer() *entities.User {
	return b.Build()
}

// Common test users for use across test files
var (
	// ValidTestUser represents a valid test user
	ValidTestUser = NewTestUserBuilder().
			WithName("John Doe").
			WithEmail("john@example.com").
			Build()

	// AnotherValidTestUser represents another valid test user
	AnotherValidTestUser = NewTestUserBuilder().
				WithID(2).
				WithName("Jane Doe").
				WithEmail("jane@example.com").
				Build()

	// InvalidNameUser represents a user with invalid name
	InvalidNameUser = NewTestUserBuilder().
			WithName("").
			Build()

	// InvalidEmailUser represents a user with invalid email
	InvalidEmailUser = NewTestUserBuilder().
				WithEmail("").
				Build()
)

// CreateTestUsers returns a slice of test users
func CreateTestUsers(count int) []*entities.User {
	users := make([]*entities.User, count)
	for i := 0; i < count; i++ {
		users[i] = NewTestUserBuilder().
			WithID(i + 1).
			WithName(fmt.Sprintf("User %d", i+1)).
			WithEmail(fmt.Sprintf("user%d@example.com", i+1)).
			Build()
	}
	return users
}

// AssertUserEqual compares two users for testing
func AssertUserEqual(expected, actual *entities.User, ignoreTimestamps bool) bool {
	if expected.ID != actual.ID {
		return false
	}
	if expected.Name != actual.Name {
		return false
	}
	if expected.Email != actual.Email {
		return false
	}
	
	if !ignoreTimestamps {
		if !expected.Created.Equal(actual.Created) {
			return false
		}
		if !expected.Updated.Equal(actual.Updated) {
			return false
		}
	}
	
	return true
}

// CloneUser creates a deep copy of a user
func CloneUser(user *entities.User) *entities.User {
	if user == nil {
		return nil
	}
	
	clone := *user
	return &clone
}

// CloneUsers creates a deep copy of a slice of users
func CloneUsers(users []*entities.User) []*entities.User {
	if users == nil {
		return nil
	}
	
	clones := make([]*entities.User, len(users))
	for i, user := range users {
		clones[i] = CloneUser(user)
	}
	return clones
}