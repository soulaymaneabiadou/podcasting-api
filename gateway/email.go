package gateway

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type EmailPayload struct {
	Receiver string
	Subject  string
	Template string
	Data     any
}

type EmailGateway interface {
	SendEmail(EmailPayload)
}

func parseTemplate(tpl string, data any) bytes.Buffer {
	t, err := parseTemplateDir("templates")
	if err != nil {
		log.Fatal("could not parse the templates dir, err: ", err)
	}

	var body bytes.Buffer
	err = t.ExecuteTemplate(&body, tpl, &data)
	if err != nil {
		log.Fatal("could not parse the template, err: ", err)
	}

	return body
}

func parseTemplateDir(dir string) (*template.Template, error) {
	var paths []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return template.ParseFiles(paths...)
}
