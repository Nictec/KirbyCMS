package main

import (
	"github.com/gorilla/mux"
	"KirbyCMS/helpers/db"
	"KirbyCMS/urls"
	"net/http"
	"fmt"
) 

const DEBUG = true
const STATIC = "bower_components/"
const MEDIA = "media/"
const PRODUCTION_HOST = "KirbyCMS"

func main(){
	//setup
	db.Init()
	defer db.Session.Close()

	//router
	router := mux.NewRouter()
	urls.Router(router)

	//staticfile server
	static := http.FileServer(http.Dir(STATIC))
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", static))
	media := http.FileServer(http.Dir(MEDIA))
	router.PathPrefix("/media/").Handler(http.StripPrefix("/media/", media)) 

	//check if server is in Debug mode and start listening
	if DEBUG == true {
		fmt.Println("EasyRent server is running on http://localhost:8000 press CTRL-C to quit")
		fmt.Println("Warning: Debug mode is enabled")
		http.ListenAndServe(":8000", router)
	} else {
		err := http.ListenAndServe(PRODUCTION_HOST+":80", router)
		if err != nil {
			panic("this port is already taken!")
		}
	}
}
