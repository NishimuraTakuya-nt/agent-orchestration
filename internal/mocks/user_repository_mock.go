package mocks

import (
	"context"
	"sync"

	"agent-orchestration/entities"
)

// Ensure, that UserRepositoryMock does implement UserRepository.
// If this is not the case, regenerate this file with moq.
//var _ repository.UserRepository = &UserRepositoryMock{}

// UserRepositoryMock is a mock implementation of UserRepository.
//
//	func TestSomethingThatUsesUserRepository(t *testing.T) {
//
//		// make and configure a mocked UserRepository
//		mockedUserRepository := &UserRepositoryMock{
//			CreateFunc: func(ctx context.Context, user *entities.User) error {
//				panic("mock out the Create method")
//			},
//			DeleteFunc: func(ctx context.Context, id int) error {
//				panic("mock out the Delete method")
//			},
//			GetByEmailFunc: func(ctx context.Context, email string) (*entities.User, error) {
//				panic("mock out the GetByEmail method")
//			},
//			GetByIDFunc: func(ctx context.Context, id int) (*entities.User, error) {
//				panic("mock out the GetByID method")
//			},
//			ListFunc: func(ctx context.Context) ([]*entities.User, error) {
//				panic("mock out the List method")
//			},
//			UpdateFunc: func(ctx context.Context, user *entities.User) error {
//				panic("mock out the Update method")
//			},
//		}
//
//		// use mockedUserRepository in code that requires UserRepository
//		// and then make assertions.
//
//	}
type UserRepositoryMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, user *entities.User) error

	// DeleteFunc mocks the Delete method.
	DeleteFunc func(ctx context.Context, id int) error

	// GetByEmailFunc mocks the GetByEmail method.
	GetByEmailFunc func(ctx context.Context, email string) (*entities.User, error)

	// GetByIDFunc mocks the GetByID method.
	GetByIDFunc func(ctx context.Context, id int) (*entities.User, error)

	// ListFunc mocks the List method.
	ListFunc func(ctx context.Context) ([]*entities.User, error)

	// UpdateFunc mocks the Update method.
	UpdateFunc func(ctx context.Context, user *entities.User) error

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// User is the user argument value.
			User *entities.User
		}
		// Delete holds details about calls to the Delete method.
		Delete []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int
		}
		// GetByEmail holds details about calls to the GetByEmail method.
		GetByEmail []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Email is the email argument value.
			Email string
		}
		// GetByID holds details about calls to the GetByID method.
		GetByID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int
		}
		// List holds details about calls to the List method.
		List []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// User is the user argument value.
			User *entities.User
		}
	}
	lockCreate     sync.RWMutex
	lockDelete     sync.RWMutex
	lockGetByEmail sync.RWMutex
	lockGetByID    sync.RWMutex
	lockList       sync.RWMutex
	lockUpdate     sync.RWMutex
}

// Create calls CreateFunc.
func (mock *UserRepositoryMock) Create(ctx context.Context, user *entities.User) error {
	if mock.CreateFunc == nil {
		panic("UserRepositoryMock.CreateFunc: method is nil but UserRepository.Create was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		User *entities.User
	}{
		Ctx:  ctx,
		User: user,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(ctx, user)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedUserRepository.CreateCalls())
func (mock *UserRepositoryMock) CreateCalls() []struct {
	Ctx  context.Context
	User *entities.User
} {
	var calls []struct {
		Ctx  context.Context
		User *entities.User
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// Delete calls DeleteFunc.
func (mock *UserRepositoryMock) Delete(ctx context.Context, id int) error {
	if mock.DeleteFunc == nil {
		panic("UserRepositoryMock.DeleteFunc: method is nil but UserRepository.Delete was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  int
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockDelete.Lock()
	mock.calls.Delete = append(mock.calls.Delete, callInfo)
	mock.lockDelete.Unlock()
	return mock.DeleteFunc(ctx, id)
}

// DeleteCalls gets all the calls that were made to Delete.
// Check the length with:
//
//	len(mockedUserRepository.DeleteCalls())
func (mock *UserRepositoryMock) DeleteCalls() []struct {
	Ctx context.Context
	ID  int
} {
	var calls []struct {
		Ctx context.Context
		ID  int
	}
	mock.lockDelete.RLock()
	calls = mock.calls.Delete
	mock.lockDelete.RUnlock()
	return calls
}

// GetByEmail calls GetByEmailFunc.
func (mock *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if mock.GetByEmailFunc == nil {
		panic("UserRepositoryMock.GetByEmailFunc: method is nil but UserRepository.GetByEmail was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Email string
	}{
		Ctx:   ctx,
		Email: email,
	}
	mock.lockGetByEmail.Lock()
	mock.calls.GetByEmail = append(mock.calls.GetByEmail, callInfo)
	mock.lockGetByEmail.Unlock()
	return mock.GetByEmailFunc(ctx, email)
}

// GetByEmailCalls gets all the calls that were made to GetByEmail.
// Check the length with:
//
//	len(mockedUserRepository.GetByEmailCalls())
func (mock *UserRepositoryMock) GetByEmailCalls() []struct {
	Ctx   context.Context
	Email string
} {
	var calls []struct {
		Ctx   context.Context
		Email string
	}
	mock.lockGetByEmail.RLock()
	calls = mock.calls.GetByEmail
	mock.lockGetByEmail.RUnlock()
	return calls
}

// GetByID calls GetByIDFunc.
func (mock *UserRepositoryMock) GetByID(ctx context.Context, id int) (*entities.User, error) {
	if mock.GetByIDFunc == nil {
		panic("UserRepositoryMock.GetByIDFunc: method is nil but UserRepository.GetByID was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  int
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGetByID.Lock()
	mock.calls.GetByID = append(mock.calls.GetByID, callInfo)
	mock.lockGetByID.Unlock()
	return mock.GetByIDFunc(ctx, id)
}

// GetByIDCalls gets all the calls that were made to GetByID.
// Check the length with:
//
//	len(mockedUserRepository.GetByIDCalls())
func (mock *UserRepositoryMock) GetByIDCalls() []struct {
	Ctx context.Context
	ID  int
} {
	var calls []struct {
		Ctx context.Context
		ID  int
	}
	mock.lockGetByID.RLock()
	calls = mock.calls.GetByID
	mock.lockGetByID.RUnlock()
	return calls
}

// List calls ListFunc.
func (mock *UserRepositoryMock) List(ctx context.Context) ([]*entities.User, error) {
	if mock.ListFunc == nil {
		panic("UserRepositoryMock.ListFunc: method is nil but UserRepository.List was just called")
	}
	callInfo := struct {
		Ctx context.Context
	}{
		Ctx: ctx,
	}
	mock.lockList.Lock()
	mock.calls.List = append(mock.calls.List, callInfo)
	mock.lockList.Unlock()
	return mock.ListFunc(ctx)
}

// ListCalls gets all the calls that were made to List.
// Check the length with:
//
//	len(mockedUserRepository.ListCalls())
func (mock *UserRepositoryMock) ListCalls() []struct {
	Ctx context.Context
} {
	var calls []struct {
		Ctx context.Context
	}
	mock.lockList.RLock()
	calls = mock.calls.List
	mock.lockList.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *UserRepositoryMock) Update(ctx context.Context, user *entities.User) error {
	if mock.UpdateFunc == nil {
		panic("UserRepositoryMock.UpdateFunc: method is nil but UserRepository.Update was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		User *entities.User
	}{
		Ctx:  ctx,
		User: user,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(ctx, user)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//
//	len(mockedUserRepository.UpdateCalls())
func (mock *UserRepositoryMock) UpdateCalls() []struct {
	Ctx  context.Context
	User *entities.User
} {
	var calls []struct {
		Ctx  context.Context
		User *entities.User
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}