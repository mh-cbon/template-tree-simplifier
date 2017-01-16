package funcmap

import (
	"fmt"
	"reflect"
	"strings"
)

//BrowsePropertyPath browse a property path ("b.c.d") on some value
func BrowsePropertyPath(some interface{}, propertypath string, args ...interface{}) interface{} {
	to := strings.Split(propertypath, ".")
	var ret interface{}
	v := reflect.ValueOf(some)
	v = indirect(v)
	for i := 0; i < len(to); i++ {
		if v.IsValid() {
			nv := v.FieldByName(to[i])
			if !v.IsValid() && !v.IsNil() {
				nv = v.MethodByName(to[i])
				if !v.IsValid() {
					err := fmt.Sprintf(
						"Field/Method %q not found at %q in value of type %v",
						strings.Join(to, "."),
						strings.Join(to[:i], "."),
						reflect.ValueOf(some),
					)
					panic(err)
				}
				// must be the last part
				reflectArgs := []reflect.Value{}
				for _, a := range args {
					reflectArgs = append(reflectArgs, reflect.ValueOf(a))
				}
				return nv.Call(reflectArgs)
			}
			v = indirect(nv)
		}
	}
	return ret
}

func indirect(some reflect.Value) reflect.Value {
	if some.Kind() == reflect.Ptr {
		some = some.Elem()
	}
	if some.Kind() == reflect.Interface {
		some = some.Elem()
	}
	return some
}
