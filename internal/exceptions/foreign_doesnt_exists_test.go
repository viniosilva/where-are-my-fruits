package exceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForeignDoesntExistsException(t *testing.T) {
	t.Run("should be an error", func(t *testing.T) {
		got := NewForeignDoesntExistsException("error")

		assert.Equal(t, ForeignDoesntExistsExceptionName, got.Name)
		assert.Equal(t, "error", got.Error())
	})

}
