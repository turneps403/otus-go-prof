package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/turneps403/otus-go-prof/hw09_struct_validator/plugins"
)

const lookingField = "validate"

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	for _, e := range v {
		fmt.Fprintf(&sb, "field: %v, err: %v; ", e.Field, e.Err.Error())
	}
	return sb.String()
}

func Validate(iface interface{}) error {
	ift := reflect.TypeOf(iface)
	ifv := reflect.ValueOf(iface)
	var errs ValidationErrors

	for i := 0; i < ift.NumField(); i++ {
		fv := ifv.Field(i)

		switch fv.Kind() {
		case reflect.Struct:
			err := Validate(fv.Interface())
			if err != nil {
				if errors.Is(err, &ValidationErrors{}) {
					errs = append(errs, err.(ValidationErrors)...)
				} else {
					return err
				}
			}
		default:
			ft := ift.Field(i)
			tagDesc := ft.Tag.Get(lookingField)
			err := plugins.IsValid(tagDesc, fv)
			if err != nil {
				errs = append(errs, ValidationError{Field: ft.Name, Err: err})
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}
