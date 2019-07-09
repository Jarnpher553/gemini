package validator

import (
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v8"
	"reflect"
)

var v = binding.Validator.Engine().(*validator.Validate)

func Register(key string, fn Func) {
	_ = v.RegisterValidation(key, func(v *validator.Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool {
		return fn(&Validate{v}, topStruct, currentStruct, field, fieldtype, fieldKind, param)
	})
}

type Validate struct {
	*validator.Validate
}

type Func func(v *Validate, topStruct reflect.Value, currentStruct reflect.Value, field reflect.Value, fieldtype reflect.Type, fieldKind reflect.Kind, param string) bool
