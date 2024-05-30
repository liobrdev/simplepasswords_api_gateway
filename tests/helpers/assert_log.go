package helpers

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

func AssertLog(t *testing.T, expected *models.Log, actual *models.Log) {
	require.Equal(t, expected.ClientIP, actual.ClientIP)
	require.Equal(t, expected.ClientOperation, actual.ClientOperation)
	require.Equal(t, expected.Detail, actual.Detail)
	require.Equal(t, expected.Extra, actual.Extra)
	require.Equal(t, expected.Level, actual.Level)
	require.Equal(t, expected.Message, actual.Message)
	require.Equal(t, expected.RequestBody, actual.RequestBody)
	require.Equal(t, expected.UserSlug, actual.UserSlug)
}
