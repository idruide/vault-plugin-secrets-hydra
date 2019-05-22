package hydra

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathUser(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.pathUserRead,
		},

		HelpSynopsis:    pathUserHelpSyn,
		HelpDescription: pathUserHelpDesc,
	}
}

// V0.9.2
func (b *backend) pathUserRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)

	// Read the Role
	role, err := getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, fmt.Errorf("error retrieving role: %s", err)
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Role '%s' not found", roleName)), nil
	}
	// Use the helper to create the secret
	return b.credentialsCreate(ctx, req.Storage, req.DisplayName, roleName, role)
}

func pathUserRollback(ctx context.Context, req *logical.Request, _kind string, data interface{}) error {

	// Get the client
	client, err := HydraClient(ctx, req.Storage)
	if err != nil {
		return err
	}

	// Delete the user
	_, err = client.DeleteOAuth2Client(data.(string))
	if err != nil {
		return err
	}

	return nil
}

const pathUserHelpSyn = `
Generate an access key pair for a specific role.
`

const pathUserHelpDesc = `
This path will generate a new, never before used key pair for
accessing Hydra.

The client credentials will have a lease associated with them. The credential
can be revoked by using the lease ID.
`
