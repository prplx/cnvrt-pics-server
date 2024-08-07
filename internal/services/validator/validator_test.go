package validator_test

import (
	"testing"

	"github.com/prplx/cnvrt/internal/services/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidator_NewValidator__should_be_valid_when_created(t *testing.T) {
	validator := validator.NewValidator()
	assert.Empty(t, validator.Errors)
	assert.True(t, validator.Valid())
}

func TestValidator_AddError__should_not_be_valid_when_called_with_falsy_argument(t *testing.T) {
	validator := validator.NewValidator()
	validator.Check(false, "key", "message")

	assert.False(t, validator.Valid())
	assert.Equal(t, "message", validator.Errors["key"])
	assert.Len(t, validator.Errors, 1)
}
