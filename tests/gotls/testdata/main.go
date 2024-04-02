package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

/*
Test File

	you can use `go build` this file
		with ARCH(amd64/arm64) and BUILDMODE(normal/pie)

	`GOOS=(linux/android) GOARCH=(amd64/arm64) go build -buildmode=xxx -o xxx/main main.go`

	when build script runing:
		you will see function `"crypto/tls.(\*Conn).Read$"` address
		and offset address of return instruction of this function.
		this offset will be useful in test case.
*/
func main() {
	rsp, err := http.Get("https://com.example.com")
	if err != nil {
		log.Fatal("http.Get error: ", err)
	}
	defer rsp.Body.Close()
	body, _ := io.ReadAll(rsp.Body)
	fmt.Printf("http.Get: \n\t%s\n", body)
}
