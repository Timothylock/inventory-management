package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Timothylock/inventory-management/config"
	"github.com/Timothylock/inventory-management/email"
	"github.com/Timothylock/inventory-management/items"
	"github.com/Timothylock/inventory-management/persistence"
	"github.com/Timothylock/inventory-management/service"
	"github.com/Timothylock/inventory-management/upc"
	"github.com/Timothylock/inventory-management/users"
	"gopkg.in/gomail.v2"
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

	emailDialer := gomail.NewPlainDialer(cfg.EmailSmtpServ, cfg.EmailSmtpPort, cfg.EmailUsername, cfg.EmailPassword)

	is := items.NewService(persister)
	us := upc.NewService(*cfg)
	user := users.NewService(persister)
	es := email.NewService(*cfg, emailDialer)

	api := service.NewAPI(is, us, user, es)

	router := service.NewRouter(&api, *cfg)

	log.Fatal(http.ListenAndServe(":9090", router))
}
