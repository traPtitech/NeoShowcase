package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    []string
		wantErr bool
	}{
		{
			name:    "returns empty if empty",
			in:      "",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "simple one command",
			in:      "npm run start",
			want:    []string{"npm", "run", "start"},
			wantErr: false,
		},
		{
			name:    "simple one command 2",
			in:      "./main",
			want:    []string{"./main"},
			wantErr: false,
		},
		{
			name:    "simple one command with extra space",
			in:      "npm run  start  ",
			want:    []string{"npm", "run", "start"},
			wantErr: false,
		},
		{
			name:    "simple one command with quoting",
			in:      "npm run start \"hello world\"",
			want:    []string{"npm", "run", "start", "hello world"},
			wantErr: false,
		},
		{
			name:    "has env shell syntax (current limitation, cannot recognize envs)",
			in:      "NODE_ENV=production npm run start",
			want:    []string{"NODE_ENV=production", "npm", "run", "start"},
			wantErr: false,
		},
		{
			name:    "has multi command shell syntax",
			in:      "npm run build && npm run start",
			want:    []string{"sh", "-c", "npm run build && npm run start"},
			wantErr: false,
		},
		{
			name:    "invalid shell line",
			in:      "hello world `",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArgs(tt.in)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equalf(t, tt.want, got, "ParseArgs(%v)", tt.in)
		})
	}
}
