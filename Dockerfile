FROM golang:1 as builder-secret
ARG DEP_VERSION=v0.5.3
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/$DEP_VERSION/dep-linux-amd64 && chmod +x /usr/local/bin/dep
WORKDIR /go/src/github.com/idruide/vault-plugin-secrets-hydra
COPY . .
RUN dep ensure -vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hydra-plugin ./hydra/cmd/hydra/main.go

FROM vault
WORKDIR /vault/plugins
COPY --from=builder-secret /go/src/github.com/idruide/vault-plugin-secrets-hydra/hydra-plugin ./hydra