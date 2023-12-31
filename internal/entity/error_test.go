package entity_test

import (
	"lucassantoss1701/bank/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorHandler_Add(t *testing.T) {
	errorHandler := &entity.ErrorHandler{}

	assert.Empty(t, errorHandler.Messages)

	errorHandler.Add("Error 1")
	assert.Equal(t, []string{"Error 1"}, errorHandler.Messages)

	errorHandler.Add("Error 2")
	errorHandler.Add("Error 3")
	assert.Equal(t, []string{"Error 1", "Error 2", "Error 3"}, errorHandler.Messages)
}

func TestErrorHandlerError_Error(t *testing.T) {
	errorHandler := &entity.ErrorHandler{}
	assert.Empty(t, errorHandler.Error())

	errorHandler.Add("Error 1")
	assert.Equal(t, "Error 1", errorHandler.Error())

	errorHandler.Add("Error 2")
	errorHandler.Add("Error 3")
	assert.Equal(t, "Error 1, Error 2, Error 3", errorHandler.Error())
}

func TestErrorHandlerError_GetTypeError(t *testing.T) {
	errorHandler := &entity.ErrorHandler{}
	errorHandler.TypeError = entity.ENTITY_ERROR
	assert.Empty(t, errorHandler.Error())

	assert.Equal(t, entity.ENTITY_ERROR, errorHandler.GetTypeError())
}

func TestErrorHandlerError_NewErrorHandler(t *testing.T) {
	errorHandler := entity.NewErrorHandler(entity.ENTITY_ERROR)
	assert.Empty(t, errorHandler.Error())

	assert.Equal(t, entity.ENTITY_ERROR, errorHandler.GetTypeError())
}
