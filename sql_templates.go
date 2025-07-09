package tmplutil

import (
	"fmt"
	"regexp"
	"strings"
)

// SQLTemplateEngine handles SQL template operations
type SQLTemplateEngine struct {
	BaseTemplateEngine
}

func NewSQLTemplateEngine(templatePath string) *SQLTemplateEngine {
	return &SQLTemplateEngine{BaseTemplateEngine{templatePath: templatePath}}
}

func (e *SQLTemplateEngine) RenderQuery(section string, params map[string]interface{}) (string, []interface{}, error) {
	queryTemplate, err := e.LoadSection(section)
	if err != nil {
		return "", nil, err
	}

	// Extract ordered parameters
	orderedParams := e.extractParamOrder(queryTemplate)
	paramValues := make([]interface{}, len(orderedParams))
	for i, paramName := range orderedParams {
		if value, ok := params[paramName]; ok {
			paramValues[i] = value
		} else {
			return "", nil, fmt.Errorf("missing parameter: %s", paramName)
		}
	}

	// Render template
	query, err := e.RenderTemplate(queryTemplate, params)
	if err != nil {
		return "", nil, err
	}

	query = remoteSemicon(query)

	return query, paramValues, nil
}

func (e *SQLTemplateEngine) extractParamOrder(queryTemplate string) []string {
	queryTemplate = regexp.MustCompile(`'[^']*'`).ReplaceAllString(queryTemplate, "")

	queryTemplate = strings.NewReplacer(
		":YYYY", "",
		":MM", "",
		":DD", "",
		":HH24", "",
		":MI", "",
		":SS", "",
	).Replace(queryTemplate)

	re := regexp.MustCompile(`:\w+`)
	matches := re.FindAllString(queryTemplate, -1)

	paramOrder := make([]string, 0, len(matches))
	for _, match := range matches {
		paramOrder = append(paramOrder, match[1:]) // Remove the ':' prefix
	}
	return paramOrder
}

func remoteSemicon(query string) string {

	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")

	re := regexp.MustCompile(`\s*;\s*`)
	query = re.ReplaceAllString(query, " ")

	return query
}
