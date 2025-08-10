package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"agent-orchestration/entities"
	httphandler "agent-orchestration/interfaces/http"
	"agent-orchestration/internal/mocks"
	"agent-orchestration/use_cases"
)

var _ = Describe("UserHandler", func() {
	var (
		handler     *httphandler.UserHandler
		mockRepo    *mocks.UserRepositoryMock
		userUseCase *use_cases.UserUseCase
		router      *chi.Mux
	)

	BeforeEach(func() {
		mockRepo = &mocks.UserRepositoryMock{}
		userUseCase = use_cases.NewUserUseCase(mockRepo)
		handler = httphandler.NewUserHandler(userUseCase)
		
		// Setup Chi router
		router = chi.NewRouter()
		router.Post("/users", handler.CreateUser)
		router.Get("/users/{id}", handler.GetUser)
		router.Put("/users/{id}", handler.UpdateUser)
		router.Delete("/users/{id}", handler.DeleteUser)
		router.Get("/users", handler.ListUsers)
	})

	Describe("CreateUser", func() {
		Context("when request is valid", func() {
			var requestBody httphandler.CreateUserRequest

			BeforeEach(func() {
				requestBody = httphandler.CreateUserRequest{
					Name:  "John Doe",
					Email: "john@example.com",
				}
				
				mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
				mockRepo.CreateFunc = func(ctx context.Context, user *entities.User) error {
					user.ID = 1
					return nil
				}
			})

			It("should create user and return 201", func() {
				body, _ := json.Marshal(requestBody)
				req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusCreated))
				
				var user entities.User
				err := json.Unmarshal(w.Body.Bytes(), &user)
				Expect(err).To(BeNil())
				Expect(user.Name).To(Equal(requestBody.Name))
				Expect(user.Email).To(Equal(requestBody.Email))
				Expect(user.ID).To(Equal(1))
			})
		})

		Context("when request body is invalid JSON", func() {
			It("should return 400 Bad Request", func() {
				req := httptest.NewRequest("POST", "/users", strings.NewReader("invalid json"))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				Expect(response["error"]).To(Equal("invalid request body"))
			})
		})

		Context("when user already exists", func() {
			BeforeEach(func() {
				existingUser := &entities.User{
					ID:    1,
					Name:  "Existing User",
					Email: "john@example.com",
				}
				mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) {
					return existingUser, nil
				}
			})

			It("should return 409 Conflict", func() {
				requestBody := httphandler.CreateUserRequest{
					Name:  "John Doe",
					Email: "john@example.com",
				}
				body, _ := json.Marshal(requestBody)
				req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusConflict))
				
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				Expect(response["error"]).To(Equal(entities.ErrUserAlreadyExists.Error()))
			})
		})

		Context("when user data is invalid", func() {
			BeforeEach(func() {
				// Setup mock to return no existing user for validation
				mockRepo.GetByEmailFunc = func(ctx context.Context, email string) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
			})
			
			DescribeTable("invalid user data scenarios",
				func(name, email string, expectedStatus int, expectedError string) {
					requestBody := httphandler.CreateUserRequest{
						Name:  name,
						Email: email,
					}
					body, _ := json.Marshal(requestBody)
					req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
					req.Header.Set("Content-Type", "application/json")
					w := httptest.NewRecorder()

					router.ServeHTTP(w, req)

					Expect(w.Code).To(Equal(expectedStatus))
					
					var response map[string]string
					json.Unmarshal(w.Body.Bytes(), &response)
					Expect(response["error"]).To(Equal(expectedError))
				},
				Entry("empty name", "", "john@example.com", http.StatusBadRequest, entities.ErrUserNameRequired.Error()),
				Entry("empty email", "John Doe", "", http.StatusBadRequest, entities.ErrUserEmailRequired.Error()),
			)
		})
	})

	Describe("GetUser", func() {
		Context("when user exists", func() {
			expectedUser := &entities.User{
				ID:      1,
				Name:    "John Doe",
				Email:   "john@example.com",
				Created: time.Now(),
				Updated: time.Now(),
			}

			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					if id == 1 {
						return expectedUser, nil
					}
					return nil, entities.ErrUserNotFound
				}
			})

			It("should return user with 200", func() {
				req := httptest.NewRequest("GET", "/users/1", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				
				var user entities.User
				err := json.Unmarshal(w.Body.Bytes(), &user)
				Expect(err).To(BeNil())
				Expect(user.ID).To(Equal(expectedUser.ID))
				Expect(user.Name).To(Equal(expectedUser.Name))
				Expect(user.Email).To(Equal(expectedUser.Email))
			})
		})

		Context("when user not found", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
			})

			It("should return 404 Not Found", func() {
				req := httptest.NewRequest("GET", "/users/999", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
				
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				Expect(response["error"]).To(Equal(entities.ErrUserNotFound.Error()))
			})
		})

		Context("when ID is invalid", func() {
			It("should return 400 Bad Request", func() {
				req := httptest.NewRequest("GET", "/users/invalid", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
				
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				Expect(response["error"]).To(Equal("invalid user ID"))
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

		Context("when update is successful", func() {
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

			It("should update user and return 200", func() {
				requestBody := httphandler.UpdateUserRequest{
					Name:  "Jane Doe",
					Email: "jane@example.com",
				}
				body, _ := json.Marshal(requestBody)
				req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				
				var user entities.User
				err := json.Unmarshal(w.Body.Bytes(), &user)
				Expect(err).To(BeNil())
				Expect(user.Name).To(Equal(requestBody.Name))
				Expect(user.Email).To(Equal(requestBody.Email))
				Expect(user.Updated).To(BeTemporally(">", existingUser.Updated))
			})
		})

		Context("when user not found", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
			})

			It("should return 404 Not Found", func() {
				requestBody := httphandler.UpdateUserRequest{Name: "New Name"}
				body, _ := json.Marshal(requestBody)
				req := httptest.NewRequest("PUT", "/users/999", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when request body is invalid", func() {
			It("should return 400 Bad Request", func() {
				req := httptest.NewRequest("PUT", "/users/1", strings.NewReader("invalid json"))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
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

			It("should delete user and return 204", func() {
				req := httptest.NewRequest("DELETE", "/users/1", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNoContent))
				Expect(w.Body.String()).To(BeEmpty())
			})
		})

		Context("when user not found", func() {
			BeforeEach(func() {
				mockRepo.GetByIDFunc = func(ctx context.Context, id int) (*entities.User, error) {
					return nil, entities.ErrUserNotFound
				}
			})

			It("should return 404 Not Found", func() {
				req := httptest.NewRequest("DELETE", "/users/999", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when ID is invalid", func() {
			It("should return 400 Bad Request", func() {
				req := httptest.NewRequest("DELETE", "/users/invalid", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusBadRequest))
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

			It("should return users with 200", func() {
				req := httptest.NewRequest("GET", "/users", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				
				var users []*entities.User
				err := json.Unmarshal(w.Body.Bytes(), &users)
				Expect(err).To(BeNil())
				Expect(users).To(HaveLen(2))
				Expect(users[0].Name).To(Equal("User 1"))
				Expect(users[1].Name).To(Equal("User 2"))
			})
		})

		Context("when no users exist", func() {
			BeforeEach(func() {
				mockRepo.ListFunc = func(ctx context.Context) ([]*entities.User, error) {
					return []*entities.User{}, nil
				}
			})

			It("should return empty array with 200", func() {
				req := httptest.NewRequest("GET", "/users", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusOK))
				
				var users []*entities.User
				err := json.Unmarshal(w.Body.Bytes(), &users)
				Expect(err).To(BeNil())
				Expect(users).To(BeEmpty())
			})
		})

		Context("when repository fails", func() {
			BeforeEach(func() {
				mockRepo.ListFunc = func(ctx context.Context) ([]*entities.User, error) {
					return nil, fmt.Errorf("database error")
				}
			})

			It("should return 500 Internal Server Error", func() {
				req := httptest.NewRequest("GET", "/users", nil)
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)

				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				
				var response map[string]string
				json.Unmarshal(w.Body.Bytes(), &response)
				Expect(response["error"]).To(Equal("failed to list users"))
			})
		})
	})
})