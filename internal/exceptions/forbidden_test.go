package exceptions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForbiddenException(t *testing.T) {
	t.Run("should be an error", func(t *testing.T) {
		got := NewForbiddenException("error")

		assert.Equal(t, ForbiddenExceptionName, got.Name)
		assert.Equal(t, "error", got.Error())
	})

}
