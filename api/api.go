package main

import (
	"fmt"
	"net/http"
)

const (
	Port = ":7880"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", DoIt)

	fmt.Println("listen port:", Port)

	err := http.ListenAndServe(Port, mux)
	if err != nil {
		panic(err)
	}
}

func DoIt(rw http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	fmt.Println(r.PostForm)
}
