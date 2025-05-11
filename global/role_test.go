package global

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRole(t *testing.T) {
	u, e := NewRole("user")
	assert.Equal(t, u, RoleUser)
	assert.Nil(t, e)
}
