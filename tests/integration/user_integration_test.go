package integration_test

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"agent-orchestration/entities"
	"agent-orchestration/infrastructure/database"
	"agent-orchestration/use_cases"
)

var _ = Describe("User Integration Tests", func() {
	var (
		userUseCase *use_cases.UserUseCase
		ctx         context.Context
	)

	BeforeEach(func() {
		// Use real repository implementation
		repo := database.NewInMemoryUserRepository()
		userUseCase = use_cases.NewUserUseCase(repo)
		ctx = context.Background()
	})

	Describe("Complete user workflow", func() {
		It("should handle full CRUD operations", func() {
			// Create first user
			user1, err := userUseCase.CreateUser(ctx, "John Doe", "john@example.com")
			Expect(err).To(BeNil())
			Expect(user1.ID).To(BeNumerically(">", 0))
			Expect(user1.Name).To(Equal("John Doe"))
			Expect(user1.Email).To(Equal("john@example.com"))

			// Create second user
			user2, err := userUseCase.CreateUser(ctx, "Jane Doe", "jane@example.com")
			Expect(err).To(BeNil())
			Expect(user2.ID).To(BeNumerically(">", user1.ID))

			// List users - should have 2 users
			users, err := userUseCase.ListUsers(ctx)
			Expect(err).To(BeNil())
			Expect(users).To(HaveLen(2))

			// Get user by ID
			retrievedUser, err := userUseCase.GetUserByID(ctx, user1.ID)
			Expect(err).To(BeNil())
			Expect(retrievedUser.Name).To(Equal(user1.Name))
			Expect(retrievedUser.Email).To(Equal(user1.Email))

			// Update user
			updatedUser, err := userUseCase.UpdateUser(ctx, user1.ID, "John Updated", "john.updated@example.com")
			Expect(err).To(BeNil())
			Expect(updatedUser.Name).To(Equal("John Updated"))
			Expect(updatedUser.Email).To(Equal("john.updated@example.com"))
			Expect(updatedUser.Updated).To(BeTemporally(">", user1.Updated))

			// Verify update persisted
			retrievedUpdated, err := userUseCase.GetUserByID(ctx, user1.ID)
			Expect(err).To(BeNil())
			Expect(retrievedUpdated.Name).To(Equal("John Updated"))
			Expect(retrievedUpdated.Email).To(Equal("john.updated@example.com"))

			// Delete user
			err = userUseCase.DeleteUser(ctx, user1.ID)
			Expect(err).To(BeNil())

			// Verify user is deleted
			_, err = userUseCase.GetUserByID(ctx, user1.ID)
			Expect(err).To(Equal(entities.ErrUserNotFound))

			// List users - should have 1 user remaining
			users, err = userUseCase.ListUsers(ctx)
			Expect(err).To(BeNil())
			Expect(users).To(HaveLen(1))
			Expect(users[0].Name).To(Equal("Jane Doe"))
		})

		It("should prevent duplicate emails", func() {
			// Create first user
			_, err := userUseCase.CreateUser(ctx, "John Doe", "duplicate@example.com")
			Expect(err).To(BeNil())

			// Try to create second user with same email
			_, err = userUseCase.CreateUser(ctx, "Jane Doe", "duplicate@example.com")
			Expect(err).To(Equal(entities.ErrUserAlreadyExists))
		})

		It("should handle concurrent operations safely", func() {
			const numGoroutines = 10
			
			// Use a channel to coordinate goroutine completion
			done := make(chan bool, numGoroutines)
			
			// Create users concurrently
			for i := 0; i < numGoroutines; i++ {
				go func(index int) {
					defer GinkgoRecover()
					
					name := fmt.Sprintf("User %d", index)
					email := fmt.Sprintf("user%d@example.com", index)
					
					user, err := userUseCase.CreateUser(ctx, name, email)
					Expect(err).To(BeNil())
					Expect(user.Name).To(Equal(name))
					Expect(user.Email).To(Equal(email))
					
					done <- true
				}(i)
			}
			
			// Wait for all goroutines to complete
			for i := 0; i < numGoroutines; i++ {
				Eventually(done).Should(Receive())
			}
			
			// Verify all users were created
			users, err := userUseCase.ListUsers(ctx)
			Expect(err).To(BeNil())
			Expect(users).To(HaveLen(numGoroutines))
		})

		It("should maintain data consistency during updates", func() {
			// Create user
			user, err := userUseCase.CreateUser(ctx, "Test User", "test@example.com")
			Expect(err).To(BeNil())

			// Update name only
			updatedUser, err := userUseCase.UpdateUser(ctx, user.ID, "Updated Name", "")
			Expect(err).To(BeNil())
			Expect(updatedUser.Name).To(Equal("Updated Name"))
			Expect(updatedUser.Email).To(Equal("test@example.com"))

			// Update email only
			updatedUser, err = userUseCase.UpdateUser(ctx, user.ID, "", "updated@example.com")
			Expect(err).To(BeNil())
			Expect(updatedUser.Name).To(Equal("Updated Name"))
			Expect(updatedUser.Email).To(Equal("updated@example.com"))

			// Update both
			updatedUser, err = userUseCase.UpdateUser(ctx, user.ID, "Final Name", "final@example.com")
			Expect(err).To(BeNil())
			Expect(updatedUser.Name).To(Equal("Final Name"))
			Expect(updatedUser.Email).To(Equal("final@example.com"))
		})
	})

	Describe("Repository integration", func() {
		Context("when repository operations are tested directly", func() {
			var repo *database.InMemoryUserRepository

			BeforeEach(func() {
				repo = database.NewInMemoryUserRepository().(*database.InMemoryUserRepository)
			})

			It("should handle repository operations correctly", func() {
				user := &entities.User{
					Name:    "Repo Test User",
					Email:   "repo@example.com",
					Created: time.Now(),
					Updated: time.Now(),
				}

				// Test Create
				err := repo.Create(ctx, user)
				Expect(err).To(BeNil())
				Expect(user.ID).To(BeNumerically(">", 0))

				// Test GetByID
				retrieved, err := repo.GetByID(ctx, user.ID)
				Expect(err).To(BeNil())
				Expect(retrieved.Name).To(Equal(user.Name))

				// Test GetByEmail
				retrieved, err = repo.GetByEmail(ctx, user.Email)
				Expect(err).To(BeNil())
				Expect(retrieved.Name).To(Equal(user.Name))

				// Test Update
				user.Name = "Updated Repo User"
				err = repo.Update(ctx, user)
				Expect(err).To(BeNil())

				retrieved, err = repo.GetByID(ctx, user.ID)
				Expect(err).To(BeNil())
				Expect(retrieved.Name).To(Equal("Updated Repo User"))

				// Test List
				users, err := repo.List(ctx)
				Expect(err).To(BeNil())
				Expect(users).To(HaveLen(1))

				// Test Delete
				err = repo.Delete(ctx, user.ID)
				Expect(err).To(BeNil())

				_, err = repo.GetByID(ctx, user.ID)
				Expect(err).To(Equal(entities.ErrUserNotFound))
			})
		})
	})
})