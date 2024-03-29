package main

import (
	"net/http"

	"github.com/JamesClonk/elephantsql-broker/broker"
	"github.com/JamesClonk/elephantsql-broker/config"
	"github.com/JamesClonk/elephantsql-broker/env"
	"github.com/JamesClonk/elephantsql-broker/log"
)

func main() {
	port := env.Get("PORT", "8080")

	log.Infoln("port:", port)
	log.Infoln("log level:", config.Get().LogLevel)
	log.Infoln("broker username:", config.Get().Username)
	log.Infoln("api url:", config.Get().API.URL)
	log.Infoln("api default region:", config.Get().API.DefaultRegion)

	// start listener
	log.Fatalln(http.ListenAndServe(":"+port, broker.NewRouter(config.Get())))
}
