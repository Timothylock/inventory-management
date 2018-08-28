package main

import (
	"net/http"
	"log"

	"inventory-management/service"
)

func main() {
	//cfg, err := config.FromEnvironment()
	//if err != nil {
	//	fmt.Printf("error initializing config %s", err.Error())
	//	os.Exit(1)
	//}

	api := service.API{}
	router := service.NewRouter(&api)

	log.Fatal(http.ListenAndServe(":9090", router))
}