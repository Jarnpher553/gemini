package validator

import (
	"github.com/Jarnpher553/micro-core/json"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"regexp"
	"sync"
	"time"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

var v *validator.Validate

func init() {
	binding.Validator = new(defaultValidator)
	v = binding.Validator.Engine().(*validator.Validate)
	v.RegisterCustomTypeFunc(func(field reflect.Value) interface{} {

		if valuer, ok := field.Interface().(json.Date); ok {

			if valuer != json.Date(time.Time{}) {
				return valuer
			}
		}

		return nil
	}, json.Date{})
}

func Register(key string, fn Func) {
	_ = v.RegisterValidation(key, func(fl validator.FieldLevel) bool {
		return fn(&Validate{v}, fl)
	})
}

type Validate struct {
	*validator.Validate
}

type FieldLevel interface {
	validator.FieldLevel
}

type Func func(v *Validate, fl FieldLevel) bool

func jsonDateRequired(
	v *validator.Validate, level validator.FieldLevel,
) bool {
	if code, ok := level.Field().Interface().(string); ok {
		reg := regexp.MustCompile(`^20\d{2}(0[1-9]|1[0-2])[0-9A-Z]{4}$`)
		ret := reg.MatchString(code)
		return ret
	}
	return true
}
