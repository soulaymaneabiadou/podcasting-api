package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"podcast/hasher"
	"strings"
)

type LocalFileHandler struct {
}

func NewLocalFileHandler() *LocalFileHandler {
	return &LocalFileHandler{}
}

func (s *LocalFileHandler) Upload(file *multipart.FileHeader) (string, error) {
	f := strings.Split(file.Filename, ".")
	name, _ := hasher.GenerateSecureToken(6)
	dst := filepath.Join("tmp", f[0]+"-"+name+"."+f[1])

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return "", err
	}

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return fmt.Sprintf("%s/%s", os.Getenv("SERVER_URL"), dst), nil
}
