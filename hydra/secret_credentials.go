package hydra

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ory/hydra/sdk/go/hydra/swagger"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

const CredentialsType = "hydra_credentials"

func credentials(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: CredentialsType,
		Fields: map[string]*framework.FieldSchema{
			"client_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Client ID",
			},

			"client_secret": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Client Secret",
			},
			"url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "URL of the server",
			},
		},
		Renew:  b.credentialsRenew,
		Revoke: b.credentialsRevoke,
	}
}

func (b *backend) credentialsCreate(ctx context.Context, s logical.Storage, displayName, roleName string, role *hydraRole) (*logical.Response, error) {

	client, err := HydraClient(ctx, s)

	oauthClient, _, err := client.CreateOAuth2Client(swagger.OAuth2Client{
		GrantTypes:    role.GrantTypes,
		RedirectUris:  role.RedirectUrls,
		Scope:         strings.Join(role.AllowedScopes, " "),
		ResponseTypes: role.ResponseTypes,
	})

	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("Error creating credentials: %v", err)), nil
	}

	// if res.StatusCode != 200 {
	// 	return logical.ErrorResponse(string(res.Payload)), nil
	// }

	resp := b.Secret(CredentialsType).Response(map[string]interface{}{
		"client_id":     oauthClient.ClientId,
		"client_secret": oauthClient.ClientSecret,
		"url":           client.Configuration.PublicURL,
	}, map[string]interface{}{
		"displayName": displayName,
		"role":        roleName,
		"id":          oauthClient.ClientId,
		"is_sts":      true,
	})

	resp.Secret.TTL = role.Lease

	return resp, nil
}

func (b *backend) credentialsRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	f := framework.LeaseExtend(0, time.Duration(1)*time.Hour, b.System())
	return f(ctx, req, d)
}

func (b *backend) credentialsRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

	// Get the clientname from the internal data
	idRaw, ok := req.Secret.InternalData["id"]
	if !ok {
		return nil, fmt.Errorf("secret is missing id internal data")
	}
	id, ok := idRaw.(string)
	if !ok {
		return nil, fmt.Errorf("secret is missing id internal data")
	}

	// Use the user rollback mechanism to delete this user
	err := pathUserRollback(ctx, req, "client", id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
