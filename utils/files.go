package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

const MB = 1048576

func checkFileType(f multipart.FileHeader, allowedTypes []string) error {
	fileType := strings.Split(f.Filename, ".")[1]

	exists := slices.Contains(allowedTypes, fileType)
	if !exists {
		return errors.New(fmt.Sprintf("type %s is not allowed", fileType))
	}

	return nil
}

func HandleFileValidation(c *gin.Context, f string, types []string, maxMb int64) (*multipart.FileHeader, error) {
	file, err := c.FormFile(f)
	if err != nil {
		return nil, err
	}

	if err := checkFileType(*file, types); err != nil {
		return nil, err
	}

	if file.Size > maxMb*MB {
		return nil, err
	}

	return file, nil
}
