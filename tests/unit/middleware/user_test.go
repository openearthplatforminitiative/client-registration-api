package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/openearthplatforminitiative/client-registration-api/middleware"
	"github.com/stretchr/testify/assert"
)

func TestUserRequired(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := map[string]struct {
		usernameHeader   string
		expectedStatus   int
		expectedResponse string
		expectUserInCtx  bool
	}{
		"User header present": {
			usernameHeader:   "test-user",
			expectedStatus:   http.StatusOK,
			expectedResponse: "success",
			expectUserInCtx:  true,
		},
		"User header missing": {
			usernameHeader:   "",
			expectedStatus:   http.StatusForbidden,
			expectedResponse: `{"error":"Not supported without user"}`,
			expectUserInCtx:  false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, router := gin.CreateTestContext(w)

			router.Use(middleware.UserRequired())
			router.GET("/test", func(c *gin.Context) {
				if tc.expectUserInCtx {
					user, exists := c.Get("user")
					if exists && user == tc.usernameHeader {
						c.String(http.StatusOK, "success")
						return
					}
				}
				c.String(http.StatusInternalServerError, "user not in context")
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tc.usernameHeader != "" {
				req.Header.Set("X-Preferred-Username", tc.usernameHeader)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedResponse, w.Body.String())
		})
	}
}
