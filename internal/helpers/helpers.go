package helpers

import (
	"fmt"
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
