package exceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForeignNotExistsException(t *testing.T) {
	t.Run("should be an error", func(t *testing.T) {
		got := NewForeignNotFoundException("error")

		assert.Equal(t, ForeignNotFoundExceptionName, got.Name)
		assert.Equal(t, "error", got.Error())
	})

}
