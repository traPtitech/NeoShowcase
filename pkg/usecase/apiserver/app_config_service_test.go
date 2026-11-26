package apiserver_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/traPtitech/neoshowcase/pkg/domain"
)

func TestSetEnvironmentVariable_SystemVariableIsNotOverwritten(t *testing.T) {
	t.Parallel()
	f := newAppFixture(t, "set-env-system", appDatabases{MariaDB: true})
	before, ok := f.env(t, domain.EnvMariaDBDatabaseKey)
	require.True(t, ok)

	err := f.svc.SetEnvironmentVariable(f.ctx, f.app.ID, domain.EnvMariaDBDatabaseKey, "nsapp_someone_else")

	_, typ, _ := domain.DecomposeError(err)
	assert.Equal(t, domain.ErrorTypeBadRequest, typ)
	after, ok := f.env(t, domain.EnvMariaDBDatabaseKey)
	require.True(t, ok)
	assert.Equal(t, before, after)
}

func TestDeleteEnvironmentVariable_SystemVariableIsNotDeleted(t *testing.T) {
	t.Parallel()
	f := newAppFixture(t, "delete-env-system", appDatabases{MariaDB: true})

	err := f.svc.DeleteEnvironmentVariable(f.ctx, f.app.ID, domain.EnvMariaDBDatabaseKey)

	_, typ, _ := domain.DecomposeError(err)
	assert.Equal(t, domain.ErrorTypeBadRequest, typ)
	_, ok := f.env(t, domain.EnvMariaDBDatabaseKey)
	assert.True(t, ok)
}

func TestEnvironmentVariable_UserDefinedNSPrefixedVariableIsEditable(t *testing.T) {
	t.Parallel()
	f := newAppFixture(t, "user-ns-env", appDatabases{MariaDB: true})
	const key = "NS_USER_DEFINED"

	require.NoError(t, f.svc.SetEnvironmentVariable(f.ctx, f.app.ID, key, "first"))
	require.NoError(t, f.svc.SetEnvironmentVariable(f.ctx, f.app.ID, key, "second"))
	got, ok := f.env(t, key)
	require.True(t, ok)
	assert.Equal(t, "second", got.Value)
	assert.False(t, got.System)

	require.NoError(t, f.svc.DeleteEnvironmentVariable(f.ctx, f.app.ID, key))
	_, ok = f.env(t, key)
	assert.False(t, ok)
}
