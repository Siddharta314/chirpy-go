package main

import "net/http"

func main(){
	serve := http.ServeMux{}
	server := http.Server{
		Addr:    ":8080",
		Handler: &serve,
	}
	server.ListenAndServe()
}
