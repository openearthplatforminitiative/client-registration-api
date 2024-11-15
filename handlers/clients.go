package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/openearthplatforminitiative/client-registration-api/keycloak"
	"github.com/openearthplatforminitiative/client-registration-api/models"
	"github.com/samber/lo"
	"net/http"
)

type ClientsHandler struct {
	Keycloak keycloak.Keycloak
}

func (cc *ClientsHandler) Clients(context *gin.Context) {
	username := context.MustGet("user").(string)
	clients, err := cc.Keycloak.GetClients(username)
	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusOK, clients)
}

func (cc *ClientsHandler) Client(context *gin.Context) {
	userName := context.MustGet("user").(string)
	clientId := context.Param("id")

	client, err := cc.Keycloak.GetClient(userName, clientId)
	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusOK, client)
}

func (cc *ClientsHandler) AddClient(context *gin.Context) {
	userName := context.MustGet("user").(string)

	var client models.Client
	if err := context.ShouldBind(&client); err != nil {
		_ = context.Error(err)
		return
	}

	clientId := userName + "-" + *client.ClientName
	clientSecret := lo.RandomString(32, lo.LettersCharset)

	client.ClientID = &clientId
	client.ClientSecret = &clientSecret

	createdClient, err := cc.Keycloak.AddClient(&client)
	if err != nil {
		_ = context.Error(err)
		return
	}
	context.JSON(http.StatusOK, createdClient)
}

func (cc *ClientsHandler) UpdateClient(context *gin.Context) {
	userName := context.MustGet("user").(string)
	clientId := context.Param("id")

	existingClient, err := cc.Keycloak.GetClient(userName, clientId)
	if err != nil {
		_ = context.Error(err)
		return
	}

	var updatedClient models.Client
	if err := context.ShouldBind(&updatedClient); err != nil {
		_ = context.Error(err)
		return
	}

	existingClient.ClientName = updatedClient.ClientName

	updatedClientId := userName + "-" + *existingClient.ClientName
	clientSecret := lo.RandomString(32, lo.LettersCharset)

	existingClient.ClientID = &updatedClientId
	existingClient.ClientSecret = &clientSecret

	client, err := cc.Keycloak.UpdateClient(existingClient)
	if err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusOK, client)

}

func (cc *ClientsHandler) DeleteClient(context *gin.Context) {
	userName := context.MustGet("user").(string)
	clientId := context.Param("id")

	if err := cc.Keycloak.DeleteClient(userName, clientId); err != nil {
		_ = context.Error(err)
		return
	}

	context.JSON(http.StatusNoContent, nil)
}
