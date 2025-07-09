package tmplutil_test

import (
	"api/pkg/tmplutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadSection_Success(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sql")
	content := `
-- START: Foo
SELECT * FROM foo;
-- END
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	engine := tmplutil.NewBaseTemplateEngine(path)

	section, err := engine.LoadSection("Foo")
	require.NoError(t, err)
	require.Contains(t, section, "SELECT")
}

func TestLoadSection_NotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.sql")
	content := `
-- START: Bar
SELECT * FROM bar;
-- END
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	engine := tmplutil.NewBaseTemplateEngine(path)

	_, err := engine.LoadSection("NonExistent")
	require.ErrorContains(t, err, "section NonExistent not found")
}

func TestLoadSection_FileNotFound(t *testing.T) {
	engine := tmplutil.NewBaseTemplateEngine("no_such_file.sql")
	_, err := engine.LoadSection("Any")
	require.ErrorContains(t, err, "failed to read template file")
}

func TestRenderTemplate_Success(t *testing.T) {
	engine := tmplutil.BaseTemplateEngine{}
	tmpl := "SELECT * FROM users WHERE id = {{.ID}} AND name = '{{.Name}}';"

	out, err := engine.RenderTemplate(tmpl, map[string]interface{}{
		"ID":   1,
		"Name": "Bob",
	})
	require.NoError(t, err)
	require.Contains(t, out, "Bob")
	require.Contains(t, out, "1")
}

func TestRenderTemplate_SyntaxError(t *testing.T) {
	engine := tmplutil.BaseTemplateEngine{}
	tmpl := "SELECT * FROM users WHERE id = {{.ID" // broken template

	_, err := engine.RenderTemplate(tmpl, map[string]interface{}{
		"ID": 1,
	})
	require.ErrorContains(t, err, "failed to parse template")
}

func TestRenderTemplate_ExecError(t *testing.T) {
	engine := tmplutil.BaseTemplateEngine{}
	tmpl := "SELECT * FROM users WHERE id = {{.ID}} AND name = '{{.Name}}';"

	_, err := engine.RenderTemplate(tmpl, map[string]interface{}{
		"ID": 1,
		// Name is missing
	})
	require.ErrorContains(t, err, "failed to execute template")
}
