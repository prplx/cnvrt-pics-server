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

func TestGetSessionCookieDomain(t *testing.T) {
	tests := []struct {
		name     string
		URL      string
		expected string
	}{
		{
			name:     "URL with subdomain",
			URL:      "https://sub.example.com",
			expected: ".example.com",
		},
		{
			name:     "URL without subdomain",
			URL:      "http://example.com",
			expected: "example.com",
		},
		{
			name:     "URL with multiple subdomains",
			URL:      "http://sub.sub.example.com",
			expected: ".sub.example.com",
		},
		{
			name:     "URL with single part",
			URL:      "http://localhost",
			expected: "localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := helpers.GetSessionCookieDomain(tt.URL)
			if result != tt.expected {
				t.Errorf("GetSessionCookieDomain(%s) = %s; expected %s", tt.URL, result, tt.expected)
			}
		})
	}
}
