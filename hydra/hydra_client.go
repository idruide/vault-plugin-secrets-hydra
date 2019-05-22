package hydra

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/vault/logical"
	"github.com/ory/hydra/sdk/go/hydra"
)

func HydraClient(ctx context.Context, s logical.Storage) (*hydra.CodeGenSDK, error) {
	var configuration *hydra.Configuration

	entry, err := s.Get(ctx, "config/root")
	if err != nil {
		return nil, err
	}
	if entry != nil {
		var config rootConfig
		if err := entry.DecodeJSON(&config); err != nil {
			return nil, fmt.Errorf("error reading root configuration: %s", err)
		}
		configuration = &hydra.Configuration{
			PublicURL:    config.PublicURL,
			AdminURL:     config.AdminURL,
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scopes:       []string{"hydra.clients", "hydra.policies"},
		}
	} else {
		return nil, errors.New("No Config found")
	}

	if configuration != nil {
		sdk, err := hydra.NewSDK(configuration)
		if err != nil {
			return nil, err
		}
		if adminURL, err := url.Parse(configuration.AdminURL); err == nil {
			if adminURL.Scheme == "http" {
				fmt.Print("Adding X-Forwarded-Proto header")
				sdk.AdminApi.Configuration.DefaultHeader["X-Forwarded-Proto"] = "https"
			}
		}
		return sdk, nil
	}

	return nil, errors.New("no valid configuration")
}
