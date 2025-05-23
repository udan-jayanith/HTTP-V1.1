package HTTPNav_test

import (
	httpNav "HTTPNav"
	"fmt"
	"testing"
)

func TestMain(m *testing.M){
	server := httpNav.GetServer()
	server.HandelFunc(httpNav.Get, "/", func() {
		fmt.Println("Hello world from handler")
	})
	fmt.Println("Running")
	server.StartServer(":8080")
}