package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Luzifer/ots/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreAttachmentCollision(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	for _, name := range []string{"secret.txt", "secret (1).txt"} {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte("existing"), storeFileMode))
	}

	require.NoError(t, storeAttachment(dir, client.SecretAttachment{
		Content: []byte("stored"),
		Name:    "secret.txt",
	}))

	assertFileContent(t, filepath.Join(dir, "secret (2).txt"), "stored")
	assertFileContent(t, filepath.Join(dir, "secret.txt"), "existing")
	assertFileContent(t, filepath.Join(dir, "secret (1).txt"), "existing")
}

func TestStoreAttachmentRejectsInvalidNames(t *testing.T) {
	t.Parallel()

	for _, name := range []string{"", ".", "/", `\`} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			err := storeAttachment(t.TempDir(), client.SecretAttachment{
				Content: []byte("stored"),
				Name:    name,
			})
			assert.Error(t, err)
		})
	}
}

func TestStoreAttachmentStripsPathComponents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filename string
		wantBase string
	}{
		{
			name:     "relative-parent",
			filename: "../outside.txt",
			wantBase: "outside.txt",
		},
		{
			name:     "relative-nested",
			filename: "nested/path/secret.txt",
			wantBase: "secret.txt",
		},
		{
			name:     "absolute",
			filename: filepath.Join(string(filepath.Separator), "tmp", "secret.txt"),
			wantBase: "secret.txt",
		},
		{
			name:     "windows-relative-parent",
			filename: `..\outside.txt`,
			wantBase: filepath.Base(`..\outside.txt`),
		},
		{
			name:     "windows-absolute",
			filename: `C:\Users\user\secret.txt`,
			wantBase: filepath.Base(`C:\Users\user\secret.txt`),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			parent := t.TempDir()
			dir := filepath.Join(parent, "downloads")
			require.NoError(t, os.Mkdir(dir, 0o700))

			require.NoError(t, storeAttachment(dir, client.SecretAttachment{
				Content: []byte(tc.name),
				Name:    tc.filename,
			}))

			assertFileContent(t, filepath.Join(dir, tc.wantBase), tc.name)
			assertNoRegularFilesOutsideDir(t, parent, dir)
		})
	}
}

func assertFileContent(t *testing.T, filename, want string) {
	t.Helper()

	content, err := os.ReadFile(filename) //#nosec:G304 // Test reads files created inside t.TempDir().
	require.NoError(t, err)
	assert.Equal(t, want, string(content))
}

func assertNoRegularFilesOutsideDir(t *testing.T, parent, allowedDir string) {
	t.Helper()

	err := filepath.WalkDir(parent, func(path string, entry os.DirEntry, err error) error {
		require.NoError(t, err)

		if !entry.Type().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(allowedDir, path)
		require.NoError(t, err)

		if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil
		}

		assert.Failf(t, "created regular file outside download dir", "%s", path)
		return nil
	})
	require.NoError(t, err)
}
