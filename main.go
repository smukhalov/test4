package main

import (
	"log"
	"net/http"
	"time"
	"github.com/urfave/negroni"
	"github.com/smukhalov/test4/common"
	"github.com/smukhalov/test4/routers"
)

//Hello gihub.com
func main() {

	//common.StartUp() - Replaced with init method
	// Get the mux router object
	router := routers.InitRoutes()
	// Create a negroni instance
	n := negroni.Classic()
	n.UseHandler(router)

	server := &http.Server{
		Addr:    common.AppConfig.Server,
		Handler: n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	log.Println("Listening...")
	server.ListenAndServe()
}
