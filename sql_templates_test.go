package tmplutil_test

import (
	"api/pkg/tmplutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderQuery_ParamOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "queries.sql")
	content := `
-- START: FindByID
SELECT * FROM users WHERE id = :UserID AND name = :Name AND created_at >= :DateFrom;
-- END
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	engine := tmplutil.NewSQLTemplateEngine(path)

	query, params, err := engine.RenderQuery("FindByID", map[string]interface{}{
		"UserID":   1,
		"Name":     "Alice",
		"DateFrom": "2023-01-01",
	})
	require.NoError(t, err)
	require.Contains(t, query, "FROM users")

	require.Equal(t, []interface{}{1, "Alice", "2023-01-01"}, params)
}

func TestRenderQuery_MissingParam(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "queries.sql")
	content := `
-- START: Test
SELECT * FROM t WHERE id = :ID AND name = :Name;
-- END
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	engine := tmplutil.NewSQLTemplateEngine(path)

	_, _, err := engine.RenderQuery("Test", map[string]interface{}{
		"ID": 10,
	})
	require.ErrorContains(t, err, "missing parameter: Name")
}

func TestRenderQuery_OracleDatePartsExcluded(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "queries.sql")
	content := `
-- START: OracleDate
SELECT TO_CHAR(sysdate, 'YYYY:MM:DD:HH24:MI:SS') FROM dual WHERE id = :ID;
-- END
`
	require.NoError(t, os.WriteFile(path, []byte(content), 0644))

	engine := tmplutil.NewSQLTemplateEngine(path)

	_, params, err := engine.RenderQuery("OracleDate", map[string]interface{}{
		"ID": 42,
	})
	require.NoError(t, err)
	require.Equal(t, []interface{}{42}, params)
}
