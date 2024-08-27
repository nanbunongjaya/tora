package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/nanbunongjaya/tora/component"
)

type (
	data struct {
		Imports    []string
		Components []string
	}
)

func Compile(comps *component.Components) (string, error) {
	const mainPkg = "main"

	data := &data{}
	imports := make(map[string]bool)

	for _, component := range comps.List() {
		compType := reflect.TypeOf(component)
		pkgPath := compType.Elem().PkgPath()
		structTypeName := removeMainPkgPrefix(fmt.Sprintf("%v", compType))

		// Ignore "main" package
		if pkgPath != mainPkg {
			imports[pkgPath] = true
		}

		// Sprint "&xxx.XXX{}"
		instanceCreation := fmt.Sprintf("%s{}", prefixStarToPrefixAmpersand(structTypeName))

		data.Components = append(data.Components, instanceCreation)
	}

	for pkg := range imports {
		data.Imports = append(data.Imports, pkg)
	}

	tmpl, err := template.New("Program").Parse(programTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var output strings.Builder
	if err := tmpl.Execute(&output, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
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

	err = os.MkdirAll(filepath.Dir(absolutePath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(absolutePath, []byte(program), 0644)
	if err != nil {
		return err
	}

	return nil
}

// Convert "*main.XXX" to "*XXX"
func removeMainPkgPrefix(s string) string {
	if strings.HasPrefix(s, "*main.") {
		return "*" + s[len("*main."):]
	}
	return s
}

// Convert "*xxx.XXX" to "&xxx.XXX"
func prefixStarToPrefixAmpersand(s string) string {
	if strings.HasPrefix(s, "*") {
		return "&" + s[1:]
	}
	return s
}
