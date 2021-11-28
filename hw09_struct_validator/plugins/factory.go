// My custom version of factory method
package plugins

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type validateFunc func(tagDesc string, rv reflect.Value) error

var plRegister = make(map[string]validateFunc)

func extendPlugins(tag string, f validateFunc) {
	if _, ok := plRegister[tag]; ok {
		panic("Double definition for validator: " + tag)
	}
	plRegister[tag] = f
}

func IsValid(tagDesc string, rv reflect.Value) error {
	for _, subTag := range strings.Split(tagDesc, "|") {
		if len(subTag) == 0 {
			continue
		}
		nv := regexp.MustCompile(":").Split(subTag, 2)
		if len(nv) == 0 {
			return fmt.Errorf("broken description")
		} else if len(nv) == 1 {
			nv = append(nv, "")
		}
		for rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if f, ok := plRegister[nv[0]]; ok {
			err := f(nv[1], rv)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("no plugin for '%v'", nv[0])
		}
	}
	return nil
}
