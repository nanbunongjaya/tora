package compiler

import (
	"log"
	"os/exec"
)

func Build(sofile, gofile string) error {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", sofile, gofile)

	var output []byte
	var err error

	if output, err = cmd.CombinedOutput(); err != nil {
		log.Printf("Output: %s", output)
		return err
	}

	log.Println("Build plugin ...success")

	return nil
}
