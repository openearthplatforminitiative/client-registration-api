package keycloak

import (
	"context"
	"errors"
	"github.com/Nerzal/gocloak/v13"
	"github.com/openearthplatforminitiative/client-registration-api/models"
	s "strings"
)

type Keycloak interface {
	GetUrl() string
	GetClients(UserName string) ([]*models.Client, error)
	GetClient(username string, id string) (*models.Client, error)
	AddClient(client *models.Client) (*models.Client, error)
	UpdateClient(client *models.Client) (*models.Client, error)
	DeleteClient(username string, id string) error
}

type KeycloakClient struct {
	KeycloakUrl          string
	KeycloakUser         string
	KeycloakPassword     string
	KeycloakMasterRealm  string
	KeycloakOpenEpiRealm string
	Ctx                  context.Context
	AdminClient          *gocloak.GoCloak
}

var (
	_                 Keycloak = (*KeycloakClient)(nil) // Ensure KeycloakClient implements Keycloak
	LoginErr                   = errors.New("failed to login to the AdminAPI")
	ClientLookupErr            = errors.New("failed to lookup clients")
	ClientNotFoundErr          = errors.New("client with given id was not found")
)

func NewKeycloak(url string, user string, pw string, masterRealm string, openEpiRealm string) Keycloak {
	return &KeycloakClient{
		KeycloakUrl:          url,
		KeycloakUser:         user,
		KeycloakPassword:     pw,
		KeycloakMasterRealm:  masterRealm,
		KeycloakOpenEpiRealm: openEpiRealm,
		Ctx:                  context.Background(),
		AdminClient:          gocloak.NewClient(url),
	}
}

func (k *KeycloakClient) GetUrl() string {
	return k.KeycloakUrl
}

func (k *KeycloakClient) GetClients(UserName string) ([]*models.Client, error) {
	token, err := k.AdminClient.LoginAdmin(k.Ctx, k.KeycloakUser, k.KeycloakPassword, k.KeycloakMasterRealm)
	if err != nil {
		return nil, LoginErr
	}

	clients, err := k.AdminClient.GetClients(k.Ctx, token.AccessToken, k.KeycloakOpenEpiRealm, gocloak.GetClientsParams{})
	if err != nil {
		return nil, ClientLookupErr
	}

	// Filter list of clients by UserName
	clientList := make([]*models.Client, 0)
	for _, cl := range clients {
		if s.HasPrefix(*cl.ClientID, UserName) {
			clientList = append(clientList, &models.Client{InternalID: cl.ID, ClientID: cl.ClientID, ClientName: cl.Name, ClientSecret: cl.Secret})
		}
	}

	return clientList, nil
}

func (k *KeycloakClient) GetClient(username string, id string) (*models.Client, error) {
	clients, err := k.GetClients(username)
	if err != nil {
		return nil, err
	}

	for _, cl := range clients {
		if *cl.ClientID == id {
			return cl, nil
		}
	}

	return nil, ClientNotFoundErr
}

func (k *KeycloakClient) AddClient(client *models.Client) (*models.Client, error) {
	token, err := k.AdminClient.LoginAdmin(k.Ctx, k.KeycloakUser, k.KeycloakPassword, k.KeycloakMasterRealm)
	if err != nil {
		return nil, LoginErr
	}

	if _, err = k.AdminClient.CreateClient(k.Ctx, token.AccessToken, k.KeycloakOpenEpiRealm, k.toGocloakClient(client)); err != nil {
		return nil, err
	}

	return client, nil
}

func (k *KeycloakClient) UpdateClient(client *models.Client) (*models.Client, error) {
	token, loginErr := k.AdminClient.LoginAdmin(k.Ctx, k.KeycloakUser, k.KeycloakPassword, k.KeycloakMasterRealm)
	if loginErr != nil {
		return nil, LoginErr
	}

	if err := k.AdminClient.UpdateClient(k.Ctx, token.AccessToken, k.KeycloakOpenEpiRealm, k.toGocloakClient(client)); err != nil {
		return nil, err
	}
	return client, nil

}

func (k *KeycloakClient) DeleteClient(username string, id string) error {
	token, err := k.AdminClient.LoginAdmin(k.Ctx, k.KeycloakUser, k.KeycloakPassword, k.KeycloakMasterRealm)
	if err != nil {
		return LoginErr
	}

	existingClient, err := k.GetClient(username, id)
	if err != nil {
		return ClientNotFoundErr
	}

	if err := k.AdminClient.DeleteClient(k.Ctx, token.AccessToken, k.KeycloakOpenEpiRealm, *existingClient.InternalID); err != nil {
		return err
	}

	return nil
}

func (k *KeycloakClient) toGocloakClient(client *models.Client) gocloak.Client {
	return gocloak.Client{
		ID:                        client.InternalID,
		Name:                      client.ClientName,
		ClientID:                  client.ClientID,                   // Required: Client ID
		Secret:                    client.ClientSecret,               // Required: Client Secret
		Enabled:                   gocloak.BoolP(true),               // Required: Enable the client
		ServiceAccountsEnabled:    gocloak.BoolP(true),               // Required: Enable Service Account
		PublicClient:              gocloak.BoolP(false),              // Required: It's not a public client
		DirectAccessGrantsEnabled: gocloak.BoolP(false),              // Disable Direct Access Grants
		StandardFlowEnabled:       gocloak.BoolP(false),              // Disable Standard Flow
		Protocol:                  gocloak.StringP("openid-connect"), // OpenID Connect Protocol
	}
}
