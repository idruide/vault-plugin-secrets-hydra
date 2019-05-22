package hydra

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/structs"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathListRoles(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "roles/?$",

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ListOperation: b.pathRoleList,
		},

		HelpSynopsis:    pathListRolesHelpSyn,
		HelpDescription: pathListRolesHelpDesc,
	}
}

type hydraRole struct {
	GrantTypes    []string      `json:"grant_types" structs:"grant_types"`
	ResponseTypes []string      `json:"response_types" structs:"response_types"`
	RedirectUrls  []string      `json:"redirect_urls" structs:"redirect_urls"`
	AllowedScopes []string      `json:"allowed_scopes" structs:"allowed_scopes"`
	Lease         time.Duration `json:"lease" structs:"lease"`
}

func pathRoles() *framework.Path {
	return &framework.Path{
		Pattern: "roles/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": &framework.FieldSchema{
				Type:        framework.TypeString,
				Description: "Name of the role",
			},
			"grant_types": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Grants supported by this role",
			},
			"response_types": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "Response types supported by this role",
			},
			"redirect_urls": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "URL allowed to redirect to",
			},
			"allowed_scopes": &framework.FieldSchema{
				Type:        framework.TypeCommaStringSlice,
				Description: "A space separated string that represent the list of supported scopes for the role",
			},
			"lease": &framework.FieldSchema{
				Type:        framework.TypeString,
				Default:     "4h",
				Description: "The lease length; defaults to 4 hours",
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.DeleteOperation: pathRolesDelete,
			logical.ReadOperation:   pathRolesRead,
			logical.UpdateOperation: pathRolesWrite,
		},

		HelpSynopsis:    pathRolesHelpSyn,
		HelpDescription: pathRolesHelpDesc,
	}
}

func (b *backend) pathRoleList(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	entries, err := req.Storage.List(ctx, "role/")
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(entries), nil
}

func pathRolesDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	err := req.Storage.Delete(ctx, "role/"+d.Get("name").(string))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func getRole(ctx context.Context, s logical.Storage, n string) (*hydraRole, error) {
	entry, err := s.Get(ctx, "role/"+n)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var result hydraRole
	if err := entry.DecodeJSON(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func pathRolesRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	role, err := getRole(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: structs.New(role).Map(),
	}, nil
}

func pathRolesWrite(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	grants := d.Get("grant_types").([]string)
	responses := d.Get("response_types").([]string)
	redirects := d.Get("redirect_urls").([]string)
	scopes := d.Get("allowed_scopes").([]string)

	leaseRaw := d.Get("lease").(string)
	lease, err := time.ParseDuration(leaseRaw)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Error parsing lease value of %s: %s", leaseRaw, err)), nil
	}

	if len(grants) == 0 {
		return logical.ErrorResponse("Role must have grants and scopes"), nil
	}

	entry := &hydraRole{
		AllowedScopes: scopes,
		GrantTypes:    grants,
		RedirectUrls:  redirects,
		ResponseTypes: responses,
		Lease:         lease,
	}

	// Store it
	entryJSON, err := logical.StorageEntryJSON("role/"+name, entry)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entryJSON); err != nil {
		return nil, err
	}

	return nil, nil
}

const pathListRolesHelpSyn = `List the existing roles in this backend`

const pathListRolesHelpDesc = `Roles will be listed by the role name.`

const pathRolesHelpSyn = `
Read, write and reference ACL policies that access keys can be made for.
`

const pathRolesHelpDesc = `
This path allows you to read and write roles that are used to
create OAuth client. These roles are associated with ACL policies that
map directly to the route to read the Client credentials.
To validate the keys, attempt to read a client credentials after writing the role.
`
