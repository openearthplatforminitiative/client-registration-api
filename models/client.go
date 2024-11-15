package models

type Client struct {
	InternalID   *string `json:"-"`
	ClientID     *string `json:"client_id,omitempty"`
	ClientName   *string `json:"client_name,omitempty" form:"client_name" binding:"required"`
	ClientSecret *string `json:"client_secret,omitempty"`
}
