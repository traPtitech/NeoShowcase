package tarfs

import "testing"

func Test_isValidRelPath(t *testing.T) {
	tests := []struct {
		name    string
		relPath string
		want    bool
	}{
		{
			name:    "valid 1",
			relPath: "./",
			want:    true,
		},
		{
			name:    "valid 2",
			relPath: "./foo/bar",
			want:    true,
		},
		{
			name:    "valid 3",
			relPath: ".",
			want:    true,
		},
		{
			name:    "valid 4",
			relPath: "./node_modules/some-module/test/../link.js",
			want:    true,
		},
		{
			name:    "valid 5",
			relPath: "./[...weird-name].txt",
			want:    true,
		},
		{
			name:    "valid 6",
			relPath: "./weird../name",
			want:    true,
		},
		{
			name:    "invalid 1",
			relPath: "../",
			want:    false,
		},
		{
			name:    "invalid 2 (same name prefix)",
			relPath: "../pathprefixed",
			want:    false,
		},
		{
			name:    "invalid 3",
			relPath: "./test/../../foo",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidRelPath(tt.relPath); got != tt.want {
				t.Errorf("isValidRelPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
