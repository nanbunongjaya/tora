package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"tora/component"
)

type (
	data struct {
		Imports    []string
		Components []string
	}
)

func Compile(comps *component.Components) (string, error) {
	data := &data{}
	imports := make(map[string]bool)
	for _, comp := range comps.List() {
		t := reflect.TypeOf(comp.Comp)

		pkgPath := t.Elem().PkgPath()
		structPtrName := removeMainPkgPrefix(fmt.Sprintf("%v", t))

		if pkgPath != "main" {
			imports[pkgPath] = true
		}

		data.Components = append(data.Components, fmt.Sprintf("(%s{}).New().(%s)", prefixStarToPrefixAmpersand(structPtrName), structPtrName))
	}

	for imp := range imports {
		data.Imports = append(data.Imports, imp)
	}

	tmpl, err := template.New("Program").Parse(programTemplate)
	if err != nil {
		return "", err
	}

	var output strings.Builder
	err = tmpl.Execute(&output, data)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func Output(path string, program string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Convert relative path to absolute path
	absolutePath := filepath.Join(currentDir, path)

	if err := os.MkdirAll(filepath.Dir(absolutePath), 0755); err != nil {
		return err
	}

	err = os.WriteFile(absolutePath, []byte(program), 0644)
	if err != nil {
		return err
	}

	return nil
}

func removeMainPkgPrefix(s string) string {
	if strings.HasPrefix(s, "*main.") {
		return "*" + s[len("*main."):]
	}
	return s
}

func prefixStarToPrefixAmpersand(s string) string {
	if strings.HasPrefix(s, "*") {
		return "&" + s[1:]
	}
	return s
}
