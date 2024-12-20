package keycloak

import (
	"context"
	"github.com/Nerzal/gocloak/v13"
	"github.com/openearthplatforminitiative/client-registration-api/models"
	"github.com/samber/lo"
	"strings"
)

type Keycloak interface {
	GetClients(UserName string) (*models.Clients, error)
	GetUrl() string
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
	_ Keycloak = (*KeycloakClient)(nil) // Ensure KeycloakClient implements Keycloak

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

func (k *KeycloakClient) GetClients(UserName string) (*models.Clients, error) {
	token, err := k.AdminClient.LoginAdmin(k.Ctx, k.KeycloakUser, k.KeycloakPassword, k.KeycloakMasterRealm)
	if err != nil {
		return nil, models.LoginErr
	}

	clients, err := k.AdminClient.GetClients(k.Ctx, token.AccessToken, k.KeycloakOpenEpiRealm, gocloak.GetClientsParams{})
	if err != nil {
		return nil, models.ClientLookupErr
	}

	clientList := lo.FilterMap(clients, func(c *gocloak.Client, index int) (*models.Client, bool) {
		if strings.HasPrefix(*c.ClientID, UserName) {
			return &models.Client{InternalID: c.ID, ClientID: c.ClientID, ClientName: c.Name, ClientSecret: c.Secret}, true
		}
		return nil, false
	})

	return &models.Clients{Clients: clientList}, nil
}

func (k *KeycloakClient) GetClient(username string, id string) (*models.Client, error) {
	clients, err := k.GetClients(username)
	if err != nil {
		return nil, err
	}

	client, found := lo.Find(clients.Clients, func(c *models.Client) bool {
		return *c.ClientID == id
	})

	if found {
		return client, nil
	}

	return nil, models.ClientNotFoundErr

}

func (k *KeycloakClient) AddClient(client *models.Client) (*models.Client, error) {
	token, err := k.AdminClient.LoginAdmin(k.Ctx, k.KeycloakUser, k.KeycloakPassword, k.KeycloakMasterRealm)
	if err != nil {
		return nil, models.LoginErr
	}

	if _, err = k.AdminClient.CreateClient(k.Ctx, token.AccessToken, k.KeycloakOpenEpiRealm, k.toGocloakClient(client)); err != nil {
		return nil, err
	}

	return client, nil
}

func (k *KeycloakClient) UpdateClient(client *models.Client) (*models.Client, error) {
	token, loginErr := k.AdminClient.LoginAdmin(k.Ctx, k.KeycloakUser, k.KeycloakPassword, k.KeycloakMasterRealm)
	if loginErr != nil {
		return nil, models.LoginErr
	}

	if err := k.AdminClient.UpdateClient(k.Ctx, token.AccessToken, k.KeycloakOpenEpiRealm, k.toGocloakClient(client)); err != nil {
		return nil, err
	}
	return client, nil

}

func (k *KeycloakClient) DeleteClient(username string, id string) error {
	token, err := k.AdminClient.LoginAdmin(k.Ctx, k.KeycloakUser, k.KeycloakPassword, k.KeycloakMasterRealm)
	if err != nil {
		return models.LoginErr
	}

	existingClient, err := k.GetClient(username, id)
	if err != nil {
		return models.ClientNotFoundErr
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
