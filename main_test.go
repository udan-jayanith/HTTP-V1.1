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
		var jsonMap map[string]any
		err := b.GetBodyAsJson(&jsonMap)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(jsonMap)
	})
	fmt.Println("Running")
	server.StartServer(":8080")
}