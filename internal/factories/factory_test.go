package factories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/viniosilva/where-are-my-fruits/internal/infra"
)

func TestFactory_Build(t *testing.T) {
	t.Run("should be successful", func(t *testing.T) {
		// when
		got, err := Build(&infra.Database{}, nil, nil)
		require.Nil(t, err)

		// then
		assert.NotNil(t, got)
	})
}
