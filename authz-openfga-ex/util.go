package main

import "reflect"

// IsNotNil - Returns `true` if the passed in value is not nil
func IsNotNil[T any](x T) bool {
	return !IsNil(x)
}

// IsNil - returns `true` if the value is nil
func IsNil[T any](x T) bool {
	v := reflect.ValueOf(x)
	k := v.Kind()
	switch k {
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		return v.IsNil()
	case reflect.Invalid: // naked nil
		return true
	default:
		return false
	}
}
