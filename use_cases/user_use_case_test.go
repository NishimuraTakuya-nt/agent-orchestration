package use_cases_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"agent-orchestration/entities"
	"agent-orchestration/internal/mocks"
	"agent-orchestration/use_cases"
)

var _ = Describe("UserUseCase", func() {
	var (
		userUseCase *use_cases.UserUseCase
		mockRepo    *mocks.UserRepositoryMock
		ctx         context.Context
	)

	BeforeEach(func() {
		mockRepo = &mocks.UserRepositoryMock{}
		userUseCase = use_cases.NewUserUseCase(mockRepo)
		ctx = context.Background()
	})

	Describe("CreateUser", func() {
		const (
			validName  = "John Doe"
			validEmail = "john@example.com"
		)

		Context("when user data is valid and user doesn't exist", func() {
			BeforeEach(func() {
				mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
				mockRepo.CreateFunc = func(ctx context.Context, user *entities.User) error {
					user.ID = 1
					return nil
				}
			})

			It("should create user successfully", func() {
				user, err := userUseCase.CreateUser(ctx, validName, validEmail)
				
				Expect(err).To(BeNil())
				Expect(user).NotTo(BeNil())
				Expect(user.Name).To(Equal(validName))
				Expect(user.Email).To(Equal(validEmail))
				Expect(user.ID).To(Equal(1))
				
				// Verify repository calls
				Expect(mockRepo.GetByEmailCalls()).To(HaveLen(1))
				Expect(mockRepo.GetByEmailCalls()[0].Email).To(Equal(validEmail))
				Expect(mockRepo.CreateCalls()).To(HaveLen(1))
			})
		})

		Context("when user already exists", func() {
			BeforeEach(func() {
				existingUser := &entities.User{
					ID:    1,
					Name:  "Existing User",
					Email: validEmail,
				}
				mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) {
					return existingUser, nil
				}
			})

			It("should return ErrUserAlreadyExists", func() {
				user, err := userUseCase.CreateUser(ctx, validName, validEmail)
				
				Expect(user).To(BeNil())
				Expect(err).To(Equal(entities.ErrUserAlreadyExists))
				Expect(mockRepo.CreateCalls()).To(HaveLen(0))
			})
		})

		Context("when user data is invalid", func() {
			DescribeTable("invalid user data scenarios",
				func(name, email string, expectedError error) {
					mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) {
						return nil, entities.ErrUserNotFound
					}
					
					user, err := userUseCase.CreateUser(ctx, name, email)
					
					Expect(user).To(BeNil())
					Expect(err).To(Equal(expectedError))
					Expect(mockRepo.CreateCalls()).To(HaveLen(0))
				},
				Entry("empty name", "", validEmail, entities.ErrUserNameRequired),
				Entry("empty email", validName, "", entities.ErrUserEmailRequired),
			)
		})

		Context("when repository create fails", func() {
			BeforeEach(func() {
				mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
				mockRepo.CreateFunc = func(ctx context.Context, user *entities.User) error {
					return errors.New("database error")
				}
			})

			It("should return the repository error", func() {
				user, err := userUseCase.CreateUser(ctx, validName, validEmail)
				
				Expect(user).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("database error"))
			})
		})
	})

	Describe("GetUserByID", func() {
		Context("when ID is valid and user exists", func() {
			expectedUser := &entities.User{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
			}

			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					if id == 1 {
						return expectedUser, nil
					}
					return nil, entities.ErrUserNotFound
				}
			})

			It("should return the user", func() {
				user, err := userUseCase.GetUserByID(ctx, 1)
				
				Expect(err).To(BeNil())
				Expect(user).To(Equal(expectedUser))
				Expect(mockRepo.GetByIDCalls()).To(HaveLen(1))
				Expect(mockRepo.GetByIDCalls()[0].ID).To(Equal(1))
			})
		})

		Context("when ID is invalid", func() {
			DescribeTable("invalid ID scenarios",
				func(id int) {
					user, err := userUseCase.GetUserByID(ctx, id)
					
					Expect(user).To(BeNil())
					Expect(err).To(Equal(entities.ErrInvalidID))
					Expect(mockRepo.GetByIDCalls()).To(HaveLen(0))
				},
				Entry("zero ID", 0),
				Entry("negative ID", -1),
			)
		})

		Context("when user not found", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
			})

			It("should return ErrUserNotFound", func() {
				user, err := userUseCase.GetUserByID(ctx, 999)
				
				Expect(user).To(BeNil())
				Expect(err).To(Equal(entities.ErrUserNotFound))
			})
		})
	})

	Describe("UpdateUser", func() {
		var existingUser *entities.User

		BeforeEach(func() {
			existingUser = &entities.User{
				ID:      1,
				Name:    "John Doe",
				Email:   "john@example.com",
				Created: time.Now().Add(-24 * time.Hour),
				Updated: time.Now().Add(-1 * time.Hour),
			}
		})

		Context("when updating name only", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					if id == 1 {
						// Return a copy to avoid modifying the original
						user := *existingUser
						return &user, nil
					}
					return nil, entities.ErrUserNotFound
				}
				mockRepo.UpdateFunc = func(ctx context.Context, user *entities.User) error {
					return nil
				}
			})

			It("should update name and timestamp", func() {
				newName := "Jane Doe"
				user, err := userUseCase.UpdateUser(ctx, 1, newName, "")
				
				Expect(err).To(BeNil())
				Expect(user.Name).To(Equal(newName))
				Expect(user.Email).To(Equal(existingUser.Email))
				Expect(user.Updated).To(BeTemporally(">", existingUser.Updated))
				
				Expect(mockRepo.GetByIDCalls()).To(HaveLen(1))
				Expect(mockRepo.UpdateCalls()).To(HaveLen(1))
			})
		})

		Context("when updating email only", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					if id == 1 {
						// Return a copy to avoid modifying the original
						user := *existingUser
						return &user, nil
					}
					return nil, entities.ErrUserNotFound
				}
				mockRepo.UpdateFunc = func(ctx context.Context, user *entities.User) error {
					return nil
				}
			})

			It("should update email and timestamp", func() {
				newEmail := "jane@example.com"
				user, err := userUseCase.UpdateUser(ctx, 1, "", newEmail)
				
				Expect(err).To(BeNil())
				Expect(user.Name).To(Equal(existingUser.Name))
				Expect(user.Email).To(Equal(newEmail))
				Expect(user.Updated).To(BeTemporally(">", existingUser.Updated))
			})
		})

		Context("when user not found", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
			})

			It("should return ErrUserNotFound", func() {
				user, err := userUseCase.UpdateUser(ctx, 999, "New Name", "")
				
				Expect(user).To(BeNil())
				Expect(err).To(Equal(entities.ErrUserNotFound))
				Expect(mockRepo.UpdateCalls()).To(HaveLen(0))
			})
		})

		Context("when ID is invalid", func() {
			It("should return ErrInvalidID", func() {
				user, err := userUseCase.UpdateUser(ctx, 0, "New Name", "")
				
				Expect(user).To(BeNil())
				Expect(err).To(Equal(entities.ErrInvalidID))
				Expect(mockRepo.GetByIDCalls()).To(HaveLen(0))
			})
		})
	})

	Describe("DeleteUser", func() {
		Context("when user exists", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					if id == 1 {
						return &entities.User{ID: 1}, nil
					}
					return nil, entities.ErrUserNotFound
				}
				mockRepo.DeleteFunc = func(ctx context.Context, id int) error {
					return nil
				}
			})

			It("should delete user successfully", func() {
				err := userUseCase.DeleteUser(ctx, 1)
				
				Expect(err).To(BeNil())
				Expect(mockRepo.GetByIDCalls()).To(HaveLen(1))
				Expect(mockRepo.DeleteCalls()).To(HaveLen(1))
				Expect(mockRepo.DeleteCalls()[0].ID).To(Equal(1))
			})
		})

		Context("when user not found", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
			})

			It("should return ErrUserNotFound", func() {
				err := userUseCase.DeleteUser(ctx, 999)
				
				Expect(err).To(Equal(entities.ErrUserNotFound))
				Expect(mockRepo.DeleteCalls()).To(HaveLen(0))
			})
		})

		Context("when ID is invalid", func() {
			It("should return ErrInvalidID", func() {
				err := userUseCase.DeleteUser(ctx, 0)
				
				Expect(err).To(Equal(entities.ErrInvalidID))
				Expect(mockRepo.GetByIDCalls()).To(HaveLen(0))
			})
		})
	})

	Describe("ListUsers", func() {
		Context("when users exist", func() {
			expectedUsers := []*entities.User{
				{ID: 1, Name: "User 1", Email: "user1@example.com"},
				{ID: 2, Name: "User 2", Email: "user2@example.com"},
			}

			BeforeEach(func() {
				mockRepo.ListFunc = func(ctx context.Context) ([]*entities.User, error) {
					return expectedUsers, nil
				}
			})

			It("should return all users", func() {
				users, err := userUseCase.ListUsers(ctx)
				
				Expect(err).To(BeNil())
				Expect(users).To(Equal(expectedUsers))
				Expect(mockRepo.ListCalls()).To(HaveLen(1))
			})
		})

		Context("when no users exist", func() {
			BeforeEach(func() {
				mockRepo.ListFunc = func(ctx context.Context) ([]*entities.User, error) {
					return []*entities.User{}, nil
				}
			})

			It("should return empty slice", func() {
				users, err := userUseCase.ListUsers(ctx)
				
				Expect(err).To(BeNil())
				Expect(users).To(BeEmpty())
			})
		})

		Context("when repository fails", func() {
			BeforeEach(func() {
				mockRepo.ListFunc = func(ctx context.Context) ([]*entities.User, error) {
					return nil, errors.New("database error")
				}
			})

			It("should return the error", func() {
				users, err := userUseCase.ListUsers(ctx)
				
				Expect(users).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("database error"))
			})
		})
	})
})