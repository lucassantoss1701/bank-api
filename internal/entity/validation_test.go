package entity_test

import (
	"lucassantoss1701/bank/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationError_Add(t *testing.T) {
	validationError := &entity.ValidationError{}

	assert.Empty(t, validationError.Messages)

	validationError.Add("Error 1")
	assert.Equal(t, []string{"Error 1"}, validationError.Messages)

	validationError.Add("Error 2")
	validationError.Add("Error 3")
	assert.Equal(t, []string{"Error 1", "Error 2", "Error 3"}, validationError.Messages)
}

func TestValidationError_Error(t *testing.T) {
	validationError := &entity.ValidationError{}
	assert.Empty(t, validationError.Error())

	validationError.Add("Error 1")
	assert.Equal(t, "Error 1", validationError.Error())

	validationError.Add("Error 2")
	validationError.Add("Error 3")
	assert.Equal(t, "Error 1, Error 2, Error 3", validationError.Error())
}
