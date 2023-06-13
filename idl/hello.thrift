// idl/hello.thrift
namespace go api

struct HelloReq {
    1: string Name (api.body="name"); // Add api annotations for easier parameter binding
}

struct HelloResp {
    1: string RespBody;
}

service HelloService {
    HelloResp HelloMethod(1: HelloReq request) (api.post="/hello");
}
