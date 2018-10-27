package backend

import (
	"context"
	"time"

	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

type kubeBackend struct {
	*framework.Backend
	testMode bool
}

// New creates and returns new instance of Kubernetes secrets manager backend
func New() *kubeBackend {
	var b kubeBackend

	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,

		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{"login"},
		},
		WALRollback:       b.walRollback,
		WALRollbackMinAge: 5 * time.Minute,
		Paths: []*framework.Path{
			pathConfig(&b),
			pathServiceAccounts(&b),
			pathSecrets(&b),
			// TODO P1 pathListRoles
			// TODO P1 pathConfigRotateToken
		},
		Secrets: []*framework.Secret{
			secretAccessTokens(&b),
		},
	}

	b.testMode = false

	return &b
}

// Factory creates and returns new backend with BackendConfig
func Factory(ctx context.Context, c *logical.BackendConfig) (logical.Backend, error) {
	b := New()
	if err := b.Setup(ctx, c); err != nil {
		return nil, err
	}
	return b, nil
}