package compiler

import (
	"log"
	"os/exec"
)

func Build(path string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "controller.so", path)

	var output []byte
	var err error

	if output, err = cmd.CombinedOutput(); err != nil {
		log.Printf("Output: %s", output)
		return err
	}

	log.Println(string(output))

	return nil
}
