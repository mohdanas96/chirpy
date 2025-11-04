package main

import (
	"fmt"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()

	server := http.Server{Addr: ":8080"}
	server.Handler = serveMux

	fmt.Println("server is listening on http://localhost:8080")
	server.ListenAndServe()

}
