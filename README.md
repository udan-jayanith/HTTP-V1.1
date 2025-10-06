# HTTPNav

HTTPNav is a lightweight HTTP/1.1 Go package designed to simplify HTTP request and response handling. It provides custom encoders and decoders to streamline working with HTTP headers and payloads.

Note: HTTPNav does not support HTTPS.

## Installation

To install HTTPNav, use go get:
```bash
go get github.com/udan-jayanith/HTTPNav
```
Then, import it into your Go code:
```go
import HTTPNav "github.com/udan-jayanith/HTTPNav"
```

## Usage

### Server
server contains server data and logic.

#### GetServer
```go
HTTPNav.GetServer()
```
GetServer() Function returns a new server.

#### HandelFunc
HandelFunc takes HTTPRequestMethod, requestTarget and a handler. If requests httpMethod and requestTarget matches the handler handler will execute.

Ex:
```go
server.HandleFunc(httpNav.Get, "/", func(response *httpNav.HTTPResponse, request *httpNav.HTTPRequest) {
		response.Write([]byte("Hello world. /"))
        //Write send a response back to the client.
})
```

#### StartServer
```go
StartServer(":8080")
```
StartServer start the server(starts listing to incoming requests). 

### HTTPRequest

Parses the HTTP request and store in it. Then pointer to the HTTPRequest is passed to matching HandleFunc()

#### GetBodyAsBytes
This return Body as a slice of bytes. If Content-Length header is not included in the header GetBodyAsBytes() returns ContentLengthHeaderNotFound.

```go
bodyBty, err := request.GetBodyAsBytes()
```
#### GetBodyAsJson
```go
var v map[string]any
err := request.GetBodyAsJson(&v)
```
GetBodyAsJson parses the JSON-encoded data and stores the result in the value pointed to by v. If Content-Type is not in the header GetBodyAsJson returns ContentTypeHeaderNotFound. Else if Content-Type != "application/json" GetBodyAsJson returns InvalidContentType. It uses GetBodyAsBytes in underlying layer.

**To see other functionalities see
[HTTP Go Doc](https://pkg.go.dev/github.com/udan-jayanith/HTTPNav) or simply see source code.**

## License
This project is licensed under the MIT License. See the [LICENSE.md](https://github.com/udan-jayanith/HTTPNav/blob/main/LICENSE.md) file for details.
