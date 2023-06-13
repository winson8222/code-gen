package main

import (
	"log"
	"os/exec"
)

func Tidy(name string) {
	cmd := exec.Command("go", "mod", "tidy")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("go mod tidy failed with %s\n", err)
	}

	cmd2 := exec.Command("go", "build", name)
	err = cmd2.Run()
	if err != nil {
		log.Fatalf("making exe failed with %s\n", err)
	}

}
