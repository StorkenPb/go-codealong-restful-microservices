package main

import "net/http"

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	http.HandleFunc("/goodbye", goodbye)

	http.ListenAndServe(":9090", nil)
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("goodbye"))
}