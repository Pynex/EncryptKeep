package appsvc

import "testing"

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
