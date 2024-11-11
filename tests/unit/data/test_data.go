package data

import (
	"encoding/json"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/openearthplatforminitiative/client-registration-api/models"
	"github.com/samber/lo"
	"reflect"
)

var Client1 = &models.Client{
	InternalID:   lo.ToPtr("internal-id-1"),
	ClientID:     lo.ToPtr("client-id-1"),
	ClientName:   lo.ToPtr("client-name-1"),
	ClientSecret: lo.ToPtr("client-secret-1"),
}

var Client2 = &models.Client{
	InternalID:   lo.ToPtr("internal-id-2"),
	ClientID:     lo.ToPtr("client-id-2"),
	ClientName:   lo.ToPtr("client-name-2"),
	ClientSecret: lo.ToPtr("client-secret-2"),
}

var TwoClients = &models.Clients{Clients: []*models.Client{Client1, Client2}}
var EmptyClients = &models.Clients{Clients: make([]*models.Client, 0)}

func MarshalJSON(v interface{}) string {
	marshaled, _ := json.Marshal(v)
	return string(marshaled)
}

type MockFieldError struct {
	field string
	tag   string
}

func (fe *MockFieldError) Tag() string                       { return fe.tag }
func (fe *MockFieldError) ActualTag() string                 { return fe.tag }
func (fe *MockFieldError) Namespace() string                 { return fe.field }
func (fe *MockFieldError) StructNamespace() string           { return fe.field }
func (fe *MockFieldError) Field() string                     { return fe.field }
func (fe *MockFieldError) StructField() string               { return fe.field }
func (fe *MockFieldError) Value() interface{}                { return nil }
func (fe *MockFieldError) Param() string                     { return "" }
func (fe *MockFieldError) Kind() reflect.Kind                { return reflect.String }
func (fe *MockFieldError) Type() reflect.Type                { return reflect.TypeOf("") }
func (fe *MockFieldError) Translate(ut ut.Translator) string { return fe.Error() }
func (fe *MockFieldError) Error() string {
	return fmt.Sprintf("Key: '%s' Error:Field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Field(), fe.Tag())
}

var RequiredFieldError = &MockFieldError{field: "ClientName", tag: "required"}
var OtherValidationError = &MockFieldError{field: "ClientName", tag: "other"}
