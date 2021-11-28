package plugins

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	extendPlugins(
		"in",
		func(desc string, rv reflect.Value) error {
			vals := make(map[string]interface{})
			for _, v := range strings.Split(desc, ",") {
				vals[v] = nil
			}

			switch rv.Kind() { // nolint:exhaustive
			case reflect.String:
				if _, ok := vals[rv.String()]; !ok {
					return fmt.Errorf("val %s outside of possible values %v", rv.String(), vals)
				}
			case reflect.Int:
				if _, ok := vals[strconv.Itoa(int(rv.Int()))]; !ok {
					return fmt.Errorf("val %s outside of possible values %v", rv.String(), vals)
				}
			}

			return nil
		},
	)
}
