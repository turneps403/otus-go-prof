package plugins

import (
	"fmt"
	"reflect"
	"strconv"
)

func init() {
	extendPlugins(
		"len",
		func(desc string, rv reflect.Value) error {
			num, err := strconv.Atoi(desc)
			if err != nil {
				return err
			}

			switch rv.Kind() {
			case reflect.Slice, reflect.String:
				if rv.Len() != num {
					return fmt.Errorf("expected %d but got %d", num, rv.Len())
				}
			}

			return nil
		},
	)
}
