package plugins

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

func init() {
	// TODO: regexp will fail in run time if we have a bad expression
	// panic in runtime possible
	extendPlugins(
		"regexp",
		func(desc string, rv reflect.Value) error {
			re := regexp.MustCompile(desc)

			switch rv.Kind() { // nolint:exhaustive
			case reflect.String:
				if !re.MatchString(rv.String()) {
					return fmt.Errorf("val %s doesnt match for a regexp '%s'", rv.String(), desc)
				}
			case reflect.Int:
				if !re.MatchString(strconv.Itoa(int(rv.Int()))) {
					return fmt.Errorf("val %s outside of possible values %v", rv.String(), desc)
				}
			}

			return nil
		},
	)
}
