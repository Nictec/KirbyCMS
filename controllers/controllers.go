package controllers

import (
	"KirbyCMS/shortcuts"
	"KirbyCMS/helpers/youtube"
	"net/http"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
)

const basefile string = "tmpl/base.tmpl" 

func Index(w http.ResponseWriter, request *http.Request) {
	//logged_in := auth.LoginCheck(store, request)
	logged_in := true
	client := youtube.GetClient("https://www.googleapis.com/auth/youtube", w, request)
	youtube.Upload(w, client)
	if logged_in {
		// type Result struct {
		// 	Users []models.User
		// }
		// var model []models.User

		//TODO: add youtube analytics data

		shortcuts.Render(w, "tmpl/index.tmpl", basefile, nil)
		return
	} else {
		http.Redirect(w, request, "/admin/login", 301)
		return
	}
} 


func OauthEndpoint(w http.ResponseWriter, request *http.Request){
	b, err := ioutil.ReadFile("client_secret.json")
    if err != nil {
        shortcuts.HttpError(w, 500, "could not read client_secret.json")
        return
    }
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/youtube")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	code := request.FormValue("code")
	tok := youtube.ExchangeToken(config, code)
	cacheFile, err := youtube.TokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	youtube.SaveToken(cacheFile, tok)
	http.Redirect(w, request, "/", 301)
}