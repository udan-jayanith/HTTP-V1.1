package HTTPNav_test

import (
	httpNav "github.com/udan-jayanith/HTTPNav"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	server := httpNav.GetServer()
	server.HandleFunc(httpNav.Get, "/", func(response *httpNav.HTTPResponse, request *httpNav.HTTPRequest) {
		response.Write([]byte("Hello world. /"))
	})

	server.HandleFunc(httpNav.Delete, "/delete", func(response *httpNav.HTTPResponse, request *httpNav.HTTPRequest) {
		response.Write([]byte("Hello world. /"))
	})

	fmt.Println("Running")
	server.StartServer(":8080")
}
