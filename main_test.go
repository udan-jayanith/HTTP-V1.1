package HTTPNav_test

import (
	httpNav "HTTPNav"
	"fmt"
	"testing"
)

func TestMain(m *testing.M){
	server := httpNav.GetServer()
	server.HandleFunc(httpNav.Get, "/", func(a string, b *httpNav.HTTPRequest) {
		fmt.Println("Hello world from handler")
		fmt.Println(b)
	})
	fmt.Println("Running")
	server.StartServer(":8080")
}