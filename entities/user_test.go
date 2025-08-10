package entities_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"agent-orchestration/entities"
)

var _ = Describe("User", func() {
	var user *entities.User

	BeforeEach(func() {
		user = &entities.User{
			ID:      1,
			Name:    "John Doe",
			Email:   "john@example.com",
			Created: time.Now(),
			Updated: time.Now(),
		}
	})

	Describe("Validate", func() {
		Context("when user has valid data", func() {
			It("should return no error", func() {
				Expect(user.Validate()).To(BeNil())
			})
		})

		Context("when name is empty", func() {
			BeforeEach(func() {
				user.Name = ""
			})

			It("should return ErrUserNameRequired", func() {
				err := user.Validate()
				Expect(err).To(Equal(entities.ErrUserNameRequired))
			})
		})

		Context("when email is empty", func() {
			BeforeEach(func() {
				user.Email = ""
			})

			It("should return ErrUserEmailRequired", func() {
				err := user.Validate()
				Expect(err).To(Equal(entities.ErrUserEmailRequired))
			})
		})

		Context("when both name and email are empty", func() {
			BeforeEach(func() {
				user.Name = ""
				user.Email = ""
			})

			It("should return ErrUserNameRequired first", func() {
				err := user.Validate()
				Expect(err).To(Equal(entities.ErrUserNameRequired))
			})
		})
	})

	Describe("IsValid", func() {
		Context("when user data is valid", func() {
			It("should return true", func() {
				Expect(user.IsValid()).To(BeTrue())
			})
		})

		Context("when user data is invalid", func() {
			BeforeEach(func() {
				user.Name = ""
			})

			It("should return false", func() {
				Expect(user.IsValid()).To(BeFalse())
			})
		})
	})

	Describe("UpdateName", func() {
		var initialUpdated time.Time

		BeforeEach(func() {
			initialUpdated = user.Updated
			time.Sleep(1 * time.Millisecond) // Ensure time difference
		})

		Context("when name is valid", func() {
			newName := "Jane Doe"

			It("should update the name and updated timestamp", func() {
				err := user.UpdateName(newName)
				Expect(err).To(BeNil())
				Expect(user.Name).To(Equal(newName))
				Expect(user.Updated).To(BeTemporally(">", initialUpdated))
			})
		})

		Context("when name is empty", func() {
			It("should return ErrUserNameRequired and not update", func() {
				originalName := user.Name
				err := user.UpdateName("")
				Expect(err).To(Equal(entities.ErrUserNameRequired))
				Expect(user.Name).To(Equal(originalName))
			})
		})
	})

	Describe("UpdateEmail", func() {
		var initialUpdated time.Time

		BeforeEach(func() {
			initialUpdated = user.Updated
			time.Sleep(1 * time.Millisecond) // Ensure time difference
		})

		Context("when email is valid", func() {
			newEmail := "jane@example.com"

			It("should update the email and updated timestamp", func() {
				err := user.UpdateEmail(newEmail)
				Expect(err).To(BeNil())
				Expect(user.Email).To(Equal(newEmail))
				Expect(user.Updated).To(BeTemporally(">", initialUpdated))
			})
		})

		Context("when email is empty", func() {
			It("should return ErrUserEmailRequired and not update", func() {
				originalEmail := user.Email
				err := user.UpdateEmail("")
				Expect(err).To(Equal(entities.ErrUserEmailRequired))
				Expect(user.Email).To(Equal(originalEmail))
			})
		})
	})

	// Table-driven tests for different user scenarios
	Describe("User validation scenarios", func() {
		DescribeTable("validation scenarios",
			func(name, email string, expectedValid bool, expectedError error) {
				user := &entities.User{Name: name, Email: email}
				err := user.Validate()
				
				if expectedValid {
					Expect(err).To(BeNil())
					Expect(user.IsValid()).To(BeTrue())
				} else {
					Expect(err).To(Equal(expectedError))
					Expect(user.IsValid()).To(BeFalse())
				}
			},
			Entry("valid user", "John Doe", "john@example.com", true, nil),
			Entry("missing name", "", "john@example.com", false, entities.ErrUserNameRequired),
			Entry("missing email", "John Doe", "", false, entities.ErrUserEmailRequired),
			Entry("missing both", "", "", false, entities.ErrUserNameRequired),
		)
	})
})