package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/openearthplatforminitiative/client-registration-api/handlers"
	"github.com/openearthplatforminitiative/client-registration-api/models"
	"github.com/openearthplatforminitiative/client-registration-api/tests/mocks/keycloak"
	"github.com/openearthplatforminitiative/client-registration-api/tests/unit/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetClients(t *testing.T) {
	userName := "test-user"

	tests := map[string]struct {
		existingClients *models.Clients
		error           error
	}{
		"no clients": {
			existingClients: data.EmptyClients,
			error:           nil,
		},
		"two clients": {
			existingClients: data.TwoClients,
			error:           nil,
		},
		"login-error": {
			existingClients: nil,
			error:           models.LoginErr,
		},
		"lookup-error": {
			existingClients: nil,
			error:           models.ClientLookupErr,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, w, router := setupTestContext(userName)
			mockedKeycloak := new(keycloak.MockKeycloak)
			mockedKeycloak.On("GetClients", userName).Return(test.existingClients, test.error)

			cc := &handlers.ClientsHandler{Keycloak: mockedKeycloak}

			router.GET("/clients", func(c *gin.Context) { cc.Clients(ctx) })

			req, _ := http.NewRequest("GET", "/clients", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			if test.error == nil {
				assert.JSONEq(t, data.MarshalJSON(test.existingClients), w.Body.String())
			} else {
				assert.Equal(t, test.error, ctx.Errors.ByType(test.error.(*gin.Error).Type).Last())
			}

			mockedKeycloak.AssertExpectations(t)
		})
	}
}

func TestGetClient(t *testing.T) {
	userName := "test-user"

	tests := map[string]struct {
		clientId       string
		existingClient *models.Client
		error          error
	}{
		"client-found": {
			clientId:       "1",
			existingClient: data.Client1,
			error:          nil,
		},
		"client-not-found": {
			clientId:       "2",
			existingClient: nil,
			error:          models.ClientNotFoundErr,
		},
		"login-error": {
			clientId:       "3",
			existingClient: nil,
			error:          models.LoginErr,
		},
		"lookup-error": {
			clientId:       "4",
			existingClient: nil,
			error:          models.ClientLookupErr,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, w, router := setupTestContext(userName)

			mockedKeycloak := new(keycloak.MockKeycloak)
			mockedKeycloak.On("GetClient", userName, test.clientId).Return(test.existingClient, test.error)

			cc := &handlers.ClientsHandler{Keycloak: mockedKeycloak}
			ctx.Params = append(ctx.Params, gin.Param{Key: "id", Value: test.clientId})

			router.GET("/client", func(c *gin.Context) { cc.Client(ctx) })

			req, _ := http.NewRequest("GET", "/client", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			if test.error == nil {
				assert.JSONEq(t, data.MarshalJSON(test.existingClient), w.Body.String())
			} else {
				assert.Equal(t, test.error, ctx.Errors.ByType(test.error.(*gin.Error).Type).Last())
			}

			mockedKeycloak.AssertExpectations(t)
		})
	}
}

func TestAddClient(t *testing.T) {
	userName := "test-user"

	tests := map[string]struct {
		client *models.Client
		error  error
	}{
		"client-added": {
			client: data.Client1,
			error:  nil,
		},
		"login-error": {
			client: data.Client1,
			error:  models.LoginErr,
		},
		"client-exists": {
			client: data.Client1,
			error:  gocloak.APIError{Code: 409, Message: "Client with id already exists"},
		},
		"required-field-error": {
			client: data.Client1,
			error: func() error {
				ve := validator.ValidationErrors{
					data.RequiredFieldError,
				}
				return ve
			}(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, w, router := setupTestContext(userName)

			mockedKeycloak := new(keycloak.MockKeycloak)
			mockedKeycloak.On("AddClient", mock.Anything).Return(test.client, test.error)

			cc := &handlers.ClientsHandler{Keycloak: mockedKeycloak}

			router.POST("/clients", func(c *gin.Context) {
				ctx.Request = c.Request
				cc.AddClient(ctx)
			})

			formData := url.Values{}
			formData.Add("client_name", *test.client.ClientName)
			req, _ := http.NewRequestWithContext(ctx, "POST", "/clients", bytes.NewBufferString(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			if test.error == nil {
				assert.JSONEq(t, data.MarshalJSON(&test.client), w.Body.String())
			} else {
				if errors.As(test.error, &gocloak.APIError{}) {
					assert.Equal(t, test.error, ctx.Errors.Last().Unwrap())
				} else if errors.As(test.error, &validator.ValidationErrors{}) {
					assert.Equal(t, test.error, ctx.Errors.Last().Unwrap())
				} else {
					assert.Equal(t, test.error, ctx.Errors.Last())
				}
			}

			mockedKeycloak.AssertExpectations(t)
		})
	}
}

func TestUpdateClient(t *testing.T) {
	userName := "test-user"

	tests := map[string]struct {
		client *models.Client
		error  error
	}{
		"client-updated": {
			client: data.Client1,
			error:  nil,
		},
		"client-not-found": {
			client: data.Client1,
			error:  models.ClientNotFoundErr,
		},
		"other-client-with-same-name-exists": {
			client: data.Client1,
			error:  gocloak.APIError{Code: 409, Message: "Client with id already exists"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, w, router := setupTestContext(userName)

			mockedKeycloak := new(keycloak.MockKeycloak)
			mockedKeycloak.On("GetClient", userName, *test.client.ClientID).Return(test.client, test.error)
			mockedKeycloak.On("UpdateClient", mock.Anything).Return(test.client, test.error)

			cc := &handlers.ClientsHandler{Keycloak: mockedKeycloak}

			router.PUT("/clients/:id", func(c *gin.Context) {
				ctx.Request = c.Request
				ctx.Params = c.Params
				cc.UpdateClient(ctx)
			})

			formData := url.Values{}
			formData.Add("client_name", *test.client.ClientName)
			req, _ := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("/clients/%s", *test.client.ClientID), bytes.NewBufferString(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			if test.error == nil {
				assert.JSONEq(t, data.MarshalJSON(&test.client), w.Body.String())
				mockedKeycloak.AssertExpectations(t)
			} else {
				if errors.As(test.error, &gocloak.APIError{}) {
					assert.Equal(t, test.error, ctx.Errors.Last().Unwrap())
				} else {
					assert.Equal(t, test.error, ctx.Errors.Last())
				}
			}

		})
	}
}

func TestDeleteClient(t *testing.T) {
	userName := "test-user"

	tests := map[string]struct {
		client     *models.Client
		statusCode int
		error      error
	}{
		"client-deleted": {
			client:     data.Client1,
			statusCode: http.StatusNoContent,
			error:      nil,
		},
		"client-not-found": {
			client:     data.Client1,
			statusCode: http.StatusOK,
			error:      gocloak.APIError{Code: 404, Message: "Client with id already exists"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, w, router := setupTestContext(userName)

			mockedKeycloak := new(keycloak.MockKeycloak)
			mockedKeycloak.On("DeleteClient", userName, *test.client.ClientID).Return(test.error)

			cc := &handlers.ClientsHandler{Keycloak: mockedKeycloak}

			router.DELETE("/clients/:id", func(c *gin.Context) {
				ctx.Request = c.Request
				ctx.Params = c.Params
				cc.DeleteClient(ctx)
			})

			req, _ := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("/clients/%s", *test.client.ClientID), nil)

			router.ServeHTTP(w, req)

			assert.Equal(t, test.statusCode, w.Code)
			if test.error == nil {
				mockedKeycloak.AssertExpectations(t)
			} else {
				if errors.As(test.error, &gocloak.APIError{}) {
					assert.Equal(t, test.error, ctx.Errors.Last().Unwrap())
				} else {
					assert.Equal(t, test.error, ctx.Errors.Last())
				}
			}

		})
	}
}

func setupTestContext(userName string) (*gin.Context, *httptest.ResponseRecorder, *gin.Engine) {
	w := httptest.NewRecorder()
	ctx, router := gin.CreateTestContext(w)
	ctx.Set("user", userName)

	return ctx, w, router
}
