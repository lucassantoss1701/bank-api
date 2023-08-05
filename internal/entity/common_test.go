package entity_test

import (
	"lucassantoss1701/bank/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommon_NewUUID(t *testing.T) {
	t.Run("Testing NewUUID with success", func(t *testing.T) {
		uuid := entity.NewUUID()
		assert.NotNil(t, uuid)
		assert.NotEmpty(t, uuid)

	})
}
