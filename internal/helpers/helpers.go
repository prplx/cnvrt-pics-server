package helpers

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func BuildPath(parts ...interface{}) string {
	var path strings.Builder
	for _, part := range parts {
		path.WriteString(fmt.Sprint(part) + "/")
	}

	return filepath.Clean(path.String())
}

func FileExtension(fileName string) string {
	return filepath.Ext(fileName)
}

func IsTest() bool {
	return os.Getenv("ENV") == "test"
}

func MustGetHostnameFromURL(input string) string {
	url, err := url.Parse(input)
	if err != nil {
		panic(err)
	}
	return strings.TrimPrefix(url.Hostname(), "www.")
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func GetMapKeys[T comparable](m map[string][]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
