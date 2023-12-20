package helpers_test

import (
	"testing"

	"github.com/prplx/cnvrt/internal/helpers"
)

func TestFileNameWithoutExtension(t *testing.T) {
	testCases := []struct {
		name     string
		fileName string
		expected string
	}{
		{
			name:     "no extension",
			fileName: "foo",
			expected: "foo",
		},
		{
			name:     "single extension",
			fileName: "foo.txt",
			expected: "foo",
		},
		{
			name:     "multiple extensions",
			fileName: "foo.tar.gz",
			expected: "foo.tar",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := helpers.FileNameWithoutExtension(tc.fileName)
			if actual != tc.expected {
				t.Errorf("expected %q, but got %q", tc.expected, actual)
			}
		})
	}
}

func TestBuildPath(t *testing.T) {
	testCases := []struct {
		name     string
		parts    []interface{}
		expected string
	}{
		{
			name:     "empty parts",
			parts:    []interface{}{},
			expected: ".",
		},
		{
			name:     "single part",
			parts:    []interface{}{"foo"},
			expected: "foo",
		},
		{
			name:     "multiple parts",
			parts:    []interface{}{"foo", "bar", "baz"},
			expected: "foo/bar/baz",
		},
		{
			name:     "trailing slash",
			parts:    []interface{}{"foo", "bar", "baz", ""},
			expected: "foo/bar/baz",
		},
		{
			name:     "leading slash",
			parts:    []interface{}{"", "foo", "bar", "baz"},
			expected: "/foo/bar/baz",
		},
		{
			name:     "mixed types",
			parts:    []interface{}{"foo", 42, "bar", 123},
			expected: "foo/42/bar/123",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := helpers.BuildPath(tc.parts...)
			if actual != tc.expected {
				t.Errorf("expected %q, but got %q", tc.expected, actual)
			}
		})
	}
}
