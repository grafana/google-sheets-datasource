package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGooglesheets(t *testing.T) {
	t.Run("column id formatting", func(t *testing.T) {
		require.Equal(t, "A", GetExcelColumnName(1))
		require.Equal(t, "B", GetExcelColumnName(2))
		require.Equal(t, "AH", GetExcelColumnName(34))
		require.Equal(t, "BN", GetExcelColumnName(66))
		require.Equal(t, "ZW", GetExcelColumnName(699))
		require.Equal(t, "AJIL", GetExcelColumnName(24582))
	})
}
