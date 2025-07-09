package nullutil_test

import (
	"testing"

	"api/pkg/nullutil"

	"github.com/stretchr/testify/require"
)

func TestJSONBMap_ScanAndValue(t *testing.T) {
	j := nullutil.JSONBMap{}

	// Simulate DB returns JSONB bytes
	err := j.Scan([]byte(`{"foo":"bar"}`))
	require.NoError(t, err)
	require.Equal(t, "bar", j["foo"])

	val, err := j.Value()
	require.NoError(t, err)
	require.Contains(t, string(val.([]byte)), "foo")
}
