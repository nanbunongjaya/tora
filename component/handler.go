package component

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

var (
	typeError = reflect.TypeOf((*error)(nil)).Elem()
	typeBytes = reflect.TypeOf(([]byte)(nil))
)

type Handler struct {
	Type     reflect.Type   // function type
	Receiver reflect.Value  // receiver
	Method   reflect.Method // handler function
	IsRawArg bool           // whether to unserialize the data or not
}

func isExported(name string) bool {
	w, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(w)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// Check function name is uppercase OR package path is empty
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isHandlerMethod(method reflect.Method) bool {
	mt := method.Type

	// Method must be exported
	if method.PkgPath != "" {
		return false
	}

	// Method needs 2 ins: (receiver, []byte)
	if mt.NumIn() != 2 {
		return false
	}

	// Method needs 1 outs: (error)
	if mt.NumOut() != 1 {
		return false
	}

	if (mt.In(1) != typeBytes) || mt.Out(0) != typeError {
		return false
	}

	return true
}
