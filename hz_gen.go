package main

import (
	"log"
	"os"
	"os/exec"
)

func Hzinstall() {
	// execute install hz
	cmd1 := exec.Command("go", "install", "github.com/cloudwego/hertz/cmd/hz@latest")
	err := cmd1.Run()
	if err != nil {
		log.Fatalf("install hertz failed with %s\n", err)
	}

	// Set the environment variable for the second command
	os.Setenv("GO111MODULE", "on")

	// execute install thriftgo
	cmd2 := exec.Command("go", "install", "github.com/cloudwego/thriftgo@latest")
	err = cmd2.Run()
	if err != nil {
		log.Fatalf("install thriftgo failed with %s\n", err)
	}
}

func Hzgen(name string) {

	//create new folder for hz
	err := os.MkdirAll("gateway", os.ModePerm)
	if err != nil {
		log.Fatalf("new folder failed with %s\n", err)
	}

	//move to directory
	desiredDir := "gateway"
	err = os.Chdir(desiredDir)
	if err != nil {
		log.Fatalf("move to folder failed with %s\n", err)
	}
	// run hz

	cmd1 := exec.Command("hz", "new", "-module", name, "-idl", "../idl/hello.thrift")
	err = cmd1.Run()
	if err != nil {
		log.Fatalf("hz gen failed with %s\n", err)
	}

	// execute install thriftgo
	cmd2 := exec.Command("go", "mod", "tidy")
	err = cmd2.Run()
	if err != nil {
		log.Fatalf("go mod tidy failed with %s\n", err)
	}
}
