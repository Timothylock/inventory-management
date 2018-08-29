package main

import (
	"net/http"
	"log"
	"fmt"
	"os"

	"github.com/Timothylock/inventory-management/service"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/persistence"
)

func main() {
	cfg, err := config.FromEnvironment()
	if err != nil {
		fmt.Printf("error initializing config %s", err.Error())
		os.Exit(1)
	}

	persister, err := persistence.NewMySQL(cfg)
	if err != nil {
		fmt.Printf("error initializing database %s", err.Error())
		os.Exit(1)
	}


	is := items.NewService(persister)
	api := service.NewAPI(is)
	router := service.NewRouter(&api)

	log.Fatal(http.ListenAndServe(":9090", router))
}