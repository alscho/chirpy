package main

import (
	"net/http"
	//"fmt"
	"log"
)

func main(){

	const port = "8080"

	mux := http.NewServeMux()

	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())

	/*
	err := server.ListenAndServe()
	if err != nil{
		fmt.Println("ListenAndServe failed: %v", err)
	}
	defer server.Close()
	*/
}