package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/leyafo/purr"
)

func demoHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, client")
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "standalone" {
		http.HandleFunc("/", demoHandle)
		go func() {
			http.ListenAndServe(":9527", nil)
		}()
		purr.RunTest("http://127.0.0.1:9527/", "./")
	} else {
		purr.RunTestWithServer(http.HandlerFunc(demoHandle), "./")
	}
}
