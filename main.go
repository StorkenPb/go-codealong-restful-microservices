package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r*http.Request) {
		d, err := io.ReadAll(r.Body)

		if(err != nil) {
			http.Error(rw, "oops", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(rw, "Hey, this was your data: %s", d)
	})

	http.ListenAndServe(":9090", nil)
}