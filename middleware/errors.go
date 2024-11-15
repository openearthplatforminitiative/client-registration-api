package middleware

import (
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/openearthplatforminitiative/client-registration-api/models"
	"net/http"
)

type ErrorMsg struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			errorCode := http.StatusInternalServerError
			errorMsgs := make([]ErrorMsg, 0)

			var apiError *gocloak.APIError
			var ve validator.ValidationErrors

			switch {
			case errors.Is(err, models.ClientNotFoundErr):
				errorCode = http.StatusNotFound
				errorMsgs = append(errorMsgs, ErrorMsg{Message: "Client with id not found"})
			case errors.As(err, &apiError):
				errorCode = apiError.Code
				errorMsgs = append(errorMsgs, ErrorMsg{Message: apiError.Message})
			case errors.As(err, &ve):
				errorCode = http.StatusBadRequest
				errorMsgs = getErrorMsg(ve)
			default:
				errorCode = http.StatusInternalServerError
			}

			c.JSON(errorCode, gin.H{
				"errors": errorMsgs,
			})
		}
	}
}

func getErrorMsg(ve validator.ValidationErrors) []ErrorMsg {
	msgs := make([]ErrorMsg, len(ve))
	for i, fe := range ve {
		tag := fe.Tag()
		switch tag {
		case "required":
			msgs[i] = ErrorMsg{fe.Field(), "This field is required"}
		default:
			msgs[i] = ErrorMsg{fe.Field(), "Validation error"}
		}
	}
	return msgs
}
