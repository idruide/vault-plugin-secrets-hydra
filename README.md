# Vault Plugin Secret for Ory Hydra

This secret engine allow you to register a client config then request credentials (clientId and clientSecret) from Hydra.
It's a simple implementation that can be enhanced by the community.

Follow instruction on [Vault website](https://www.vaultproject.io/docs/plugin/) to activate the plugin. 

```
# config.hcl
api_addr = "http://127.0.0.1:8200"
plugin_directory = "/vault/plugins"
...
```
Then register the plugin in the catalog and activate the secret engine.
```
sha256sum /vault/plugins/hydra
vault plugin register -sha256=ed7059de1557294e8a03c1bd0a3e2c89b0d10d4c1f8f10bc3d9c6ce979651887 -command=hydra secret hydra
vault secrets enable -path=hydra hydra
```

First, we need to tell our plugin where its admin service is.

```
$ vault write hydra/config/root admin_url=http://hydra.svc:4445 [skip_tls_verify=true/false] [public_url=https://hydra.localhost.svc] [client_id=xxxx] [client_secret=xxxx]
```

Register a client config

```
$ vault write hydra/roles/oauth_client_test \
  grant_types=authorization_code,refresh_token,client_credentials,implicit \
  response_types=token,code,id_token \
  allowed_scopes=openid,offline,profile,email,phone,address \
  redirect_urls=https://localhost.svc/callback
  [lease=24h]
```

Finally read creds will generate a new credentials based on the config you define previously

```
$ vault read hydra/creds/oauth_client_test
Key                Value
---                -----
lease_id           hydra/creds/oauth_client_test/nELoZiaHAhNXtAU8B8yU5TJM
lease_duration     24h
lease_renewable    true
client_id          2ddbc3fc-4410-4fe6-a5a3-99281736dd49
client_secret      4zQNSO3piU8o
url                https://hydra.localhost.svc

```

You can now use this credentials to talk to your hydra endpoint. 