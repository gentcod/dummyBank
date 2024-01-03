package api

import (
	"github.com/gentcod/DummyBank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSuppoertedCurrency(currency)
	}
	return false
}