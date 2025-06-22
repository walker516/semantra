package tmplutil

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"text/template"
)

// TemplateEngine defines the common interface for all template engines
type TemplateEngine interface {
	LoadSection(section string) (string, error)
	RenderTemplate(tmpl string, params map[string]interface{}) (string, error)
}

// BaseTemplateEngine provides common template functionality
type BaseTemplateEngine struct {
	templatePath string
}

func (e *BaseTemplateEngine) LoadSection(section string) (string, error) {
	data, err := os.ReadFile(e.templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %w", err)
	}

	re := regexp.MustCompile(fmt.Sprintf(`(?s)-- START: %s\b(.*?)-- END`, section))
	match := re.FindStringSubmatch(string(data))
	if len(match) < 2 {
		return "", fmt.Errorf("section %s not found in file", section)
	}

	return match[1], nil
}

func (e *BaseTemplateEngine) RenderTemplate(tmpl string, params map[string]interface{}) (string, error) {
	t, err := template.New("query").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
