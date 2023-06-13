package main

import (
	"net/http"
)

func main() {
	//GatewayURL (customised)
	//FilepathtoService (fixed to ..idl/xxx.thrift)
	//EtcdURL (customised)
	//Servicename: genericlient method name, constants. (IDL)
	//HandlerFile: Service name but to lowercase (IDL file name with _service)
	//ServiceUpstreamURL (customised)
	//Methods: array of of methods of the service (IDL)
	//MethodName: Name of method, used in handler, and generic call (IDL)
	//MethodType: req type, used in generic call (IDL)
	//IDLName: Package name, import address, handler folder name (get from IDL namespac)
	//GatewayName: import address, hz gen command (customised)
	//RequestStruct: used in handler (IDL)
	gatewayURL := "0.0.0.0:8888"
	filepathtoservice := "..idl/hello.thrift"
	etcdurl := "0.0.0.0:20000"
	servicename := "HelloService"
	handlerfile := "hello"
	serviceupstreamurl := "0.0.0.0:8000"
	methodname := "HelloMethod"
	methodtype := http.MethodPost
	idlname := "api"
	gatewayname := "example"
	requeststruct := "HelloReq"

	Hzinstall()
	Hzgen(gatewayname)
	method := Method{MethodName: methodname, MethodType: methodtype}
	methods := []Method{method}
	constants := Constants{
		GatewayURL:         gatewayURL,
		FilepathToService:  filepathtoservice,
		EtcdURL:            etcdurl,
		ServiceName:        servicename, //to be changed into Hello
		ServiceUpstreamURL: serviceupstreamurl,
		Methods:            methods,
		IDLName:            idlname,
		GatewayName:        gatewayname,
	}

	handler := Handler{
		MethodName:    methodname,
		ServiceName:   servicename,
		IDLName:       idlname,
		RequestStruct: requeststruct,
	}

	handlers := []Handler{handler}

	services := Service{
		IDLName:     idlname,
		GatewayName: gatewayname,
		HandlerFile: handlerfile,
		Handlers:    handlers,
	}

	CreateConstant(constants)
	Creategencli(constants)
	Createhandler(services)
	Tidy(gatewayname)
}
