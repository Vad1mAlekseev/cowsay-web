package main

import (
	"cowsay-web/apiserver"
	"log"
)

func main() {
	s := apiserver.ApiServer{}
	s.SetupDefaultRoutes()
	if err := s.ListenAndServe(":8080"); err != nil {
		log.Fatal(err)
	}
}
