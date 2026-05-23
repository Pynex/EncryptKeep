package appsvc

import (
	"os"
	"path/filepath"
	"testing"

	"encryptkeep-backend/internal/keymanager"
)

func TestNormalizeMasterPassword(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "unchanged plain password",
			in:   "strong-pass-123",
			want: "strong-pass-123",
		},
		{
			name: "trim trailing newline",
			in:   "strong-pass-123\n",
			want: "strong-pass-123",
		},
		{
			name: "trim trailing windows newline",
			in:   "strong-pass-123\r\n",
			want: "strong-pass-123",
		},
		{
			name: "preserve spaces",
			in:   "  strong-pass-123  ",
			want: "  strong-pass-123  ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeMasterPassword(tt.in)
			if got != tt.want {
				t.Fatalf("normalizeMasterPassword() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResetStoredKeys(t *testing.T) {
	configDir := t.TempDir()
	keyFilePath := filepath.Join(configDir, "keys.json")
	if err := os.WriteFile(keyFilePath, []byte(`{"sample":"data"}`), 0o600); err != nil {
		t.Fatalf("write test key file: %v", err)
	}

	svc := &Service{
		km: keymanager.NewKeyManager(keymanager.KeyManagerConfig{
			ConfigDir: configDir,
		}),
	}

	if err := svc.ResetStoredKeys(); err != nil {
		t.Fatalf("ResetStoredKeys() error = %v", err)
	}
	if _, err := os.Stat(keyFilePath); !os.IsNotExist(err) {
		t.Fatalf("expected key file removed, stat err = %v", err)
	}
}
