package main

import (
	"log"
	"os"
	"strings"
	"text/template"
)

type Constants struct {
	GatewayURL         string
	FilepathToService  string
	EtcdURL            string
	ServiceName        string
	ServiceUpstreamURL string
	Methods            []Method
	IDLName            string
	GatewayName        string
}

type Method struct {
	MethodName string
	MethodType string
}

func CreateConstant(constants Constants) {
	// Define the values for the constants

	// Define the template string
	templateString :=
		`package constants

	 import (
		"strings"
	 )
 
	func ToConstant(s string) string {
		return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
	}

	// info about API Gateway
	const (
		GATEWAY_URL              = "{{ .GatewayURL }}"
		FILEPATH_TO_{{ .ServiceName | ToConstant }} = "{{ .FilepathToService }}" //relative to root directory of module!!
	)
	
	// info about etcd instance
	const (
		ETCD_URL = "{{ .EtcdURL }}" //connects to a single etcd instance in the cluster
	)
	
	// info about the service
	const (
		{{ .ServiceName | ToConstant }}_NAME         = "{{ .ServiceName }}" //name registered with svc registry as rpcendpoint
		{{ .ServiceName | ToConstant }}_UPSTREAM_URL = "{{ .ServiceUpstreamURL }}"
	)
	
	// ToConstant function converts the input string to a constant format
	`

	// Create a new template
	tmpl := template.Must(template.New("constants").Funcs(template.FuncMap{"ToConstant": ToConstant}).Parse(templateString))

	// Create the output file
	err := os.MkdirAll("constants", os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output folder: %s\n", err)
	}
	outputFile, err := os.Create("constants/constants.go")
	if err != nil {
		log.Fatalf("Error creating output file: %s\n", err)
	}
	defer outputFile.Close()

	// Execute the template with the constants and write the output to the file
	err = tmpl.Execute(outputFile, constants)
	if err != nil {
		log.Fatalf("Error executing template: %s\n", err)
	}

	log.Println("Template generation completed successfully.")
}

func ToConstant(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}
