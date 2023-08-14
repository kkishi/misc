package main

import (
	"fmt"
	"net/http"
	"os"
)

type handler struct {
}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%v\n", *r)
	fmt.Fprintln(w, "Hello, world!")
}

func run() error {
	s := &http.Server{
		Addr:    ":8888",
		Handler: handler{},
	}
	return s.ListenAndServe()
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
