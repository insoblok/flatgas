package internal

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveAddressOrAlias(t *testing.T) {
	tempDir := t.TempDir()
	aliasesPath := filepath.Join(tempDir, "wallet")
	require.NoError(t, os.MkdirAll(aliasesPath, 0755))

	aliasFile := filepath.Join(aliasesPath, "aliases.json")
	aliasesJSON := `{
	  "alice": "0x1111111111111111111111111111111111111111",
	  "bob":   "0x2222222222222222222222222222222222222222"
	}`
	require.NoError(t, os.WriteFile(aliasFile, []byte(aliasesJSON), 0644))

	t.Run("resolve known alias", func(t *testing.T) {
		addr, err := ResolveAddressOrAlias(tempDir, "alice")
		require.NoError(t, err)
		require.Equal(t, "0x1111111111111111111111111111111111111111", addr)
	})

	t.Run("resolve valid address", func(t *testing.T) {
		addr, err := ResolveAddressOrAlias(tempDir, "0x3333333333333333333333333333333333333333")
		require.NoError(t, err)
		require.Equal(t, "0x3333333333333333333333333333333333333333", addr)
	})

	t.Run("unknown input", func(t *testing.T) {
		_, err := ResolveAddressOrAlias(tempDir, "not-found")
		require.Error(t, err)
	})
}
