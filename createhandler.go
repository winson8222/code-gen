package main

import (
	"log"
	"os"
	"text/template"
)

func Createhandler(service Service) {
	headertamplate := `// Code generated by hertz generator.

	package {{ .IDLName }}
	
	import (
		"context"
	
		"{{ .GatewayName }}/biz/model/{{ .IDLName }}"
	
		"github.com/cloudwego/hertz/pkg/app"
		"github.com/cloudwego/hertz/pkg/protocol/consts"
	)
	`

	methodtemplate := `
	func {{ .MethodName }}(ctx context.Context, c *app.RequestContext) {
		var req {{ .IDLName }}.{{ .RequestStruct }}
		err := c.BindAndValidate(&req)
		if err != nil {
			c.String(consts.StatusBadRequest, err.Error())
			return
		}
	
		cli := {{ .ServiceName }}GenericClient()
		resp, err := Do{{ .MethodName }}(cli, c) // Pass the generic client and requestContext
		if err != nil {
			c.String(consts.StatusInternalServerError, "Failed to perform generic call")
			return
		}
	
		c.JSON(consts.StatusOK, resp)
	}
	`

	// Create the output file
	outputFile, err := os.Create("biz/handler/" + service.IDLName + "/" + service.HandlerFile + "_service.go")
	if err != nil {
		log.Fatalf("Error creating output file: %s\n", err)
	}
	defer outputFile.Close()

	// Create a new template for the generic client
	headerTmpl := template.Must(template.New("header").Parse(headertamplate))

	err = headerTmpl.Execute(outputFile, service)
	if err != nil {
		log.Fatalf("Error executing generic client template: %s\n", err)
	}

	// Create a new template for the method
	methodTmpl := template.Must(template.New("method").Parse(methodtemplate))

	// Generate code for each method
	for _, method := range service.Handlers {
		// Execute the method template with the current method and write the output to the file
		err = methodTmpl.Execute(outputFile, method)
		if err != nil {
			log.Fatalf("Error executing method template: %s\n", err)
		}
	}

	log.Println("Handler code Generated successfully.")

}

type Handler struct {
	MethodName    string
	ServiceName   string
	IDLName       string
	RequestStruct string
}

type Service struct {
	IDLName     string
	GatewayName string
	HandlerFile string
	Handlers    []Handler
}
