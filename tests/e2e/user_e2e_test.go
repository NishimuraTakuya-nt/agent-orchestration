package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"agent-orchestration/entities"
	httphandler "agent-orchestration/interfaces/http"
)

const (
	serverURL = "http://localhost:8080"
	timeout   = 30 * time.Second
)

var _ = Describe("User E2E Tests", Ordered, func() {
	var (
		serverCmd *exec.Cmd
		httpClient *http.Client
	)

	BeforeAll(func() {
		httpClient = &http.Client{
			Timeout: 10 * time.Second,
		}

		// Start the server
		serverCmd = exec.Command("go", "run", "../../cmd/server/main.go")
		serverCmd.Dir = "."
		serverCmd.Stdout = GinkgoWriter
		serverCmd.Stderr = GinkgoWriter
		
		err := serverCmd.Start()
		Expect(err).To(BeNil())

		// Wait for server to be ready
		Eventually(func() error {
			resp, err := httpClient.Get(serverURL + "/health")
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("server not ready, status: %d", resp.StatusCode)
			}
			return nil
		}, timeout, 1*time.Second).Should(Succeed())

		DeferCleanup(func() {
			if serverCmd != nil && serverCmd.Process != nil {
				serverCmd.Process.Kill()
				serverCmd.Wait()
			}
		})
	})

	Describe("Health Check", func() {
		It("should return server health status", func() {
			resp, err := httpClient.Get(serverURL + "/health")
			Expect(err).To(BeNil())
			defer resp.Body.Close()

			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var health map[string]string
			err = json.NewDecoder(resp.Body).Decode(&health)
			Expect(err).To(BeNil())
			Expect(health["status"]).To(Equal("ok"))
		})
	})

	Describe("User API E2E", func() {
		var createdUserIDs []int

		AfterEach(func() {
			// Cleanup created users
			for _, id := range createdUserIDs {
				req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", serverURL, id), nil)
				httpClient.Do(req)
			}
			createdUserIDs = nil
		})

		Context("when creating users", func() {
			It("should create a user successfully", func() {
				createReq := httphandler.CreateUserRequest{
					Name:  "E2E Test User",
					Email: "e2e@example.com",
				}

				body, _ := json.Marshal(createReq)
				resp, err := httpClient.Post(serverURL+"/users", "application/json", bytes.NewReader(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				var user entities.User
				err = json.NewDecoder(resp.Body).Decode(&user)
				Expect(err).To(BeNil())
				Expect(user.Name).To(Equal(createReq.Name))
				Expect(user.Email).To(Equal(createReq.Email))
				Expect(user.ID).To(BeNumerically(">", 0))

				createdUserIDs = append(createdUserIDs, user.ID)
			})

			It("should reject duplicate emails", func() {
				createReq := httphandler.CreateUserRequest{
					Name:  "First User",
					Email: "duplicate@example.com",
				}

				// Create first user
				body, _ := json.Marshal(createReq)
				resp, err := httpClient.Post(serverURL+"/users", "application/json", bytes.NewReader(body))
				Expect(err).To(BeNil())
				resp.Body.Close()
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				var firstUser entities.User
				json.NewDecoder(resp.Body).Decode(&firstUser)
				createdUserIDs = append(createdUserIDs, firstUser.ID)

				// Try to create second user with same email
				createReq.Name = "Second User"
				body, _ = json.Marshal(createReq)
				resp, err = httpClient.Post(serverURL+"/users", "application/json", bytes.NewReader(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusConflict))

				var errorResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				Expect(err).To(BeNil())
				Expect(errorResp["error"]).To(Equal(entities.ErrUserAlreadyExists.Error()))
			})
		})

		Context("when retrieving users", func() {
			var testUser entities.User

			BeforeEach(func() {
				// Create a test user
				createReq := httphandler.CreateUserRequest{
					Name:  "Get Test User",
					Email: "gettest@example.com",
				}

				body, _ := json.Marshal(createReq)
				resp, err := httpClient.Post(serverURL+"/users", "application/json", bytes.NewReader(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				err = json.NewDecoder(resp.Body).Decode(&testUser)
				Expect(err).To(BeNil())

				createdUserIDs = append(createdUserIDs, testUser.ID)
			})

			It("should get user by ID", func() {
				resp, err := httpClient.Get(fmt.Sprintf("%s/users/%d", serverURL, testUser.ID))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var user entities.User
				err = json.NewDecoder(resp.Body).Decode(&user)
				Expect(err).To(BeNil())
				Expect(user.ID).To(Equal(testUser.ID))
				Expect(user.Name).To(Equal(testUser.Name))
				Expect(user.Email).To(Equal(testUser.Email))
			})

			It("should return 404 for non-existent user", func() {
				resp, err := httpClient.Get(fmt.Sprintf("%s/users/99999", serverURL))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})

			It("should list all users", func() {
				resp, err := httpClient.Get(serverURL + "/users")
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var users []*entities.User
				err = json.NewDecoder(resp.Body).Decode(&users)
				Expect(err).To(BeNil())
				Expect(users).To(HaveLen(1))
				Expect(users[0].ID).To(Equal(testUser.ID))
			})
		})

		Context("when updating users", func() {
			var testUser entities.User

			BeforeEach(func() {
				// Create a test user
				createReq := httphandler.CreateUserRequest{
					Name:  "Update Test User",
					Email: "updatetest@example.com",
				}

				body, _ := json.Marshal(createReq)
				resp, err := httpClient.Post(serverURL+"/users", "application/json", bytes.NewReader(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				err = json.NewDecoder(resp.Body).Decode(&testUser)
				Expect(err).To(BeNil())

				createdUserIDs = append(createdUserIDs, testUser.ID)
			})

			It("should update user successfully", func() {
				updateReq := httphandler.UpdateUserRequest{
					Name:  "Updated Name",
					Email: "updated@example.com",
				}

				body, _ := json.Marshal(updateReq)
				req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/users/%d", serverURL, testUser.ID), bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := httpClient.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var updatedUser entities.User
				err = json.NewDecoder(resp.Body).Decode(&updatedUser)
				Expect(err).To(BeNil())
				Expect(updatedUser.ID).To(Equal(testUser.ID))
				Expect(updatedUser.Name).To(Equal(updateReq.Name))
				Expect(updatedUser.Email).To(Equal(updateReq.Email))
				Expect(updatedUser.Updated).To(BeTemporally(">", testUser.Updated))
			})

			It("should return 404 for non-existent user update", func() {
				updateReq := httphandler.UpdateUserRequest{
					Name: "Should Not Work",
				}

				body, _ := json.Marshal(updateReq)
				req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/users/99999", serverURL), bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")

				resp, err := httpClient.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("when deleting users", func() {
			var testUser entities.User

			BeforeEach(func() {
				// Create a test user
				createReq := httphandler.CreateUserRequest{
					Name:  "Delete Test User",
					Email: "deletetest@example.com",
				}

				body, _ := json.Marshal(createReq)
				resp, err := httpClient.Post(serverURL+"/users", "application/json", bytes.NewReader(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				err = json.NewDecoder(resp.Body).Decode(&testUser)
				Expect(err).To(BeNil())

				createdUserIDs = append(createdUserIDs, testUser.ID)
			})

			It("should delete user successfully", func() {
				req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/%d", serverURL, testUser.ID), nil)
				resp, err := httpClient.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))

				// Verify user is deleted
				getResp, err := httpClient.Get(fmt.Sprintf("%s/users/%d", serverURL, testUser.ID))
				Expect(err).To(BeNil())
				defer getResp.Body.Close()

				Expect(getResp.StatusCode).To(Equal(http.StatusNotFound))

				// Remove from cleanup list since it's already deleted
				for i, id := range createdUserIDs {
					if id == testUser.ID {
						createdUserIDs = append(createdUserIDs[:i], createdUserIDs[i+1:]...)
						break
					}
				}
			})

			It("should return 404 for non-existent user deletion", func() {
				req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/users/99999", serverURL), nil)
				resp, err := httpClient.Do(req)
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("when handling edge cases", func() {
			It("should handle invalid JSON in request body", func() {
				resp, err := httpClient.Post(serverURL+"/users", "application/json", bytes.NewReader([]byte("invalid json")))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should handle invalid user ID in URL", func() {
				resp, err := httpClient.Get(serverURL + "/users/invalid")
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			})

			It("should handle missing required fields", func() {
				createReq := httphandler.CreateUserRequest{
					Name: "", // Missing name
					Email: "test@example.com",
				}

				body, _ := json.Marshal(createReq)
				resp, err := httpClient.Post(serverURL+"/users", "application/json", bytes.NewReader(body))
				Expect(err).To(BeNil())
				defer resp.Body.Close()

				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

				var errorResp map[string]string
				err = json.NewDecoder(resp.Body).Decode(&errorResp)
				Expect(err).To(BeNil())
				Expect(errorResp["error"]).To(Equal(entities.ErrUserNameRequired.Error()))
			})
		})
	})
})