package hydra

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
		},

		Paths: []*framework.Path{
			pathConfigRoot(),
			pathRoles(),
			pathListRoles(&b),
			pathUser(&b),
		},

		Secrets: []*framework.Secret{
			credentials(&b),
		},

		WALRollback:       walRollback,
		WALRollbackMinAge: time.Duration(5) * time.Minute,
		BackendType:       logical.TypeLogical,
	}

	return &b
}

type backend struct {
	*framework.Backend
}

const backendHelp = `
The Hydra backend dynamically generates credentials for a set of
service to request client id and client secret in order to provide or access api. 

After mounting this backend, credentials to generate keys must
be configured with the "root" path and policies must be written using
the "roles/" endpoints before any access keys can be generated.

Each service must be map to a roles/name in order to specify name and common service configuration (e.g: redirect URI, scope, token response,etc.)
`
