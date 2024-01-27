package exceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundException(t *testing.T) {
	t.Run("should be an error", func(t *testing.T) {
		got := NewNotFoundException("error")

		assert.Equal(t, NotFoundExceptionName, got.Name)
		assert.Equal(t, "error", got.Error())
	})

}
