package semconview

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalyzeSemconvDependencies(t *testing.T) {
	ctx := context.Background()
	t.Run("Test1", func(t *testing.T) {
		got, err := AnalyzeSemconvDependencies(ctx, []string{"../../example/test1/*.go"})
		require.NoError(t, err)
		assert.Len(t, got.Attributes, 4)
		assert.Contains(t, got.Attributes, SemconvAttribute{
			Key:     "http.status_code",
			Version: "v1.20.0",
		})
		assert.Contains(t, got.Attributes, SemconvAttribute{
			Key:     "http.method",
			Version: "v1.20.0",
		})
		assert.Contains(t, got.Attributes, SemconvAttribute{
			Key:     "http.scheme",
			Version: "v1.20.0",
		})
		assert.Contains(t, got.Attributes, SemconvAttribute{
			Key:     "user_agent.original",
			Version: "v1.30.0",
		})
	})
}
