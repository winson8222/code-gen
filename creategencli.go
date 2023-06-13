package main

import (
	"log"
	"os"
	"text/template"
)

func Creategencli(constants Constants) {
	// Generic client generation template
	genericClientTemplate := `
package {{ .IDLName }}

import (
	"bytes"
	"context"
	"{{ .GatewayName }}/constants"
	"log"
	"net/http"
	"strings"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/generic/descriptor"
	"github.com/cloudwego/kitex/pkg/klog"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func ToConstant(s string) string {
	return strings.ToUpper(strings.ReplaceAll(s, " ", "_"))
}


// Creates generic client "[ServiceName]GenericClient"
func {{ .ServiceName }}GenericClient() genericclient.Client {
	r, err := etcd.NewEtcdResolver([]string{constants.ETCD_URL})
	if err != nil {
		log.Fatal(err)
	}

	path := constants.FILEPATH_TO_{{ .ServiceName | ToConstant}}
	p, err := generic.NewThriftFileProvider(path)
	if err != nil {
		klog.Fatalf("new thrift file provider failed: %v", err)
	}
	g, err := generic.HTTPThriftGeneric(p)
	if err != nil {
		klog.Fatalf("new http thrift generic failed: %v", err)
	}

	cli, err := genericclient.NewClient(constants.{{ .ServiceName | ToConstant }}_NAME, g, client.WithResolver(r)) //should use dns resolver

	if err != nil {
		klog.Fatalf("new http generic client failed: %v", err)
	}
	return cli
}

`

	// Method template
	methodTemplate := `
func Do{{ .MethodName }}(cli genericclient.Client, c *app.RequestContext) (*descriptor.HTTPResponse, error) {
	req, err := http.NewRequest("{{ .MethodType }}", c.URI().String(), bytes.NewBuffer(c.Request.BodyBytes()))

	if err != nil {
		klog.Fatalf("new http request failed: %v", err)
	}

	customReq, err := generic.FromHTTPRequest(req)
	if err != nil {
		klog.Fatalf("convert request failed: %v", err)
	}

	resp, err := cli.GenericCall(context.Background(), "", customReq)
	if err != nil {
		klog.Fatalf("generic call failed: %v", err)
	}

	realResp := resp.(*generic.HTTPResponse)
	klog.Infof("{{ .MethodName }} response, status code: %v, headers: %v, body: %v\n",
		realResp.StatusCode, realResp.Header, realResp.Body)

	return realResp, err
}

`

	// Create the output file
	outputFile, err := os.Create("biz/handler/" + constants.IDLName + "/gen_client.go")
	if err != nil {
		log.Fatalf("Error creating output file: %s\n", err)
	}
	defer outputFile.Close()

	// Create a new template for the generic client
	clientTmpl := template.Must(template.New("genericClient").Funcs(template.FuncMap{"ToConstant": ToConstant}).Parse(genericClientTemplate))

	// Execute the generic client template with the constants and write the output to the file
	err = clientTmpl.Execute(outputFile, constants)
	if err != nil {
		log.Fatalf("Error executing generic client template: %s\n", err)
	}

	// Create a new template for the method
	methodTmpl := template.Must(template.New("method").Parse(methodTemplate))

	// Generate code for each method
	for _, method := range constants.Methods {
		// Execute the method template with the current method and write the output to the file
		err = methodTmpl.Execute(outputFile, method)
		if err != nil {
			log.Fatalf("Error executing method template: %s\n", err)
		}
	}

	log.Println("Generic Client code Generated successfully.")
}
