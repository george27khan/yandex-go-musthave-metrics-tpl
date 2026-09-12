package main

import "net/http"

func main() {
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		panic(err)
	}
}
