package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/persistence"
	"github.com/Timothylock/inventory-management/service"
	"github.com/Timothylock/inventory-management/upc"
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
	us := upc.NewService(*cfg)
	api := service.NewAPI(is, us)
	router := service.NewRouter(&api)

	log.Fatal(http.ListenAndServe(":9090", router))
}
