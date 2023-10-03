package upload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalFileHandler struct {
}

func NewLocalFileHandler() *LocalFileHandler {
	return &LocalFileHandler{}
}

func (s *LocalFileHandler) Upload(file *multipart.FileHeader) (string, error) {
	dst := randomizeFileName(file.Filename)

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
