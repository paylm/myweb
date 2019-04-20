package utils

import (
	"errors"
	"log"
	"os/exec"
)

func convert(file string, size string, done chan error) {
	if len(file) == 0 {
		done <- errors.New("filename can't be null")
		return
	}
	cmd := exec.Command("/bin/bash", "../../script/convert.sh", file, size)
	log.Printf("Running command and waiting for it to finish...")
	err := cmd.Run()
	log.Printf("Command finished with error: %v", err)
	done <- err
}
