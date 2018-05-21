package urls 

import (
	"github.com/gorilla/mux"
	"KirbyCMS/controllers"
	) 

func Router(router *mux.Router){
	router.HandleFunc("/", controllers.Index)
	router.HandleFunc("/oauth", controllers.OauthEndpoint)
	router.HandleFunc("/new-video", controllers.NewVideo)
}