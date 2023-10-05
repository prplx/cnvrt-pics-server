package helpers

import (
	"path/filepath"
	"strings"
)

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
