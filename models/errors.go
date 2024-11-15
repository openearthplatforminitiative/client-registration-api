package models

import (
	"errors"
	"github.com/gin-gonic/gin"
)

var (
	LoginErr          = newGinError("failed to login to the AdminAPI", gin.ErrorTypePrivate)
	ClientLookupErr   = newGinError("failed to lookup clients", gin.ErrorTypePrivate)
	ClientNotFoundErr = newGinError("Client with id not found", gin.ErrorTypePublic)
)

func newGinError(message string, errorType gin.ErrorType) *gin.Error {
	return &gin.Error{
		Err:  errors.New(message),
		Type: errorType,
	}
}
