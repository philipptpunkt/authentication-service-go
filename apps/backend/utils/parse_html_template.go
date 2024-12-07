package utils

import (
	"bytes"
	"html/template"
)

func ParseHtmlTemplate(templatePath string, data map[string]interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
