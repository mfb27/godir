package pathutil

import (
	"path"
	"strings"
)

func Ext(filename string) string {
	return strings.ToLower(path.Ext(filename))
}

func BaseWithoutExt(filename string) string {
	return strings.TrimPrefix(filename, Ext(filename))
}
