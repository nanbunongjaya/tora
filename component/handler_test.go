package component

import (
	"reflect"
	"testing"
)

func TestIsExported(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"A", true},
		{"a", false},
		{"Test", true},
		{"test", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test
			exported := isExported(tt.name)

			// Compare
			if exported != tt.expected {
				t.Errorf("Got %v, expected %v", exported, tt.expected)
			}
		})
	}
}

func TestIsExportedOrBuiltinType(t *testing.T) {
	tests := []struct {
		typ      reflect.Type
		expected bool
	}{
		{
			// int - built-in type
			typ:      reflect.TypeOf(1),
			expected: true,
		},
		{
			// string - built-in type
			typ:      reflect.TypeOf("string"),
			expected: true,
		},
		{
			// &test{} - unexported struct type
			typ:      reflect.TypeOf(&test{}),
			expected: false,
		},
		{
			// &data{} - unexported struct type
			typ:      reflect.TypeOf(&data{}),
			expected: false,
		},
		{
			// &Components{} - exported struct type
			typ:      reflect.TypeOf(&Components{}),
			expected: true,
		},
		{
			// []byte{} - built-in type
			typ:      reflect.TypeOf([]byte{}),
			expected: true,
		},
		{
			// map[string]int{} - built-in type
			typ:      reflect.TypeOf(map[string]int{}),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.typ.String(), func(t *testing.T) {
			// Test
			result := isExportedOrBuiltinType(tt.typ)

			// Compare
			if result != tt.expected {
				t.Errorf("Got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestIsHandlerMethod(t *testing.T) {
	tests := []struct {
		method reflect.Method
		valid  bool
	}{
		{
			// A - Valid
			method: reflect.TypeOf(&test{}).Method(0),
			valid:  true,
		},
		{
			// B - Invalid
			method: reflect.TypeOf(&test{}).Method(1),
			valid:  false,
		},
		{
			// C - Invalid
			method: reflect.TypeOf(&test{}).Method(2),
			valid:  false,
		},
		{
			// D - Valid
			method: reflect.TypeOf(&test{}).Method(3),
			valid:  false,
		},
		{
			// E - Invalid
			method: reflect.TypeOf(&test{}).Method(4),
			valid:  false,
		},
		{
			// F - Invalid
			method: reflect.TypeOf(&test{}).Method(5),
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.method.Name, func(t *testing.T) {
			// Test
			valid := isHandlerMethod(tt.method)

			// Compare
			if valid != tt.valid {
				t.Errorf("Got %v, expected %v", valid, tt.valid)
			}
		})
	}
}

// Test struct
type test struct{ Base }
type data struct{}

func (t *test) A(data []byte) error {
	return nil
}

// In[1] with non-bytes type (invalid)
func (t *test) B(data string) error {
	return nil
}

// Out[0] with non-error type (invalid)
func (t *test) C(data []byte) string {
	return ""
}

// In[1] with non-bytes type (invalid)
func (t *test) D(data *data) error {
	return nil
}

// In[1] with non-bytes type (invalid)
func (t *test) E(data data) error {
	return nil
}

// In[1] not exist (invalid)
func (t *test) F() error {
	return nil
}
