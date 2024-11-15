package middleware

import (
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/openearthplatforminitiative/client-registration-api/middleware"
	"github.com/openearthplatforminitiative/client-registration-api/models"
	"github.com/openearthplatforminitiative/client-registration-api/tests/unit/data"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrors(t *testing.T) {
	tests := map[string]struct {
		error            error
		expectedStatus   int
		expectedResponse string
	}{
		"error-not-found": {
			error:            models.ClientNotFoundErr,
			expectedStatus:   http.StatusNotFound,
			expectedResponse: `{"errors":[{"message":"Client with id not found"}]}`,
		},
		"gocloak-api-error": {
			error: &gocloak.APIError{
				Code:    http.StatusConflict,
				Message: "A conflict occurred",
			},
			expectedStatus:   http.StatusConflict,
			expectedResponse: `{"errors":[{"message":"A conflict occurred"}]}`,
		},
		"field-required-error": {
			error: func() error {
				ve := validator.ValidationErrors{
					data.RequiredFieldError,
				}
				return ve
			}(),
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"errors":[{"field":"ClientName","message":"This field is required"}]}`,
		},
		"other-validation-error": {
			error: func() error {
				ve := validator.ValidationErrors{
					data.OtherValidationError,
				}
				return ve
			}(),
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"errors":[{"field":"ClientName","message":"Validation error"}]}`,
		},
		"unexpected-error": {
			error:            errors.New("unknown error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"errors":[]}`,
		},
		"no-error": {
			error:            nil,
			expectedStatus:   http.StatusOK,
			expectedResponse: `success`,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, router := gin.CreateTestContext(w)

			router.Use(middleware.ErrorHandler())
			router.GET("/test", func(c *gin.Context) {
				if tc.error != nil {
					_ = c.Error(tc.error)
				} else {
					c.String(http.StatusOK, "success")
				}

			})

			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedResponse, w.Body.String())
		})

	}
}
