package plugins

import (
	"fmt"
	"reflect"
	"strconv"
)

func init() {
	extendPlugins(
		"min",
		func(desc string, rv reflect.Value) error {
			num, err := strconv.Atoi(desc)
			if err != nil {
				return err
			}

			switch rv.Kind() {
			case reflect.Slice, reflect.String:
				if rv.Len() < num {
					return fmt.Errorf("expected %d but got %d", num, rv.Len())
				}
			case reflect.Int:
				val := rv.Int()
				if val < int64(num) {
					return fmt.Errorf("expected %d but got %d", num, val)
				}
			case reflect.Float64:
				val := rv.Float()
				if val < float64(num) {
					return fmt.Errorf("expected %d but got %f", num, val)
				}
			}

			return nil
		},
	)
}
