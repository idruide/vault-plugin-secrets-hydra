package hydra

import (
	"context"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathConfigRoot() *framework.Path {
	return &framework.Path{
		Pattern: "config/root",
		Fields: map[string]*framework.FieldSchema{
			"admin_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Address of Hydra Server Admin URL",
			},
			"public_url": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Address of Hydra Server Public URL",
			},
			"client_id": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Hydra client id",
			},
			"client_secret": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Hydra client secret",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: pathConfigRootWrite,
		},

		HelpSynopsis:    pathConfigRootHelpSyn,
		HelpDescription: pathConfigRootHelpDesc,
	}
}

func pathConfigRootWrite(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	entry, err := logical.StorageEntryJSON("config/root", rootConfig{
		AdminURL:     data.Get("admin_url").(string),
		PublicURL:    data.Get("public_url").(string),
		ClientID:     data.Get("client_id").(string),
		ClientSecret: data.Get("client_secret").(string),
	})
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

type rootConfig struct {
	AdminURL     string `json:"admin_url"`
	PublicURL    string `json:"public_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

const pathConfigRootHelpSyn = `
Configure the root credentials that are used to manage Hydra.
`

const pathConfigRootHelpDesc = `
Before doing anything, the Hydra backend needs configuration about endpoints location and 
credentials to manage roles etc. This endpoint is used
to configure those credentials.
`
