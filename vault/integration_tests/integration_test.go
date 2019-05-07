package integration_tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gitlab.com/gitlab-org/gitlab-runner/vault"
	"gitlab.com/gitlab-org/gitlab-runner/vault/config"
	"gitlab.com/gitlab-org/gitlab-runner/vault/secret"
)

func TestTokenLogin(t *testing.T) {
	s := NewService(t)

	conf := config.Vault{
		Server: s.GetVaultServerConfig(ServiceProxyPort),
		Auth: config.VaultAuth{
			Token: s.GetVaultTokenAuthConfig(),
		},
	}

	v := vault.New(nil)

	err := v.Connect(conf.Server)
	assert.NoError(t, err)

	err = v.Authenticate(conf.Auth)
	assert.NoError(t, err)
}

func TestUserpassLogin(t *testing.T) {
	s := NewService(t)

	conf := config.Vault{
		Server: s.GetVaultServerConfig(ServiceProxyPort),
		Auth: config.VaultAuth{
			Userpass: s.GetVaultUserpassAuthConfig(),
		},
	}

	v := vault.New(nil)

	err := v.Connect(conf.Server)
	assert.NoError(t, err)

	err = v.Authenticate(conf.Auth)
	assert.NoError(t, err)
}

func TestTLSLogin(t *testing.T) {
	s := NewService(t)

	conf := config.Vault{
		Server: s.GetVaultServerConfig(ServiceDirectPort),
		Auth: config.VaultAuth{
			TLS: s.GetVaultTLSAuthConfig(),
		},
	}

	v := vault.New(nil)

	err := v.Connect(conf.Server)
	assert.NoError(t, err)

	err = v.Authenticate(conf.Auth)
	assert.NoError(t, err)
}

func TestSecretRead(t *testing.T) {
	s := NewService(t)

	conf := config.Vault{
		Server: s.GetVaultServerConfig(ServiceProxyPort),
		Auth: config.VaultAuth{
			Token: s.GetVaultTokenAuthConfig(),
		},
		Secrets: s.GetVaultSecretsConfig(),
	}

	builderMock := new(secret.MockBuilder)
	defer builderMock.AssertExpectations(t)

	for _, sec := range conf.Secrets {
		for _, key := range sec.Keys {
			builderMock.On("BuildSecret", key, mock.Anything).
				Return(nil).
				Once()
		}
	}

	v := vault.New(builderMock)

	err := v.Connect(conf.Server)
	assert.NoError(t, err)

	err = v.Authenticate(conf.Auth)
	assert.NoError(t, err)

	err = v.ReadSecrets(conf.Secrets)
	assert.NoError(t, err)
}
