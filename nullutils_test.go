package nullutil_test

import (
	"database/sql"
	"testing"

	"api/pkg/nullutil"

	"github.com/stretchr/testify/require"
)

func TestNullStringToString(t *testing.T) {
	require.Equal(t, "", nullutil.NullStringToString(sql.NullString{}))
	require.Equal(t, "abc", nullutil.NullStringToString(sql.NullString{String: "abc", Valid: true}))
}

func TestNullInt64ToInt64(t *testing.T) {
	require.Equal(t, int64(0), nullutil.NullInt64ToInt64(sql.NullInt64{}))
	require.Equal(t, int64(42), nullutil.NullInt64ToInt64(sql.NullInt64{Int64: 42, Valid: true}))
}

func TestNullInt64ToPointer(t *testing.T) {
	require.Nil(t, nullutil.NullInt64ToPointer(sql.NullInt64{}))
	ptr := nullutil.NullInt64ToPointer(sql.NullInt64{Int64: 55, Valid: true})
	require.NotNil(t, ptr)
	require.Equal(t, 55, *ptr)
}
