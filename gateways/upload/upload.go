package upload

import (
	"mime/multipart"
	"path/filepath"
	"strings"

	"podcast/hasher"
)

type FileHandler interface {
	Upload(file *multipart.FileHeader) (string, error)
}

func randomizeFileName(filename string) string {
	f := strings.Split(filename, ".")
	name, _ := hasher.GenerateSecureToken(6)
	dst := filepath.Join("tmp", f[0]+"-"+name+"."+f[1])

	return dst
}
